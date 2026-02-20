package services

import (
	"fmt"
	"math"
	"math/rand"
	"strings" // Added for strings.Join

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
)

type LootService struct {
	itemRepo      database.ItemRepository
	weaponRepo    database.WeaponRepository
	characterRepo database.CharacterRepository
}

type LootResult struct {
	Items      []models.Item   `json:"items"`
	Weapons    []models.Weapon `json:"weapons"`
	GoldEarned int64           `json:"gold_earned"`
	TotalMoney int64           `json:"total_money"` // Added: Character's total money after loot
	Message    string          `json:"message"`     // Added: A user-friendly message
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

// rarityWeights maps each rarity to its weight for a given tier.
type rarityWeights struct {
	Common    int
	Uncommon  int
	Rare      int
	VeryRare  int
	Legendary int
}

var lootTable = map[string]map[string]lootParams{
	"common_enemy": {
		"low":    {0, 1, 0, 0, 50, 200},
		"medium": {1, 1, 0, 1, 200, 500},
		"high":   {1, 2, 0, 1, 500, 1000},
	},
	"boss_enemy": {
		"low":    {1, 1, 0, 1, 500, 1000},
		"medium": {1, 2, 1, 1, 1000, 3000},
		"high":   {2, 3, 1, 2, 3000, 8000},
	},
	"room": {
		"low":    {0, 1, 0, 0, 100, 300},
		"medium": {1, 1, 0, 0, 300, 800},
		"high":   {1, 2, 0, 1, 800, 2000},
	},
}

var tierWeights = map[string]rarityWeights{
	"low":    {70, 25, 4, 1, 0},
	"medium": {45, 35, 15, 4, 1},
	"high":   {20, 30, 30, 14, 6},
}

// applyLevelShift shifts weight from Common to higher rarities every 5 levels above 1.
func applyLevelShift(w rarityWeights, level int) rarityWeights {
	shifts := (level - 1) / 5
	if shifts <= 0 {
		return w
	}
	totalShift := shifts * 5
	if totalShift > w.Common {
		totalShift = w.Common
	}
	w.Common -= totalShift
	// Distribute evenly to Uncommon, Rare, VeryRare, Legendary
	perBucket := totalShift / 4
	remainder := totalShift % 4
	w.Uncommon += perBucket
	w.Rare += perBucket
	w.VeryRare += perBucket
	w.Legendary += perBucket + remainder
	return w
}

func rollRarity(w rarityWeights) string {
	total := w.Common + w.Uncommon + w.Rare + w.VeryRare + w.Legendary
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
	if roll < w.VeryRare {
		return "Very Rare"
	}
	return "Legendary"
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

// GenerateLoot modified to update character inventory and return new fields
func (s *LootService) GenerateLoot(characterID uuid.UUID, source, tier string) (*LootResult, error) {
	// Fetch character to get level and update inventory
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

	systemItems := make([]*models.Item, 0)
	for _, item := range allItems {
		if item.UserID == nil {
			systemItems = append(systemItems, item)
		}
	}
	systemWeapons := make([]models.Weapon, 0)
	for _, w := range allWeapons {
		if w.UserID == nil {
			systemWeapons = append(systemWeapons, w)
		}
	}

	itemsByRarity := make(map[string][]*models.Item)
	for _, item := range systemItems {
		itemsByRarity[item.Rarity] = append(itemsByRarity[item.Rarity], item)
	}
	weaponsByRarity := make(map[string][]models.Weapon)
	for _, w := range systemWeapons {
		weaponsByRarity[w.Rarity] = append(weaponsByRarity[w.Rarity], w)
	}

	itemCount := randRange(params.ItemMin, params.ItemMax)
	resultItems := make([]models.Item, 0, itemCount)
	for i := 0; i < itemCount; i++ {
		rarity := rollRarity(weights)
		pool := itemsByRarity[rarity]
		if len(pool) == 0 {
			pool = itemsByRarity["Common"]
		}
		if len(pool) == 0 {
			continue
		}
		selectedItem := *pool[rand.Intn(len(pool))]
		resultItems = append(resultItems, selectedItem)
		character.Items = append(character.Items, selectedItem)
	}

	weaponCount := randRange(params.WeaponMin, params.WeaponMax)
	resultWeapons := make([]models.Weapon, 0, weaponCount)
	for i := 0; i < weaponCount; i++ {
		rarity := rollRarity(weights)
		pool := weaponsByRarity[rarity]
		if len(pool) == 0 {
			pool = weaponsByRarity["Common"]
		}
		if len(pool) == 0 {
			continue
		}
		selectedWeapon := pool[rand.Intn(len(pool))]
		resultWeapons = append(resultWeapons, selectedWeapon)
		character.Weapons = append(character.Weapons, selectedWeapon)
	}

	baseGold := randRange64(params.GoldMin, params.GoldMax)
	multiplier := 1.0 + float64(level-1)*0.15
	goldEarned := int64(math.Round(float64(baseGold) * multiplier))
	character.Money += goldEarned

	err = s.characterRepo.Update(character)
	if err != nil {
		return nil, fmt.Errorf("failed to update character with new loot: %w", err)
	}

	// Construct user-friendly message
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

	message := "You found: "
	if len(lootMessages) > 0 {
		message += strings.Join(lootMessages, ", ") + "."
	} else {
		message += "nothing of value."
	}

	return &LootResult{
		Items:      resultItems,
		Weapons:    resultWeapons,
		GoldEarned: goldEarned,
		TotalMoney: character.Money, // Return the updated total money
		Message:    message,         // Return the generated message
	}, nil
}


