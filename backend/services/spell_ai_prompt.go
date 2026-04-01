package services

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rpupo63/unified-personal-site-backend/models"
)

// SpellAIOpinionResponseSchema is the JSON schema passed to the LLM for structured spell review output.
// Keep in sync with [SpellOpinion] fields and the DATABASE CONSTRAINTS block in spellAIOpinionPromptTemplate.
const SpellAIOpinionResponseSchema = `{
		"type": "object",
		"properties": {
			"description_opinion": {"type": "string"},
			"damage_opinion": {"type": "string"},
			"effect_opinion": {"type": "string"},
			"overall_verdict": {"type": "string"},
			"recommended_name": {"type": ["string", "null"]},
			"recommended_description": {"type": ["string", "null"]},
			"recommended_type": {
				"type": ["string", "null"],
				"description": "Exactly one of: Attack, Save, Effect, Healing, Utility (exact capitalization)."
			},
			"recommended_range": {
				"type": ["string", "null"],
				"description": "Digits only: non-negative integer feet as a string, e.g. 120 or 0. No units or prose."
			},
			"recommended_duration": {
				"type": ["string", "null"],
				"description": "Timed (e.g. 1 min, 2 hours), rounds (e.g. 1 round), or keyword: concentration, instantaneous, instant, until dispelled, until triggered, special, permanent."
			},
			"recommended_damage_dice_count": {
				"type": ["integer", "null"],
				"description": "Number of dice (e.g. 2 for 2d6). Must be used with recommended_damage_die_size."
			},
			"recommended_damage_die_size": {
				"type": ["integer", "null"],
				"description": "Die faces: 4, 6, 8, 10, 12, 20, or 100 (must match recommended_damage_dice_count)."
			},
			"recommended_damage_type": {
				"type": ["string", "null"],
				"description": "Exactly one of: Slashing, Piercing, Bludgeoning, Fire, Cold, Lightning, Thunder, Poison, Acid, Necrotic, Radiant, Force, Psychic."
			},
			"recommended_save_attr": {
				"type": ["string", "null"],
				"description": "Exactly one of: STR, DEX, CON, INT, WIS, CHA."
			}
		},
		"required": ["description_opinion", "damage_opinion", "effect_opinion", "overall_verdict"]
	}`

// spellAIOpinionPromptTemplate is the full rubric + instructions for spell AI review.
// Placeholders: %s spell block, %s components block, %d level, %s damage summary, %s type, %s range, %s duration.
// The DATABASE CONSTRAINTS subsection must match server rules: [models.ValidateSpellDuration], [models.ParseSpellRangeFeet], [models.SpellType].
const spellAIOpinionPromptTemplate = `
	You are an expert magical theorist and Game Master for the Faradhaven TTRPG.
	Your task is to provide a technical and balance-focused opinion on a player-crafted spell, AND propose specific edits to improve it.

	### FARADHAVEN SPELL SCALING RULES (CRITICAL CONTEXT):
	- Spells are built from modular components.
	- The 'Level' of a spell is equal to the TOTAL number of components used.
	- SCALE: A Level 2 or 3 spell in Faradhaven is extremely weak, comparable to a D&D 5e Cantrip. 
	- SCALE: A Level 5 spell is roughly equivalent to a 1st or 2nd level D&D spell.
	- SCALE: A Level 10 spell is a powerful mid-tier effect.
	- Do NOT compare Faradhaven levels directly to D&D levels (e.g., a Faradhaven Level 3 spell is NOT a D&D 3rd-level spell like Fireball; it is a Cantrip).

	### COMPONENT PILLARS (what each category means):
	- Forma: Shape — geometric delivery (Projectile, Beam, Nova, Wall, Zone, Cone, Aura, Touch, etc.).
	- Scopus: Targeting — anchor point (Target, Self, Ground, Chain, etc.).
	- Essentia: Domain / matter — substance or abstract being manipulated (Ignis, Umbra, Vita, Arcanum, pathos tones, etc.).
	- Actio: Kinetic verb — what the magic does to Essentia or the target (Push, Pull, Bind, Create, etc.).
	- Magnitudo: Scale modifiers — power dials (Strong, Increase, Extreme, etc.).

	### FIELDS ON EACH COMPONENT (use all of them when judging the spell):
	- Symbol: Short alchemical-style code for the setting's notation; optional flavor, can help disambiguate similar names.
	- Tier: 1 = basic building block; 2 = advanced — often implies broader or rarer capability within the same component count.
	- Element: Usually set on Essentia — the game's suggested damage or energy type (e.g. fire, necrotic, force). Align recommended damage_type and flavor with Essentia Element when present.
	- Description: Canonical meaning of that component in Faradhaven; the spell's name, description, range, duration, and damage should not contradict these texts.

	### DATABASE CONSTRAINTS FOR RECOMMENDED EDITS (MANDATORY — SERVER VALIDATION):
	The API persists recommended_* fields under strict rules (invalid values are dropped or rejected). Your JSON recommendations MUST conform or they will not be saved as-is.
	- **recommended_type**: EXACTLY one of these five strings (this capitalization): Attack, Save, Effect, Healing, Utility. No other labels (e.g. not "Buff", "Debuff", "Enchantment").
	- **recommended_range**: A non-negative integer as a STRING of digits only, meaning distance in **feet** (e.g. "120", "60", "0"). Use "0" for self-centered or touch-adjacent effects where range is effectively zero. Do NOT write prose ("120 feet away"), compound phrases ("Self (15 ft)"), or units inside the string — digits only so the database integer column accepts the value.
	- **recommended_duration**: MUST match one of these patterns (trim whitespace; typical casing is fine for keywords):
	  • Timed: "<positive integer> <unit>" where unit is one of: min, minute, minutes, hour, hours, day, days, week, weeks, month, months, year, years. Examples: "1 min", "2 hours", "3 days".
	  • Rounds: "<positive integer> round" or "<positive integer> rounds". Example: "1 round".
	  • Keywords (exact set): concentration, instantaneous, instant, until dispelled, until triggered, special, permanent.
	  Do NOT output free text like "Until you dismiss it", "See text", or D&D book phrases that do not match the above — use "special" or "permanent" when appropriate, or a timed form.
	- **recommended_damage_type**: EXACTLY one of these strings (capitalization as shown): Slashing, Piercing, Bludgeoning, Fire, Cold, Lightning, Thunder, Poison, Acid, Necrotic, Radiant, Force, Psychic. Align with Essentia Element when possible.
	- **recommended_damage_dice_count** and **recommended_damage_die_size**: Integers for damage dice (e.g. 2 and 6 for 2d6). Must both be null or both set. Die size must be one of: 4, 6, 8, 10, 12, 20, 100. Count must be at least 1. Do NOT use a single string like "2d6" — use two integers.
	- **recommended_save_attr**: If suggesting a save, EXACTLY one of: STR, DEX, CON, INT, WIS, CHA.

	When the spell has **no** If/Then/Therefore (Logica) components, the component list below is reordered Forma → Scopus → Essentia → Actio → Magnitudo (then by name) for readability.
	When the spell **uses Logica**, a **LOGIC FLOW** section preserves **exact cast order**: each **Phase** is one narrative beat; If/Then/Therefore mark boundaries (e.g. water pooled in phase 1, then freeze in phase 2). Same Essentia name in two phases means two different story beats, not one pooled multiset.

	### SPELL TO REVIEW:
	%s

	### COMPONENTS USED (full records):
	%s

	### YOUR INSTRUCTIONS:
	Provide your opinion on the following three aspects in a structured format:
	1. Description: Does the flavor text match the components? Is it evocative?
	2. Damage: Based on the level (%d components) and the scaling rules above, is the damage (%s) appropriate, too high, or too low?
	3. Effect: Does the spell type (%s) and its mechanics (Range: %s, Duration: %s) align with the intended Actio and Essentia?

	IMPORTANT: Propose specific "Recommended" values only when a field needs improvement. If a field is fine, use null for that key. Every non-null recommended_type, recommended_range, recommended_duration, recommended_damage_* (type + dice pair + save) MUST satisfy the **DATABASE CONSTRAINTS** section above — the server enforces the same rules.

	Return a JSON object with the following fields:
	- description_opinion (string)
	- damage_opinion (string)
	- effect_opinion (string)
	- overall_verdict (string: A final summary recommendation for the GM)
	- recommended_name (string or null)
	- recommended_description (string or null)
	- recommended_type (string or null)
	- recommended_range (string or null)
	- recommended_duration (string or null)
	- recommended_damage_dice_count (integer or null)
	- recommended_damage_die_size (integer or null)
	- recommended_damage_type (string or null)
	- recommended_save_attr (string or null)
	`

// BuildSpellAIOpinionPrompt is the single source of truth for spell review LLM prompts.
// Use this from [SpellAIService.GetSpellOpinion]; batch tools should call GetSpellOpinion rather than reimplementing prompts.
func BuildSpellAIOpinionPrompt(spell *models.Spell, components []models.Component) string {
	spellContext := FormatSpellContextForSpellAI(spell)
	componentContext := FormatComponentContextForSpellAI(components)
	return fmt.Sprintf(spellAIOpinionPromptTemplate,
		spellContext,
		componentContext,
		spell.Level,
		formatSpellDamageForPrompt(spell),
		spell.Type,
		formatSpellRangeFeet(spell.Range),
		derefStr(spell.Duration, "N/A"),
	)
}

// FormatSpellContextForSpellAI renders spell fields for [BuildSpellAIOpinionPrompt].
func FormatSpellContextForSpellAI(spell *models.Spell) string {
	if spell == nil {
		return "(nil spell)\n"
	}
	return fmt.Sprintf("Name: %s\nDescription: %s\nLevel: %d\nType: %s\nRange (feet): %s\nDamage: %s (%s)",
		spell.Name, spell.Description, spell.Level, spell.Type,
		formatSpellRangeFeet(spell.Range),
		formatSpellDamageForPrompt(spell), derefDamageType(spell.DamageType))
}

func formatSpellDamageForPrompt(spell *models.Spell) string {
	if spell.DamageDiceCount != nil && spell.DamageDieSize != nil {
		return models.FormatSpellDamageDice(*spell.DamageDiceCount, *spell.DamageDieSize)
	}
	return "None"
}

func derefDamageType(dt *models.DamageType) string {
	if dt == nil {
		return "None"
	}
	return string(*dt)
}

// FormatComponentContextForSpellAI renders component records for [BuildSpellAIOpinionPrompt].
func FormatComponentContextForSpellAI(components []models.Component) string {
	if len(components) == 0 {
		return "(no components linked to this spell)\n"
	}
	if spellSequenceHasLogica(components) {
		return formatComponentContextWithLogicPhases(components)
	}
	var sb strings.Builder
	for _, c := range sortComponentsForContext(components) {
		sb.WriteString(formatOneComponentLine(c))
	}
	return sb.String()
}

func spellSequenceHasLogica(components []models.Component) bool {
	for _, c := range components {
		if c.Category == models.CategoryLogica {
			return true
		}
	}
	return false
}

// formatComponentContextWithLogicPhases groups non-Logica runs into numbered phases separated by If/Then/Therefore in chain order.
func formatComponentContextWithLogicPhases(components []models.Component) string {
	var sb strings.Builder
	sb.WriteString("LOGIC FLOW (cast order — phases are separate narrative beats; connectors separate earlier magic from later magic):\n\n")
	phaseNum := 0
	var phase []models.Component
	flushPhase := func() {
		if len(phase) == 0 {
			return
		}
		phaseNum++
		sb.WriteString(fmt.Sprintf("**Phase %d** (components for this beat, in order):\n", phaseNum))
		for _, c := range phase {
			sb.WriteString(formatOneComponentLine(c))
		}
		sb.WriteString("\n")
		phase = nil
	}
	for _, c := range components {
		if c.Category == models.CategoryLogica {
			flushPhase()
			sb.WriteString(fmt.Sprintf("**── %s** (Logica — separates phases; order relative to surrounding phases matters)\n", c.Name))
			sb.WriteString(formatOneComponentLine(c))
			sb.WriteString("\n")
			continue
		}
		phase = append(phase, c)
	}
	flushPhase()
	return sb.String()
}

func derefStr(val *string, def string) string {
	if val == nil {
		return def
	}
	return *val
}

func formatSpellRangeFeet(p *int) string {
	if p == nil {
		return "N/A"
	}
	return fmt.Sprintf("%d ft", *p)
}

// categoryPriority orders components for prompt readability (spell grammar).
func categoryPriority(c models.ComponentCategory) int {
	switch c {
	case models.CategoryForma:
		return 0
	case models.CategoryScopus:
		return 1
	case models.CategoryEssentia:
		return 2
	case models.CategoryActio:
		return 3
	case models.CategoryMagnitudo:
		return 4
	case models.CategoryLogica:
		return 5
	default:
		return 99
	}
}

func sortComponentsForContext(components []models.Component) []models.Component {
	if len(components) <= 1 {
		out := make([]models.Component, len(components))
		copy(out, components)
		return out
	}
	for _, c := range components {
		if c.Category == models.CategoryLogica {
			out := make([]models.Component, len(components))
			copy(out, components)
			return out
		}
	}
	out := make([]models.Component, len(components))
	copy(out, components)
	sort.Slice(out, func(i, j int) bool {
		pi, pj := categoryPriority(out[i].Category), categoryPriority(out[j].Category)
		if pi != pj {
			return pi < pj
		}
		if out[i].Tier != out[j].Tier {
			return out[i].Tier < out[j].Tier
		}
		return out[i].Name < out[j].Name
	})
	return out
}

func formatOneComponentLine(c models.Component) string {
	tierLabel := "basic (tier 1)"
	if c.Tier >= 2 {
		tierLabel = fmt.Sprintf("advanced (tier %d)", c.Tier)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("- **%s — %s**", c.Category, c.Name))
	if c.Symbol != "" {
		sb.WriteString(fmt.Sprintf(" · symbol %q", c.Symbol))
	}
	sb.WriteString(fmt.Sprintf(" · %s\n", tierLabel))
	if c.Element != "" {
		sb.WriteString(fmt.Sprintf("  - Element: %s\n", c.Element))
	}
	desc := strings.TrimSpace(c.Description)
	if desc == "" {
		desc = "(no description stored for this component)"
	}
	sb.WriteString(fmt.Sprintf("  - Description: %s\n", desc))
	return sb.String()
}
