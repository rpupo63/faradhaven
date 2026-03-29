package models

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	legacyReFeet     = regexp.MustCompile(`(?i)(\d+)\s*(?:ft|feet)\b`)
	legacyReDigits   = regexp.MustCompile(`^\s*(\d+)\s*$`)
	legacyReHasDigit = regexp.MustCompile(`\d`)
)

// LegacyRangeTextToFeet parses free-text range (e.g. "120 ft", "Self (15 ft radius)", "Touch")
// into a feet value. Used by data migration and AI recommendation normalization.
// Returns ok false when the text should be treated as NULL (unparseable).
func LegacyRangeTextToFeet(s string) (int, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	if n, err := ParseSpellRangeFeet(s); err == nil {
		return n, true
	}

	ls := strings.ToLower(s)
	switch ls {
	case "touch":
		return 5, true
	}
	if ls == "self" || ls == "self-centered" {
		return 0, true
	}
	if strings.Contains(ls, "touch") && !legacyReHasDigit.MatchString(s) {
		return 5, true
	}
	if strings.Contains(ls, "self") && !legacyReHasDigit.MatchString(s) {
		return 0, true
	}

	best := 0
	for _, m := range legacyReFeet.FindAllStringSubmatch(s, -1) {
		if len(m) >= 2 {
			n, _ := strconv.Atoi(m[1])
			if n > best {
				best = n
			}
		}
	}
	if best > 0 {
		return best, true
	}

	if m := legacyReDigits.FindStringSubmatch(s); len(m) >= 2 {
		n, err := strconv.Atoi(m[1])
		if err == nil && n >= 0 {
			return n, true
		}
	}

	return 0, false
}
