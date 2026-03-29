package models

import (
	"regexp"
	"strings"
)

var regexpCollapseSpaces = regexp.MustCompile(`\s+`)

// NormalizeSpellDurationCandidate trims whitespace, collapses spaces, and maps a few legacy
// abbreviations to strings that pass [ValidateSpellDuration]. See also frontend spellMechanics.
func NormalizeSpellDurationCandidate(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	s = regexpCollapseSpaces.ReplaceAllString(s, " ")
	low := strings.ToLower(s)
	switch low {
	case "inst", "inst.":
		return "instant"
	case "perm", "perm.":
		return "permanent"
	case "conc", "conc.":
		return "concentration"
	}
	return s
}

// ParseSpellTypeRecommendation parses a single LLM/user string into a [SpellType]. Returns ok false if empty or unknown.
func ParseSpellTypeRecommendation(s string) (SpellType, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", false
	}
	lower := strings.ToLower(s)
	switch lower {
	case "attack":
		return SpellTypeAttack, true
	case "save":
		return SpellTypeSave, true
	case "effect":
		return SpellTypeEffect, true
	case "healing":
		return SpellTypeHealing, true
	case "utility":
		return SpellTypeUtility, true
	}
	t := SpellType(s)
	if t.IsValid() {
		return t, true
	}
	return "", false
}

// NormalizeSpellAIRecommendations coerces raw LLM strings into typed, validated mechanic fields.
// Nil outputs mean "do not store a recommendation" (invalid or empty input).
func NormalizeSpellAIRecommendations(recommendedType, recommendedRange, recommendedDuration *string) (
	outType *SpellType,
	outRange *int,
	outDuration *string,
) {
	if recommendedType != nil {
		if st, ok := ParseSpellTypeRecommendation(*recommendedType); ok {
			outType = &st
		}
	}
	if recommendedRange != nil {
		s := strings.TrimSpace(*recommendedRange)
		if s != "" {
			if n, err := ParseSpellRangeFeet(s); err == nil {
				outRange = &n
			} else if n, ok := LegacyRangeTextToFeet(s); ok {
				outRange = &n
			}
		}
	}
	if recommendedDuration != nil {
		s := strings.TrimSpace(*recommendedDuration)
		if s != "" {
			cand := s
			if ValidateSpellDuration(cand) != nil {
				cand = NormalizeSpellDurationCandidate(s)
			}
			if ValidateSpellDuration(cand) == nil {
				outDuration = &cand
			}
		}
	}
	return
}

// NormalizeSpellAIRecommendationsExtras coerces save, damage type, and dice pair from LLM output.
func NormalizeSpellAIRecommendationsExtras(
	recommendedSave, recommendedDamageType *string,
	recommendedDamageDiceCount, recommendedDamageDieSize *int,
) (
	outSave *SaveAttribute,
	outDamageType *DamageType,
	outDiceCount, outDieSize *int,
) {
	if recommendedSave != nil {
		if a, ok := ParseSaveAttribute(*recommendedSave); ok {
			outSave = &a
		}
	}
	if recommendedDamageType != nil {
		if d, ok := ParseDamageType(*recommendedDamageType); ok {
			outDamageType = &d
		}
	}
	if recommendedDamageDiceCount != nil && recommendedDamageDieSize != nil {
		c := *recommendedDamageDiceCount
		sz := *recommendedDamageDieSize
		if err := ValidateSpellDamageDicePair(&c, &sz); err == nil {
			outDiceCount = &c
			outDieSize = &sz
		}
	}
	return
}

// ApplyNormalizedAIRecommendations sets AI mechanic recommendation columns from raw LLM strings.
func ApplyNormalizedAIRecommendations(spell *Spell, recommendedType, recommendedRange, recommendedDuration *string) {
	tt, rr, dd := NormalizeSpellAIRecommendations(recommendedType, recommendedRange, recommendedDuration)
	spell.AIRecommendedType = tt
	spell.AIRecommendedRange = rr
	spell.AIRecommendedDuration = dd
}

// ApplyNormalizedAIRecommendationsExtras sets AI columns for save, damage type, and dice from LLM output.
func ApplyNormalizedAIRecommendationsExtras(spell *Spell,
	recommendedSave, recommendedDamageType *string,
	recommendedDamageDiceCount, recommendedDamageDieSize *int,
) {
	s, dt, dc, ds := NormalizeSpellAIRecommendationsExtras(recommendedSave, recommendedDamageType, recommendedDamageDiceCount, recommendedDamageDieSize)
	spell.AIRecommendedSaveAttr = s
	spell.AIRecommendedDamageType = dt
	spell.AIRecommendedDamageDiceCount = dc
	spell.AIRecommendedDamageDieSize = ds
}
