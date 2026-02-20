package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// MockLLMClient is a mock implementation of the LLMClient interface.
// It simulates structured content generation by returning a hardcoded JSON
// that is vaguely consistent with the requested schema and prompt.
type MockLLMClient struct {
	// Potentially store loaded schemas for more sophisticated mocking
}

// NewMockLLMClient creates a new instance of MockLLMClient.
func NewMockLLMClient() *MockLLMClient {
	return &MockLLMClient{}
}

// GenerateStructuredContent simulates an LLM generating structured JSON content.
// In a real scenario, this would involve calling a generative AI model.
// For this mock, it returns a predefined JSON structure based on the prompt's keywords.
func (m *MockLLMClient) GenerateStructuredContent(ctx context.Context, prompt string, schema string) (string, error) {
	// In a real LLM integration, 'schema' would be used by the LLM to format its output.
	// For this mock, we'll ignore the schema content for simplicity and return a fixed example
	// that we know adheres to our Monster schema.

	// Example output for a "fire goblin" with CR "1/2"
	// This should be dynamic based on the prompt for a better mock,
	// but for now, we'll keep it static and assume a successful LLM call.

	// Attempt to load the schema from file
	schemaPath := "docs/schemas/monster_schema.json"
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return "", fmt.Errorf("mock LLM client: failed to read monster schema from %s: %w", schemaPath, err)
	}

	var parsedSchema map[string]interface{}
	if err := json.Unmarshal(schemaBytes, &parsedSchema); err != nil {
		return "", fmt.Errorf("mock LLM client: failed to parse monster schema: %w", err)
	}

	// Dynamic values based on prompt or a simple keyword check for better mocking
	name := "Generated Monster"
	cr := "1"
	visualDescription := "A generic fantasy creature."

	if strings.Contains(strings.ToLower(prompt), "goblin") && strings.Contains(strings.ToLower(prompt), "fire") {
		name = "Fiery Goblin"
		visualDescription = "A small, scrawny goblin with unusually red skin, glowing yellow eyes, and small flames flickering around its hands. It carries a crude, charred scimitar."
	} else if strings.Contains(strings.ToLower(prompt), "dragon") && strings.Contains(strings.ToLower(prompt), "ice") {
		name = "Ice Dragon Wyrmling"
		visualDescription = "A small, reptilian creature with iridescent blue scales and sharp claws, exhaling a puff of frosty air. Its eyes gleam with ancient, cold wisdom."
	}
	if strings.Contains(strings.ToLower(prompt), "cr 1/2") {
		cr = "1/2"
	} else if strings.Contains(strings.ToLower(prompt), "cr 5") {
		cr = "5"
	}

	// Use json.Marshal to properly escape dynamic strings for safe JSON embedding
	escapedName, _ := json.Marshal(name)
	escapedCR, _ := json.Marshal(cr)
	escapedNotes, _ := json.Marshal(prompt)
	escapedVisual, _ := json.Marshal(visualDescription)

	mockMonsterJSON := fmt.Sprintf(`{
		"name": %s,
		"size": "Medium",
		"type": "Aberration",
		"alignment": "Chaotic Neutral",
		"challenge_rating": %s,
		"proficiency_bonus": 2,
		"armor_class": 13,
		"hit_points": 22,
		"hit_dice": "4d8 + 4",
		"speed": "30 ft",
		"strength": 12,
		"dexterity": 10,
		"constitution": 14,
		"intelligence": 8,
		"wisdom": 10,
		"charisma": 6,
		"saving_throws": ["Con"],
		"skills": ["Perception +2"],
		"damage_immunities": [],
		"damage_resistances": [],
		"condition_immunities": [],
		"senses": "darkvision 60 ft, passive Perception 12",
		"languages": "Deep Speech",
		"notes": %s,
		"image_url": null,
		"attacks": [
			{
				"name": "Tentacle Lash",
				"description": "Melee Weapon Attack: +3 to hit, reach 5 ft., one target. Hit: 1d6 + 1 bludgeoning damage."
			}
		],
		"actions": [],
		"special_traits": ["Amorphous: The monster can move through a space as narrow as 1 inch wide without squeezing."],
		"legendary_actions": [],
		"reactions": [],
		"lair_actions": [],
		"environments": ["Underdark"],
		"source": "Generated",
		"visual_description": %s
	}`, escapedName, escapedCR, escapedNotes, escapedVisual)

	return mockMonsterJSON, nil
}
