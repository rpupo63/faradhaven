package api

import (
	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/services"
)

// buildClassResources aggregates resource definitions, level values, and character state
// into a response-ready slice. Shared between characterHandler and levelHandler.
// character may be nil; when set, Vapor Blade shadow_points Value reflects Dex + proficiency.
func buildClassResources(
	classRepo database.ClassRepository,
	characterResourceRepo database.CharacterResourceRepository,
	classID uuid.UUID,
	level int,
	characterID uuid.UUID,
	character *models.Character,
) []ClassResourceResponse {
	defs, resourceMap, err := classRepo.FindResourceInfo(classID, level)
	if err != nil || len(defs) == 0 {
		return nil
	}

	charResources, _ := characterResourceRepo.FindByCharacterID(characterID)
	charResMap := make(map[string]*models.CharacterResource, len(charResources))
	for _, cr := range charResources {
		charResMap[cr.ResourceKey] = cr
	}

	result := make([]ClassResourceResponse, 0, len(defs))
	for _, def := range defs {
		val := resourceMap[def.ResourceKey]
		if def.ResourceKey == "shadow_points" && character != nil {
			val = services.VaporBladeShadowPointsMax(character.Dexterity, level)
		}
		resp := ClassResourceResponse{
			Key:          def.ResourceKey,
			DisplayName:  def.DisplayName,
			Category:     string(def.Category),
			Description:  def.Description,
			Value:        val,
			IsTrackable:  def.IsTrackable,
			DisplayOrder: def.DisplayOrder,
		}
		if def.IsTrackable {
			if cr, ok := charResMap[def.ResourceKey]; ok {
				resp.CurrentValue = &cr.CurrentValue
				resp.MaxValue = cr.MaxValue
			}
		}
		result = append(result, resp)
	}
	return result
}
