package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/database"
	"github.com/rpupo63/faradhaven/backend/models"
	"github.com/rpupo63/faradhaven/backend/seed/faradhaven_classes"
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
	eventRepo   database.MonsterGenerationEventRepository
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
	eventRepo database.MonsterGenerationEventRepository,
	s3Service *S3Service,
	llmClient LLMClient, // New: LLM Client
) *MonsterGenerationService {
	return &MonsterGenerationService{
		monsterRepo: monsterRepo,
		eventRepo:   eventRepo,
		s3Service:   s3Service,
		llmClient:   llmClient,
		schemaPath:  "docs/schemas/monster_schema.json",
	}
}

type GenerationContext struct {
	Role                string `json:"role,omitempty"`
	Environment         string `json:"environment,omitempty"`
	Temperament         string `json:"temperament,omitempty"`
	EncounterGoal       string `json:"encounter_goal,omitempty"`
	PartyLevel          int    `json:"party_level,omitempty"`
	TemplateID          string `json:"template_id,omitempty"`
	ClassThemeIntensity string `json:"class_theme_intensity,omitempty"` // light,strong
}

// GenerateMonsterFromPromptRequest encapsulates the input for monster generation.
type GenerateMonsterFromPromptRequest struct {
	UserID          uuid.UUID `json:"user_id"`
	Description     string    `json:"description"`
	ChallengeRating string    `json:"challenge_rating"` // e.g., "1", "1/2", "CR 5"
	// FaradhavenClassName, when non-empty, selects a seeded Faradhaven class from
	// seed/faradhaven_classes. The LLM receives that class's full seed JSON as context
	// and builds a monster stat block themed as an enemy embodying that class (not a PC sheet).
	FaradhavenClassName string            `json:"faradhaven_class_name,omitempty"`
	GenerationContext   GenerationContext `json:"generation_context,omitempty"`
}

// GenerateMonsterFromPrompt orchestrates the monster generation process.
func (s *MonsterGenerationService) GenerateMonsterFromPrompt(ctx context.Context, req GenerateMonsterFromPromptRequest) (*models.Monster, error) {
	monster, retries, err := s.generateMonsterStruct(ctx, req)
	if err != nil {
		s.logEvent(ctx, req.UserID, nil, "create", retries, false, "llm_or_parse_error", req)
		return nil, err
	}
	if err := s.monsterRepo.Add(monster); err != nil {
		s.logEvent(ctx, req.UserID, nil, "create", retries, false, "db_save_error", req)
		return nil, fmt.Errorf("failed to save generated monster to database: %w", err)
	}
	s.logEvent(ctx, req.UserID, &monster.ID, "create", retries, true, "", req)
	return monster, nil
}

func (s *MonsterGenerationService) PreviewMonsterFromPrompt(ctx context.Context, req GenerateMonsterFromPromptRequest) (*models.Monster, error) {
	monster, retries, err := s.generateMonsterStruct(ctx, req)
	if err != nil {
		s.logEvent(ctx, req.UserID, nil, "preview", retries, false, "llm_or_parse_error", req)
		return nil, err
	}
	s.logEvent(ctx, req.UserID, nil, "preview", retries, true, "", req)
	return monster, nil
}

func (s *MonsterGenerationService) RegenerateSection(ctx context.Context, base *models.Monster, userID uuid.UUID, section string) (*models.Monster, error) {
	desc := fmt.Sprintf("Regenerate only the %s for this monster.\n\nMonster JSON:\n%s", section, mustJSON(base))
	req := GenerateMonsterFromPromptRequest{
		UserID:              userID,
		Description:         desc,
		ChallengeRating:     base.ChallengeRating,
		FaradhavenClassName: deref(base.GenerationClassName),
	}
	regen, retries, err := s.generateMonsterStruct(ctx, req)
	if err != nil {
		s.logEvent(ctx, userID, &base.ID, "regenerate", retries, false, "regen_error", req)
		return nil, err
	}
	switch strings.ToLower(strings.TrimSpace(section)) {
	case "attacks":
		base.Attacks = regen.Attacks
	case "traits", "special_traits":
		base.SpecialTraits = regen.SpecialTraits
	case "lore", "description":
		base.VisualDescription = regen.VisualDescription
		base.Notes = regen.Notes
	default:
		base.Actions = regen.Actions
	}
	s.logEvent(ctx, userID, &base.ID, "regenerate", retries, true, "", req)
	return base, nil
}

func (s *MonsterGenerationService) BuildVariant(ctx context.Context, base *models.Monster, userID uuid.UUID, variant string) (*models.Monster, error) {
	desc := fmt.Sprintf("Create a %s variant of this monster while preserving fantasy and core identity.\n\nMonster JSON:\n%s", variant, mustJSON(base))
	req := GenerateMonsterFromPromptRequest{
		UserID:              userID,
		Description:         desc,
		ChallengeRating:     base.ChallengeRating,
		FaradhavenClassName: deref(base.GenerationClassName),
	}
	variantMonster, retries, err := s.generateMonsterStruct(ctx, req)
	if err != nil {
		s.logEvent(ctx, userID, &base.ID, "variant", retries, false, "variant_error", req)
		return nil, err
	}
	v := strings.TrimSpace(variant)
	variantMonster.GenerationMode = "variant"
	variantMonster.GenerationTemplate = &v
	s.logEvent(ctx, userID, &base.ID, "variant", retries, true, "", req)
	return variantMonster, nil
}

func (s *MonsterGenerationService) generateMonsterStruct(ctx context.Context, req GenerateMonsterFromPromptRequest) (*models.Monster, int, error) {
	className := strings.TrimSpace(req.FaradhavenClassName)

	var llmOutput string
	var err error
	retries := 0
	if className != "" {
		llmOutput, err = s.generateMonsterFromFaradhavenClass(ctx, className, req.Description, req.ChallengeRating, req.GenerationContext)
	} else {
		llmOutput, err = s.generateMonsterDataFromLLM(ctx, req.Description, req.ChallengeRating, req.GenerationContext)
	}
	if err != nil {
		return nil, retries, fmt.Errorf("LLM text generation failed: %w", err)
	}

	// Unmarshal LLM output into a Monster struct
	var monster models.Monster
	if err := json.Unmarshal([]byte(llmOutput), &monster); err != nil {
		// Retry once with an explicit correction instruction.
		retries++
		retryPrompt := req.Description + "\n\nIMPORTANT: Return strictly valid JSON that matches the schema exactly."
		if className != "" {
			llmOutput, err = s.generateMonsterFromFaradhavenClass(ctx, className, retryPrompt, req.ChallengeRating, req.GenerationContext)
		} else {
			llmOutput, err = s.generateMonsterDataFromLLM(ctx, retryPrompt, req.ChallengeRating, req.GenerationContext)
		}
		if err != nil || json.Unmarshal([]byte(llmOutput), &monster) != nil {
			return nil, retries, fmt.Errorf("failed to unmarshal LLM output into Monster struct: %w", err)
		}
	}

	monster.NormalizeCreatureFields()

	// Overwrite some fields from the request/system
	monster.UserID = req.UserID
	monster.Notes = req.Description // Store original prompt in Notes
	if className != "" {
		monster.Source = fmt.Sprintf("User Generated (Faradhaven class: %s)", className)
		monster.GenerationMode = "class-themed"
		monster.GenerationClassName = &className
	} else {
		monster.Source = "User Generated"
		monster.GenerationMode = "custom"
	}
	monster.ChallengeRating = req.ChallengeRating // Ensure CR from request is used
	if req.GenerationContext.TemplateID != "" {
		templateID := req.GenerationContext.TemplateID
		monster.GenerationTemplate = &templateID
	}
	ctxJSON, _ := json.Marshal(req.GenerationContext)
	monster.GenerationContext = ctxJSON

	// Step B: Image Generation
	imageURL, err := s.generateImage(ctx, monster.VisualDescription)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Image generation failed, proceeding without image.")
		monster.ImageURL = nil // Set to nil if image generation fails
	} else {
		monster.ImageURL = &imageURL
	}

	return &monster, retries, nil
}

// generateMonsterDataFromLLM uses the LLMClient to get structured monster data.
func (s *MonsterGenerationService) generateMonsterDataFromLLM(ctx context.Context, prompt, cr string, generationContext GenerationContext) (string, error) {
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
	Generation context: %s

	JSON Schema:
	%s
	`, prompt, cr, mustJSON(generationContext), schemaString)

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
func (s *MonsterGenerationService) generateMonsterFromFaradhavenClass(ctx context.Context, className, flavor, cr string, generationContext GenerationContext) (string, error) {
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
Generation context: %s

JSON Schema:
%s
`, string(classJSON), flavor, cr, mustJSON(generationContext), schemaString)

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

func (s *MonsterGenerationService) logEvent(ctx context.Context, userID uuid.UUID, monsterID *uuid.UUID, eventType string, retries int, success bool, failureType string, req GenerateMonsterFromPromptRequest) {
	if s.eventRepo == nil {
		return
	}
	start := time.Now()
	meta, _ := json.Marshal(map[string]any{
		"challenge_rating": req.ChallengeRating,
		"class_name":       req.FaradhavenClassName,
		"generation_mode":  req.GenerationContext,
	})
	event := &models.MonsterGenerationEvent{
		UserID:        userID,
		MonsterID:     monsterID,
		EventType:     eventType,
		PromptVersion: "v2",
		ModelName:     llmMonsterGenerationModel,
		LatencyMs:     time.Since(start).Milliseconds(),
		RetryCount:    retries,
		Success:       success,
		FailureType:   failureType,
		Metadata:      meta,
	}
	if err := s.eventRepo.Add(event); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("failed to persist monster generation event")
	}
}

func mustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func deref(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
