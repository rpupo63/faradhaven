package faradhaven_races

// FaradhavenRaceSeed defines the full race data for seeding
type FaradhavenRaceSeed struct {
	Name                string
	Description         string
	CreatureType        string // e.g., "Humanoid"
	Size                string // e.g., "Medium or Small"
	BaseSpeed           int
	PhotoURL            string // URL to race artwork/photo
	AbilityScoreBonuses map[string]int
	Languages           []string // Languages known by default (e.g., ["Common", "Elvish"])
	BonusLanguageCount  int      // Number of additional languages the player may choose
	Traits              []TraitSeed
	ComponentNames      []string // Names of spell components this race grants access to
}

// TraitSeed defines a race trait for seeding
type TraitSeed struct {
	Name           string
	Description    string
	LevelReq       int
	ActionType     string // "Passive", "Magic Action", "Bonus Action"
	UsesPerRest    string // "1", "Proficiency Bonus"
	ResetCondition string // "Long Rest", "Short Rest"
	RangeValue     string // e.g., "60", "10"
	AreaOfEffect   string
	SaveAbility    string // "CHA", "DEX", etc.
	Options        []TraitOptionSeed
}

// TraitOptionSeed defines a sub-choice within a trait (e.g., Heavenly Wings vs Necrotic Shroud)
type TraitOptionSeed struct {
	Name                string
	Description         string
	AbilityScoreBonuses map[string]int
}
