package services

import "testing"

func TestVaporBladeShadowPointsMax(t *testing.T) {
	if got := VaporBladeShadowPointsMax(14, 1); got != 4 { // +2 dex + 2 prof
		t.Fatalf("dex 14 level 1: got %d want 4", got)
	}
	if got := VaporBladeShadowPointsMax(8, 1); got != 1 { // -1 mod +2 prof → min 1
		t.Fatalf("dex 8 level 1: got %d want 1", got)
	}
	if got := VaporBladeShadowPointsMax(20, 20); got != 11 { // +5 + 6 prof
		t.Fatalf("dex 20 level 20: got %d want 11", got)
	}
}
