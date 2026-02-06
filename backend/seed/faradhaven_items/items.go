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
		{
			Name:         "Aetheric Revivifier",
			Description:  "A bubbling blue fluid that smells of ozone. It tastes like copper and electricity.",
			Category:     "Potion",
			Rarity:       "Uncommon",
			Cost:         "150 gp",
			Weight:       "0.5 lb.",
			Effects:      "Restore 10 Spell Points, 5 Stability, or 5 Blood Ichor.",
			IsConsumable: true,
		},
		{
			Name:         "Smog-Filter Lozenge",
			Description:  "A hard, charcoal-colored candy designed for the soot-heavy air of the industrial districts.",
			Category:     "Potion",
			Rarity:       "Common",
			Cost:         "10 gp",
			Weight:       "0.1 lb.",
			Effects:      "Advantage on saving throws against inhaled gases and poisons for 1 hour.",
			IsConsumable: true,
		},
		{
			Name:         "Alchemical Accelerant",
			Description:  "A silver, fast-moving liquid that causes the heart to race and the world to slow down.",
			Category:     "Potion",
			Rarity:       "Uncommon",
			Cost:         "200 gp",
			Weight:       "0.5 lb.",
			Effects:      "Gain +5 to Initiative and +10ft movement speed for 1 minute.",
			IsConsumable: true,
		},
		{
			Name:         "Stabilizer Serum",
			Description:  "A thick, viscous grey fluid that tastes of ash and helps calm the mind during transformation.",
			Category:     "Potion",
			Rarity:       "Uncommon",
			Cost:         "120 gp",
			Weight:       "0.5 lb.",
			Effects:      "The next Madness Save for a Mutagen mutation is made with Advantage.",
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
		{
			Name:         "Anatomist's Kit",
			Description:  "Precision scalpels, silver tweezers, and magnifying lenses for delicate biological study.",
			Category:     "Tool",
			Rarity:       "Uncommon",
			Cost:         "50 gp",
			Weight:       "5 lb.",
			Effects:      "Provides advantage on Lorewright Anatomical Insight and Visceral Psychometry checks.",
			IsConsumable: false,
		},
		{
			Name:         "Pneumatic Lockpicks",
			Description:  "Steam-pressurized tools that hiss as they force pins into place.",
			Category:     "Tool",
			Rarity:       "Uncommon",
			Cost:         "100 gp",
			Weight:       "2 lb.",
			Effects:      "Grants +2 to Thieves' Tools checks but makes a loud hiss audible within 20 feet.",
			IsConsumable: false,
		},
		{
			Name:         "Sanguine Extraction Pump",
			Description:  "A brass-and-glass device with fine needles for drawing Ichor from willing or unconscious subjects.",
			Category:     "Tool",
			Rarity:       "Uncommon",
			Cost:         "75 gp",
			Weight:       "3 lb.",
			Effects:      "Allows a Sanguinist to use Siphon on a target without causing a level of Exhaustion if used over 10 minutes.",
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
		{
			Name: "Clockwork Chronometer", Description: "A gold-plated pocket watch that ticks with unnatural precision.", Category: "Gear", Rarity: "Uncommon", Cost: "150 gp", Weight: "1 lb.", Effects: "Adds +1 second to Powder Mage timer durations.", IsConsumable: false,
		},
		{
			Name: "Aether-Lantern", Description: "A heavy brass lantern fueled by a glowing blue crystal.", Category: "Gear", Rarity: "Uncommon", Cost: "100 gp", Weight: "3 lb.", Effects: "Reveals magical auras and footprints within 10ft when focused as an action.", IsConsumable: false,
		},
		{
			Name: "Faraday Mesh Cloak", Description: "A cloak woven with fine copper wire to ground electrical energy.", Category: "Gear", Rarity: "Uncommon", Cost: "250 gp", Weight: "4 lb.", Effects: "Grants resistance to Lightning damage.", IsConsumable: false,
		},
		{
			Name: "Vapor-Smoke Pellet", Description: "A small brass sphere that releases a thick, soot-heavy cloud when crushed.", Category: "Gear", Rarity: "Common", Cost: "25 gp", Weight: "0.2 lb.", Effects: "10ft radius Heavily Obscured for 1 round when crushed as a bonus action.", IsConsumable: true,
		},
		{
			Name: "Stability Pressure Gauge", Description: "A small dial that clips onto a weapon to monitor kinetic buildup.", Category: "Gear", Rarity: "Uncommon", Cost: "100 gp", Weight: "0.5 lb.", Effects: "Advantage on Piston Brawler checks to prevent Malfunction.", IsConsumable: false,
		},
		{
			Name: "Specimen Preservation Jar", Description: "A lead-lined glass jar filled with alchemical brine.", Category: "Gear", Rarity: "Common", Cost: "15 gp", Weight: "2 lb.", Effects: "Keeps a creature's liver fresh for Lorewright use for up to 24 hours.", IsConsumable: false,
		},
		{
			Name: "Piston Core Assembly Kit", Description: "A specialized toolkit for maintaining and recalibrating Piston Cores.", Category: "Tool", Rarity: "Uncommon", Cost: "100 gp", Weight: "15 lb.", Effects: "Required for Piston Brawler weapon maintenance.", IsConsumable: false,
		},
		{
			Name: "Pneumatic Grappling Hook", Description: "A steam-powered launcher that fires a four-pronged hook with a high-tensile wire.", Category: "Gear", Rarity: "Uncommon", Cost: "150 gp", Weight: "8 lb.", Effects: "Grants a climbing speed of 30ft for 1 minute as an action.", IsConsumable: false,
		},
		{
			Name: "Industrial Respirator", Description: "A heavy leather mask with brass filters and glass goggles.", Category: "Gear", Rarity: "Common", Cost: "25 gp", Weight: "2 lb.", Effects: "Immunity to the effects of smoke and non-magical toxic fumes.", IsConsumable: false,
		},
		{
			Name: "Aetheric Tuning Fork", Description: "A vibrating fork that hums when near high concentrations of mana.", Category: "Gear", Rarity: "Uncommon", Cost: "120 gp", Weight: "1 lb.", Effects: "Can be used to locate the nearest elemental rift within 500 feet.", IsConsumable: false,
		},
		{
			Name: "Clockwork Lockbox", Description: "A small box with a shifting gear-based combination lock.", Category: "Gear", Rarity: "Uncommon", Cost: "50 gp", Weight: "3 lb.", Effects: "Requires a DC 20 Sleight of Hand check to open without the code.", IsConsumable: false,
		},
		{
			Name: "Smelter's Gloves", Description: "Heavy, insulated gloves designed for handling superheated metal.", Category: "Gear", Rarity: "Common", Cost: "15 gp", Weight: "2 lb.", Effects: "Resistance to Fire damage for the hands only.", IsConsumable: false,
		},
		{
			Name: "Aetheric Component Pouch", Description: "A pouch containing various alchemical salts and conductive filaments.", Category: "Gear", Rarity: "Common", Cost: "25 gp", Weight: "2 lb.", Effects: "Required for Powder Mage casting.", IsConsumable: false,
		},
		{
			Name: "Empty Blood Vial", Description: "A sterile glass tube with a rubber stopper.", Category: "Gear", Rarity: "Common", Cost: "1 gp", Weight: "0.1 lb.", Effects: "Can store 1 unit of Blood Ichor.", IsConsumable: false,
		},
		{
			Name: "Galvanic Battery", Description: "A heavy lead-acid cell with copper terminals.", Category: "Gear", Rarity: "Uncommon", Cost: "80 gp", Weight: "12 lb.", Effects: "Can be spent to restore 10 Stability or 2 Components.", IsConsumable: true,
		},
		{
			Name: "Mechanized Oil Can", Description: "A long-spouted brass can filled with high-viscosity lubricating oil.", Category: "Gear", Rarity: "Common", Cost: "10 gp", Weight: "2 lb.", Effects: "Heals a construct for 1d8 HP as an action.", IsConsumable: false,
		},
		{
			Name: "Dark Cloak", Description: "A heavy, soot-colored cloak designed to blend into the shadows of the city.", Category: "Gear", Rarity: "Common", Cost: "10 gp", Weight: "3 lb.", Effects: "Provides advantage on Stealth checks in dim light or darkness.", IsConsumable: false,
		},
		{
			Name: "Protective Goggles", Description: "Brass-rimmed goggles with interchangeable lenses for different lighting conditions.", Category: "Gear", Rarity: "Common", Cost: "5 gp", Weight: "0.5 lb.", Effects: "Immunity to blindness caused by bright light or steam.", IsConsumable: false,
		},
		{
			Name: "Bag of Gears and Springs", Description: "A collection of miscellaneous clockwork parts scavenged from the city.", Category: "Gear", Rarity: "Common", Cost: "5 gp", Weight: "5 lb.", Effects: "Can be used to repair simple mechanical objects.", IsConsumable: false,
		},
		{
			Name: "Bone Talisman", Description: "A string of small, bleached bones carved with protective runes.", Category: "Gear", Rarity: "Common", Cost: "10 gp", Weight: "0.5 lb.", Effects: "Can be used as a focus for Lorewright abilities.", IsConsumable: false,
		},
		{
			Name: "Traveler's Clothes", Description: "Sturdy, multi-layered clothing suitable for the unpredictable London weather.", Category: "Gear", Rarity: "Common", Cost: "2 gp", Weight: "4 lb.", Effects: "Provides protection against mild cold and rain.", IsConsumable: false,
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
