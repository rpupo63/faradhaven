package faradhaven_items

type LootTags struct {
	Themes        []string
	Locations     []string
	Sources       []string
	Tiers         []string
	RewardAmounts []string
	LevelBands    []string
	Weight        float64
}

// WeaponSeed defines the data structure for seeding weapons
type WeaponSeed struct {
	Name                string
	Description         string
	Category            string // e.g. "Simple Melee", "Martial Ranged"
	Rarity              string // e.g. "Common", "Rare"
	RangeType           string // "Melee" or "Ranged"
	Cost                string
	Weight              string
	AttackModifier      string   // "Strength", "Dexterity"
	Properties          []string // e.g. "Finesse", "Light"
	RangeNormal         int
	RangeLong           int
	VersatileDamageDice *string
	SecondaryEffect     string
	Damages             []WeaponDamageSeed
	LootTags            LootTags
}

// WeaponDamageSeed defines the damage data for a weapon seed
type WeaponDamageSeed struct {
	DamageDice     string // e.g. "1d8"
	DamageType     string // e.g. "Slashing"
	DamageCategory string // "Base", "Bonus"
}
