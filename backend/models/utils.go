package models

import (
	"regexp"
	"strconv"
	"strings"
)

// ParseWeightString parses a weight string (e.g., "1 lb.", "2.5 kg") and returns
// its value as a float64 in pounds. It assumes "lb" as the default unit if no unit is specified
// or if an unrecognized unit is found.
// This function needs to be robust for various inputs.
func ParseWeightString(weightStr string) float64 {
	weightStr = strings.ToLower(strings.TrimSpace(weightStr))

	// Regex to capture number and optional unit
	re := regexp.MustCompile(`^(\d*\.?\d+)\s*(lb|lbs|kg)?\.?$`)
	matches := re.FindStringSubmatch(weightStr)

	if len(matches) < 2 {
		// No valid number found, or empty string. Return 0.
		return 0.0
	}

	valueStr := matches[1]
	unit := ""
	if len(matches) >= 3 {
		unit = matches[2]
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0.0 // Could not parse number, return 0
	}

	// Convert to pounds
	switch unit {
	case "kg":
		return value * 2.20462 // 1 kg = 2.20462 lbs
	case "lb", "lbs", "": // Default to pounds if unit is lb, lbs, or not specified
		return value
	default:
		// If an unrecognized unit is found, treat the value as pounds
		return value
	}
}
