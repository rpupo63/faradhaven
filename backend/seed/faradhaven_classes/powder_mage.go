package faradhaven_classes

// Fusilier returns the Powder Mage class seed
func PowderMage() FaradhavenClassSeed {
	return FaradhavenClassSeed{
		Name:           "The Powder Mage",
		Description:    "Etch spells directly onto bullet casings and fire magic through your Caster Gun. Cycle through ice shots, piercing rounds, and explosive shells to dominate at range with deadly precision.",
		HitDie:         10,
		PrimaryAbility: "dexterity",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/powder_mage.jpg",
		Archetype:      "Ranged DPS / Sniper",
		Concept:        "The gun-mage. While the Tinkerer builds turrets, the Fusilier is a master of the Caster Gun. They etch spells onto bullet casings.",
		ClassFeatures: []string{
			"Spell-Shot: You do not cast spells; you load them. Casting a spell is part of the Attack action with a firearm.",
			"Chamber Roulette: You can prepare different shells (Ice shot, Net shot, Piercing shot) and cycle through them instantly.",
		},
		DnDSkillFocus: []string{"Perception", "Investigation"},
		Proficiencies: "Simple Weapons, Firearms (Pistols, Rifles, Blunderbusses), Medium Armor",
		SkillChoice:   []string{"Perception", "Insight", "Survival", "Investigation"},
		Tools:         []string{"Gunsmith's Tools"},
		SavingThrows:  []string{"Dexterity", "Wisdom"},
		StartingEquip: []string{"Heavy coat (padded armor)", "Flintlock Rifle or two Pistols", "Gunsmith's kit", "Ammo pouch"},
		LevelFeatures: powderMageLevelFeatures(),
	}
}

func powderMageLevelFeatures() map[int]string {
	return map[int]string{
		1:  "Spell-Shot, Chamber Roulette — You gain the core class features. When you take the Attack action with a firearm, you can spend spell points to imbue your shot with a component (Ignite, Frost, Pierce, Seeker). The component effect triggers on a hit. You have advantage on Perception checks to spot hidden creatures beyond 30 feet.",
		2:  "Steady Barrel — When you haven't moved this turn, you can spend 2 spell points to add your Dexterity modifier to the damage of your next firearm attack. You have advantage on that attack roll if the target is more than 30 feet away.",
		3:  "Quick Reload — You can reload a firearm as a bonus action. When you use Chamber Roulette to switch shells, it no longer costs a bonus action; you can switch as part of the Attack action.",
		4:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Dexterity cannot exceed 20.",
		5:  "Extra Attack — When you take the Attack action, you can attack twice instead of once. Your second attack can use a different component shell.",
		6:  "Ignite Mastery — When you hit with an Ignite-augmented shot, the target takes 1d4 fire damage at the start of their next turn. Frost-augmented shots reduce the target's speed by 10 feet until the end of their next turn.",
		7:  "Pierce Through — When you use the Pierce modifier, your shot ignores half cover and three-quarters cover. You can spend 3 additional spell points to also ignore full cover if you can trace a line to the target.",
		8:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Dexterity cannot exceed 20.",
		9:  "Seeker Lock — When you use the Seeker modifier, the target has disadvantage on the saving throw (if any). Once per short rest, you can spend 5 spell points to make a Seeker shot that automatically hits (no attack roll) but deals half damage.",
		10: "Overcharged Round — Once per short rest, you can spend 10 spell points as part of an attack to fire an overcharged shot. On a hit, the target takes an extra 2d6 damage of the component's type (fire for Ignite, cold for Frost, etc.) and must succeed on a Constitution save or be knocked prone.",
		11: "Rapid Volley — When you use the Attack action and hit with both attacks, you can spend 5 spell points as a bonus action to make one additional firearm attack. You can use this feature a number of times per long rest equal to your Dexterity modifier.",
		12: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Dexterity cannot exceed 20.",
		13: "Dual Component — When you fire a shot, you can spend 5 additional spell points to apply a second modifier component (Pierce or Seeker) to the same attack. Both effects trigger on a hit.",
		14: "Dead Eye — Your firearm attacks score a critical hit on a roll of 19 or 20. When you score a critical hit with a component-augmented shot, you regain spell points equal to half the spell points spent on that shot (rounded down).",
		15: "Evasive Reload — You can take the Disengage or Dodge action as a bonus action when you reload a firearm. You have resistance to opportunity attacks made against you when you move after firing.",
		16: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Dexterity cannot exceed 20.",
		17: "Devastating Shot — Once per long rest, when you hit with a firearm attack, you can spend 15 spell points to maximize the damage dice of that attack (including component effects). The target must succeed on a Constitution save or be stunned until the end of your next turn.",
		18: "Overcharged Mastery — Your Overcharged Round no longer requires a short rest. You can use it once per long rest. Additionally, your standard Overcharged Round damage increases to 3d6.",
		19: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Dexterity cannot exceed 20.",
		20: "Fusilier's Precision — You are a master of the Caster Gun. You regain 20 spell points when you roll initiative and have none remaining. Your Steady Barrel no longer requires you to remain still, and you can use it once per turn without spending spell points. When you use Rapid Volley, the additional attack costs no spell points.",
	}
}
