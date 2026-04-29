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

	// Armor-specific (only for Category containing "Armor" or "Shield")
	ArmorType           *string
	BaseAC              *int
	StrengthRequirement *int
	StealthDisadvantage *bool
	LootTags            LootTags
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
		// Standard fantasy / PHB-style potions
		{Name: "Superior Healing Potion", Description: "A ruby liquid that closes wounds with visible speed.", Category: "Potion", Rarity: "Rare", Cost: "500 gp", Weight: "0.5 lb.", Effects: "Restore 8d4+8 hit points.", IsConsumable: true},
		{Name: "Supreme Healing Potion", Description: "A brilliant crimson elixir that thrums with restorative magic.", Category: "Potion", Rarity: "Very Rare", Cost: "2000 gp", Weight: "0.5 lb.", Effects: "Restore 10d4+20 hit points.", IsConsumable: true},
		{Name: "Potion of Climbing", Description: "A sticky, almost tar-like potion that smells of wet stone.", Category: "Potion", Rarity: "Uncommon", Cost: "75 gp", Weight: "0.5 lb.", Effects: "Climbing speed equals walking speed for 1 hour.", IsConsumable: true},
		{Name: "Potion of Animal Friendship", Description: "A murky green liquid with bits of leaf suspended in it.", Category: "Potion", Rarity: "Uncommon", Cost: "100 gp", Weight: "0.5 lb.", Effects: "For 1 hour, cast Animal Friendship at will (save DC 13).", IsConsumable: true},
		{Name: "Potion of Growth", Description: "A cloudy potion that swirls with red and silver.", Category: "Potion", Rarity: "Uncommon", Cost: "90 gp", Weight: "0.5 lb.", Effects: "Become Large if Medium or smaller for 10 minutes; +1d4 weapon damage (if applicable).", IsConsumable: true},
		{Name: "Potion of Water Breathing", Description: "A smooth oil that smells of the sea.", Category: "Potion", Rarity: "Uncommon", Cost: "100 gp", Weight: "0.5 lb.", Effects: "Breathe underwater for 1 hour.", IsConsumable: true},
		{Name: "Potion of Resistance", Description: "A metallic-tasting liquid attuned to one damage type (acid, cold, fire, force, lightning, necrotic, poison, psychic, radiant, or thunder).", Category: "Potion", Rarity: "Uncommon", Cost: "300 gp", Weight: "0.5 lb.", Effects: "Resistance to that damage type for 1 hour.", IsConsumable: true},
		{Name: "Potion of Heroism", Description: "A bright blue liquid that sparkles when swirled.", Category: "Potion", Rarity: "Rare", Cost: "500 gp", Weight: "0.5 lb.", Effects: "Gain 10 temporary hit points; Bless on yourself for 1 minute.", IsConsumable: true},
		{Name: "Potion of Invisibility", Description: "Looks like nothing at all until you uncork it.", Category: "Potion", Rarity: "Very Rare", Cost: "5000 gp", Weight: "0.5 lb.", Effects: "Invisible for 1 hour or until you attack or cast a spell.", IsConsumable: true},
		{Name: "Potion of Speed", Description: "A yellow fluid streaked with black.", Category: "Potion", Rarity: "Very Rare", Cost: "5000 gp", Weight: "0.5 lb.", Effects: "Haste effect on yourself for 1 minute (no lethargy).", IsConsumable: true},
		{Name: "Potion of Vitality", Description: "A deep crimson potion that leaves you feeling invincible.", Category: "Potion", Rarity: "Very Rare", Cost: "5000 gp", Weight: "0.5 lb.", Effects: "Remove exhaustion, end disease/poison, max HP for 24 hours.", IsConsumable: true},
		{Name: "Potion of Mind Reading", Description: "An opalescent liquid that tastes faintly of copper.", Category: "Potion", Rarity: "Rare", Cost: "500 gp", Weight: "0.5 lb.", Effects: "Detect Thoughts (save DC 13) at will for 10 minutes.", IsConsumable: true},
		{Name: "Elixir of Health", Description: "A clear liquid with a single golden bubble floating in the center.", Category: "Potion", Rarity: "Rare", Cost: "120 gp", Weight: "0.5 lb.", Effects: "Cures any disease, neutralizes poison in your system.", IsConsumable: true},
		{Name: "Philter of Love", Description: "A rose-hued elixir that smells of summer flowers.", Category: "Potion", Rarity: "Uncommon", Cost: "90 gp", Weight: "0.5 lb.", Effects: "Charmed by the first creature you see for 1 hour (DC 13 Wis save).", IsConsumable: true},
		{Name: "Potion of Poison", Description: "Looks identical to a Potion of Healing.", Category: "Potion", Rarity: "Uncommon", Cost: "100 gp", Weight: "0.5 lb.", Effects: "DC 13 Con save or take 3d6 poison damage.", IsConsumable: true},
		{Name: "Oil of Sharpness", Description: "A clear, viscous oil that clings to metal.", Category: "Potion", Rarity: "Very Rare", Cost: "5000 gp", Weight: "0.5 lb.", Effects: "Coat one slashing weapon: +3 to attack and damage for 1 hour.", IsConsumable: true},
		{Name: "Potion of Fire Giant Strength", Description: "A deep red potion that feels warm to the touch.", Category: "Potion", Rarity: "Rare", Cost: "500 gp", Weight: "0.5 lb.", Effects: "Strength becomes 25 for 1 hour.", IsConsumable: true},
		{Name: "Potion of Gaseous Form", Description: "A cloudy gray fluid that never settles.", Category: "Potion", Rarity: "Rare", Cost: "500 gp", Weight: "0.5 lb.", Effects: "Gaseous Form for 1 hour.", IsConsumable: true},
		{Name: "Potion of Diminution", Description: "A red liquid that shrinks as you drink.", Category: "Potion", Rarity: "Rare", Cost: "500 gp", Weight: "0.5 lb.", Effects: "Reduce (Enlarge/Reduce) for 1 hour with no concentration.", IsConsumable: true},
		{Name: "Potion of Longevity", Description: "A clear liquid with a single hair-thin streak of silver.", Category: "Potion", Rarity: "Very Rare", Cost: "9000 gp", Weight: "0.5 lb.", Effects: "Reduce apparent age by 1d6+6 years (max 13 years younger).", IsConsumable: true},
	}
}

// Shields returns the shield seeds
func Shields() []ItemSeed {
	ac2 := 2
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
			BaseAC:       &ac2,
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
			BaseAC:       &ac2,
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
	all = append(all, StandardAdventuringGear()...)
	all = append(all, Armor()...)
	all = append(all, VendorThematicItems()...)
	all = append(all, ThemedLootItems()...)
	return all
}

func ThemedLootItems() []ItemSeed {
	return []ItemSeed{
		{
			Name:         "Rusted Reliquary Key",
			Description:  "A heavy iron key stamped with an ancient chapel sigil.",
			Category:     "Gear",
			Rarity:       "Uncommon",
			Cost:         "65 gp",
			Weight:       "1 lb.",
			Effects:      "Can unlock forgotten reliquary locks and grants advantage on checks to identify old dungeon mechanisms.",
			IsConsumable: false,
			LootTags: LootTags{
				Themes:        []string{"dungeon"},
				Locations:     []string{"underground"},
				Sources:       []string{"room", "boss_enemy"},
				Tiers:         []string{"medium", "high"},
				RewardAmounts: []string{"standard", "bountiful"},
				LevelBands:    []string{"adventurer", "veteran"},
				Weight:        1.35,
			},
		},
		{
			Name:         "Executive Seal Stamp",
			Description:  "A polished brass stamp used to authorize guild payroll and permits.",
			Category:     "Tool",
			Rarity:       "Common",
			Cost:         "55 gp",
			Weight:       "0.5 lb.",
			Effects:      "Provides leverage in office and bureaucracy social checks.",
			IsConsumable: false,
			LootTags: LootTags{
				Themes:        []string{"office"},
				Locations:     []string{"indoor", "urban"},
				Sources:       []string{"room", "common_enemy"},
				Tiers:         []string{"low", "medium"},
				RewardAmounts: []string{"scarce", "standard"},
				LevelBands:    []string{"novice", "adventurer"},
				Weight:        1.15,
			},
		},
		{
			Name:         "Velvet Coin Ledger",
			Description:  "A lacquered account book that includes hidden noble debt transactions.",
			Category:     "Gear",
			Rarity:       "Rare",
			Cost:         "220 gp",
			Weight:       "2 lb.",
			Effects:      "Can be traded for favors with wealthy factions.",
			IsConsumable: false,
			LootTags: LootTags{
				Themes:        []string{"rich"},
				Locations:     []string{"estate", "indoor"},
				Sources:       []string{"room", "boss_enemy"},
				Tiers:         []string{"medium", "high"},
				RewardAmounts: []string{"bountiful", "jackpot"},
				LevelBands:    []string{"adventurer", "veteran", "legend"},
				Weight:        1.55,
			},
		},
		{
			Name:         "Patchwork Ration Bundle",
			Description:  "A tightly wrapped bundle of dried roots, stale bread, and preserved fish.",
			Category:     "Gear",
			Rarity:       "Common",
			Cost:         "18 gp",
			Weight:       "2 lb.",
			Effects:      "Counts as 5 days of rations.",
			IsConsumable: true,
			LootTags: LootTags{
				Themes:        []string{"poor"},
				Locations:     []string{"slums", "urban"},
				Sources:       []string{"room", "common_enemy"},
				Tiers:         []string{"low", "medium"},
				RewardAmounts: []string{"scarce", "standard"},
				LevelBands:    []string{"novice", "adventurer"},
				Weight:        1.2,
			},
		},
		{
			Name:         "Contraband Cipher Notebook",
			Description:  "A cramped notebook filled with turf maps and coded extortion records.",
			Category:     "Tool",
			Rarity:       "Uncommon",
			Cost:         "95 gp",
			Weight:       "1 lb.",
			Effects:      "Grants advantage on checks for underground contacts and street intel.",
			IsConsumable: false,
			LootTags: LootTags{
				Themes:        []string{"gangster"},
				Locations:     []string{"street", "urban"},
				Sources:       []string{"common_enemy", "boss_enemy", "room"},
				Tiers:         []string{"medium", "high"},
				RewardAmounts: []string{"standard", "bountiful"},
				LevelBands:    []string{"adventurer", "veteran"},
				Weight:        1.4,
			},
		},
		{
			Name:         "Fractured Mana Prism",
			Description:  "A cracked crystal prism that leaks motes of arcane static.",
			Category:     "Gear",
			Rarity:       "Rare",
			Cost:         "310 gp",
			Weight:       "1 lb.",
			Effects:      "Can be consumed to restore 12 spell points.",
			IsConsumable: true,
			LootTags: LootTags{
				Themes:        []string{"arcane"},
				Locations:     []string{"indoor", "underground"},
				Sources:       []string{"room", "boss_enemy"},
				Tiers:         []string{"medium", "high"},
				RewardAmounts: []string{"bountiful", "jackpot"},
				LevelBands:    []string{"adventurer", "veteran", "legend"},
				Weight:        1.6,
			},
		},
		{
			Name:         "Predator Scent Satchel",
			Description:  "A leather pouch packed with wild herbs, resin, and beast musk.",
			Category:     "Gear",
			Rarity:       "Uncommon",
			Cost:         "55 gp",
			Weight:       "0.8 lb.",
			Effects:      "Provides advantage on tracking checks in wilderness terrain.",
			IsConsumable: false,
			LootTags: LootTags{
				Themes:        []string{"wilderness"},
				Locations:     []string{"wilds"},
				Sources:       []string{"common_enemy", "room", "boss_enemy"},
				Tiers:         []string{"low", "medium", "high"},
				RewardAmounts: []string{"standard", "bountiful"},
				LevelBands:    []string{"novice", "adventurer", "veteran"},
				Weight:        1.3,
			},
		},
	}
}
