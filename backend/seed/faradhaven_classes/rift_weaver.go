package faradhaven_classes

import (
	"github.com/rpupo63/unified-personal-site-backend/seed/seedmedia"
)

// RiftWeaver returns the Rift Weaver class seed
func RiftWeaver() FaradhavenClassSeed {
	return FaradhavenClassSeed{
		Name:                "The Rift Weaver",
		Description:         "Open rifts to elemental planes and channel raw fire, ice, lightning, and earth into devastating evocations. You're the party's artillery—stand back and obliterate the battlefield with area destruction.",
		HitDie:              8,
		PrimaryAbility:      "intelligence",
		PhotoURL:            seedmedia.URL("rift_weaver.jpg"),
		Archetype:           "Full Caster / Evoker / Damage",
		Concept:             "A scholar who learned to open rifts to elemental planes. They draw power from fire, ice, lightning, and earth—commanding raw elemental forces for devastating evocation magic. You are the party's artillery.",
		DnDSkillFocus:       []string{"Arcana", "Nature"},
		Proficiencies:       "Simple weapons, Daggers, Quarterstaffs, Light Armor",
		SkillChoice:         []string{"Arcana", "History", "Nature", "Religion"},
		Tools:               []string{"Arcane Focus (orb or wand)"},
		SavingThrows:        []string{"Intelligence", "Wisdom"},
		AutomaticEquipNames: []string{"Spellbook (grimoire of elemental formulae)"},
		AutomaticItemNames:  []string{"Scholar's robes (light armor)", "Explorer's pack"},
		EquipmentChoices: []EquipmentChoiceSeed{
			{
				Instruction: "Choose your focus",
				Options: []EquipmentOptionSeed{
					{Description: "Aetheric Tuning Fork", ItemNames: []string{"Aetheric Tuning Fork"}},
					{Description: "Arcane Focus (Orb)", ItemNames: []string{"Arcane Focus (Orb)"}},
				},
			},
			{
				Instruction: "Choose your elemental protection",
				Options: []EquipmentOptionSeed{
					{Description: "Smelter's Gloves (Fire Resist)", ItemNames: []string{"Smelter's Gloves"}},
					{Description: "Protective Goggles (Vision)", ItemNames: []string{"Protective Goggles"}},
				},
			},
		},
		LevelFeatures:    riftWeaverLevelFeatures(),
		LevelProgression: riftWeaverLevelProgression(),
		ComponentPool: []string{
			// Forma (Shape)
			"Projectile", "Beam", "Nova", "Wall", "Zone", "Cone", "Aura",
			// Scopus (Targeting)
			"Target", "Self", "Ground", "Chain",
			// Essentia (Matter & Energy)
			"Ignis", "Aqua", "Terra", "Aer", "Fulgur", "Spatium", "Arcanum",
			// Actio (Verbs)
			"Create", "Destroy", "Push", "Pull", "Crush", "Pierce",
			// Magnitudo (Modifiers)
			"Increase", "Decrease", "Strong", "Weak", "Extreme",
			// Logica (sequential links)
			"If", "Then", "Therefore",
		},
		ResourceDefinitions: []ResourceDefinitionSeed{
			{Key: "spell_points", DisplayName: "Spell Points", Category: "pool", Description: "A pool of magical energy for casting spells.", IsTrackable: true, RestoreOnShortRest: true, RestoreOnLongRest: true, DisplayOrder: 1},
			{Key: "overchannel_uses", DisplayName: "Overchannel", Category: "pool", Description: "1/Short Rest: Spend 3 extra spell points to add your INT modifier to a damage roll.", IsTrackable: true, RestoreOnShortRest: true, RestoreOnLongRest: true, DisplayOrder: 2},
			{Key: "elemental_surge_uses", DisplayName: "Elemental Surge", Category: "pool", Description: "1/Short Rest: Double a spell's range/area, or change its damage type.", IsTrackable: true, RestoreOnShortRest: true, RestoreOnLongRest: true, DisplayOrder: 3},
			{Key: "elemental_overload_uses", DisplayName: "Elemental Overload", Category: "pool", Description: "1/Long Rest: Spend all remaining spell points (min 20) for a 30ft elemental burst.", IsTrackable: true, RestoreOnLongRest: true, DisplayOrder: 4},
		},
	}
}

func riftWeaverLevelProgression() map[int]ClassLevelSeed {
	// Rift Weaver is a full caster - cantrips and spells known scale with level.
	// spell_points must match ClassLevel.MaxSpellPoints (maxSpellPointsByLevel) for CharacterResource rows.
	cantrips2, cantrips3, cantrips4, cantrips5 := 2, 3, 4, 5
	spells2, spells3, spells4, spells5, spells6, spells7, spells8 := 2, 3, 4, 5, 6, 7, 8
	spells9, spells10, spells11, spells12, spells13, spells14, spells15 := 9, 10, 11, 12, 13, 14, 15
	// overchannel unlocks at level 3, elemental_surge at level 5, elemental_overload at level 17
	sp := func(level int) map[string]int {
		overchannelUses := 0
		if level >= 3 {
			overchannelUses = 1
		}
		elementalSurgeUses := 0
		if level >= 5 {
			elementalSurgeUses = 1
		}
		elementalOverloadUses := 0
		if level >= 17 {
			elementalOverloadUses = 1
		}
		return map[string]int{
			"spell_points":            maxSpellPointsByLevel(level),
			"overchannel_uses":        overchannelUses,
			"elemental_surge_uses":    elementalSurgeUses,
			"elemental_overload_uses": elementalOverloadUses,
		}
	}
	return map[int]ClassLevelSeed{
		1:  {CantripsKnown: &cantrips2, SpellsKnown: &spells2, Resources: sp(1)},
		2:  {CantripsKnown: &cantrips2, SpellsKnown: &spells3, Resources: sp(2)},
		3:  {CantripsKnown: &cantrips2, SpellsKnown: &spells4, Resources: sp(3)},
		4:  {CantripsKnown: &cantrips3, SpellsKnown: &spells5, Resources: sp(4)},
		5:  {CantripsKnown: &cantrips3, SpellsKnown: &spells6, Resources: sp(5)},
		6:  {CantripsKnown: &cantrips3, SpellsKnown: &spells7, Resources: sp(6)},
		7:  {CantripsKnown: &cantrips3, SpellsKnown: &spells8, Resources: sp(7)},
		8:  {CantripsKnown: &cantrips3, SpellsKnown: &spells9, Resources: sp(8)},
		9:  {CantripsKnown: &cantrips3, SpellsKnown: &spells10, Resources: sp(9)},
		10: {CantripsKnown: &cantrips4, SpellsKnown: &spells11, Resources: sp(10)},
		11: {CantripsKnown: &cantrips4, SpellsKnown: &spells12, Resources: sp(11)},
		12: {CantripsKnown: &cantrips4, SpellsKnown: &spells12, Resources: sp(12)},
		13: {CantripsKnown: &cantrips4, SpellsKnown: &spells13, Resources: sp(13)},
		14: {CantripsKnown: &cantrips4, SpellsKnown: &spells13, Resources: sp(14)},
		15: {CantripsKnown: &cantrips4, SpellsKnown: &spells14, Resources: sp(15)},
		16: {CantripsKnown: &cantrips4, SpellsKnown: &spells14, Resources: sp(16)},
		17: {CantripsKnown: &cantrips5, SpellsKnown: &spells15, Resources: sp(17)},
		18: {CantripsKnown: &cantrips5, SpellsKnown: &spells15, Resources: sp(18)},
		19: {CantripsKnown: &cantrips5, SpellsKnown: &spells15, Resources: sp(19)},
		20: {CantripsKnown: &cantrips5, SpellsKnown: &spells15, Resources: sp(20)},
	}
}

func riftWeaverLevelFeatures() map[int][]FeatureSeed {
	return map[int][]FeatureSeed{
		1: {
			{Name: "Elemental Channeling", Description: "You spend spell points to channel universal components. Cost: 2 Spell Points per Component. Damage: 1d10 per Elemental Component (stacking identical components yields diminishing returns). Area: 5ft radius per Shape Component. Save DC: 8 + Proficiency + Intelligence modifier. Attack: Proficiency + Intelligence modifier."},
			{Name: "Elemental Sight", Description: "You can perceive elemental auras within 30 feet. Creatures or objects with elemental affinity appear as a faint glow. You have advantage on Arcana checks to identify elemental effects."},
		},
		2:  {{Name: "Ritual Attunement", Description: "As a ritual (10 minutes), you can cast a spell using components you know without spending spell points. The effect's duration is halved when cast this way."}},
		3:  {{Name: "Overchannel", Description: "When you spend spell points on a damage-dealing component, you can spend 3 additional spell points to add your Intelligence modifier to the damage. Once per short rest."}},
		5:  {{Name: "Elemental Surge", Description: "Once per short rest, when you cast a spell, you can double its range or area of effect. Alternatively, you can change the damage type of one component to another element you know."}},
		6:  {{Name: "Potent Evocation", Description: "Your elemental damage ignores resistance to one damage type of your choice (Fire [Ignis], Cold [Aqua], Lightning [Fulgur], or Acid/Earth [Acidum/Terra]). You can change this choice after a long rest."}},
		7:  {{Name: "Elemental Resilience", Description: "You gain resistance to one damage type of your choice (Fire, Cold, Lightning, or Acid/Earth). As a reaction when you take that damage, you can spend 5 spell points to reduce the damage taken by half again (after resistance)."}},
		9:  {{Name: "Sculpt Spells", Description: "When you cast a spell that affects an area, you can choose a number of creatures equal to 1 + your Intelligence modifier. Those creatures are unaffected by the spell: they automatically succeed on their saving throw and take no damage."}},
		10: {{Name: "Elemental Amplification", Description: "When you deal elemental damage, you can spend 5 spell points to add 2d6 damage of the same type to one target. Additionally, elemental buffs you cast on allies last for 1 hour."}},
		11: {{Name: "Elemental Barrage", Description: "When you cast a single-target spell, you can spend 5 additional spell points to target a second creature within 30 feet of the first with the same effect."}},
		13: {{Name: "Elemental Veil", Description: "As an action, you spend 15 spell points to wreathe yourself in elemental energy for 1 minute. Creatures that hit you with a melee attack take 2d8 damage of your chosen element (Fire, Cold, Lightning, or Acid/Earth)."}},
		14: {{Name: "Dual Element", Description: "When you cast a spell, you can spend 10 additional spell points to add a second elemental component to the effect for free. Both elemental damage types apply simultaneously."}},
		15: {{Name: "Elemental Mastery", Description: "You can maintain concentration on two elemental effects simultaneously. Additionally, targets have disadvantage on the first saving throw they make against your spells."}},
		17: {{Name: "Elemental Overload", Description: "Once per long rest, you can spend all remaining spell points (minimum 20) to create a 30ft radius burst. Each creature of your choice takes 1d10 damage of your choice per 2 spell points spent (Dexterity save for half)."}},
		18: {{Name: "Rift Dominion", Description: "When an enemy fails a save against your spell, you can spend 10 spell points to extend the duration by 1 minute or cause the damage to deal its maximum possible value. Once per target per long rest."}},
		20: {{Name: "Archmage of the Rifts", Description: "You regain 20 spell points when you roll initiative and have none remaining. Your Intelligence score increases by 4, to a maximum of 24. Your elemental damage ignores immunity once per long rest per creature."}},
	}
}
