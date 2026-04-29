package services

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/database"
	"github.com/rpupo63/faradhaven/backend/models"
	"gorm.io/gorm"
)

type LootService struct {
	db            *gorm.DB
	itemRepo      database.ItemRepository
	weaponRepo    database.WeaponRepository
	characterRepo database.CharacterRepository
	pendingMu     sync.Mutex
	pendingLoot   map[uuid.UUID]*pendingLootSession
}

// LootDrop records one resolved drop for API/UI display.
type LootDrop struct {
	Kind   string `json:"kind"` // "item" or "weapon"
	Name   string `json:"name"`
	Rarity string `json:"rarity"`
}

type LootResult struct {
	Items          []models.Item   `json:"items"`
	Weapons        []models.Weapon `json:"weapons"`
	GoldEarned     int64           `json:"gold_earned"`
	TotalMoney     int64           `json:"total_money"`
	Message        string          `json:"message"`
	ItemsRolled    int             `json:"items_rolled"`
	WeaponsRolled  int             `json:"weapons_rolled"`
	Drops          []LootDrop      `json:"drops"`
	RoomTheme      string          `json:"room_theme"`
	RewardAmount   string          `json:"reward_amount"`
	LevelBand      string          `json:"level_band"`
	ProfileNotes   []string        `json:"profile_notes"`
	ExpectedBudget int64           `json:"expected_budget"`
	SessionBudget  int64           `json:"session_budget"`
	EndingBudget   int64           `json:"ending_budget"`
	DebtUsed       bool            `json:"debt_used"`
	DebtAmount     int64           `json:"debt_amount"`
}

type LootPartyMember struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type LootPreview struct {
	SessionID    uuid.UUID
	Loot         *LootResult
	PartyMembers []LootPartyMember
}

type pendingLootSession struct {
	SessionID      uuid.UUID
	SourceCharacterID uuid.UUID
	PartyID        *uuid.UUID
	CreatedAt      time.Time
	Result         *LootResult
	Chosen         []lootCandidate
}

func NewLootService(db *gorm.DB, itemRepo database.ItemRepository, weaponRepo database.WeaponRepository, characterRepo database.CharacterRepository) *LootService {
	return &LootService{
		db: db, itemRepo: itemRepo, weaponRepo: weaponRepo, characterRepo: characterRepo,
		pendingLoot: make(map[uuid.UUID]*pendingLootSession),
	}
}

// lootParams defines the count ranges and gold for a (source, tier) combo.
type lootParams struct {
	ItemMin, ItemMax     int
	WeaponMin, WeaponMax int
	GoldMin, GoldMax     int64
}

type lootModifier struct {
	ItemDelta      int
	WeaponDelta    int
	GoldMultiplier float64
	RarityShift    int
}

// rarityWeights maps each rarity to its weight for a given tier (Legendary is never used for loot).
type rarityWeights struct {
	Common    int
	Uncommon  int
	Rare      int
	VeryRare  int
	Legendary int
}

var lootTable = map[string]map[string]lootParams{
	"common_enemy": {
		"low":    {0, 2, 0, 1, 50, 200},
		"medium": {1, 3, 0, 2, 200, 500},
		"high":   {2, 4, 1, 2, 500, 1000},
	},
	"boss_enemy": {
		"low":    {1, 3, 0, 2, 500, 1000},
		"medium": {2, 4, 1, 2, 1000, 3000},
		"high":   {3, 5, 1, 3, 3000, 8000},
	},
	"room": {
		"low":    {0, 2, 0, 1, 100, 300},
		"medium": {1, 3, 0, 1, 300, 800},
		"high":   {2, 4, 1, 2, 800, 2000},
	},
}

// Legendary weight is always zero; former Legendary weight is folded into Very Rare.
var tierWeights = map[string]rarityWeights{
	"low":    {70, 25, 4, 1, 0},
	"medium": {45, 35, 15, 5, 0},
	"high":   {20, 30, 30, 20, 0},
}

var themeModifiers = map[string]lootModifier{
	"dungeon":    {ItemDelta: 1, WeaponDelta: 0, GoldMultiplier: 1.05, RarityShift: 0},
	"office":     {ItemDelta: 1, WeaponDelta: -1, GoldMultiplier: 0.95, RarityShift: -2},
	"rich":       {ItemDelta: 1, WeaponDelta: 1, GoldMultiplier: 1.35, RarityShift: 5},
	"poor":       {ItemDelta: 0, WeaponDelta: -1, GoldMultiplier: 0.75, RarityShift: -4},
	"gangster":   {ItemDelta: 0, WeaponDelta: 1, GoldMultiplier: 1.15, RarityShift: 2},
	"arcane":     {ItemDelta: 2, WeaponDelta: 0, GoldMultiplier: 1.1, RarityShift: 3},
	"wilderness": {ItemDelta: 1, WeaponDelta: 1, GoldMultiplier: 1.0, RarityShift: 1},
}

var rewardModifiers = map[string]lootModifier{
	"scarce":    {ItemDelta: -1, WeaponDelta: -1, GoldMultiplier: 0.7, RarityShift: -5},
	"standard":  {ItemDelta: 0, WeaponDelta: 0, GoldMultiplier: 1.0, RarityShift: 0},
	"bountiful": {ItemDelta: 1, WeaponDelta: 1, GoldMultiplier: 1.25, RarityShift: 4},
	"jackpot":   {ItemDelta: 2, WeaponDelta: 1, GoldMultiplier: 1.6, RarityShift: 8},
}

var levelBandModifiers = map[string]lootModifier{
	"novice":     {ItemDelta: 0, WeaponDelta: 0, GoldMultiplier: 0.95, RarityShift: -2},
	"adventurer": {ItemDelta: 0, WeaponDelta: 0, GoldMultiplier: 1.0, RarityShift: 0},
	"veteran":    {ItemDelta: 1, WeaponDelta: 1, GoldMultiplier: 1.1, RarityShift: 3},
	"legend":     {ItemDelta: 1, WeaponDelta: 1, GoldMultiplier: 1.2, RarityShift: 5},
}

var ErrInvalidLootOption = errors.New("invalid loot option")
var ErrInvalidCostString = errors.New("invalid cost string")

const (
	budgetSigmaRatio      = 0.22
	minDebtAllowanceCP    = int64(50)
	debtAllowanceRatio    = 0.2
	maxBudgetFillAttempts = 400
	minSpendFloorRatio    = 0.08
	minSpendFloorCP       = int64(50)
	profileSpendFloorRate = 0.045
)

type lootCandidate struct {
	Kind       string
	Item       *models.Item
	Weapon     *models.Weapon
	Name       string
	Rarity     string
	Category   string
	PriceCP    int64
	BaseWeight float64
}

type budgetFillOutcome struct {
	Chosen       []lootCandidate
	RemainingCP  int64
	DebtUsed     bool
	DebtAmountCP int64
}

// applyLevelShift shifts weight from Common to Uncommon, Rare, and Very Rare every 5 levels above 1.
// Legendary never receives shifted weight (loot never drops Legendary items).
func applyLevelShift(w rarityWeights, level int) rarityWeights {
	w.Legendary = 0
	shifts := (level - 1) / 5
	if shifts <= 0 {
		return w
	}
	totalShift := shifts * 5
	if totalShift > w.Common {
		totalShift = w.Common
	}
	w.Common -= totalShift
	perBucket := totalShift / 3
	remainder := totalShift % 3
	w.Uncommon += perBucket
	w.Rare += perBucket
	w.VeryRare += perBucket + remainder
	return w
}

func rollRarity(w rarityWeights) string {
	w.Legendary = 0
	total := w.Common + w.Uncommon + w.Rare + w.VeryRare
	if total <= 0 {
		return "Common"
	}
	roll := rand.Intn(total)
	if roll < w.Common {
		return "Common"
	}
	roll -= w.Common
	if roll < w.Uncommon {
		return "Uncommon"
	}
	roll -= w.Uncommon
	if roll < w.Rare {
		return "Rare"
	}
	roll -= w.Rare
	return "Very Rare"
}

func randRange(min, max int) int {
	if min >= max {
		return min
	}
	return min + rand.Intn(max-min+1)
}

func randRange64(min, max int64) int64 {
	if min >= max {
		return min
	}
	return min + rand.Int63n(max-min+1)
}

func parseCostToCopper(cost string) (int64, error) {
	s := strings.TrimSpace(strings.ToLower(cost))
	if s == "" {
		return 0, fmt.Errorf("%w: empty", ErrInvalidCostString)
	}
	parts := strings.Fields(s)
	if len(parts) < 2 {
		return 0, fmt.Errorf("%w: %q", ErrInvalidCostString, cost)
	}
	value, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", ErrInvalidCostString, cost)
	}
	unit := parts[1]
	switch unit {
	case "cp":
		return int64(math.Round(value)), nil
	case "sp":
		return int64(math.Round(value * 10)), nil
	case "gp":
		return int64(math.Round(value * 100)), nil
	case "pp":
		return int64(math.Round(value * 1000)), nil
	default:
		return 0, fmt.Errorf("%w: unsupported unit %q", ErrInvalidCostString, unit)
	}
}

func rarityScarcityWeight(rarity string) float64 {
	switch strings.ToLower(strings.TrimSpace(rarity)) {
	case "very rare":
		return 6.0
	case "rare":
		return 3.2
	case "uncommon":
		return 1.8
	default:
		return 1.0
	}
}

func sampleGaussianBudget(rng *rand.Rand, expected int64) int64 {
	if expected <= 0 {
		return 0
	}
	sigma := float64(expected) * budgetSigmaRatio
	rolled := float64(expected) + rng.NormFloat64()*sigma
	if rolled < float64(expected)/4 {
		rolled = float64(expected) / 4
	}
	return int64(math.Round(rolled))
}

func debtAllowanceForRemaining(remaining int64) int64 {
	scaled := int64(math.Round(float64(remaining) * debtAllowanceRatio))
	if scaled < minDebtAllowanceCP {
		return minDebtAllowanceCP
	}
	return scaled
}

func weightedPickIndex(rng *rand.Rand, candidates []lootCandidate) int {
	if len(candidates) == 0 {
		return -1
	}
	total := 0.0
	for i := range candidates {
		if candidates[i].BaseWeight > 0 {
			total += candidates[i].BaseWeight
		}
	}
	if total <= 0 {
		return rng.Intn(len(candidates))
	}
	r := rng.Float64() * total
	acc := 0.0
	for i := range candidates {
		w := candidates[i].BaseWeight
		if w <= 0 {
			continue
		}
		acc += w
		if r <= acc {
			return i
		}
	}
	return len(candidates) - 1
}

func spendFloorForRemaining(remaining, expectedBudget int64, picksSoFar int) int64 {
	if remaining <= 0 {
		return 0
	}
	scaled := int64(math.Round(float64(remaining) * minSpendFloorRatio))
	profileScaled := int64(math.Round(float64(expectedBudget) * profileSpendFloorRate))
	if profileScaled > scaled {
		scaled = profileScaled
	}
	// Gentle escalation avoids long tails of tiny purchases.
	if picksSoFar > 0 {
		escalation := int64(math.Round(float64(expectedBudget) * 0.01 * float64(picksSoFar)))
		if escalation > 0 {
			scaled += escalation
		}
	}
	if scaled < minSpendFloorCP {
		scaled = minSpendFloorCP
	}
	if scaled > remaining {
		return remaining
	}
	return scaled
}

func rollBudgetFromProfile(rng *rand.Rand, profile models.LootLevelBudgetProfile) (sessionBudgetCP int64, expectedBudgetCP int64) {
	expectedBudgetCP = int64(math.Round(profile.MeanGP * 100))
	sigmaCP := profile.SigmaGP * 100
	rolled := float64(expectedBudgetCP) + rng.NormFloat64()*sigmaCP
	minCP := int64(math.Round(profile.MinGP * 100))
	maxCP := int64(math.Round(profile.MaxGP * 100))
	sessionBudgetCP = int64(math.Round(rolled))
	if sessionBudgetCP < minCP {
		sessionBudgetCP = minCP
	}
	if maxCP > 0 && sessionBudgetCP > maxCP {
		sessionBudgetCP = maxCP
	}
	if sessionBudgetCP < 0 {
		sessionBudgetCP = 0
	}
	return sessionBudgetCP, expectedBudgetCP
}

func selectByRarity(candidates []lootCandidate, rarity string) []lootCandidate {
	filtered := make([]lootCandidate, 0)
	for _, c := range candidates {
		if strings.EqualFold(c.Rarity, rarity) {
			filtered = append(filtered, c)
		}
	}
	if len(filtered) > 0 {
		return filtered
	}
	for _, c := range candidates {
		if strings.EqualFold(c.Rarity, "Common") {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

func runBudgetFillLoop(rng *rand.Rand, candidates []lootCandidate, weights rarityWeights, startingBudget int64, expectedBudget int64) budgetFillOutcome {
	out := budgetFillOutcome{
		Chosen:      make([]lootCandidate, 0),
		RemainingCP: startingBudget,
	}
	if startingBudget <= 0 || len(candidates) == 0 {
		return out
	}
	available := make([]lootCandidate, len(candidates))
	copy(available, candidates)

	for attempts := 0; out.RemainingCP > 0 && attempts < maxBudgetFillAttempts; attempts++ {
		minSpend := spendFloorForRemaining(out.RemainingCP, expectedBudget, len(out.Chosen))
		affordable := make([]lootCandidate, 0)
		affordableIndexes := make([]int, 0)
		for i := range available {
			c := available[i]
			if c.PriceCP <= out.RemainingCP && c.PriceCP >= minSpend {
				affordable = append(affordable, c)
				affordableIndexes = append(affordableIndexes, i)
			}
		}
		if len(affordable) == 0 {
			break
		}
		pickedAffordableIndex := weightedPickIndex(rng, affordable)
		if pickedAffordableIndex < 0 {
			break
		}
		pickedIndex := affordableIndexes[pickedAffordableIndex]
		picked := available[pickedIndex]
		out.Chosen = append(out.Chosen, picked)
		out.RemainingCP -= picked.PriceCP
		available = append(available[:pickedIndex], available[pickedIndex+1:]...)
	}

	return out
}

func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func deriveLevelBand(level int) string {
	switch {
	case level <= 4:
		return "novice"
	case level <= 9:
		return "adventurer"
	case level <= 14:
		return "veteran"
	default:
		return "legend"
	}
}

func shiftedWeights(w rarityWeights, shift int) rarityWeights {
	if shift == 0 {
		return w
	}
	w.Legendary = 0
	if shift > 0 {
		move := clampInt(shift, 0, w.Common)
		w.Common -= move
		w.Uncommon += move / 2
		w.Rare += move / 3
		w.VeryRare += move - ((move / 2) + (move / 3))
		return w
	}

	move := clampInt(-shift, 0, w.Uncommon+w.Rare+w.VeryRare)
	fromVeryRare := clampInt(move/2, 0, w.VeryRare)
	w.VeryRare -= fromVeryRare
	w.Common += fromVeryRare

	remaining := move - fromVeryRare
	fromRare := clampInt(remaining/2, 0, w.Rare)
	w.Rare -= fromRare
	w.Common += fromRare

	remaining -= fromRare
	fromUncommon := clampInt(remaining, 0, w.Uncommon)
	w.Uncommon -= fromUncommon
	w.Common += fromUncommon

	return w
}

func applyLootModifiers(base lootParams, mods ...lootModifier) (lootParams, int, float64) {
	adjusted := base
	totalShift := 0
	goldMultiplier := 1.0
	for _, mod := range mods {
		adjusted.ItemMin = clampInt(adjusted.ItemMin+mod.ItemDelta, 0, 10)
		adjusted.ItemMax = clampInt(adjusted.ItemMax+mod.ItemDelta, adjusted.ItemMin, 12)
		adjusted.WeaponMin = clampInt(adjusted.WeaponMin+mod.WeaponDelta, 0, 10)
		adjusted.WeaponMax = clampInt(adjusted.WeaponMax+mod.WeaponDelta, adjusted.WeaponMin, 12)
		totalShift += mod.RarityShift
		if mod.GoldMultiplier > 0 {
			goldMultiplier *= mod.GoldMultiplier
		}
	}
	adjusted.GoldMin = int64(math.Max(0, math.Round(float64(adjusted.GoldMin)*goldMultiplier)))
	adjusted.GoldMax = int64(math.Max(float64(adjusted.GoldMin), math.Round(float64(adjusted.GoldMax)*goldMultiplier)))
	return adjusted, totalShift, goldMultiplier
}

func validateLootInputs(source, roomTheme string) error {
	if _, ok := models.ParseLootSource(source); !ok {
		return fmt.Errorf("%w: source %q", ErrInvalidLootOption, source)
	}
	if _, ok := models.ParseLootTheme(roomTheme); !ok {
		return fmt.Errorf("%w: room_theme %q", ErrInvalidLootOption, roomTheme)
	}
	return nil
}

func validateLootLocation(location string) error {
	if _, ok := models.ParseLootLocation(location); !ok {
		return fmt.Errorf("%w: location %q", ErrInvalidLootOption, location)
	}
	return nil
}

func toWeightMap[K comparable](rows []struct {
	ID     K
	Key    string
	Weight float64
}) map[K]map[string]float64 {
	out := make(map[K]map[string]float64)
	for _, row := range rows {
		if out[row.ID] == nil {
			out[row.ID] = make(map[string]float64)
		}
		if row.Weight <= 0 {
			out[row.ID][row.Key] = 1.0
		} else {
			out[row.ID][row.Key] = row.Weight
		}
	}
	return out
}

func dimensionWeight[K comparable](m map[K]map[string]float64, id K, selected string) (float64, bool) {
	dim, exists := m[id]
	if !exists || len(dim) == 0 {
		return 1.0, true
	}
	w, ok := dim[selected]
	if !ok {
		return 0, false
	}
	if w <= 0 {
		return 1.0, true
	}
	return w, true
}

func isLegendaryRarity(r string) bool {
	return strings.EqualFold(strings.TrimSpace(r), "Legendary")
}

// pickRarityCategoryMap returns the category→pool map for rarity, falling back to Common if empty.
func pickRarityCategoryMap[T any](byRarity map[string]map[string][]T, rarity string) map[string][]T {
	if cm, ok := byRarity[rarity]; ok && categoryMapHasItems(cm) {
		return cm
	}
	if cm, ok := byRarity["Common"]; ok && categoryMapHasItems(cm) {
		return cm
	}
	return nil
}

func categoryMapHasItems[T any](cm map[string][]T) bool {
	for _, items := range cm {
		if len(items) > 0 {
			return true
		}
	}
	return false
}

// pickStratified chooses a category uniformly among categories that have items, then an item in that category.
func pickStratifiedItem(byRarity map[string]map[string][]*models.Item, rarity string) *models.Item {
	catMap := pickRarityCategoryMap(byRarity, rarity)
	if catMap == nil {
		return nil
	}
	categories := make([]string, 0)
	for cat, pool := range catMap {
		if len(pool) > 0 {
			categories = append(categories, cat)
		}
	}
	if len(categories) == 0 {
		return nil
	}
	cat := categories[rand.Intn(len(categories))]
	pool := catMap[cat]
	return pool[rand.Intn(len(pool))]
}

func pickStratifiedWeapon(byRarity map[string]map[string][]models.Weapon, rarity string) *models.Weapon {
	catMap := pickRarityCategoryMap(byRarity, rarity)
	if catMap == nil {
		return nil
	}
	categories := make([]string, 0)
	for cat, pool := range catMap {
		if len(pool) > 0 {
			categories = append(categories, cat)
		}
	}
	if len(categories) == 0 {
		return nil
	}
	cat := categories[rand.Intn(len(categories))]
	pool := catMap[cat]
	chosen := pool[rand.Intn(len(pool))]
	return &chosen
}

func (s *LootService) generateLootResult(characterID uuid.UUID, source, roomTheme string, locationOverride *string, lootLevel int) (*LootResult, []lootCandidate, *uuid.UUID, []LootPartyMember, error) {
	rng := rand.New(rand.NewSource(rand.Int63()))
	character, err := s.characterRepo.FindByIDWithInventory(characterID)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to find character %s: %w", characterID, err)
	}
	if lootLevel < 1 || lootLevel > 20 {
		return nil, nil, nil, nil, fmt.Errorf("%w: loot_level %d must be between 1 and 20", ErrInvalidLootOption, lootLevel)
	}
	levelBand := deriveLevelBand(lootLevel)
	roomTheme = strings.ToLower(strings.TrimSpace(roomTheme))
	location := "indoor"
	if locationOverride != nil && strings.TrimSpace(*locationOverride) != "" {
		location = strings.ToLower(strings.TrimSpace(*locationOverride))
	}
	if err := validateLootLocation(location); err != nil {
		return nil, nil, nil, nil, err
	}

	if err := validateLootInputs(source, roomTheme); err != nil {
		return nil, nil, nil, nil, err
	}

	allItems, err := s.itemRepo.FindAll()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to fetch items: %w", err)
	}
	allWeapons, err := s.weaponRepo.FindAll()
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to fetch weapons: %w", err)
	}

	itemsByRarityCategory := make(map[string]map[string][]*models.Item)
	for _, item := range allItems {
		if item.UserID != nil || isLegendaryRarity(item.Rarity) {
			continue
		}
		r := item.Rarity
		if itemsByRarityCategory[r] == nil {
			itemsByRarityCategory[r] = make(map[string][]*models.Item)
		}
		itemsByRarityCategory[r][item.Category] = append(itemsByRarityCategory[r][item.Category], item)
	}

	weaponsByRarityCategory := make(map[string]map[string][]models.Weapon)
	for i := range allWeapons {
		w := allWeapons[i]
		if w.UserID != nil || isLegendaryRarity(w.Rarity) {
			continue
		}
		r := w.Rarity
		if weaponsByRarityCategory[r] == nil {
			weaponsByRarityCategory[r] = make(map[string][]models.Weapon)
		}
		weaponsByRarityCategory[r][w.Category] = append(weaponsByRarityCategory[r][w.Category], w)
	}
	profile := models.LootLevelBudgetProfile{}
	if err := s.db.Where("loot_level = ?", lootLevel).First(&profile).Error; err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to load loot level budget profile for level %d: %w", lootLevel, err)
	}
	sessionBudget, expectedBudget := rollBudgetFromProfile(rng, profile)

	candidates := make([]lootCandidate, 0, len(allItems)+len(allWeapons))
	for rarity, categoryMap := range itemsByRarityCategory {
		for category, pool := range categoryMap {
			if len(pool) == 0 {
				continue
			}
			scarcity := rarityScarcityWeight(rarity) * (1.0 + 1.0/float64(len(pool)))
			for _, item := range pool {
				price, perr := parseCostToCopper(item.Cost)
				if perr != nil || price <= 0 {
					continue
				}
				candidates = append(candidates, lootCandidate{
					Kind:       "item",
					Item:       item,
					Name:       item.Name,
					Rarity:     rarity,
					Category:   category,
					PriceCP:    price,
					BaseWeight: scarcity,
				})
			}
		}
	}
	for rarity, categoryMap := range weaponsByRarityCategory {
		for category, pool := range categoryMap {
			if len(pool) == 0 {
				continue
			}
			scarcity := rarityScarcityWeight(rarity) * (1.0 + 1.0/float64(len(pool)))
			for i := range pool {
				w := pool[i]
				price, perr := parseCostToCopper(w.Cost)
				if perr != nil || price <= 0 {
					continue
				}
				candidates = append(candidates, lootCandidate{
					Kind:       "weapon",
					Weapon:     &w,
					Name:       w.Name,
					Rarity:     rarity,
					Category:   category,
					PriceCP:    price,
					BaseWeight: scarcity,
				})
			}
		}
	}

	// Apply enum-backed relational mapping filters/weights when rows exist.
	itemIDs := make([]uuid.UUID, 0, len(allItems))
	for _, item := range allItems {
		itemIDs = append(itemIDs, item.ID)
	}
	weaponIDs := make([]uuid.UUID, 0, len(allWeapons))
	for i := range allWeapons {
		weaponIDs = append(weaponIDs, allWeapons[i].ID)
	}

	itemThemeRows := make([]struct {
		ID     uuid.UUID
		Key    string
		Weight float64
	}, 0)
	if err := s.db.Model(&models.ItemLootTheme{}).Where("item_id IN ?", itemIDs).Select("item_id as id, theme as key, weight").Scan(&itemThemeRows).Error; err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to load item loot themes: %w", err)
	}
	itemSourceRows := make([]struct {
		ID     uuid.UUID
		Key    string
		Weight float64
	}, 0)
	if err := s.db.Model(&models.ItemLootSource{}).Where("item_id IN ?", itemIDs).Select("item_id as id, source as key, weight").Scan(&itemSourceRows).Error; err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to load item loot sources: %w", err)
	}
	itemBandRows := make([]struct {
		ID     uuid.UUID
		Key    string
		Weight float64
	}, 0)
	if err := s.db.Model(&models.ItemLootLevelBand{}).Where("item_id IN ?", itemIDs).Select("item_id as id, level_band as key, weight").Scan(&itemBandRows).Error; err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to load item loot level bands: %w", err)
	}
	itemLocationRows := make([]struct {
		ID     uuid.UUID
		Key    string
		Weight float64
	}, 0)
	if err := s.db.Model(&models.ItemLootLocation{}).Where("item_id IN ?", itemIDs).Select("item_id as id, location as key, weight").Scan(&itemLocationRows).Error; err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to load item loot locations: %w", err)
	}

	weaponThemeRows := make([]struct {
		ID     uuid.UUID
		Key    string
		Weight float64
	}, 0)
	if err := s.db.Model(&models.WeaponLootTheme{}).Where("weapon_id IN ?", weaponIDs).Select("weapon_id as id, theme as key, weight").Scan(&weaponThemeRows).Error; err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to load weapon loot themes: %w", err)
	}
	weaponSourceRows := make([]struct {
		ID     uuid.UUID
		Key    string
		Weight float64
	}, 0)
	if err := s.db.Model(&models.WeaponLootSource{}).Where("weapon_id IN ?", weaponIDs).Select("weapon_id as id, source as key, weight").Scan(&weaponSourceRows).Error; err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to load weapon loot sources: %w", err)
	}
	weaponBandRows := make([]struct {
		ID     uuid.UUID
		Key    string
		Weight float64
	}, 0)
	if err := s.db.Model(&models.WeaponLootLevelBand{}).Where("weapon_id IN ?", weaponIDs).Select("weapon_id as id, level_band as key, weight").Scan(&weaponBandRows).Error; err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to load weapon loot level bands: %w", err)
	}
	weaponLocationRows := make([]struct {
		ID     uuid.UUID
		Key    string
		Weight float64
	}, 0)
	if err := s.db.Model(&models.WeaponLootLocation{}).Where("weapon_id IN ?", weaponIDs).Select("weapon_id as id, location as key, weight").Scan(&weaponLocationRows).Error; err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to load weapon loot locations: %w", err)
	}

	itemThemeMap := toWeightMap(itemThemeRows)
	itemSourceMap := toWeightMap(itemSourceRows)
	itemBandMap := toWeightMap(itemBandRows)
	itemLocationMap := toWeightMap(itemLocationRows)
	weaponThemeMap := toWeightMap(weaponThemeRows)
	weaponSourceMap := toWeightMap(weaponSourceRows)
	weaponBandMap := toWeightMap(weaponBandRows)
	weaponLocationMap := toWeightMap(weaponLocationRows)

	filteredCandidates := make([]lootCandidate, 0, len(candidates))
	for _, c := range candidates {
		weightScale := 1.0
		allowed := true
		if c.Kind == "item" && c.Item != nil {
			if w, ok := dimensionWeight(itemThemeMap, c.Item.ID, roomTheme); !ok {
				allowed = false
			} else {
				weightScale *= w
			}
			if w, ok := dimensionWeight(itemSourceMap, c.Item.ID, source); !ok {
				allowed = false
			} else {
				weightScale *= w
			}
			if w, ok := dimensionWeight(itemBandMap, c.Item.ID, levelBand); !ok {
				allowed = false
			} else {
				weightScale *= w
			}
			if w, ok := dimensionWeight(itemLocationMap, c.Item.ID, location); !ok {
				allowed = false
			} else {
				weightScale *= w
			}
		}
		if c.Kind == "weapon" && c.Weapon != nil {
			if w, ok := dimensionWeight(weaponThemeMap, c.Weapon.ID, roomTheme); !ok {
				allowed = false
			} else {
				weightScale *= w
			}
			if w, ok := dimensionWeight(weaponSourceMap, c.Weapon.ID, source); !ok {
				allowed = false
			} else {
				weightScale *= w
			}
			if w, ok := dimensionWeight(weaponBandMap, c.Weapon.ID, levelBand); !ok {
				allowed = false
			} else {
				weightScale *= w
			}
			if w, ok := dimensionWeight(weaponLocationMap, c.Weapon.ID, location); !ok {
				allowed = false
			} else {
				weightScale *= w
			}
		}
		if !allowed {
			continue
		}
		c.BaseWeight *= weightScale
		filteredCandidates = append(filteredCandidates, c)
	}
	candidates = filteredCandidates

	fill := runBudgetFillLoop(rng, candidates, rarityWeights{}, sessionBudget, expectedBudget)
	resultItems := make([]models.Item, 0)
	resultWeapons := make([]models.Weapon, 0)
	drops := make([]LootDrop, 0, len(fill.Chosen))
	for _, chosen := range fill.Chosen {
		if chosen.Kind == "item" && chosen.Item != nil {
			resultItems = append(resultItems, *chosen.Item)
			drops = append(drops, LootDrop{Kind: "item", Name: chosen.Item.Name, Rarity: chosen.Item.Rarity})
			continue
		}
		if chosen.Kind == "weapon" && chosen.Weapon != nil {
			resultWeapons = append(resultWeapons, *chosen.Weapon)
			drops = append(drops, LootDrop{Kind: "weapon", Name: chosen.Weapon.Name, Rarity: chosen.Weapon.Rarity})
		}
	}
	itemCount := len(resultItems)
	weaponCount := len(resultWeapons)
	goldEarned := int64(0)

	var lootMessages []string
	if goldEarned > 0 {
		lootMessages = append(lootMessages, fmt.Sprintf("%d copper pieces", goldEarned))
	}
	for _, item := range resultItems {
		lootMessages = append(lootMessages, item.Name)
	}
	for _, weapon := range resultWeapons {
		lootMessages = append(lootMessages, weapon.Name)
	}

	countPrefix := fmt.Sprintf("You purchased %d item(s) and %d weapon(s) from a level %d %s-themed loot budget. ", itemCount, weaponCount, lootLevel, roomTheme)
	message := countPrefix + "You found: "
	if len(lootMessages) > 0 {
		message += strings.Join(lootMessages, ", ") + "."
	} else {
		message += "nothing of value."
	}

	result := &LootResult{
		Items:         resultItems,
		Weapons:       resultWeapons,
		GoldEarned:    goldEarned,
		TotalMoney:    character.Money,
		Message:       message,
		ItemsRolled:   itemCount,
		WeaponsRolled: weaponCount,
		Drops:         drops,
		RoomTheme:     roomTheme,
		RewardAmount:  "",
		LevelBand:     levelBand,
		ProfileNotes: []string{
			fmt.Sprintf("Theme profile: %s", roomTheme),
			fmt.Sprintf("Location profile: %s", location),
			fmt.Sprintf("Loot level: %d", lootLevel),
			fmt.Sprintf("Level band: %s", levelBand),
			fmt.Sprintf("Expected budget: %d cp", expectedBudget),
			fmt.Sprintf("Session budget roll: %d cp", sessionBudget),
		},
		ExpectedBudget: expectedBudget,
		SessionBudget:  sessionBudget,
		EndingBudget:   fill.RemainingCP,
		DebtUsed:       fill.DebtUsed,
		DebtAmount:     fill.DebtAmountCP,
	}
	partyMembers := make([]LootPartyMember, 0)
	if character.PartyID != nil {
		members := make([]models.Character, 0)
		if err := s.db.Select("id, name").Where("party_id = ?", *character.PartyID).Find(&members).Error; err != nil {
			return nil, nil, nil, nil, fmt.Errorf("failed to load party members: %w", err)
		}
		for _, member := range members {
			partyMembers = append(partyMembers, LootPartyMember{ID: member.ID, Name: member.Name})
		}
	}
	if len(partyMembers) == 0 {
		partyMembers = append(partyMembers, LootPartyMember{ID: character.ID, Name: character.Name})
	}
	return result, fill.Chosen, character.PartyID, partyMembers, nil
}

func (s *LootService) GenerateLootPreview(characterID uuid.UUID, source, roomTheme string, locationOverride *string, lootLevel int) (*LootPreview, error) {
	result, chosen, partyID, members, err := s.generateLootResult(characterID, source, roomTheme, locationOverride, lootLevel)
	if err != nil {
		return nil, err
	}
	sessionID := uuid.New()
	s.pendingMu.Lock()
	s.pendingLoot[sessionID] = &pendingLootSession{
		SessionID:      sessionID,
		SourceCharacterID: characterID,
		PartyID:        partyID,
		CreatedAt:      time.Now(),
		Result:         result,
		Chosen:         chosen,
	}
	s.pendingMu.Unlock()
	return &LootPreview{SessionID: sessionID, Loot: result, PartyMembers: members}, nil
}

// GenerateLoot keeps backward-compatible behavior: preview then auto-assign to caller.
func (s *LootService) GenerateLoot(characterID uuid.UUID, source, roomTheme string, locationOverride *string, lootLevel int) (*LootResult, error) {
	preview, err := s.GenerateLootPreview(characterID, source, roomTheme, locationOverride, lootLevel)
	if err != nil {
		return nil, err
	}
	assignments := make([]struct {
		dropIndex int
		charID    uuid.UUID
	}, 0, len(preview.Loot.Drops))
	for i := range preview.Loot.Drops {
		assignments = append(assignments, struct {
			dropIndex int
			charID    uuid.UUID
		}{dropIndex: i, charID: characterID})
	}
	req := make([]LootAssignmentRequest, 0, len(assignments))
	for _, a := range assignments {
		req = append(req, LootAssignmentRequest{DropIndex: a.dropIndex, CharacterID: a.charID})
	}
	return s.ConfirmLootPickup(characterID, preview.SessionID, req)
}

type LootAssignmentRequest struct {
	DropIndex   int
	CharacterID uuid.UUID
}

func (s *LootService) ConfirmLootPickup(sourceCharacterID uuid.UUID, sessionID uuid.UUID, assignments []LootAssignmentRequest) (*LootResult, error) {
	s.pendingMu.Lock()
	session, ok := s.pendingLoot[sessionID]
	if ok && time.Since(session.CreatedAt) > 30*time.Minute {
		delete(s.pendingLoot, sessionID)
		ok = false
	}
	if ok {
		delete(s.pendingLoot, sessionID)
	}
	s.pendingMu.Unlock()
	if !ok || session == nil {
		return nil, fmt.Errorf("%w: invalid or expired loot session", ErrInvalidLootOption)
	}
	if session.SourceCharacterID != sourceCharacterID {
		return nil, fmt.Errorf("%w: loot session does not belong to this character", ErrInvalidLootOption)
	}
	if len(assignments) != len(session.Chosen) {
		return nil, fmt.Errorf("%w: every drop must be assigned before pickup", ErrInvalidLootOption)
	}
	recipients := make(map[uuid.UUID]*models.Character)
	allowed := make(map[uuid.UUID]struct{})
	if session.PartyID != nil {
		members := make([]models.Character, 0)
		if err := s.db.Where("party_id = ?", *session.PartyID).Find(&members).Error; err != nil {
			return nil, fmt.Errorf("failed to load party members for pickup: %w", err)
		}
		for i := range members {
			member := members[i]
			allowed[member.ID] = struct{}{}
		}
	}
	allowed[sourceCharacterID] = struct{}{}
	for _, assignment := range assignments {
		if assignment.DropIndex < 0 || assignment.DropIndex >= len(session.Chosen) {
			return nil, fmt.Errorf("%w: invalid drop index %d", ErrInvalidLootOption, assignment.DropIndex)
		}
		if _, ok := allowed[assignment.CharacterID]; !ok {
			return nil, fmt.Errorf("%w: character %s is not in party", ErrInvalidLootOption, assignment.CharacterID)
		}
		if _, ok := recipients[assignment.CharacterID]; !ok {
			char, err := s.characterRepo.FindByIDWithInventory(assignment.CharacterID)
			if err != nil {
				return nil, fmt.Errorf("failed to load recipient character %s: %w", assignment.CharacterID, err)
			}
			recipients[assignment.CharacterID] = char
		}
	}
	for _, assignment := range assignments {
		char := recipients[assignment.CharacterID]
		chosen := session.Chosen[assignment.DropIndex]
		if chosen.Kind == "item" && chosen.Item != nil {
			char.Items = append(char.Items, *chosen.Item)
		} else if chosen.Kind == "weapon" && chosen.Weapon != nil {
			char.Weapons = append(char.Weapons, *chosen.Weapon)
		}
	}
	for _, char := range recipients {
		if err := s.characterRepo.Update(char); err != nil {
			return nil, fmt.Errorf("failed to persist loot pickup for %s: %w", char.ID, err)
		}
	}
	return session.Result, nil
}
