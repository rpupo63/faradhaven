package models

import "testing"

func TestNormalizeSpellAIRecommendations(t *testing.T) {
	t.Run("all nil", func(t *testing.T) {
		tt, rr, dd := NormalizeSpellAIRecommendations(nil, nil, nil)
		if tt != nil || rr != nil || dd != nil {
			t.Fatalf("expected all nil")
		}
	})
	t.Run("valid triple", func(t *testing.T) {
		ts := "attack"
		rs := "120 ft"
		ds := "1 min"
		tt, rr, dd := NormalizeSpellAIRecommendations(&ts, &rs, &ds)
		if tt == nil || *tt != SpellTypeAttack {
			t.Fatalf("type: got %v", tt)
		}
		if rr == nil || *rr != 120 {
			t.Fatalf("range: got %v", rr)
		}
		if dd == nil || *dd != "1 min" {
			t.Fatalf("duration: got %v", dd)
		}
	})
	t.Run("invalid type dropped", func(t *testing.T) {
		ts := "not_a_type"
		tt, _, _ := NormalizeSpellAIRecommendations(&ts, nil, nil)
		if tt != nil {
			t.Fatalf("expected nil type")
		}
	})
}

func TestNormalizeSpellAIRecommendationsExtras(t *testing.T) {
	s := "dex"
	dt := "fire"
	c, sz := 2, 6
	sa, dtp, dc, ds := NormalizeSpellAIRecommendationsExtras(&s, &dt, &c, &sz)
	if sa == nil || *sa != SaveAttrDEX || dtp == nil || *dtp != DamageFire || dc == nil || *dc != 2 || ds == nil || *ds != 6 {
		t.Fatalf("extras: sa=%v dt=%v dc=%v ds=%v", sa, dtp, dc, ds)
	}
	_, _, _, _ = NormalizeSpellAIRecommendationsExtras(nil, nil, nil, nil)
}
