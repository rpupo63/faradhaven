package faradhaven_items

// TransformativeWeapons returns the trick weapon seeds
func TransformativeWeapons() []WeaponSeed {
	sawCleaverVersatile := "1d10"
	return []WeaponSeed{
		{
			Name:        "Saw Cleaver",
			Description: "A jagged blade on a hinged handle. It can be used as a short-range hack-and-slash weapon or flicked open for a long-reach, high-momentum saw.",
			Category:    "Martial Melee",
			Rarity:      "Uncommon",
			RangeType:   "Melee",
			Cost:        "50 gp",
			Weight:      "6 lb.",
			AttackModifier: "Dexterity",
			Properties:     []string{"Transformative", "Versatile", "Finesse"},
			RangeNormal:    5,
			VersatileDamageDice: &sawCleaverVersatile,
			SecondaryEffect: "In extended mode, gains Reach property and deals 1d10 damage.",
			Damages: []WeaponDamageSeed{
				{DamageDice: "1d6", DamageType: "Slashing", DamageCategory: "Base"},
			},
		},
		{
			Name:        "Whirligig Saw",
			Description: "A long pole that attaches to a motorized, circular serrated disc. When activated, the disc spins at high speeds via a steam engine.",
			Category:    "Martial Melee",
			Rarity:      "Rare",
			RangeType:   "Melee",
			Cost:        "250 gp",
			Weight:      "15 lb.",
			AttackModifier: "Strength",
			Properties:     []string{"Heavy", "Two-Handed", "Motorized"},
			RangeNormal:    5,
			SecondaryEffect: "As a bonus action, rev the engine to deal an extra 1d6 slashing damage on your next hit.",
			Damages: []WeaponDamageSeed{
				{DamageDice: "2d6", DamageType: "Slashing", DamageCategory: "Base"},
			},
		},
		{
			Name:        "Folding Blade",
			Description: "A masterwork of clockwork and tension springs. The blade is divided into segments that collapse into the hilt for total concealment.",
			Category:    "Simple Melee",
			Rarity:      "Common",
			RangeType:   "Melee",
			Cost:        "50 gp",
			Weight:      "1 lb.",
			AttackModifier: "Dexterity",
			Properties:     []string{"Finesse", "Light", "Concealed"},
			RangeNormal:    5,
			SecondaryEffect: "Advantage on checks to conceal this weapon. Deploys instantly.",
			Damages: []WeaponDamageSeed{
				{DamageDice: "1d4", DamageType: "Piercing", DamageCategory: "Base"},
			},
		},
		{
			Name:        "Switchglaive",
			Description: "Uses a complex gear-and-lock system to transition between a folding scythe (heavy), a short sword (fast), and a glaive (mid-range).",
			Category:    "Martial Melee",
			Rarity:      "Rare",
			RangeType:   "Melee",
			Cost:        "150 gp",
			Weight:      "8 lb.",
			AttackModifier: "Strength",
			Properties:     []string{"Transformative", "Reach", "Two-Handed"},
			RangeNormal:    10,
			SecondaryEffect: "Can transform into Scythe (2d4 Slashing) or Sword (1d6 Slashing, Finesse) as a bonus action.",
			Damages: []WeaponDamageSeed{
				{DamageDice: "1d10", DamageType: "Slashing", DamageCategory: "Base"},
			},
		},
	}
}

// NaturalWeapons returns innate or class-specific natural weapons
func NaturalWeapons() []WeaponSeed {
	return []WeaponSeed{
		{
			Name:           "Bite",
			Description:    "A vampiric bite that extracts blood ichor from victims.",
			Category:       "Natural Melee",
			Rarity:         "Common",
			RangeType:      "Melee",
			Cost:           "0 gp",
			Weight:         "0 lb.",
			AttackModifier: "Charisma",
			Properties:     []string{"Finesse"},
			RangeNormal:    5,
			SecondaryEffect: "On hit, regain Ichor equal to Proficiency Bonus. Yields 1 Unstable Component.",
			Damages: []WeaponDamageSeed{
				{DamageDice: "2d8", DamageType: "Necrotic", DamageCategory: "Base"},
			},
		},
	}
}

// HydraulicWeapons returns the pressure-based weapon seeds
func HydraulicWeapons() []WeaponSeed {
	return []WeaponSeed{
		{
			Name:        "Atlas Gauntlets",
			Description: "Oversized metal fists powered by Hextech. They use hydraulic pistons to deliver strikes that can shatter stone and generate defensive shields.",
			Category:    "Martial Melee",
			Rarity:      "Rare",
			RangeType:   "Melee",
			Cost:        "400 gp",
			Weight:      "20 lb.",
			AttackModifier: "Strength",
			Properties:     []string{"Heavy", "Siege"},
			RangeNormal:    5,
			SecondaryEffect: "Deals double damage to objects and structures. Can generate a temporary shield (Reaction, +2 AC).",
			Damages: []WeaponDamageSeed{
				{DamageDice: "1d8", DamageType: "Bludgeoning", DamageCategory: "Base"},
			},
		},
		{
			Name:        "Tsuranuki Zutsu",
			Description: "A 'piercing gun' that uses a backpack-mounted steam boiler to fire a metal bolt with enough pressure to penetrate industrial plating.",
			Category:    "Martial Ranged",
			Rarity:      "Very Rare",
			RangeType:   "Ranged",
			Cost:        "800 gp",
			Weight:      "45 lb.",
			AttackModifier: "Dexterity",
			Properties:     []string{"Heavy", "Loading", "Two-Handed", "Armor Piercing"},
			RangeNormal:    100,
			RangeLong:      300,
			SecondaryEffect: "+2 bonus to attack rolls against armored targets.",
			Damages: []WeaponDamageSeed{
				{DamageDice: "2d12", DamageType: "Piercing", DamageCategory: "Base"},
			},
		},
		{
			Name:        "Steam-Pressure Rifle",
			Description: "A firearm connected to a portable boiler via hoses. It uses high-pressure steam instead of gunpowder to propel lead shots.",
			Category:    "Martial Ranged",
			Rarity:      "Uncommon",
			RangeType:   "Ranged",
			Cost:        "200 gp",
			Weight:      "12 lb.",
			AttackModifier: "Dexterity",
			Properties:     []string{"Two-Handed", "Ammunition"},
			RangeNormal:    80,
			RangeLong:      240,
			SecondaryEffect: "Silent firing compared to gunpowder weapons.",
			Damages: []WeaponDamageSeed{
				{DamageDice: "1d10", DamageType: "Piercing", DamageCategory: "Base"},
			},
		},
		{
			Name:        "Pneumatic Crossbow",
			Description: "A silent weapon that uses pre-charged air canisters to fire bolts, often featuring modular lenses for long-range targeting.",
			Category:    "Simple Ranged",
			Rarity:      "Uncommon",
			RangeType:   "Ranged",
			Cost:        "150 gp",
			Weight:      "8 lb.",
			AttackModifier: "Dexterity",
			Properties:     []string{"Ammunition", "Loading", "Two-Handed", "Silent"},
			RangeNormal:    100,
			RangeLong:      400,
			SecondaryEffect: "Shots are completely silent and leave no smoke trail.",
			Damages: []WeaponDamageSeed{
				{DamageDice: "1d8", DamageType: "Piercing", DamageCategory: "Base"},
			},
		},
		{
			Name:        "Drill Arm",
			Description: "A massive, fuel-injected industrial drill mounted to the arm. It features a heavy motor and cooling vents to prevent seizing during use.",
			Category:    "Martial Melee",
			Rarity:      "Rare",
			RangeType:   "Melee",
			Cost:        "300 gp",
			Weight:      "25 lb.",
			AttackModifier: "Strength",
			Properties:     []string{"Heavy", "Special"},
			RangeNormal:    5,
			SecondaryEffect: "On a hit, you can choose to grapple the target. While grappling, you deal automatic drill damage at the start of your turn.",
			Damages: []WeaponDamageSeed{
				{DamageDice: "1d12", DamageType: "Piercing", DamageCategory: "Base"},
			},
		},
	}
}

// AlchemicalWeapons returns the alchemical/elemental weapon seeds
func AlchemicalWeapons() []WeaponSeed {
	return []WeaponSeed{
		{
			Name:        "Mercury Hammer",
			Description: "A heavy warhammer with an internal energy core. It can transform into a long-range cannon that fires 'shocks' of unstable energy.",
			Category:    "Martial Melee",
			Rarity:      "Legendary",
			RangeType:   "Melee",
			Cost:        "2500 gp",
			Weight:      "18 lb.",
			AttackModifier: "Strength",
			Properties:     []string{"Versatile", "Transformative", "Thrown"},
			RangeNormal:    5,
			RangeLong:      60,
			SecondaryEffect: "Can transform into cannon mode to fire lightning blasts (Range 60/120, 2d8 Lightning).",
			Damages: []WeaponDamageSeed{
				{DamageDice: "1d10", DamageType: "Bludgeoning", DamageCategory: "Base"},
				{DamageDice: "1d6", DamageType: "Lightning", DamageCategory: "Bonus"},
			},
		},
		{
			Name:        "Chemical Thrower",
			Description: "A series of pressurized tanks and nozzles designed to vent liquid nitrogen, napalm, or electrified gel.",
			Category:    "Martial Ranged",
			Rarity:      "Rare",
			RangeType:   "Ranged",
			Cost:        "500 gp",
			Weight:      "30 lb.",
			AttackModifier: "Dexterity",
			Properties:     []string{"Heavy", "Two-Handed", "Special"},
			RangeNormal:    15,
			SecondaryEffect: "Fires in a 15-foot cone. Targets must make a DC 15 Dex save or take damage.",
			Damages: []WeaponDamageSeed{
				{DamageDice: "2d6", DamageType: "Acid", DamageCategory: "Base"},
			},
		},
		{
			Name:        "The Big Daddy Harpoon",
			Description: "A heavy, spring-loaded launcher that fires oversized harpoons. Some variants include 'drill tips' or 'incendiary heads.'",
			Category:    "Martial Ranged",
			Rarity:      "Rare",
			RangeType:   "Ranged",
			Cost:        "350 gp",
			Weight:      "18 lb.",
			AttackModifier: "Dexterity",
			Properties:     []string{"Heavy", "Loading", "Two-Handed"},
			RangeNormal:    40,
			RangeLong:      120,
			SecondaryEffect: "On a hit, target is tethered if using a line. Can reel in Medium or smaller creatures.",
			Damages: []WeaponDamageSeed{
				{DamageDice: "1d12", DamageType: "Piercing", DamageCategory: "Base"},
			},
		},
		{
			Name:        "Pepperbox Revolver",
			Description: "A multi-barreled pistol with a hand-cranked rotating assembly. In steampunk settings, these often use alchemical cartridges (fire, frost, or acid).",
			Category:    "Martial Ranged",
			Rarity:      "Uncommon",
			RangeType:   "Ranged",
			Cost:        "175 gp",
			Weight:      "4 lb.",
			AttackModifier: "Dexterity",
			Properties:     []string{"Ammunition", "Reload (6)", "Light"},
			RangeNormal:    60,
			RangeLong:      180,
			SecondaryEffect: "Can load different elemental cartridges. Misfire score 2.",
			Damages: []WeaponDamageSeed{
				{DamageDice: "1d8", DamageType: "Piercing", DamageCategory: "Base"},
			},
		},
	}
}

// ElectricWeapons returns the tesla/electric weapon seeds
func ElectricWeapons() []WeaponSeed {
	return []WeaponSeed{
		{
			Name:        "Tesla Rifle",
			Description: "A weapon composed of copper coils and glass tubes. It projects arcs of high-voltage electricity that can 'chain' between multiple targets.",
			Category:    "Martial Ranged",
			Rarity:      "Rare",
			RangeType:   "Ranged",
			Cost:        "600 gp",
			Weight:      "10 lb.",
			AttackModifier: "Dexterity",
			Properties:     []string{"Two-Handed", "Energy"},
			RangeNormal:    60,
			RangeLong:      120,
			SecondaryEffect: "On hit, arcs to a second target within 10 ft for 1d4 Lightning damage.",
			Damages: []WeaponDamageSeed{
				{DamageDice: "1d8", DamageType: "Lightning", DamageCategory: "Base"},
			},
		},
		{
			Name:        "Lightning Gun",
			Description: "A sleek, silver rifle that fires concentrated electrical bursts, fueled by early chemical battery cells.",
			Category:    "Martial Ranged",
			Rarity:      "Very Rare",
			RangeType:   "Ranged",
			Cost:        "1200 gp",
			Weight:      "15 lb.",
			AttackModifier: "Dexterity",
			Properties:     []string{"Heavy", "Two-Handed"},
			RangeNormal:    100,
			RangeLong:      300,
			SecondaryEffect: "Advantage on attack rolls against targets wearing metal armor.",
			Damages: []WeaponDamageSeed{
				{DamageDice: "3d6", DamageType: "Lightning", DamageCategory: "Base"},
			},
		},
		{
			Name:        "The Amrad Cannon",
			Description: "A heavy, shoulder-mounted weapon that fires superheated steam or concussive air blasts, representing the pinnacle of pure steam power.",
			Category:    "Martial Ranged",
			Rarity:      "Legendary",
			RangeType:   "Ranged",
			Cost:        "5000 gp",
			Weight:      "60 lb.",
			AttackModifier: "Strength",
			Properties:     []string{"Heavy", "Siege", "Two-Handed"},
			RangeNormal:    150,
			RangeLong:      600,
			SecondaryEffect: "Requires a turn to charge before firing. Deals massive structural damage.",
			Damages: []WeaponDamageSeed{
				{DamageDice: "2d10", DamageType: "Thunder", DamageCategory: "Base"},
				{DamageDice: "2d10", DamageType: "Fire", DamageCategory: "Base"},
			},
		},
	}
}

// StandardWeapons returns basic D&D 5e weapon seeds
func StandardWeapons() []WeaponSeed {
	return []WeaponSeed{
		{
			Name: "Dagger", Description: "A small, double-edged blade.", Category: "Simple Melee", Rarity: "Common", RangeType: "Melee", Cost: "2 gp", Weight: "1 lb.",
			AttackModifier: "Dexterity", Properties: []string{"Finesse", "Light", "Thrown"}, RangeNormal: 20, RangeLong: 60,
			Damages: []WeaponDamageSeed{{DamageDice: "1d4", DamageType: "Piercing", DamageCategory: "Base"}},
		},
		{
			Name: "Shortsword", Description: "A light, agile blade.", Category: "Simple Melee", Rarity: "Common", RangeType: "Melee", Cost: "10 gp", Weight: "2 lb.",
			AttackModifier: "Dexterity", Properties: []string{"Finesse", "Light"}, RangeNormal: 5,
			Damages: []WeaponDamageSeed{{DamageDice: "1d6", DamageType: "Piercing", DamageCategory: "Base"}},
		},
		{
			Name: "Longsword", Description: "A versatile, double-edged blade.", Category: "Martial Melee", Rarity: "Common", RangeType: "Melee", Cost: "15 gp", Weight: "3 lb.",
			AttackModifier: "Strength", Properties: []string{"Versatile"}, RangeNormal: 5,
			VersatileDamageDice: stringPtr("1d10"),
			Damages: []WeaponDamageSeed{{DamageDice: "1d8", DamageType: "Slashing", DamageCategory: "Base"}},
		},
		{
			Name: "Greataxe", Description: "A massive, two-handed axe.", Category: "Martial Melee", Rarity: "Common", RangeType: "Melee", Cost: "30 gp", Weight: "7 lb.",
			AttackModifier: "Strength", Properties: []string{"Heavy", "Two-Handed"}, RangeNormal: 5,
			Damages: []WeaponDamageSeed{{DamageDice: "1d12", DamageType: "Slashing", DamageCategory: "Base"}},
		},
		{
			Name: "Greatsword", Description: "A massive, two-handed sword.", Category: "Martial Melee", Rarity: "Common", RangeType: "Melee", Cost: "50 gp", Weight: "6 lb.",
			AttackModifier: "Strength", Properties: []string{"Heavy", "Two-Handed"}, RangeNormal: 5,
			Damages: []WeaponDamageSeed{{DamageDice: "2d6", DamageType: "Slashing", DamageCategory: "Base"}},
		},
		{
			Name: "Mace", Description: "A heavy, bludgeoning weapon with a flanged head.", Category: "Simple Melee", Rarity: "Common", RangeType: "Melee", Cost: "5 gp", Weight: "4 lb.",
			AttackModifier: "Strength", Properties: []string{}, RangeNormal: 5,
			Damages: []WeaponDamageSeed{{DamageDice: "1d6", DamageType: "Bludgeoning", DamageCategory: "Base"}},
		},
		{
			Name: "Light Hammer", Description: "A small hammer that can be thrown.", Category: "Simple Melee", Rarity: "Common", RangeType: "Melee", Cost: "2 gp", Weight: "2 lb.",
			AttackModifier: "Strength", Properties: []string{"Light", "Thrown"}, RangeNormal: 20, RangeLong: 60,
			Damages: []WeaponDamageSeed{{DamageDice: "1d4", DamageType: "Bludgeoning", DamageCategory: "Base"}},
		},
		{
			Name: "Handaxe", Description: "A small axe designed for throwing or close combat.", Category: "Simple Melee", Rarity: "Common", RangeType: "Melee", Cost: "5 gp", Weight: "2 lb.",
			AttackModifier: "Strength", Properties: []string{"Light", "Thrown"}, RangeNormal: 20, RangeLong: 60,
			Damages: []WeaponDamageSeed{{DamageDice: "1d6", DamageType: "Slashing", DamageCategory: "Base"}},
		},
		{
			Name: "Spear", Description: "A simple polearm with a pointed tip.", Category: "Simple Melee", Rarity: "Common", RangeType: "Melee", Cost: "1 gp", Weight: "3 lb.",
			AttackModifier: "Strength", Properties: []string{"Thrown", "Versatile"}, RangeNormal: 20, RangeLong: 60,
			VersatileDamageDice: stringPtr("1d8"),
			Damages: []WeaponDamageSeed{{DamageDice: "1d6", DamageType: "Piercing", DamageCategory: "Base"}},
		},
		{
			Name: "Quarterstaff", Description: "A simple wooden staff.", Category: "Simple Melee", Rarity: "Common", RangeType: "Melee", Cost: "2 sp", Weight: "4 lb.",
			AttackModifier: "Strength", Properties: []string{"Versatile"}, RangeNormal: 5,
			VersatileDamageDice: stringPtr("1d8"),
			Damages: []WeaponDamageSeed{{DamageDice: "1d6", DamageType: "Bludgeoning", DamageCategory: "Base"}},
		},
		{
			Name: "Rapier", Description: "A slender, pointed sword designed for precision.", Category: "Martial Melee", Rarity: "Common", RangeType: "Melee", Cost: "25 gp", Weight: "2 lb.",
			AttackModifier: "Dexterity", Properties: []string{"Finesse"}, RangeNormal: 5,
			Damages: []WeaponDamageSeed{{DamageDice: "1d8", DamageType: "Piercing", DamageCategory: "Base"}},
		},
		{
			Name: "Scimitar", Description: "A curved, single-edged blade.", Category: "Martial Melee", Rarity: "Common", RangeType: "Melee", Cost: "25 gp", Weight: "3 lb.",
			AttackModifier: "Dexterity", Properties: []string{"Finesse", "Light"}, RangeNormal: 5,
			Damages: []WeaponDamageSeed{{DamageDice: "1d6", DamageType: "Slashing", DamageCategory: "Base"}},
		},
		{
			Name: "Warhammer", Description: "A heavy, two-handed hammer.", Category: "Martial Melee", Rarity: "Common", RangeType: "Melee", Cost: "15 gp", Weight: "2 lb.",
			AttackModifier: "Strength", Properties: []string{"Versatile"}, RangeNormal: 5,
			VersatileDamageDice: stringPtr("1d10"),
			Damages: []WeaponDamageSeed{{DamageDice: "1d8", DamageType: "Bludgeoning", DamageCategory: "Base"}},
		},
		{
			Name: "Battleaxe", Description: "A heavy, double-edged axe.", Category: "Martial Melee", Rarity: "Common", RangeType: "Melee", Cost: "10 gp", Weight: "4 lb.",
			AttackModifier: "Strength", Properties: []string{"Versatile"}, RangeNormal: 5,
			VersatileDamageDice: stringPtr("1d10"),
			Damages: []WeaponDamageSeed{{DamageDice: "1d8", DamageType: "Slashing", DamageCategory: "Base"}},
		},
		{
			Name: "Shortbow", Description: "A simple ranged weapon.", Category: "Simple Ranged", Rarity: "Common", RangeType: "Ranged", Cost: "25 gp", Weight: "2 lb.",
			AttackModifier: "Dexterity", Properties: []string{"Ammunition", "Two-Handed"}, RangeNormal: 80, RangeLong: 320,
			Damages: []WeaponDamageSeed{{DamageDice: "1d6", DamageType: "Piercing", DamageCategory: "Base"}},
		},
		{
			Name: "Longbow", Description: "A powerful ranged weapon.", Category: "Martial Ranged", Rarity: "Common", RangeType: "Ranged", Cost: "50 gp", Weight: "2 lb.",
			AttackModifier: "Dexterity", Properties: []string{"Ammunition", "Heavy", "Two-Handed"}, RangeNormal: 150, RangeLong: 600,
			Damages: []WeaponDamageSeed{{DamageDice: "1d8", DamageType: "Piercing", DamageCategory: "Base"}},
		},
		{
			Name: "Light Crossbow", Description: "A simple crossbow.", Category: "Simple Ranged", Rarity: "Common", RangeType: "Ranged", Cost: "25 gp", Weight: "5 lb.",
			AttackModifier: "Dexterity", Properties: []string{"Ammunition", "Loading", "Two-Handed"}, RangeNormal: 80, RangeLong: 320,
			Damages: []WeaponDamageSeed{{DamageDice: "1d8", DamageType: "Piercing", DamageCategory: "Base"}},
		},
		{
			Name: "Hand Crossbow", Description: "A small, one-handed crossbow.", Category: "Martial Ranged", Rarity: "Common", RangeType: "Ranged", Cost: "75 gp", Weight: "3 lb.",
			AttackModifier: "Dexterity", Properties: []string{"Ammunition", "Light", "Loading"}, RangeNormal: 30, RangeLong: 120,
			Damages: []WeaponDamageSeed{{DamageDice: "1d6", DamageType: "Piercing", DamageCategory: "Base"}},
		},
	}
}

func stringPtr(s string) *string { return &s }

// AllWeapons collects all weapon categories
func AllWeapons() []WeaponSeed {
	var all []WeaponSeed
	all = append(all, StandardWeapons()...)
	all = append(all, TransformativeWeapons()...)
	all = append(all, NaturalWeapons()...)
	all = append(all, HydraulicWeapons()...)
	all = append(all, AlchemicalWeapons()...)
	all = append(all, ElectricWeapons()...)
	return all
}
