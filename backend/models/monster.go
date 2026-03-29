package models

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Monster represents a creature or NPC in the game world.
type Monster struct {
	Base
	UserID              uuid.UUID       `json:"user_id"`
	Name                string          `json:"name"`
	Size                CreatureSize    `json:"size"`
	Type                CreatureType    `json:"type"`
	Alignment           string          `json:"alignment"`
	ChallengeRating     string          `json:"challenge_rating"` // e.g., "1/2", "5", "CR 10"
	ProficiencyBonus    int             `json:"proficiency_bonus"`
	ArmorClass          int             `json:"armor_class"`
	HitPoints           int             `json:"hit_points"`
	HitDice             string          `json:"hit_dice"` // e.g., "5d8 + 10"
	Speed               string          `json:"speed"`    // e.g., "30 ft, fly 60 ft"
	Strength            int             `json:"strength"`
	Dexterity           int             `json:"dexterity"`
	Constitution        int             `json:"constitution"`
	Intelligence        int             `json:"intelligence"`
	Wisdom              int             `json:"wisdom"`
	Charisma            int             `json:"charisma"`
	SavingThrows        pq.StringArray  `json:"saving_throws" gorm:"type:text[]"` // e.g., ["Str", "Dex"]
	Skills              pq.StringArray  `json:"skills" gorm:"type:text[]"`        // e.g., ["Perception +5", "Stealth +7"]
	DamageImmunities    pq.StringArray  `json:"damage_immunities" gorm:"type:text[]"`
	DamageResistances   pq.StringArray  `json:"damage_resistances" gorm:"type:text[]"`
	ConditionImmunities pq.StringArray  `json:"condition_immunities" gorm:"type:text[]"`
	Senses              string          `json:"senses"` // e.g., "darkvision 60 ft, passive Perception 12"
	Languages           string          `json:"languages"`
	Notes               string          `json:"notes"`     // User's original text description/LLM notes
	ImageURL            *string         `json:"image_url"` // S3 URL for the monster image
	Attacks             []MonsterAttack `json:"attacks" gorm:"foreignKey:MonsterID"`
	Actions             []MonsterAction `json:"actions" gorm:"many2many:monster_actions;"` // Custom actions not tied to attacks
	SpecialTraits       pq.StringArray  `json:"special_traits" gorm:"type:text[]"`         // Passive abilities, e.g., "Magic Resistance"
	LegendaryActions    pq.StringArray  `json:"legendary_actions" gorm:"type:text[]"`
	Reactions           pq.StringArray  `json:"reactions" gorm:"type:text[]"`
	LairActions         pq.StringArray  `json:"lair_actions" gorm:"type:text[]"`
	Environments        pq.StringArray  `json:"environments" gorm:"type:text[]"` // e.g., "Forest", "Desert", for filtering/search
	Source              string          `json:"source"`                          // e.g., "User Generated", "Monster Manual"

	// LLM Generated fields
	VisualDescription string `json:"visual_description"` // Optimized for image generation AI
}

// MonsterAttack represents a single attack a monster can make.
// This is intentionally simplified from CharacterAttack for monster purposes.
type MonsterAttack struct {
	Base
	MonsterID   uuid.UUID  `json:"monster_id" gorm:"type:uuid;not null"` // Foreign key to Monster
	Name        string     `json:"name"`
	Description string     `json:"description"` // Full description, e.g., "Melee Weapon Attack: +5 to hit, reach 5 ft., one target. Hit: 1d6 + 3 piercing damage."
	IsLegendary bool       `json:"is_legendary"`
	WeaponID    *uuid.UUID `json:"weapon_id"` // Optional: link to a specific weapon if applicable
	Weapon      *Weapon    `json:"weapon,omitempty"`
}

// MonsterAction represents a non-attack action a monster can take.
type MonsterAction struct {
	Base
	Name        string `json:"name"`
	Description string `json:"description"`
	IsLegendary bool   `json:"is_legendary"`
}

// NormalizeCreatureFields maps LLM or legacy casing to canonical [CreatureSize] / [CreatureType].
func (m *Monster) NormalizeCreatureFields() {
	if sz, ok := ParseCreatureSize(string(m.Size)); ok {
		m.Size = sz
	}
	if ty, ok := ParseCreatureType(string(m.Type)); ok {
		m.Type = ty
	}
}

// ParseChallengeRating converts the CR string (e.g., "1/2", "5") into a float64.
func (m *Monster) ParseChallengeRating() float64 {
	if strings.Contains(m.ChallengeRating, "/") {
		parts := strings.Split(m.ChallengeRating, "/")
		numerator, _ := strconv.ParseFloat(parts[0], 64)
		denominator, _ := strconv.ParseFloat(parts[1], 64)
		if denominator == 0 { // Avoid division by zero
			return 0
		}
		return numerator / denominator
	}
	cr, _ := strconv.ParseFloat(m.ChallengeRating, 64)
	return cr
}

// AbilityModifier calculates the standard D&D 5e ability modifier.
func (m *Monster) AbilityModifier(score int) int {
	return int(math.Floor(float64(score-10) / 2))
}

// calculateProficiencyBonus determines the proficiency bonus based on CR.
// This is a simplified version; official rules have more granular steps.
func calculateProficiencyBonus(cr float64) int {
	if cr < 1 { // CR 0, 1/8, 1/4, 1/2
		return 2
	} else if cr < 5 { // CR 1-4
		return 2
	} else if cr < 9 { // CR 5-8
		return 3
	} else if cr < 13 { // CR 9-12
		return 4
	} else if cr < 17 { // CR 13-16
		return 5
	} else { // CR 17+
		return 6
	}
}

// MarshalJSON customizes JSON marshalling to handle nested structs if needed.
// Example: If you want to customize how Attacks or Actions are represented.
func (m Monster) MarshalJSON() ([]byte, error) {
	type Alias Monster
	return json.Marshal(&struct {
		*Alias
		// Add any custom fields or transformations here
		ParsedCR float64 `json:"parsed_cr"`
	}{
		Alias:    (*Alias)(&m),
		ParsedCR: m.ParseChallengeRating(),
	})
}

// UnmarshalJSON customizes JSON unmarshalling.
func (m *Monster) UnmarshalJSON(data []byte) error {
	type Alias Monster
	aux := &struct {
		Attacks []json.RawMessage `json:"attacks"`
		Actions []json.RawMessage `json:"actions"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// For Attacks and Actions, assume they come as simple string descriptions or full objects.
	// We'll primarily rely on the `Attacks` and `Actions` fields for the LLM output directly.
	// If the LLM provides just strings, we map them to the Name field of Attack/Action struct.
	// If it provides full objects, then direct unmarshalling works.
	// This part might need adjustment based on final LLM output format expectations.
	for _, raw := range aux.Attacks {
		var att MonsterAttack
		if err := json.Unmarshal(raw, &att); err != nil {
			var attStr string
			if err := json.Unmarshal(raw, &attStr); err == nil {
				att.Name = attStr        // If it's just a string, use it as name
				att.Description = attStr // and description
			} else {
				return fmt.Errorf("failed to unmarshal monster attack: %w", err)
			}
		}
		m.Attacks = append(m.Attacks, att)
	}

	for _, raw := range aux.Actions {
		var act MonsterAction
		if err := json.Unmarshal(raw, &act); err != nil {
			var actStr string
			if err := json.Unmarshal(raw, &actStr); err == nil {
				act.Name = actStr        // If it's just a string, use it as name
				act.Description = actStr // and description
			} else {
				return fmt.Errorf("failed to unmarshal monster action: %w", err)
			}
		}
		m.Actions = append(m.Actions, act)
	}

	return nil
}
