package models

import "testing"

func TestValidateSpellDuration(t *testing.T) {
	good := []string{
		"1 min",
		"2 hours",
		"3 days",
		"1 round",
		"2 rounds",
		"concentration",
		"Instantaneous",
		"until dispelled",
		"special",
	}
	for _, s := range good {
		if err := ValidateSpellDuration(s); err != nil {
			t.Errorf("ValidateSpellDuration(%q): want nil, got %v", s, err)
		}
	}
	bad := []string{"", "  ", "up to 1 hour", "1", "foo"}
	for _, s := range bad {
		if err := ValidateSpellDuration(s); err == nil {
			t.Errorf("ValidateSpellDuration(%q): want error, got nil", s)
		}
	}
}

func TestParseSpellRangeFeet(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"120", 120},
		{"120 ft", 120},
		{" 90 FT ", 90},
		{"0", 0},
	}
	for _, tc := range cases {
		got, err := ParseSpellRangeFeet(tc.in)
		if err != nil || got != tc.want {
			t.Errorf("ParseSpellRangeFeet(%q) = (%d, %v), want (%d, nil)", tc.in, got, err, tc.want)
		}
	}
	if _, err := ParseSpellRangeFeet(""); err == nil {
		t.Error("ParseSpellRangeFeet(\"\"): want error")
	}
	if _, err := ParseSpellRangeFeet("x"); err == nil {
		t.Error("ParseSpellRangeFeet(\"x\"): want error")
	}
}
