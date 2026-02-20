package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/seed/faradhaven_classes"
	"github.com/rs/zerolog/log"
)

// getCharacterSheet returns the fully calculated character sheet (Class + ClassLevel joined)
func (h *characterHandler) getCharacterSheet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "characterID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		character, err := h.characterRepo.FindByID(id)
		if err != nil {
			respondError(w, http.StatusNotFound, "Character not found")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		if authUserID != character.UserID.String() {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		sheet, err := h.getCharacterSheetData(id)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get character sheet data")
			respondError(w, http.StatusInternalServerError, "Failed to get character sheet")
			return
		}

		respondJSON(w, http.StatusOK, sheet)
	}
}

func (h *characterHandler) getCharacterSheetData(id uuid.UUID) (*CharacterSheetResponse, error) {
	character, err := h.characterRepo.FindByIDWithSkills(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %w", err)
	}

	// Preload weapons with modifiers and items for the sheet
	if err := h.characterRepo.GetDB().
		Preload("CharacterWeapons.Weapon.Damages").
		Preload("CharacterWeapons.Modifiers").
		Preload("Items").
		Preload("EquippedArmor").
		Preload("EquippedShield").
		First(character, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("failed to preload equipment: %w", err)
	}

	// Fetch class with components for spell crafting
	class, err := h.classRepo.FindByIDWithLevels(character.ClassID)
	if err != nil {
		return nil, fmt.Errorf("class not found for character: %w", err)
	}

	// Fetch race with components for spell crafting
	race, err := h.raceRepo.FindByID(character.RaceID)
	if err != nil {
		return nil, fmt.Errorf("race not found for character: %w", err)
	}

	classLevel, err := h.classRepo.FindLevelByClassAndLevel(character.ClassID, character.Level)
	if err != nil {
		return nil, fmt.Errorf("level data not found for class: %w", err)
	}

	// Assign fully-loaded class (with Levels + LevelFeatures) so ComputeStats
	// can resolve Unarmored Defense and other level-gated features.
	character.Class = *class
	character.CurrentClassLevel = classLevel
	character.ComputeStats()
	stats := character.ComputedStats

	maxSP := classLevel.MaxSpellPoints
	currentSP := character.CurrentSpellPoints
	if currentSP > maxSP {
		currentSP = maxSP
	}

	// Map class.SavingThrows (ability names) to ability ids for frontend
	saveProfs := make([]string, 0, len(class.SavingThrows))
	for _, name := range class.SavingThrows {
		saveProfs = append(saveProfs, strings.ToLower(strings.TrimSpace(name)))
	}

	// Combine class and race components for spell crafting (deduplicated by ID)
	componentMap := make(map[uuid.UUID]models.Component)
	for _, comp := range class.Components {
		componentMap[comp.ID] = comp
	}
	for _, comp := range race.Components {
		componentMap[comp.ID] = comp
	}
	availableComponents := make([]models.Component, 0, len(componentMap))
	for _, comp := range componentMap {
		availableComponents = append(availableComponents, comp)
	}

	// Race Traits
	raceTraits, err := h.raceRepo.FindByIDWithTraits(character.RaceID)
	var traits []models.Trait
	if err == nil {
		traits = raceTraits.Traits
	}

	// Lineage Traits
	var lineage *models.Lineage
	if character.LineageID != nil {
		lineage, err = h.raceRepo.FindLineageByIDWithTraits(*character.LineageID)
		if err == nil && lineage != nil {
			traits = append(traits, lineage.LineageTraits...)
		}
	}

	// Map CharacterWeapons to response format
	inventoryWeapons := make([]CharacterWeaponResponse, 0, len(character.CharacterWeapons))
	for _, cw := range character.CharacterWeapons {
		activeModifiers := make([]WeaponModifierResponse, 0, len(cw.Modifiers))
		for _, m := range cw.Modifiers {
			bonusDamage := []BonusDamageInfo{}

			// Compute bonus damage based on modifier type and character stats
			if m.ModifierType == models.ModifierTypePistonCore && m.IsActive {
				intMod := abilityMod(character.Intelligence)
				if intMod > 0 {
					bonusDamage = append(bonusDamage, BonusDamageInfo{
						Dice:       fmt.Sprintf("%+d", intMod),
						DamageType: "Fixed", // Adds to base damage
					})
				}
			} else if m.ModifierType == models.ModifierTypeVenomCoating && m.IsActive {
				dice := "1d4"
				if character.Level >= 6 {
					dice = "2d4"
				}
				bonusDamage = append(bonusDamage, BonusDamageInfo{
					Dice:       dice,
					DamageType: "Poison",
				})
			}

			activeModifiers = append(activeModifiers, WeaponModifierResponse{
				ModifierType: string(m.ModifierType),
				IsActive:     m.IsActive,
				BonusDamage:  bonusDamage,
				Metadata:     m.Metadata,
			})
		}

		inventoryWeapons = append(inventoryWeapons, CharacterWeaponResponse{
			CharacterWeaponID: cw.ID.String(),
			Weapon:            cw.Weapon,
			IsPrimary:         cw.IsPrimary,
			IsEquipped:        cw.IsEquipped,
			CustomName:        cw.CustomName,
			ActiveModifiers:   activeModifiers,
		})
	}

	componentsList := character.Components

	var harvestedAbilities models.HarvestedAbilities
	if len(character.HarvestedAbilities) > 0 {
		json.Unmarshal(character.HarvestedAbilities, &harvestedAbilities)
	}

	sheet := &CharacterSheetResponse{
		Character:                character,
		Class:                    class,
		ClassLevel:               classLevel,
		MaxHP:                    stats.MaxHP,
		CurrentHP:                stats.CurrentHP,
		TempHP:                   stats.TempHP,
		AC:                       stats.ArmorClass,
		SaveDC:                   stats.SpellSaveDC,
		MaxSpellPoints:           maxSP,
		CurrentSpellPoints:       currentSP,
		SavingThrowProficiencies: saveProfs,
		AvailableComponents:      availableComponents,
		HitDiceTotal:             stats.HitDiceTotal,
		HitDiceRemaining:         stats.HitDiceRemaining,
		HitDie:                   stats.HitDie,
		Money:                    character.Money,
		MeleeAttackBonus:         stats.MeleeAttackBonus,
		RangedAttackBonus:        stats.RangedAttackBonus,
		RaceTraits:               traits,
		Lineage:                  lineage,
		InventoryWeapons:         inventoryWeapons,
		InventoryItems:           character.Items,
		Components:               componentsList,
		HarvestedAbilities:       harvestedAbilities,
		ClassResources:           h.buildClassResources(character.ClassID, character.Level, character.ID),
	}

	if class.Name == "The Lorewright" {
		sheet.MadnessTable = faradhaven_classes.LorewrightMadnessTable()
	}

	if class.Name == "The Mutagen" {
		sheet.MadnessTable = faradhaven_classes.MutagenFeralTable()
	}

	return sheet, nil
}
