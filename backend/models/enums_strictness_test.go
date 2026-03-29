package models

import "testing"

func TestParseCreatureSize(t *testing.T) {
	tests := []struct {
		in   string
		want CreatureSize
		ok   bool
	}{
		{"Medium", SizeMedium, true},
		{"medium", SizeMedium, true},
		{"Large", SizeLarge, true},
		{"invalid", "", false},
	}
	for _, tt := range tests {
		got, ok := ParseCreatureSize(tt.in)
		if ok != tt.ok || (tt.ok && got != tt.want) {
			t.Errorf("ParseCreatureSize(%q) = (%v, %v), want (%v, %v)", tt.in, got, ok, tt.want, tt.ok)
		}
	}
}

func TestParseCreatureType(t *testing.T) {
	got, ok := ParseCreatureType("humanoid")
	if !ok || got != TypeHumanoid {
		t.Fatalf("ParseCreatureType(humanoid) = (%v, %v), want (Humanoid, true)", got, ok)
	}
	_, ok = ParseCreatureType("not-a-type")
	if ok {
		t.Fatal("expected failure for invalid creature type")
	}
}

func TestParseMapTokenType(t *testing.T) {
	got, ok := ParseMapTokenType("PC")
	if !ok || got != MapTokenPC {
		t.Fatalf("ParseMapTokenType(PC) = (%v, %v), want (pc, true)", got, ok)
	}
}

func TestParseClassResourceCategory(t *testing.T) {
	got, ok := ParseClassResourceCategory("die_size")
	if !ok || got != ClassResourceCategoryDieSize {
		t.Fatalf("got %v ok=%v", got, ok)
	}
}

func TestParseWeaponDamageCategory(t *testing.T) {
	got, ok := ParseWeaponDamageCategory("base")
	if !ok || got != WeaponDamageCategoryBase {
		t.Fatalf("got %v ok=%v", got, ok)
	}
}
