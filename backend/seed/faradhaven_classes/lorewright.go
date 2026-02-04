package faradhaven_classes

// Lorewright returns the Lorewright class seed.
// Inspired by Milo Thatch / Sherlock Holmes: a planner who deciphers ancient languages,
// finds weak points in enemy machinery, records findings, and observes the environment.
// Low damage output, high information value.
func Lorewright() FaradhavenClassSeed {
	return FaradhavenClassSeed{
		Name:           "The Lorewright",
		Description:    "Decipher ancient languages, find weak points in enemy machinery, and record everything in your Field Journal. You don't deal damage—you make sure your party hunts smarter.",
		HitDie:         8,
		PrimaryAbility: "intelligence",
		PhotoURL:       "https://photos-for-apps.s3.us-east-2.amazonaws.com/lorewright.jpg",
		Archetype:      "Support / Information / Controller",
		Concept:        "A scholar who has been around long enough to know the monsters. Max INT and WIS, dump STR. You decipher ancient languages, find weak points in enemy machinery, and record everything. Like a planner who observes—not powerful in a fight, but your party hunts smarter because of you.",
		ClassFeatures: []string{
			"Field Journal: You use components to extract and record information. Decipher reads ancient scripts and obscure languages; Weak Point reveals structural flaws in creatures and machinery. Your findings go into your journal and benefit the whole party.",
			"Keen Observer: You are acutely aware of your environment. You have advantage on Perception and Investigation checks to notice hidden details, track creatures, or spot inconsistencies. Your passive Perception and Investigation scores increase by 2.",
		},
		DnDSkillFocus: []string{"Investigation", "Perception"},
		Proficiencies: "Simple weapons, Light armor, Crossbows (hand or light)",
		SkillChoice:   []string{"History", "Nature", "Insight", "Survival"},
		Tools:         []string{"Cartographer's Tools", "Thieves' Tools (for locks and traps)"},
		SavingThrows:  []string{"Intelligence", "Wisdom"},
		StartingEquip: []string{"Light crossbow or quarterstaff", "Leather armor", "Field journal (blank)", "Magnifying lens", "Ink and quill", "Explorer's pack"},
		LevelFeatures: lorewrightLevelFeatures(),
	}
}

func lorewrightLevelFeatures() map[int]string {
	return map[int]string{
		1:  "Field Journal, Keen Observer — You gain the core class features. You can spend 2 spell points to use Decipher on a written source within 30 feet: you understand its meaning and can relay it to your party. You can spend 3 spell points to use Weak Point on a creature or object you can see: you learn one vulnerability (resistance, immunity, or a flaw that grants allies advantage on the next attack).",
		2:  "Record — When you use Decipher or Weak Point, you can spend 2 additional spell points to add the finding to your Field Journal. Any ally who has seen your journal in the last 24 hours gains the benefit of that finding (e.g., advantage on the next attack against a creature whose weak point you recorded). Your journal can hold a number of entries equal to your Intelligence modifier + your Field Scholar level.",
		3:  "Observe — You can spend 5 spell points as an action to enhance your perception for 1 minute. You automatically detect hidden doors, traps, and creatures within 30 feet (you know they exist, not exact location). You can maintain one Observe effect at a time.",
		4:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Intelligence and Wisdom cannot exceed 20.",
		5:  "Linguist's Eye — Your Decipher component works on spoken languages you do not know. You can spend 5 spell points to understand and respond in any language for 10 minutes. You can also use Decipher to identify magical writings and their school of magic.",
		6:  "Weak Point Mastery — When you use Weak Point on a creature, you learn two vulnerabilities instead of one. When an ally attacks a creature whose weak point you have recorded, they add your Intelligence modifier to the damage of the first hit.",
		7:  "Environmental Awareness — Your Observe component extends to 60 feet. Additionally, you cannot be surprised while conscious, and you have advantage on initiative rolls.",
		8:  "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Intelligence and Wisdom cannot exceed 20.",
		9:  "Machinery Savant — When you use Weak Point on a construct, vehicle, or mechanical object, you learn how to disable or sabotage it. An ally can use your finding to impose disadvantage on the target's next attack or cause it to malfunction (DM's discretion). You have advantage on Arcana checks related to machinery.",
		10: "Chronicler's Blessing — When you add a finding to your Field Journal with the Record modifier, you can choose one ally. That ally gains advantage on all checks and saves related to that finding for 1 hour. You can maintain this on a number of allies equal to your Wisdom modifier.",
		11: "Decipher Mastery — Your Decipher component can extract information from non-written sources: carvings, murals, or even the scars on a creature's body. You learn one piece of useful information (history, weakness, or recent activity) per 5 spell points spent.",
		12: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Intelligence and Wisdom cannot exceed 20.",
		13: "Preternatural Deduction — When you observe a creature for at least 1 minute (or have recorded its weak point), you can spend 10 spell points as an action to deduce its next likely action. You and your allies have advantage on saving throws against that creature's next attack or ability.",
		14: "Living Ledger — Your Field Journal entries never expire. You can recall any recorded finding as a bonus action without spending spell points. Your journal can hold a number of entries equal to twice your Field Scholar level.",
		15: "Observe Mastery — Your Observe component reveals exact locations of hidden creatures and traps within range. You can maintain two Observe effects simultaneously. While Observing, you have advantage on Insight checks against creatures in the area.",
		16: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Intelligence and Wisdom cannot exceed 20.",
		17: "Monster Lore — You have studied the beasts of Faradhaven. When you use Weak Point on a creature type you have encountered before (recorded or not), you learn all of its resistances, immunities, and vulnerabilities. You can spend 5 spell points to share this knowledge with your party as a bonus action.",
		18: "Eidetic Recall — You no longer need to add findings to your journal to remember them. Any creature or object you have used Decipher or Weak Point on is permanently known to you. You can relay any past finding to an ally as a bonus action.",
		19: "Ability Score Improvement — Increase one ability score by 2, or two ability scores by 1. Intelligence and Wisdom cannot exceed 20.",
		20: "Master Inquisitive — You are the party's eyes and mind. You regain 20 spell points when you roll initiative and have none remaining. Your Keen Observer bonus to passive Perception and Investigation increases to 5. When you use Weak Point, you learn all vulnerabilities and the creature has disadvantage on its next attack roll. Your Field Journal entries are considered magical and cannot be altered or destroyed without your consent.",
	}
}
