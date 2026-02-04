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
			AttackModifier: "Strength",
			Properties:     []string{"Transformative", "Versatile"},
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

// AllWeapons collects all weapon categories
func AllWeapons() []WeaponSeed {
	var all []WeaponSeed
	all = append(all, TransformativeWeapons()...)
	all = append(all, HydraulicWeapons()...)
	all = append(all, AlchemicalWeapons()...)
	all = append(all, ElectricWeapons()...)
	return all
}
