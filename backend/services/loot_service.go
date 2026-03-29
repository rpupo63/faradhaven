package services

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
)

type LootService struct {
	itemRepo      database.ItemRepository
	weaponRepo    database.WeaponRepository
	characterRepo database.CharacterRepository
}

// LootDrop records one resolved drop for API/UI display.
type LootDrop struct {
	Kind   string `json:"kind"` // "item" or "weapon"
	Name   string `json:"name"`
	Rarity string `json:"rarity"`
}

type LootResult struct {
	Items         []models.Item   `json:"items"`
	Weapons       []models.Weapon `json:"weapons"`
	GoldEarned    int64           `json:"gold_earned"`
	TotalMoney    int64           `json:"total_money"`
	Message       string          `json:"message"`
	ItemsRolled   int             `json:"items_rolled"`
	WeaponsRolled int             `json:"weapons_rolled"`
	Drops         []LootDrop      `json:"drops"`
}

func NewLootService(itemRepo database.ItemRepository, weaponRepo database.WeaponRepository, characterRepo database.CharacterRepository) *LootService {
	return &LootService{itemRepo: itemRepo, weaponRepo: weaponRepo, characterRepo: characterRepo}
}

// lootParams defines the count ranges and gold for a (source, tier) combo.
type lootParams struct {
	ItemMin, ItemMax     int
	WeaponMin, WeaponMax int
	GoldMin, GoldMax     int64
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

// GenerateLoot updates character inventory and returns loot details.
func (s *LootService) GenerateLoot(characterID uuid.UUID, source, tier string) (*LootResult, error) {
	character, err := s.characterRepo.FindByIDWithInventory(characterID)
	if err != nil {
		return nil, fmt.Errorf("failed to find character %s: %w", characterID, err)
	}

	level := character.Level

	sourceTable, ok := lootTable[source]
	if !ok {
		return nil, fmt.Errorf("invalid source: %s (must be common_enemy, boss_enemy, or room)", source)
	}
	params, ok := sourceTable[tier]
	if !ok {
		return nil, fmt.Errorf("invalid tier: %s (must be low, medium, or high)", tier)
	}

	baseWeights, ok := tierWeights[tier]
	if !ok {
		return nil, fmt.Errorf("invalid tier: %s", tier)
	}
	weights := applyLevelShift(baseWeights, level)

	allItems, err := s.itemRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch items: %w", err)
	}
	allWeapons, err := s.weaponRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weapons: %w", err)
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

	itemCount := randRange(params.ItemMin, params.ItemMax)
	weaponCount := randRange(params.WeaponMin, params.WeaponMax)

	resultItems := make([]models.Item, 0, itemCount)
	resultWeapons := make([]models.Weapon, 0, weaponCount)
	drops := make([]LootDrop, 0, itemCount+weaponCount)

	for i := 0; i < itemCount; i++ {
		rarity := rollRarity(weights)
		selected := pickStratifiedItem(itemsByRarityCategory, rarity)
		if selected == nil {
			continue
		}
		resultItems = append(resultItems, *selected)
		character.Items = append(character.Items, *selected)
		drops = append(drops, LootDrop{Kind: "item", Name: selected.Name, Rarity: selected.Rarity})
	}

	for i := 0; i < weaponCount; i++ {
		rarity := rollRarity(weights)
		selected := pickStratifiedWeapon(weaponsByRarityCategory, rarity)
		if selected == nil {
			continue
		}
		resultWeapons = append(resultWeapons, *selected)
		character.Weapons = append(character.Weapons, *selected)
		drops = append(drops, LootDrop{Kind: "weapon", Name: selected.Name, Rarity: selected.Rarity})
	}

	baseGold := randRange64(params.GoldMin, params.GoldMax)
	multiplier := 1.0 + float64(level-1)*0.15
	goldEarned := int64(math.Round(float64(baseGold) * multiplier))
	character.Money += goldEarned

	err = s.characterRepo.Update(character)
	if err != nil {
		return nil, fmt.Errorf("failed to update character with new loot: %w", err)
	}

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

	countPrefix := fmt.Sprintf("You rolled %d item drop(s) and %d weapon drop(s). ", itemCount, weaponCount)
	message := countPrefix + "You found: "
	if len(lootMessages) > 0 {
		message += strings.Join(lootMessages, ", ") + "."
	} else {
		message += "nothing of value."
	}

	return &LootResult{
		Items:         resultItems,
		Weapons:       resultWeapons,
		GoldEarned:    goldEarned,
		TotalMoney:    character.Money,
		Message:       message,
		ItemsRolled:   itemCount,
		WeaponsRolled: weaponCount,
		Drops:         drops,
	}, nil
}
