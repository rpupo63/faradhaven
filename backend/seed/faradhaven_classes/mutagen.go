package faradhaven_classes

// Mutagen returns the Mutagen class seed
func Mutagen() FaradhavenClassSeed {
	return FaradhavenClassSeed{
		Name:           "The Mutagen",
		Description:    "Transform into a monstrous beast by metabolizing alchemical reagents, trading spell slots for raw biological mutations. You're the front-line juggernaut who grows claws, hardens skin, and regenerates mid-combat.",
		HitDie:         12,
		PrimaryAbility: "wisdom",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/mutagen.jpg",
		Archetype:      "Tank / Shock Trooper",
		Concept:        "A survivor of industrial pollution or experimental alchemy who has learned to control the mutagenic beast within. They don't cast spells; they metabolize components to trigger biological mutations.",
		ClassFeatures: []string{
			"Mutagenesis: You do not have spell slots. Instead, you ingest components (spending spell points) to grow claws, harden skin, heighten senses, or accelerate healing for a duration. Each mutation has its own duration based on the component used.",
			"Primal Form: At level 1, you design a specific Beast Shape (Wolf, Bear, Rat, Insect). Your form influences your appearance in combat and grants thematic bonuses when you activate mutagens.",
			"Metabolic Surge: Your body burns through alchemical reagents quickly. Mutagens are potent but shorter-lived than traditional magic—you trade longevity for raw physiological power.",
		},
		DnDSkillFocus: []string{"Survival", "Perception"},
		Proficiencies: "Simple weapons, Medium Armor (Hide/Leather only), Shields",
		SkillChoice:   []string{"Nature", "Athletics", "Intimidation", "Insight"},
		Tools:         []string{"Alchemist's Supplies (for stabilizing mutagens)"},
		SavingThrows:  []string{"Constitution", "Wisdom"},
		StartingEquip: []string{"Tattered traveler's clothes", "Heavy wooden shield", "Handaxe", "Iron flask (for mixing)"},
		LevelFeatures: mutagenLevelFeatures(),
	}
}

func mutagenLevelFeatures() map[int]string {
	return map[int]string{
		1:  "Primal Form — Choose your Beast Shape: Wolf (pack tactics, keen smell), Bear (raw power, thick fur), Rat (nimble, disease resistance), or Insect (chitin plating, climbing). When you use Mutagenesis, you partially assume traits of your form. Your unarmed strikes in Primal Form deal 1d6 + Strength modifier slashing damage (claws or fangs).",
		2:  "Feral Senses — Your heightened senses from the mutagen grant advantage on Perception (Wisdom) and Survival (Wisdom) checks to track creatures by scent or sound. When using the Adrenaline component, you also gain darkvision out to 60 feet for the duration.",
		3:  "Bestial Fury — When you have Feral Muscle active and are in Primal Form, you can use a bonus action to make one additional claw or bite attack. This attack deals 1d6 + Strength modifier damage.",
		4:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Wisdom cannot exceed 20.",
		5:  "Extra Attack — When you take the Attack action, you can attack twice instead of once. Your beast form is now fully combat-ready; you gain a +1 bonus to AC while Chitin/Hide is active.",
		6:  "Thick Hide — When you use the Chitin/Hide component, you gain resistance to nonmagical bludgeoning, piercing, and slashing damage for the duration. Bear-form Chimeras gain an additional +1 temporary HP per level while this mutagen is active.",
		7:  "Primal Resilience — Your body has adapted to the worst toxins. You have advantage on saving throws against poison and disease. When you use Regenerate, you can end one disease or the poisoned condition on yourself as part of the same action.",
		8:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Wisdom cannot exceed 20.",
		9:  "Hybrid Form — You can manifest partial beast traits without full transformation. As a bonus action, you can enter a 'half-form' that lasts 1 hour (no spell point cost). In half-form, you retain your claws (1d6 damage) and one minor trait from your Beast Shape, but you can still speak and use equipment. You can still layer full Mutagenesis effects on top.",
		10: "Predator's Stride — When you use the Adrenaline component, difficult terrain no longer costs you extra movement. You can also take the Dash action as a bonus action once per turn while Adrenaline is active. Wolf-form Chimeras can move through an ally's space when charging.",
		11: "Greater Mutagen — Your mutagens last twice as long. You can maintain two different mutagens simultaneously (e.g., Feral Muscle + Chitin/Hide). Each costs spell points separately. Your body has grown more efficient at processing reagents.",
		12: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Wisdom cannot exceed 20.",
		13: "Alpha's Roar — Once per short rest, you can use an action to unleash a bestial roar. Each creature of your choice within 30 feet must succeed on a Wisdom saving throw or be frightened of you for 1 minute. A creature can repeat the save at the end of each of its turns. Creatures that succeed are immune to your Alpha's Roar for 24 hours.",
		14: "Adaptive Biology — Your body has evolved from repeated exposure to mutagens. When you enter Primal Form, choose one damage type: acid, cold, fire, lightning, or poison. You have resistance to that damage type until you leave Primal Form. You can change this choice when you finish a long rest.",
		15: "True Hybrid — You can remain in Hybrid Form indefinitely when not in combat. There is no duration limit for your half-form. When you use Primal Form in combat, your mutagen durations are not reduced by taking damage. Your claws and fangs are now considered magical for the purpose of overcoming resistance.",
		16: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Wisdom cannot exceed 20.",
		17: "Primal Surge — When you use the Regenerate component on yourself, you can extend the healing to one willing creature you touch. That creature gains half the healing you receive (rounded down) for the duration. You can use this feature a number of times equal to your Wisdom modifier per long rest.",
		18: "Indomitable Beast — When you are reduced to 0 hit points but not killed outright, you can make a Constitution saving throw (DC 10). On a success, you drop to 1 hit point instead and cannot be reduced below 1 hit point until the start of your next turn. Once you use this feature, you cannot use it again until you finish a long rest.",
		19: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Wisdom cannot exceed 20.",
		20: "Master of the Form — You are the apex predator. Your mutagen durations are tripled (stacking with Greater Mutagen). When you enter Primal Form, you gain +1 AC and +1 to attack and damage rolls with natural weapons. Your natural weapon damage dice increase by one size (1d6 becomes 1d8). You have mastered the beast within.",
	}
}
