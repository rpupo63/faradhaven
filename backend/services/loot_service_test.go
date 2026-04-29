package services

import (
	"math"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/models"
)

func TestDeriveLevelBand(t *testing.T) {
	tests := []struct {
		level int
		want  string
	}{
		{1, "novice"},
		{4, "novice"},
		{5, "adventurer"},
		{9, "adventurer"},
		{10, "veteran"},
		{14, "veteran"},
		{15, "legend"},
	}

	for _, tt := range tests {
		if got := deriveLevelBand(tt.level); got != tt.want {
			t.Fatalf("deriveLevelBand(%d)=%s, want %s", tt.level, got, tt.want)
		}
	}
}

func TestValidateLootInputs(t *testing.T) {
	if err := validateLootInputs("room", "dungeon"); err != nil {
		t.Fatalf("validateLootInputs valid case returned err: %v", err)
	}

	cases := []struct {
		name      string
		source    string
		roomTheme string
	}{
		{name: "invalid source", source: "crate", roomTheme: "dungeon"},
		{name: "invalid room theme", source: "room", roomTheme: "space"},
	}

	for _, tc := range cases {
		if err := validateLootInputs(tc.source, tc.roomTheme); err == nil {
			t.Fatalf("%s: expected error, got nil", tc.name)
		}
	}
}

func TestApplyLootModifiers(t *testing.T) {
	base := lootParams{
		ItemMin:   1,
		ItemMax:   3,
		WeaponMin: 0,
		WeaponMax: 1,
		GoldMin:   100,
		GoldMax:   300,
	}

	adjusted, rarityShift, goldMult := applyLootModifiers(
		base,
		lootModifier{ItemDelta: 1, WeaponDelta: 1, GoldMultiplier: 1.25, RarityShift: 3},
		lootModifier{ItemDelta: -1, WeaponDelta: 0, GoldMultiplier: 0.8, RarityShift: -1},
	)

	if adjusted.ItemMin != 1 || adjusted.ItemMax != 3 {
		t.Fatalf("unexpected item range: %+v", adjusted)
	}
	if adjusted.WeaponMin != 1 || adjusted.WeaponMax != 2 {
		t.Fatalf("unexpected weapon range: %+v", adjusted)
	}
	if adjusted.GoldMin != 100 || adjusted.GoldMax != 300 {
		t.Fatalf("unexpected gold range after multipliers: %+v", adjusted)
	}
	if rarityShift != 2 {
		t.Fatalf("unexpected rarity shift %d", rarityShift)
	}
	if goldMult <= 0 {
		t.Fatalf("gold multiplier must stay positive, got %f", goldMult)
	}
}

func TestParseCostToCopper(t *testing.T) {
	tests := []struct {
		cost    string
		want    int64
		wantErr bool
	}{
		{"50 cp", 50, false},
		{"5 sp", 50, false},
		{"3 gp", 300, false},
		{"1 pp", 1000, false},
		{"bad", 0, true},
	}

	for _, tt := range tests {
		got, err := parseCostToCopper(tt.cost)
		if tt.wantErr && err == nil {
			t.Fatalf("expected error for %q", tt.cost)
		}
		if !tt.wantErr && (err != nil || got != tt.want) {
			t.Fatalf("parseCostToCopper(%q)=%d, err=%v; want=%d", tt.cost, got, err, tt.want)
		}
	}
}

func TestSampleGaussianBudgetNonNegative(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 50; i++ {
		got := sampleGaussianBudget(rng, 1000)
		if got < 250 {
			t.Fatalf("sampleGaussianBudget floor violated: %d", got)
		}
	}
}

func TestBudgetFillStopsAtNoAffordable(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	candidates := []lootCandidate{
		{Kind: "item", Name: "Light", Rarity: "Rare", PriceCP: 100, BaseWeight: 5},
		{Kind: "item", Name: "Heavy", Rarity: "Rare", PriceCP: 240, BaseWeight: 5},
	}
	weights := rarityWeights{}
	out := runBudgetFillLoop(rng, candidates, weights, 200, 200)
	if len(out.Chosen) != 1 {
		t.Fatalf("expected one chosen item, got %d", len(out.Chosen))
	}
	if out.RemainingCP != 100 {
		t.Fatalf("expected remaining budget 100, got %d", out.RemainingCP)
	}
}

func TestBudgetFillWithoutReplacement(t *testing.T) {
	rng := rand.New(rand.NewSource(3))
	candidates := []lootCandidate{
		{Kind: "item", Name: "Single", Rarity: "Common", PriceCP: 100, BaseWeight: 1},
	}
	out := runBudgetFillLoop(rng, candidates, rarityWeights{}, 500, 500)
	if len(out.Chosen) != 1 {
		t.Fatalf("expected exactly one pick without replacement, got %d", len(out.Chosen))
	}
	if out.RemainingCP != 400 {
		t.Fatalf("expected remaining budget 400, got %d", out.RemainingCP)
	}
}

func TestBudgetFillMinimumSpendFloor(t *testing.T) {
	rng := rand.New(rand.NewSource(11))
	candidates := []lootCandidate{
		{Kind: "item", Name: "Cheap A", Rarity: "Common", PriceCP: 30, BaseWeight: 1},
		{Kind: "item", Name: "Cheap B", Rarity: "Common", PriceCP: 40, BaseWeight: 1},
		{Kind: "item", Name: "Meaningful", Rarity: "Common", PriceCP: 400, BaseWeight: 1},
	}
	out := runBudgetFillLoop(rng, candidates, rarityWeights{}, 500, 500)
	if len(out.Chosen) != 1 {
		t.Fatalf("expected only one meaningful pick due to spend floor, got %d", len(out.Chosen))
	}
	if out.Chosen[0].Name != "Meaningful" {
		t.Fatalf("expected meaningful pick first, got %s", out.Chosen[0].Name)
	}
}

func TestQualityQuantityInversion(t *testing.T) {
	weights := rarityWeights{Common: 100, Uncommon: 0, Rare: 0, VeryRare: 0}

	commonCandidates := []lootCandidate{
		{Kind: "item", Name: "CommonA", Rarity: "Common", PriceCP: 80, BaseWeight: 1},
		{Kind: "item", Name: "CommonB", Rarity: "Common", PriceCP: 90, BaseWeight: 1},
	}
	expensiveCandidates := []lootCandidate{
		{Kind: "item", Name: "RareBig", Rarity: "Common", PriceCP: 210, BaseWeight: 1},
	}

	rngA := rand.New(rand.NewSource(7))
	outCommon := runBudgetFillLoop(rngA, commonCandidates, weights, 220, 220)
	rngB := rand.New(rand.NewSource(7))
	outExpensive := runBudgetFillLoop(rngB, expensiveCandidates, weights, 220, 220)

	if len(outCommon.Chosen) <= len(outExpensive.Chosen) {
		t.Fatalf("expected common pool to yield more picks: common=%d expensive=%d", len(outCommon.Chosen), len(outExpensive.Chosen))
	}
}

func TestDimensionWeightFiltering(t *testing.T) {
	id := uuid.New()
	unmapped := uuid.New()
	weights := map[uuid.UUID]map[string]float64{
		id: {
			"dungeon": 1.7,
			"arcane":  1.1,
		},
	}
	if w, ok := dimensionWeight(weights, id, "dungeon"); !ok || w != 1.7 {
		t.Fatalf("expected mapped weight 1.7, got %f ok=%v", w, ok)
	}
	if _, ok := dimensionWeight(weights, id, "office"); ok {
		t.Fatalf("expected non-matching mapped value to fail")
	}
	if w, ok := dimensionWeight(weights, unmapped, "office"); !ok || w != 1.0 {
		t.Fatalf("expected unmapped id fallback to 1.0, got %f ok=%v", w, ok)
	}
}

type economySimulationStats struct {
	avgPicks     float64
	zeroRollRate float64
}

func buildBandCandidates(avgCostCP int64, size int) []lootCandidate {
	out := make([]lootCandidate, 0, size)
	spread := int64(math.Max(10, float64(avgCostCP)*0.25))
	for i := 0; i < size; i++ {
		cost := avgCostCP - spread + int64((2*int(spread))*i/size)
		if cost < 25 {
			cost = 25
		}
		out = append(out, lootCandidate{
			Kind:       "item",
			Name:       "Sim",
			Rarity:     "Common",
			PriceCP:    cost,
			BaseWeight: 1.0,
		})
	}
	return out
}

func simulateEconomyBand(profile models.LootLevelBudgetProfile, candidates []lootCandidate, trials int, seed int64) economySimulationStats {
	rng := rand.New(rand.NewSource(seed))
	totalPicks := 0
	zero := 0
	for i := 0; i < trials; i++ {
		session, expected := rollBudgetFromProfile(rng, profile)
		if session == 0 {
			zero++
		}
		out := runBudgetFillLoop(rng, candidates, rarityWeights{}, session, expected)
		totalPicks += len(out.Chosen)
	}
	return economySimulationStats{
		avgPicks:     float64(totalPicks) / float64(trials),
		zeroRollRate: float64(zero) / float64(trials),
	}
}

func TestEconomyBaselineSimulation(t *testing.T) {
	lowProfile := models.LootLevelBudgetProfile{LootLevel: 3, MeanGP: 12, SigmaGP: 6.5, MinGP: 0, MaxGP: 30}
	midProfile := models.LootLevelBudgetProfile{LootLevel: 9, MeanGP: 35, SigmaGP: 13, MinGP: 0, MaxGP: 80}
	highProfile := models.LootLevelBudgetProfile{LootLevel: 17, MeanGP: 110, SigmaGP: 30, MinGP: 0, MaxGP: 220}

	low := simulateEconomyBand(lowProfile, buildBandCandidates(500, 64), 1000, 101)
	mid := simulateEconomyBand(midProfile, buildBandCandidates(900, 64), 1000, 102)
	high := simulateEconomyBand(highProfile, buildBandCandidates(1500, 64), 1000, 103)

	t.Logf("economy baseline low(avg=%.2f zero=%.3f) mid(avg=%.2f zero=%.3f) high(avg=%.2f zero=%.3f)",
		low.avgPicks, low.zeroRollRate, mid.avgPicks, mid.zeroRollRate, high.avgPicks, high.zeroRollRate)
}

func TestEconomyGuardrailDropBands(t *testing.T) {
	lowProfile := models.LootLevelBudgetProfile{LootLevel: 3, MeanGP: 12, SigmaGP: 6.5, MinGP: 0, MaxGP: 30}
	midProfile := models.LootLevelBudgetProfile{LootLevel: 9, MeanGP: 35, SigmaGP: 13, MinGP: 0, MaxGP: 80}
	highProfile := models.LootLevelBudgetProfile{LootLevel: 17, MeanGP: 110, SigmaGP: 30, MinGP: 0, MaxGP: 220}

	low := simulateEconomyBand(lowProfile, buildBandCandidates(500, 64), 1200, 201)
	mid := simulateEconomyBand(midProfile, buildBandCandidates(900, 64), 1200, 202)
	high := simulateEconomyBand(highProfile, buildBandCandidates(1500, 64), 1200, 203)

	if low.avgPicks < 1.0 || low.avgPicks > 3.0 {
		t.Fatalf("low band avg picks out of target: %.2f", low.avgPicks)
	}
	if mid.avgPicks < 2.0 || mid.avgPicks > 5.0 {
		t.Fatalf("mid band avg picks out of target: %.2f", mid.avgPicks)
	}
	if high.avgPicks < 3.0 || high.avgPicks > 7.0 {
		t.Fatalf("high band avg picks out of target: %.2f", high.avgPicks)
	}
}

func TestEconomyGuardrailZeroRollBounded(t *testing.T) {
	profiles := []models.LootLevelBudgetProfile{
		{LootLevel: 3, MeanGP: 12, SigmaGP: 6.5, MinGP: 0, MaxGP: 30},
		{LootLevel: 9, MeanGP: 35, SigmaGP: 13, MinGP: 0, MaxGP: 80},
		{LootLevel: 17, MeanGP: 110, SigmaGP: 30, MinGP: 0, MaxGP: 220},
	}
	anyNonZero := false
	for i, p := range profiles {
		stats := simulateEconomyBand(p, buildBandCandidates(900, 64), 2000, int64(300+i))
		if stats.zeroRollRate > 0 {
			anyNonZero = true
		}
		if stats.zeroRollRate > 0.06 {
			t.Fatalf("zero-roll rate too high for level %d: %.3f", p.LootLevel, stats.zeroRollRate)
		}
	}
	if !anyNonZero {
		t.Fatalf("expected at least one level band to retain non-zero dead-roll probability")
	}
}
