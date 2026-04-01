package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/models"
	"github.com/rpupo63/unified-personal-site-backend/seed/faradhaven_classes"
	"github.com/rs/zerolog/log"
	// Vertex AI for image generation (assuming it's the "Google image service")
	// If it's another service like Imagen or a custom solution, this import might change.
	// Placeholder: vertexaipb "cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	// For now, I'll assume a simpler HTTP call to a generic image generation service or mock it.
	// Assuming a client for LLM interaction will be available or created.
	// For demonstration, I'll use a placeholder for LLM interaction.
)

const (
	llmMonsterGenerationModel = "gemini-1.5-flash-preview-0514"     // Example LLM model
	imageGenerationEndpoint   = "https://image.google.com/generate" // Placeholder
)

// MonsterGenerationService handles the orchestration of monster generation.
type MonsterGenerationService struct {
	monsterRepo database.MonsterRepository
	s3Service   *S3Service
	llmClient   LLMClient // New: LLM Client
	schemaPath  string
}

// LLMClient interface for structured content generation.
type LLMClient interface {
	GenerateStructuredContent(ctx context.Context, prompt string, schema string) (string, error)
}

// NewMonsterGenerationService creates a new MonsterGenerationService.
func NewMonsterGenerationService(
	monsterRepo database.MonsterRepository,
	s3Service *S3Service,
	llmClient LLMClient, // New: LLM Client
) *MonsterGenerationService {
	return &MonsterGenerationService{
		monsterRepo: monsterRepo,
		s3Service:   s3Service,
		llmClient:   llmClient,
		schemaPath:  "docs/schemas/monster_schema.json",
	}
}

// GenerateMonsterFromPromptRequest encapsulates the input for monster generation.
type GenerateMonsterFromPromptRequest struct {
	UserID          uuid.UUID `json:"user_id"`
	Description     string    `json:"description"`
	ChallengeRating string    `json:"challenge_rating"` // e.g., "1", "1/2", "CR 5"
	// FaradhavenClassName, when non-empty, selects a seeded Faradhaven class from
	// seed/faradhaven_classes. The LLM receives that class's full seed JSON as context
	// and builds a monster stat block themed as an enemy embodying that class (not a PC sheet).
	FaradhavenClassName string `json:"faradhaven_class_name,omitempty"`
}

// GenerateMonsterFromPrompt orchestrates the monster generation process.
func (s *MonsterGenerationService) GenerateMonsterFromPrompt(ctx context.Context, req GenerateMonsterFromPromptRequest) (*models.Monster, error) {
	className := strings.TrimSpace(req.FaradhavenClassName)

	var llmOutput string
	var err error
	if className != "" {
		llmOutput, err = s.generateMonsterFromFaradhavenClass(ctx, className, req.Description, req.ChallengeRating)
	} else {
		llmOutput, err = s.generateMonsterDataFromLLM(ctx, req.Description, req.ChallengeRating)
	}
	if err != nil {
		return nil, fmt.Errorf("LLM text generation failed: %w", err)
	}

	// Unmarshal LLM output into a Monster struct
	var monster models.Monster
	if err := json.Unmarshal([]byte(llmOutput), &monster); err != nil {
		return nil, fmt.Errorf("failed to unmarshal LLM output into Monster struct: %w", err)
	}

	monster.NormalizeCreatureFields()

	// Overwrite some fields from the request/system
	monster.UserID = req.UserID
	monster.Notes = req.Description // Store original prompt in Notes
	if className != "" {
		monster.Source = fmt.Sprintf("User Generated (Faradhaven class: %s)", className)
	} else {
		monster.Source = "User Generated"
	}
	monster.ChallengeRating = req.ChallengeRating // Ensure CR from request is used

	// Step B: Image Generation
	imageURL, err := s.generateImage(ctx, monster.VisualDescription)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Image generation failed, proceeding without image.")
		monster.ImageURL = nil // Set to nil if image generation fails
	} else {
		monster.ImageURL = &imageURL
	}

	// Step D: Assembly and Save to DB
	if err := s.monsterRepo.Add(&monster); err != nil {
		return nil, fmt.Errorf("failed to save generated monster to database: %w", err)
	}

	return &monster, nil
}

// generateMonsterDataFromLLM uses the LLMClient to get structured monster data.
func (s *MonsterGenerationService) generateMonsterDataFromLLM(ctx context.Context, prompt, cr string) (string, error) {
	if s.llmClient == nil {
		return "", fmt.Errorf("LLM client not initialized")
	}

	// 1. Read the monster_schema.json file.
	schemaBytes, err := os.ReadFile(s.schemaPath)
	if err != nil {
		return "", fmt.Errorf("failed to read monster schema from %s: %w", s.schemaPath, err)
	}
	schemaString := string(schemaBytes)

	// 2. Construct a prompt.
	// The prompt instructs the LLM to generate a monster based on the description and CR,
	// and to ensure the output strictly adheres to the provided JSON schema.
	fullPrompt := fmt.Sprintf(`
	Generate a D&D 5e-style monster based on the following description and Challenge Rating (CR).
	The output MUST be a JSON object that strictly adheres to the provided JSON schema.
	Populate all fields relevant to the monster, including attributes, combat stats, attacks, actions, and a detailed visual_description suitable for an image generation AI.
	For attacks and actions, ensure names and descriptions are detailed.
	Ensure the 'challenge_rating' field in the output matches the requested CR.
	If any field is a list, provide an empty list if there are no items.

	User Description: "%s"
	Target Challenge Rating: "%s"

	JSON Schema:
	%s
	`, prompt, cr, schemaString)

	// 3. Call the LLM with the constructed prompt and the schema.
	// The `GenerateStructuredContent` method is expected to handle the interaction with the LLM
	// and return a JSON string that conforms to the schema.
	llmResponse, err := s.llmClient.GenerateStructuredContent(ctx, fullPrompt, schemaString)
	if err != nil {
		return "", fmt.Errorf("LLM structured content generation failed: %w", err)
	}

	return llmResponse, nil
}

func findFaradhavenClassSeedByName(name string) (*faradhaven_classes.FaradhavenClassSeed, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return nil, fmt.Errorf("empty Faradhaven class name")
	}
	all := faradhaven_classes.AllClasses()
	if len(all) == 0 {
		return nil, fmt.Errorf("no Faradhaven classes configured")
	}
	for i := range all {
		if strings.EqualFold(all[i].Name, trimmed) {
			seed := all[i]
			return &seed, nil
		}
	}
	return nil, fmt.Errorf("unknown Faradhaven class %q (example: %q)", trimmed, all[0].Name)
}

// generateMonsterFromFaradhavenClass builds the LLM prompt from faradhaven_classes seed data
// instead of a generic creature brief.
func (s *MonsterGenerationService) generateMonsterFromFaradhavenClass(ctx context.Context, className, flavor, cr string) (string, error) {
	if s.llmClient == nil {
		return "", fmt.Errorf("LLM client not initialized")
	}

	seed, err := findFaradhavenClassSeedByName(className)
	if err != nil {
		return "", err
	}

	schemaBytes, err := os.ReadFile(s.schemaPath)
	if err != nil {
		return "", fmt.Errorf("failed to read monster schema from %s: %w", s.schemaPath, err)
	}
	schemaString := string(schemaBytes)

	classJSON, err := json.MarshalIndent(seed, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal class seed: %w", err)
	}

	fullPrompt := fmt.Sprintf(`
Generate a D&D 5e-style monster stat block for a hostile NPC that embodies the Faradhaven player class below.
The JSON is authoritative game design from our seed data (features, components, resources, archetypes). Translate its themes and signature mechanics into monster traits, actions, and attacks appropriate for an enemy—not a player character sheet.
Scale difficulty to the requested Challenge Rating. Respect the class's fantasy (concept, proficiencies, component pool, resource definitions) when naming abilities.
The output MUST be a JSON object that strictly adheres to the provided JSON schema.
Populate all relevant fields including a detailed visual_description suited to image generation.
If any field is a list, provide an empty list when there are no items.

Faradhaven class seed (JSON):
%s

Additional scene or flavor direction from the user: "%s"
Target Challenge Rating: "%s"

JSON Schema:
%s
`, string(classJSON), flavor, cr, schemaString)

	return s.llmClient.GenerateStructuredContent(ctx, fullPrompt, schemaString)
}

// generateImage calls an image generation service and uploads to S3.
func (s *MonsterGenerationService) generateImage(ctx context.Context, visualDescription string) (string, error) {
	if s.s3Service == nil {
		return "", fmt.Errorf("S3 service not initialized for image upload")
	}

	// --- Placeholder Image Generation API Call ---
	// In a real scenario, this would call a service like Vertex AI Imagen, DALL-E, etc.
	// For now, it will mock fetching an image and uploading it.

	// Use a slightly dynamic placeholder based on visualDescription
	placeholderKeyword := "monster"
	if strings.Contains(strings.ToLower(visualDescription), "dragon") {
		placeholderKeyword = "dragon"
	} else if strings.Contains(strings.ToLower(visualDescription), "goblin") {
		placeholderKeyword = "goblin"
	} else if strings.Contains(strings.ToLower(visualDescription), "beast") {
		placeholderKeyword = "beast"
	}

	// Example: Fetch a dynamic placeholder image from a service like placehold.co or picsum.photos
	// For simplicity, let's use a fixed size and just vary text
	placeholderImageURL := fmt.Sprintf("https://placehold.co/300x300/FF0000/FFFFFF/png?text=%s", strings.ReplaceAll(strings.ToTitle(placeholderKeyword), " ", "+"))

	log.Ctx(ctx).Info().Msgf("Simulating image generation for: '%s', fetching from %s", visualDescription, placeholderImageURL)

	resp, err := http.Get(placeholderImageURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch placeholder image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch placeholder image, status: %s", resp.Status)
	}

	// Step C: Image Processing & S3 Upload
	filename := fmt.Sprintf("monster-%s.png", uuid.New().String())
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/png" // Default if not provided by the placeholder
	}

	uploadedURL, err := s.s3Service.UploadStream(ctx, resp.Body, filename, contentType)
	if err != nil {
		return "", fmt.Errorf("failed to upload generated image to S3: %w", err)
	}

	return uploadedURL, nil
}
