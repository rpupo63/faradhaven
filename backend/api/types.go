package api

import (
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
)

// routeHandlers contains all the handlers for the API routes
type routeHandlers struct {
	authHandler      *authHandler
	userHandler      *userHandler
	characterHandler *characterHandler
	spellHandler     *spellHandler
	beastHandler     *beastHandler
	levelHandler     *levelHandler
}

// ErrorResponse represents an error response from the API
type ErrorResponse struct {
	Error   string `json:"error" example:"Internal Server Error"`
	Status  string `json:"status" example:"error"`
	Field   string `json:"field,omitempty" example:"title"`
	Details string `json:"details,omitempty" example:"Additional error details"`
	Cause   string `json:"cause,omitempty" example:"Underlying error cause"`
}

// User Request/Response Types

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

// Character Request/Response Types

type CreateCharacterRequest struct {
	UserID             uuid.UUID `json:"user_id"`
	Name               string    `json:"name"`
	RaceID             uuid.UUID `json:"race_id"`
	ClassID            uuid.UUID `json:"class_id"`
	Level              int       `json:"level"`
	Spellbook          []string  `json:"spellbook"`
	Strength           int       `json:"strength"`
	Dexterity          int       `json:"dexterity"`
	Constitution       int       `json:"constitution"`
	Intelligence       int       `json:"intelligence"`
	Wisdom             int       `json:"wisdom"`
	Charisma           int       `json:"charisma"`
	CurrentSpellPoints int       `json:"current_spell_points"`
	SkillProficiencies []string  `json:"skill_proficiencies"` // D&D 5e skill ids (e.g. "persuasion", "stealth")
}

type UpdateCharacterRequest struct {
	Name               *string    `json:"name,omitempty"`
	RaceID             *uuid.UUID `json:"race_id,omitempty"`
	ClassID            *uuid.UUID `json:"class_id,omitempty"`
	Level              *int       `json:"level,omitempty"`
	Spellbook          []string   `json:"spellbook,omitempty"`
	Strength           *int       `json:"strength,omitempty"`
	Dexterity          *int       `json:"dexterity,omitempty"`
	Constitution       *int       `json:"constitution,omitempty"`
	Intelligence       *int       `json:"intelligence,omitempty"`
	Wisdom             *int       `json:"wisdom,omitempty"`
	Charisma           *int       `json:"charisma,omitempty"`
	CurrentSpellPoints *int       `json:"current_spell_points,omitempty"`
	SkillProficiencies []string   `json:"skill_proficiencies,omitempty"` // D&D 5e skill ids
}

// ClassWithLevelsResponse is a class with all levels 1-20 for the compendium/book view
type ClassWithLevelsResponse struct {
	*models.Class
	Levels []models.ClassLevel `json:"levels"`
}

// CharacterSheetResponse is the fully calculated character sheet (Class + ClassLevel joined)
type CharacterSheetResponse struct {
	Character                *models.Character  `json:"character"`
	Class                    *models.Class      `json:"class"`
	ClassLevel               *models.ClassLevel `json:"class_level"`
	TotalHP                  int                `json:"total_hp"`
	AC                       int                `json:"ac"`      // 8 + ProficiencyBonus + DexMod
	SaveDC                   int                `json:"save_dc"` // 8 + ProficiencyBonus + PrimaryAbilityMod
	MaxSpellPoints           int                `json:"max_spell_points"`
	CurrentSpellPoints       int                `json:"current_spell_points"`
	SavingThrowProficiencies []string           `json:"saving_throw_proficiencies"` // ability ids from class.SavingThrows
}

// Spell Request/Response Types

type CreateSpellRequest struct {
	UserID       uuid.UUID   `json:"user_id"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	ComponentIDs []uuid.UUID `json:"component_ids"`
	SlotLevel    int         `json:"slot_level"`
}

type UpdateSpellRequest struct {
	Name         *string     `json:"name,omitempty"`
	Description  *string     `json:"description,omitempty"`
	ComponentIDs []uuid.UUID `json:"component_ids,omitempty"`
	SlotLevel    *int        `json:"slot_level,omitempty"`
}

// Beast Request/Response Types

type CreateBeastRequest struct {
	UserID           uuid.UUID             `json:"user_id"`
	Name             string                `json:"name"`
	ImageURL         *string               `json:"image_url,omitempty"`
	Size             models.CreatureSize   `json:"size"`
	Type             models.CreatureType   `json:"type"`
	Alignment        string                `json:"alignment"`
	ArmorClass       int                   `json:"armor_class"`
	HitPoints        int                   `json:"hit_points"`
	HitDice          string                `json:"hit_dice"`
	Speed            string                `json:"speed"`
	Strength         int                   `json:"strength"`
	Dexterity        int                   `json:"dexterity"`
	Constitution     int                   `json:"constitution"`
	Intelligence     int                   `json:"intelligence"`
	Wisdom           int                   `json:"wisdom"`
	Charisma         int                   `json:"charisma"`
	ChallengeRating  string                `json:"challenge_rating"`
	Abilities        *string               `json:"abilities,omitempty"`
	Actions          *string               `json:"actions,omitempty"`
	LegendaryActions *string               `json:"legendary_actions,omitempty"`
	Description      string                `json:"description"`
	Attacks          []CreateAttackRequest `json:"attacks,omitempty"`
}

type UpdateBeastRequest struct {
	Name             *string              `json:"name,omitempty"`
	ImageURL         *string              `json:"image_url,omitempty"`
	Size             *models.CreatureSize `json:"size,omitempty"`
	Type             *models.CreatureType `json:"type,omitempty"`
	Alignment        *string              `json:"alignment,omitempty"`
	ArmorClass       *int                 `json:"armor_class,omitempty"`
	HitPoints        *int                 `json:"hit_points,omitempty"`
	HitDice          *string              `json:"hit_dice,omitempty"`
	Speed            *string              `json:"speed,omitempty"`
	Strength         *int                 `json:"strength,omitempty"`
	Dexterity        *int                 `json:"dexterity,omitempty"`
	Constitution     *int                 `json:"constitution,omitempty"`
	Intelligence     *int                 `json:"intelligence,omitempty"`
	Wisdom           *int                 `json:"wisdom,omitempty"`
	Charisma         *int                 `json:"charisma,omitempty"`
	ChallengeRating  *string              `json:"challenge_rating,omitempty"`
	Abilities        *string              `json:"abilities,omitempty"`
	Actions          *string              `json:"actions,omitempty"`
	LegendaryActions *string              `json:"legendary_actions,omitempty"`
	Description      *string              `json:"description,omitempty"`
}

// Attack Request/Response Types

type CreateAttackRequest struct {
	Name        string            `json:"name"`
	AttackBonus int               `json:"attack_bonus"`
	DamageType  models.DamageType `json:"damage_type"`
	DamageDice  string            `json:"damage_dice"`
	Reach       *string           `json:"reach,omitempty"`
	Description *string           `json:"description,omitempty"`
}

// API Response Types

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
