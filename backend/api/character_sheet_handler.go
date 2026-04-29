package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/models"
	"github.com/rpupo63/faradhaven/backend/seed/faradhaven_classes"
	"github.com/rs/zerolog/log"
)

// getCharacterSheet returns the fully calculated character sheet
func (h *characterHandler) getCharacterSheet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "characterID")
		id, err := uuid.Parse(idStr)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid character ID")
			return
		}

		authUserID, err := ctxGetUserID(r.Context())
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		sheet, err := h.getCharacterSheetData(id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				respondError(w, http.StatusNotFound, "Character not found")
				return
			}
			log.Error().Err(err).Msg("Failed to get character sheet data")
			respondError(w, http.StatusInternalServerError, "Failed to get character sheet")
			return
		}

		if authUserID != sheet.Character.UserID.String() {
			respondError(w, http.StatusForbidden, "Forbidden")
			return
		}

		respondJSON(w, http.StatusOK, sheet)
	}
}

func (h *characterHandler) getCharacterSheetData(id uuid.UUID) (*CharacterSheetResponse, error) {
	// Single query: loads all relations needed for the sheet including equipment,
	// class levels with features, race traits, and components.
	character, err := h.characterRepo.FindByIDForSheet(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get character: %w", err)
	}

	if h.resourceService != nil {
		if ensureErr := h.resourceService.EnsureTrackableClassResources(character); ensureErr != nil {
			log.Warn().Err(ensureErr).Str("characterID", id.String()).Msg("EnsureTrackableClassResources")
		}
	}

	// If character is part of a party, load party members and identified beasts in ONE call
	if character.PartyID != nil {
		party, err := h.partyRepo.FindByIDWithMembersAndBeasts(*character.PartyID)
		if err != nil {
			log.Error().Err(err).Str("partyID", character.PartyID.String()).Msg("Failed to get party for character sheet")
		} else {
			character.Party = party
		}
	}

	// Derive classLevel from the preloaded class levels — no extra DB call needed
	class := &character.Class
	var classLevel *models.ClassLevel
	for i := range class.Levels {
		if class.Levels[i].Level == character.Level {
			classLevel = &class.Levels[i]
			break
		}
	}
	if classLevel == nil {
		return nil, fmt.Errorf("level data not found for class at level %d", character.Level)
	}

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

	// Spell pool: always from current class + race join tables (updates when seeds/DB change).
	availableComponents := MergeSpellPoolComponents(class, &character.Race)

	// Race traits are already preloaded via FindByIDForSheet (Race.Traits.Options)
	traits := character.Race.Traits

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

	// Inventory-only rows (foraging, scavenging, sanguine extraction, etc.); spell pool is available_components.
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
		ClassResources:           h.buildClassResources(character.ClassID, character.Level, character.ID, character),
	}

	if class.Name == "The Lorewright" {
		sheet.MadnessTable = faradhaven_classes.LorewrightMadnessTable()
	}

	if class.Name == "The Mutagen" {
		sheet.MadnessTable = faradhaven_classes.MutagenFeralTable()
	}

	// Build trait use states for limited-use race/lineage traits
	traitUseStates := make(map[string]int)
	traitMaxUses := make(map[string]int)

	profBonus := profBonusByLevel(character.Level)
	_ = initTraitUseResources(id, traits, profBonus, character.Level, h.characterResourceRepo)

	for _, t := range traits {
		if t.UsesPerRest == "" {
			continue
		}
		key := "trait_" + t.ID.String()
		if cr, lookupErr := h.characterResourceRepo.FindByCharacterAndKey(id, key); lookupErr == nil {
			traitUseStates[t.ID.String()] = cr.CurrentValue
			if cr.MaxValue != nil {
				traitMaxUses[t.ID.String()] = *cr.MaxValue
			}
		}
	}
	if len(traitUseStates) > 0 {
		sheet.TraitUseStates = traitUseStates
	}
	if len(traitMaxUses) > 0 {
		sheet.TraitMaxUses = traitMaxUses
	}

	return sheet, nil
}
