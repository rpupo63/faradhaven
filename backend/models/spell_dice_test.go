package models

import "testing"

func TestParseLegacyDamageDiceString(t *testing.T) {
	c, s, ok := ParseLegacyDamageDiceString("2d6")
	if !ok || c != 2 || s != 6 {
		t.Fatalf("2d6: got %d,%d ok=%v", c, s, ok)
	}
	_, _, ok = ParseLegacyDamageDiceString("2d7")
	if ok {
		t.Fatal("invalid die size should fail")
	}
}

func TestValidateSpellDamageDicePair(t *testing.T) {
	a, b := 2, 6
	if err := ValidateSpellDamageDicePair(&a, &b); err != nil {
		t.Fatal(err)
	}
	if err := ValidateSpellDamageDicePair(nil, nil); err != nil {
		t.Fatal(err)
	}
	x := 2
	if err := ValidateSpellDamageDicePair(&x, nil); err == nil {
		t.Fatal("expected error for partial pair")
	}
}
