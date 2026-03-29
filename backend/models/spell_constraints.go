package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Duration validation must stay in sync with frontend isValidSpellDuration in
// ../../frontend/src/lib/spellMechanics.ts (same patterns).
var (
	spellDurationTimed   = regexp.MustCompile(`(?i)^\d+\s+(min|minute|minutes|hour|hours|day|days|week|weeks|month|months|year|years)$`)
	spellDurationRounds  = regexp.MustCompile(`(?i)^\d+\s+rounds?$`)
	spellDurationLiteral = regexp.MustCompile(`(?i)^(concentration|instantaneous|instant|until dispelled|until triggered|special|permanent)$`)
)

// ValidateSpellDuration returns an error if s is not a permitted duration string:
// either a timed value like "1 min" or "2 hours", a round count like "1 round",
// or a standard D&D-style keyword (concentration, instantaneous, etc.).
func ValidateSpellDuration(s string) error {
	t := strings.TrimSpace(s)
	if t == "" {
		return fmt.Errorf("duration is empty")
	}
	if spellDurationTimed.MatchString(t) || spellDurationRounds.MatchString(t) || spellDurationLiteral.MatchString(t) {
		return nil
	}
	return fmt.Errorf("duration must be like \"1 min\" / \"2 hours\", \"1 round\", or concentration, instantaneous, until dispelled, until triggered, special, or permanent")
}

// ParseSpellRangeFeet parses AI or user text such as "120", "120 ft", " 90 FT " into a non-negative integer (feet).
func ParseSpellRangeFeet(s string) (int, error) {
	lower := strings.ToLower(strings.TrimSpace(s))
	lower = strings.TrimSpace(strings.TrimSuffix(lower, " feet"))
	lower = strings.TrimSpace(strings.TrimSuffix(lower, " ft"))
	lower = strings.TrimSpace(lower)
	if lower == "" {
		return 0, fmt.Errorf("range is empty")
	}
	n, err := strconv.Atoi(lower)
	if err != nil {
		return 0, fmt.Errorf("range must be a non-negative integer (feet)")
	}
	if n < 0 {
		return 0, fmt.Errorf("range must be non-negative")
	}
	return n, nil
}
