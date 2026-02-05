package faradhaven_items

// ItemSeed defines the item data for seeding
type ItemSeed struct {
	Name         string
	Description  string
	Category     string
	Rarity       string
	Cost         string
	Weight       string
	Effects      string
	IsConsumable bool
}

// Potions returns the potion seeds
func Potions() []ItemSeed {
	return []ItemSeed{
		{
			Name:         "Healing Potion",
			Description:  "A vial of shimmering red liquid that restores vitality.",
			Category:     "Potion",
			Rarity:       "Common",
			Cost:         "50 gp",
			Weight:       "0.5 lb.",
			Effects:      "Restore 2d4+2 hit points.",
			IsConsumable: true,
		},
		{
			Name:         "Greater Healing Potion",
			Description:  "A larger vial of deep red liquid with floating golden flecks.",
			Category:     "Potion",
			Rarity:       "Uncommon",
			Cost:         "150 gp",
			Weight:       "0.5 lb.",
			Effects:      "Restore 4d4+4 hit points.",
			IsConsumable: true,
		},
		{
			Name:         "Steam-Breath Elixir",
			Description:  "This potion bubbles constantly. When drunk, your breath becomes visible and hot.",
			Category:     "Potion",
			Rarity:       "Uncommon",
			Cost:         "100 gp",
			Weight:       "0.5 lb.",
			Effects:      "Gain a breath weapon (15ft cone, 3d6 Fire) for 1 minute.",
			IsConsumable: true,
		},
	}
}

// Shields returns the shield seeds
func Shields() []ItemSeed {
	return []ItemSeed{
		{
			Name:         "Standard Shield",
			Description:  "A sturdy metal or wood shield.",
			Category:     "Shield",
			Rarity:       "Common",
			Cost:         "10 gp",
			Weight:       "6 lb.",
			Effects:      "AC +2",
			IsConsumable: false,
		},
		{
			Name:         "Clockwork Buckler",
			Description:  "A small, gear-driven shield that can expand its surface area.",
			Category:     "Shield",
			Rarity:       "Uncommon",
			Cost:         "75 gp",
			Weight:       "4 lb.",
			Effects:      "AC +2. As a reaction, expand to grant +4 AC against a single attack.",
			IsConsumable: false,
		},
	}
}

// Tools returns the tool seeds
func Tools() []ItemSeed {
	return []ItemSeed{
		{
			Name:         "Thieves' Tools",
			Description:  "A set of small files, lock picks, and specialized mirrors.",
			Category:     "Tool",
			Rarity:       "Common",
			Cost:         "25 gp",
			Weight:       "1 lb.",
			Effects:      "Enables proficiency bonus on checks to pick locks or disarm traps.",
			IsConsumable: false,
		},
		{
			Name:         "Tinker's Tools",
			Description:  "A collection of wrenches, hammers, and magnifying lenses.",
			Category:     "Tool",
			Rarity:       "Common",
			Cost:         "50 gp",
			Weight:       "10 lb.",
			Effects:      "Required for creating and repairing clockwork devices.",
			IsConsumable: false,
		},
	}
}

// Gear returns general adventuring gear seeds
func Gear() []ItemSeed {
	return []ItemSeed{
		{
			Name:         "Explorer's Pack",
			Description:  "Includes a backpack, bedroll, mess kit, tinderbox, 10 torches, 10 days of rations, and a waterskin.",
			Category:     "Gear",
			Rarity:       "Common",
			Cost:         "10 gp",
			Weight:       "59 lb.",
			Effects:      "Essential survival equipment.",
			IsConsumable: false,
		},
		{
			Name:         "Burglar's Pack",
			Description:  "Includes a backpack, 1000 ball bearings, 10ft of string, a bell, 5 candles, a crowbar, a hammer, 10 pitons, a hooded lantern, 2 flasks of oil, 5 days of rations, a tinderbox, and a waterskin.",
			Category:     "Gear",
			Rarity:       "Common",
			Cost:         "16 gp",
			Weight:       "47.5 lb.",
			Effects:      "Equipment for infiltration.",
			IsConsumable: false,
		},
		{
			Name: "Arcane Focus (Orb)", Description: "An orb of glass or crystal.", Category: "Gear", Rarity: "Common", Cost: "20 gp", Weight: "3 lb.", Effects: "A focus for spellcasting.", IsConsumable: false,
		},
		{
			Name: "Arcane Focus (Wand)", Description: "A wand of wood or metal.", Category: "Gear", Rarity: "Common", Cost: "10 gp", Weight: "1 lb.", Effects: "A focus for spellcasting.", IsConsumable: false,
		},
		{
			Name: "Arcane Focus (Crystal)", Description: "A large crystal.", Category: "Gear", Rarity: "Common", Cost: "10 gp", Weight: "1 lb.", Effects: "A focus for spellcasting.", IsConsumable: false,
		},
		{
			Name: "Alchemist's Supplies", Description: "All the beakers and chemicals needed for basic alchemy.", Category: "Gear", Rarity: "Common", Cost: "50 gp", Weight: "8 lb.", Effects: "Required for alchemical crafting.", IsConsumable: false,
		},
		{
			Name: "Scholar's Pack", Description: "Includes a backpack, a book of lore, a bottle of ink, an ink pen, 10 sheets of parchment, a little bag of sand, and a small knife.", Category: "Gear", Rarity: "Common", Cost: "40 gp", Weight: "10 lb.", Effects: "Equipment for scholars.", IsConsumable: false,
		},
		{
			Name: "Dungeoneer's Pack", Description: "Includes a backpack, a crowbar, a hammer, 10 pitons, 10 torches, a tinderbox, 10 days of rations, and a waterskin. The pack also has 50 feet of hempen rope strapped to the side of it.", Category: "Gear", Rarity: "Common", Cost: "12 gp", Weight: "61.5 lb.", Effects: "Equipment for dungeon delving.", IsConsumable: false,
		},
	}
}

// AllItems collects all item categories
func AllItems() []ItemSeed {
	var all []ItemSeed
	all = append(all, Potions()...)
	all = append(all, Shields()...)
	all = append(all, Tools()...)
	all = append(all, Gear()...)
	return all
}
