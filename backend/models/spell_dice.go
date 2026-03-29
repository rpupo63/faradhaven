package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Standard die sizes allowed for spell damage dice (D&D-style).
// Keep in sync with STANDARD_SPELL_DIE_SIZES in frontend/src/lib/spellMechanics.ts.
var allowedSpellDieSizes = map[int]struct{}{
	4: {}, 6: {}, 8: {}, 10: {}, 12: {}, 20: {}, 100: {},
}

// IsAllowedSpellDieSize reports whether n is a permitted die face count for spells.
func IsAllowedSpellDieSize(n int) bool {
	_, ok := allowedSpellDieSizes[n]
	return ok
}

// FormatSpellDamageDice returns notation like "2d6" for display.
func FormatSpellDamageDice(count, size int) string {
	return fmt.Sprintf("%dd%d", count, size)
}

var legacyDamageDiceRe = regexp.MustCompile(`(?i)^\s*(\d+)\s*d\s*(\d+)\s*$`)

// ParseLegacyDamageDiceString parses a single "XdY" token (e.g. "2d6", "1d8") from legacy storage or LLM text.
func ParseLegacyDamageDiceString(s string) (count int, size int, ok bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, 0, false
	}
	m := legacyDamageDiceRe.FindStringSubmatch(s)
	if len(m) != 3 {
		return 0, 0, false
	}
	c, err1 := strconv.Atoi(m[1])
	f, err2 := strconv.Atoi(m[2])
	if err1 != nil || err2 != nil || c < 1 {
		return 0, 0, false
	}
	if !IsAllowedSpellDieSize(f) {
		return 0, 0, false
	}
	return c, f, true
}

// ValidateSpellDamageDicePair enforces: both nil or both set; count >= 1; allowed die size.
func ValidateSpellDamageDicePair(count, size *int) error {
	if count == nil && size == nil {
		return nil
	}
	if count == nil || size == nil {
		return fmt.Errorf("damage_dice_count and damage_die_size must both be set or both omitted")
	}
	if *count < 1 {
		return fmt.Errorf("damage_dice_count must be at least 1")
	}
	if !IsAllowedSpellDieSize(*size) {
		return fmt.Errorf("damage_die_size must be one of 4, 6, 8, 10, 12, 20, 100")
	}
	return nil
}
