package seed

import (
	"log"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/seed/batch"
	"github.com/rpupo63/unified-personal-site-backend/seed/uuids"
	"gorm.io/gorm"
)

// SpellSeed defines a spell for seeding
type SpellSeed struct {
	ID          uuid.UUID
	Name        string
	Description string
	SlotLevel   int
	Type        string
	Range       *string
	Duration    *string
	Concentration bool
	SaveAttr    *string
	DamageDice  *string
	DamageType  *models.DamageType
	AddModifier bool
	Components  []string // Names of components
}

// ExistingSpells returns a list of pre-defined spells
func ExistingSpells() []SpellSeed {
	return []SpellSeed{
		{
			ID:         uuids.ComponentUUID("Thaumaturgy"),
			Name:       "Thaumaturgy",
			SlotLevel:  0,
			Type:       "Utility",
			Components: []string{"Arcanum", "Self"},
		},
		{
			ID:         uuids.ComponentUUID("lightning goes boom"),
			Name:       "lightning goes boom",
			SlotLevel:  1,
			Type:       "Attack",
			Components: []string{"Fulgur", "Increase", "Beam"}, // Expand -> Increase
		},
		{
			ID:         uuids.ComponentUUID("blue man"),
			Name:       "blue man",
			SlotLevel:  1,
			Type:       "Utility",
			Components: []string{"Fulgur", "Bind", "Increase"}, // Fuse -> Bind, Expand -> Increase
		},
		{
			ID:        uuids.ComponentUUID("Calm Down"),
			Name:      "Calm Down",
			SlotLevel: 1,
			Type:      "Utility",
			Components: []string{"Arcanum", "Strong", "Decrease"}, // Psi -> Arcanum, Focus -> Strong, Cool -> Decrease
		},
		{
			ID:        uuids.ComponentUUID("Can't Sleep Yet"),
			Name:      "Can't Sleep Yet",
			SlotLevel: 1,
			Type:      "Utility",
			Components: []string{"Vita", "Create", "Target"}, // Revive -> Vita, Create; Touch -> Target
		},
		{
			ID:        uuids.ComponentUUID("Splash"),
			Name:      "Splash",
			SlotLevel: 1,
			Type:      "Attack",
			Components: []string{"Push", "Decrease", "Aqua"}, // Cool -> Decrease
		},
		{
			ID:        uuids.ComponentUUID("Flood"),
			Name:      "Flood",
			SlotLevel: 2,
			Type:      "Attack",
			Components: []string{"Chain", "Push", "Decrease", "Aqua"}, // Cool -> Decrease
		},
		{
			ID:        uuids.ComponentUUID("Massage Therapy 2.0"),
			Name:      "Massage Therapy 2.0",
			SlotLevel: 1,
			Type:      "Utility",
			Components: []string{"Vita", "Increase", "Target"}, // Mend -> Vita, Increase; Touch -> Target
		},
		{
			ID:        uuids.ComponentUUID("What big teeth you have!"),
			Name:      "What big teeth you have!",
			SlotLevel: 1,
			Type:      "Utility",
			Components: []string{"Arcanum", "Crush", "Self"}, // Psi -> Arcanum
		},
		{
			ID:        uuids.ComponentUUID("test"),
			Name:      "test",
			SlotLevel: 0,
			Type:      "Utility",
			Components: []string{"Fulgur", "Arcanum"},
		},
		{
			ID:        uuids.ComponentUUID("Bear Hug"),
			Name:      "Bear Hug",
			SlotLevel: 1,
			Type:      "Effect",
			Components: []string{"Arcanum", "Increase", "Target"},
		},
		{
			ID:        uuids.ComponentUUID("Steep"),
			Name:      "Steep",
			SlotLevel: 2,
			Type:      "Attack",
			Components: []string{"Strong", "Increase", "Aqua"},
		},
		{
			ID:        uuids.ComponentUUID("Hot Shot"),
			Name:      "Hot Shot",
			SlotLevel: 2,
			Type:      "Attack",
			Components: []string{"Ignis", "Zone"},
		},
		{
			ID:        uuids.ComponentUUID("Drip"),
			Name:      "Drip",
			SlotLevel: 1,
			Type:      "Attack",
			Components: []string{"Aqua"},
		},
		{
			ID:        uuids.ComponentUUID("Magnet Hug 2.0"),
			Name:      "Magnet Hug 2.0",
			SlotLevel: 2,
			Type:      "Save",
			Components: []string{"Bind", "Ferrum", "Pull", "Grab"},
		},
		{
			ID:        uuids.ComponentUUID("Tainted Love"),
			Name:      "Tainted Love",
			SlotLevel: 1,
			Type:      "Effect",
			Components: []string{"Arcanum", "Decrease"},
		},
		{
			ID:        uuids.ComponentUUID("Forced Entry"),
			Name:      "Forced Entry",
			SlotLevel: 1,
			Type:      "Utility",
			Components: []string{"Vita", "Push"},
		},
		{
			ID:        uuids.ComponentUUID("2112"),
			Name:      "2112",
			SlotLevel: 2,
			Type:      "Effect",
			Components: []string{"Anger", "Extreme", "Sonus", "Cone"},
		},
		{
			ID:        uuids.ComponentUUID("Can I Kick it?"),
			Name:      "Can I Kick it?",
			SlotLevel: 2,
			Type:      "Attack",
			Components: []string{"Fulgur", "Increase", "Zone", "Push"},
		},
	}
}

// SeedExistingSpells seeds the spells from ExistingSpells
func SeedExistingSpells(tx *gorm.DB) error {
	// Ensure system user exists
	systemUserID := uuids.SystemUserUUID("system")
	systemUser := models.User{
		ID:       systemUserID,
		Email:    "system@faradhaven.com",
		Name:  "system",
	}
	if err := tx.FirstOrCreate(&systemUser, "id = ?", systemUserID).Error; err != nil {
		return err
	}

	spellSeeds := ExistingSpells()
	spells := make([]models.Spell, 0, len(spellSeeds))
	spellComponents := make([]models.SpellComponent, 0)

	// Build a map of component names to their UUIDs for quick lookup
	var componentsList []models.Component
	if err := tx.Find(&componentsList).Error; err != nil {
		return err
	}
	componentMap := make(map[string]uuid.UUID)
	for _, c := range componentsList {
		componentMap[c.Name] = c.ID
	}

	for _, ss := range spellSeeds {
		spells = append(spells, models.Spell{
			ID:          ss.ID,
			UserID:      systemUserID, // Associate with system user
			Name:        ss.Name,
			Description: ss.Description,
			SlotLevel:   ss.SlotLevel,
			Type:        ss.Type,
			Range:       ss.Range,
			Duration:    ss.Duration,
			Concentration: ss.Concentration,
			SaveAttr:    ss.SaveAttr,
			DamageDice:  ss.DamageDice,
			DamageType:  ss.DamageType,
			AddModifier: ss.AddModifier,
		})

		for _, compName := range ss.Components {
			if compID, ok := componentMap[compName]; ok {
				spellComponents = append(spellComponents, models.SpellComponent{
					SpellID:     ss.ID,
					ComponentID: compID,
				})
			} else {
				log.Printf("Warning: component '%s' not found for spell '%s'", compName, ss.Name)
			}
		}
	}

	if err := batch.UpsertBatchUpdateAll(tx, spells, 100); err != nil {
		return err
	}

	// Clear existing spell-component associations for these spells before re-adding
	spellIDs := make([]uuid.UUID, len(spells))
	for i, s := range spells {
		spellIDs[i] = s.ID
	}
	if err := tx.Where("spell_id IN ?", spellIDs).Delete(&models.SpellComponent{}).Error; err != nil {
		return err
	}

	if err := batch.InsertBatch(tx, spellComponents, 100); err != nil {
		return err
	}

	log.Printf("Upserted %d existing spells and their components", len(spells))
	return nil
}
