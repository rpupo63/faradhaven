package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rpupo63/faradhaven/backend/models"
)

// SpellOpinion represents the LLM's opinion on a crafted spell.
type SpellOpinion struct {
	DescriptionOpinion string `json:"description_opinion"`
	DamageOpinion      string `json:"damage_opinion"`
	EffectOpinion      string `json:"effect_opinion"`
	OverallVerdict     string `json:"overall_verdict"`

	// Recommended Edits
	RecommendedName            *string `json:"recommended_name,omitempty"`
	RecommendedDescription     *string `json:"recommended_description,omitempty"`
	RecommendedType            *string `json:"recommended_type,omitempty"`
	RecommendedRange           *string `json:"recommended_range,omitempty"`
	RecommendedDuration        *string `json:"recommended_duration,omitempty"`
	RecommendedDamageDiceCount *int    `json:"recommended_damage_dice_count,omitempty"`
	RecommendedDamageDieSize   *int    `json:"recommended_damage_die_size,omitempty"`
	RecommendedDamageType      *string `json:"recommended_damage_type,omitempty"`
	RecommendedSaveAttr        *string `json:"recommended_save_attr,omitempty"`
}

// SpellAIService provides AI-powered opinions on spells.
type SpellAIService struct {
	llmClient LLMClient
}

// NewSpellAIService creates a new SpellAIService.
func NewSpellAIService(llmClient LLMClient) *SpellAIService {
	return &SpellAIService{llmClient: llmClient}
}

// GetSpellOpinion generates an opinion on a spell based on its components and characteristics.
func (s *SpellAIService) GetSpellOpinion(ctx context.Context, spell *models.Spell, components []models.Component) (*SpellOpinion, string, error) {
	if s.llmClient == nil {
		return nil, "", fmt.Errorf("LLM client not initialized")
	}

	prompt := BuildSpellAIOpinionPrompt(spell, components)

	response, err := s.llmClient.GenerateStructuredContent(ctx, prompt, SpellAIOpinionResponseSchema)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate AI opinion: %w", err)
	}

	var opinion SpellOpinion
	if err := json.Unmarshal([]byte(response), &opinion); err != nil {
		return nil, response, fmt.Errorf("failed to parse AI opinion: %w", err)
	}

	return &opinion, response, nil
}
