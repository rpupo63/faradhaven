package models

import "strings"

type LootTheme string
type LootLocation string
type LootSource string
type LootTier string
type LootRewardAmount string
type LootLevelBand string

const (
	LootThemeDungeon    LootTheme = "dungeon"
	LootThemeOffice     LootTheme = "office"
	LootThemeRich       LootTheme = "rich"
	LootThemePoor       LootTheme = "poor"
	LootThemeGangster   LootTheme = "gangster"
	LootThemeArcane     LootTheme = "arcane"
	LootThemeWilderness LootTheme = "wilderness"
)

const (
	LootLocationIndoor      LootLocation = "indoor"
	LootLocationUnderground LootLocation = "underground"
	LootLocationUrban       LootLocation = "urban"
	LootLocationSlums       LootLocation = "slums"
	LootLocationEstate      LootLocation = "estate"
	LootLocationStreet      LootLocation = "street"
	LootLocationWilds       LootLocation = "wilds"
)

const (
	LootSourceCommonEnemy LootSource = "common_enemy"
	LootSourceBossEnemy   LootSource = "boss_enemy"
	LootSourceRoom        LootSource = "room"
)

const (
	LootTierLow    LootTier = "low"
	LootTierMedium LootTier = "medium"
	LootTierHigh   LootTier = "high"
)

const (
	LootRewardScarce    LootRewardAmount = "scarce"
	LootRewardStandard  LootRewardAmount = "standard"
	LootRewardBountiful LootRewardAmount = "bountiful"
	LootRewardJackpot   LootRewardAmount = "jackpot"
)

const (
	LootLevelNovice     LootLevelBand = "novice"
	LootLevelAdventurer LootLevelBand = "adventurer"
	LootLevelVeteran    LootLevelBand = "veteran"
	LootLevelLegend     LootLevelBand = "legend"
)

func AllLootThemes() []LootTheme {
	return []LootTheme{
		LootThemeDungeon, LootThemeOffice, LootThemeRich, LootThemePoor,
		LootThemeGangster, LootThemeArcane, LootThemeWilderness,
	}
}

func AllLootLocations() []LootLocation {
	return []LootLocation{
		LootLocationIndoor, LootLocationUnderground, LootLocationUrban,
		LootLocationSlums, LootLocationEstate, LootLocationStreet, LootLocationWilds,
	}
}

func AllLootSources() []LootSource {
	return []LootSource{LootSourceCommonEnemy, LootSourceBossEnemy, LootSourceRoom}
}

func AllLootTiers() []LootTier {
	return []LootTier{LootTierLow, LootTierMedium, LootTierHigh}
}

func AllLootRewardAmounts() []LootRewardAmount {
	return []LootRewardAmount{LootRewardScarce, LootRewardStandard, LootRewardBountiful, LootRewardJackpot}
}

func AllLootLevelBands() []LootLevelBand {
	return []LootLevelBand{LootLevelNovice, LootLevelAdventurer, LootLevelVeteran, LootLevelLegend}
}

func ParseLootTheme(s string) (LootTheme, bool) {
	v := LootTheme(strings.ToLower(strings.TrimSpace(s)))
	for _, candidate := range AllLootThemes() {
		if v == candidate {
			return v, true
		}
	}
	return "", false
}

func ParseLootLocation(s string) (LootLocation, bool) {
	v := LootLocation(strings.ToLower(strings.TrimSpace(s)))
	for _, candidate := range AllLootLocations() {
		if v == candidate {
			return v, true
		}
	}
	return "", false
}

func ParseLootSource(s string) (LootSource, bool) {
	v := LootSource(strings.ToLower(strings.TrimSpace(s)))
	for _, candidate := range AllLootSources() {
		if v == candidate {
			return v, true
		}
	}
	return "", false
}

func ParseLootTier(s string) (LootTier, bool) {
	v := LootTier(strings.ToLower(strings.TrimSpace(s)))
	for _, candidate := range AllLootTiers() {
		if v == candidate {
			return v, true
		}
	}
	return "", false
}

func ParseLootRewardAmount(s string) (LootRewardAmount, bool) {
	v := LootRewardAmount(strings.ToLower(strings.TrimSpace(s)))
	for _, candidate := range AllLootRewardAmounts() {
		if v == candidate {
			return v, true
		}
	}
	return "", false
}

func ParseLootLevelBand(s string) (LootLevelBand, bool) {
	v := LootLevelBand(strings.ToLower(strings.TrimSpace(s)))
	for _, candidate := range AllLootLevelBands() {
		if v == candidate {
			return v, true
		}
	}
	return "", false
}
