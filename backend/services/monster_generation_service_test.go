package services

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/rpupo63/faradhaven/backend/models"
)

// --- Test mocks ---

// testLLMClient returns a fixed JSON string, no disk reads needed.
type testLLMClient struct {
	response  string
	responses []string
	calls     int
	err       error
}

func (t *testLLMClient) GenerateStructuredContent(ctx context.Context, prompt string, schema string) (string, error) {
	t.calls++
	if len(t.responses) >= t.calls {
		return t.responses[t.calls-1], t.err
	}
	return t.response, t.err
}

// inMemoryMonsterRepo stores monsters in a slice for assertions.
type inMemoryMonsterRepo struct {
	monsters []*models.Monster
}

func (r *inMemoryMonsterRepo) Add(m *models.Monster) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	r.monsters = append(r.monsters, m)
	return nil
}
func (r *inMemoryMonsterRepo) FindByID(id uuid.UUID) (*models.Monster, error) { return nil, nil }
func (r *inMemoryMonsterRepo) FindByUserID(userID uuid.UUID) ([]*models.Monster, error) {
	return nil, nil
}
func (r *inMemoryMonsterRepo) FindAll() ([]*models.Monster, error) { return nil, nil }
func (r *inMemoryMonsterRepo) Update(m *models.Monster) error      { return nil }
func (r *inMemoryMonsterRepo) Delete(id uuid.UUID) error           { return nil }

// --- Tests ---

const validMonsterJSON = `{
	"name": "Fiery Goblin",
	"size": "Small",
	"type": "Humanoid",
	"alignment": "Chaotic Evil",
	"challenge_rating": "1/2",
	"proficiency_bonus": 2,
	"armor_class": 13,
	"hit_points": 18,
	"hit_dice": "4d6 + 4",
	"speed": "30 ft",
	"strength": 8,
	"dexterity": 14,
	"constitution": 12,
	"intelligence": 10,
	"wisdom": 8,
	"charisma": 8,
	"saving_throws": ["Dex"],
	"skills": ["Stealth +4"],
	"damage_immunities": ["fire"],
	"damage_resistances": [],
	"condition_immunities": [],
	"senses": "darkvision 60 ft, passive Perception 9",
	"languages": "Goblin, Common",
	"notes": "",
	"image_url": null,
	"attacks": [
		{
			"name": "Flame Dagger",
			"description": "Melee Weapon Attack: +4 to hit, reach 5 ft., one target. Hit: 1d4 + 2 fire damage."
		}
	],
	"actions": [
		{
			"name": "Fire Breath (Recharge 5-6)",
			"description": "The goblin exhales fire in a 15-foot cone. Each creature in the area must make a DC 11 Dexterity saving throw, taking 2d6 fire damage on a failed save, or half as much on a successful one."
		}
	],
	"special_traits": ["Fire Absorption: Whenever the goblin takes fire damage, it regains hit points equal to half the damage dealt."],
	"legendary_actions": [],
	"reactions": [],
	"lair_actions": [],
	"environments": ["Underground", "Volcanic"],
	"source": "Generated",
	"visual_description": "A small goblin wreathed in crackling flames."
}`

func TestGenerateMonsterFromPrompt_Success(t *testing.T) {
	repo := &inMemoryMonsterRepo{}
	llm := &testLLMClient{response: validMonsterJSON}
	svc := NewMonsterGenerationService(repo, nil, nil, llm) // nil s3 → image gen skipped gracefully
	svc.schemaPath = "../docs/schemas/monster_schema.json"

	userID := uuid.New()
	monster, err := svc.GenerateMonsterFromPrompt(context.Background(), GenerateMonsterFromPromptRequest{
		UserID:          userID,
		Description:     "A fire-breathing goblin",
		ChallengeRating: "1/2",
	})

	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Request-level fields are applied on top of LLM output
	if monster.UserID != userID {
		t.Errorf("UserID = %v, want %v", monster.UserID, userID)
	}
	if monster.ChallengeRating != "1/2" {
		t.Errorf("ChallengeRating = %q, want %q", monster.ChallengeRating, "1/2")
	}
	if monster.Source != "User Generated" {
		t.Errorf("Source = %q, want %q", monster.Source, "User Generated")
	}
	if monster.Notes != "A fire-breathing goblin" {
		t.Errorf("Notes = %q, want original description", monster.Notes)
	}

	// LLM-generated fields are preserved
	if monster.Name != "Fiery Goblin" {
		t.Errorf("Name = %q, want %q", monster.Name, "Fiery Goblin")
	}
	if monster.ArmorClass != 13 {
		t.Errorf("ArmorClass = %d, want 13", monster.ArmorClass)
	}
	if len(monster.Attacks) != 1 {
		t.Errorf("Attacks count = %d, want 1", len(monster.Attacks))
	} else if monster.Attacks[0].Name != "Flame Dagger" {
		t.Errorf("Attack name = %q, want %q", monster.Attacks[0].Name, "Flame Dagger")
	}
	if len(monster.Actions) != 1 {
		t.Errorf("Actions count = %d, want 1", len(monster.Actions))
	}

	// Image should be nil since s3Service is nil
	if monster.ImageURL != nil {
		t.Errorf("ImageURL = %v, want nil (no s3 service)", *monster.ImageURL)
	}

	// Monster was persisted in the repo
	if len(repo.monsters) != 1 {
		t.Fatalf("repo has %d monsters, want 1", len(repo.monsters))
	}
}

func TestGenerateMonsterFromPrompt_LLMError(t *testing.T) {
	repo := &inMemoryMonsterRepo{}
	llm := &testLLMClient{err: fmt.Errorf("model overloaded")}
	svc := NewMonsterGenerationService(repo, nil, nil, llm)

	_, err := svc.GenerateMonsterFromPrompt(context.Background(), GenerateMonsterFromPromptRequest{
		UserID:          uuid.New(),
		Description:     "anything",
		ChallengeRating: "1",
	})

	if err == nil {
		t.Fatal("expected error when LLM fails, got nil")
	}
}

func TestGenerateMonsterFromPrompt_InvalidJSON(t *testing.T) {
	repo := &inMemoryMonsterRepo{}
	llm := &testLLMClient{response: `{"name": "broken`} // malformed JSON
	svc := NewMonsterGenerationService(repo, nil, nil, llm)

	_, err := svc.GenerateMonsterFromPrompt(context.Background(), GenerateMonsterFromPromptRequest{
		UserID:          uuid.New(),
		Description:     "anything",
		ChallengeRating: "1",
	})

	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestGenerateMonsterFromPrompt_NilLLMClient(t *testing.T) {
	repo := &inMemoryMonsterRepo{}
	svc := NewMonsterGenerationService(repo, nil, nil, nil) // nil LLM client

	_, err := svc.GenerateMonsterFromPrompt(context.Background(), GenerateMonsterFromPromptRequest{
		UserID:          uuid.New(),
		Description:     "anything",
		ChallengeRating: "1",
	})

	if err == nil {
		t.Fatal("expected error when LLM client is nil, got nil")
	}
}

func TestGenerateMonsterFromPrompt_FaradhavenClass(t *testing.T) {
	repo := &inMemoryMonsterRepo{}
	llm := &testLLMClient{response: validMonsterJSON}
	svc := NewMonsterGenerationService(repo, nil, nil, llm)
	svc.schemaPath = "../docs/schemas/monster_schema.json"

	userID := uuid.New()
	monster, err := svc.GenerateMonsterFromPrompt(context.Background(), GenerateMonsterFromPromptRequest{
		UserID:              userID,
		Description:         "patrolling the alchemy ward",
		ChallengeRating:     "3",
		FaradhavenClassName: "The Mutagen",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if monster.UserID != userID {
		t.Errorf("UserID = %v, want %v", monster.UserID, userID)
	}
	if !strings.Contains(monster.Source, "Mutagen") {
		t.Errorf("Source = %q, want substring Mutagen", monster.Source)
	}
	if len(repo.monsters) != 1 {
		t.Fatalf("repo has %d monsters, want 1", len(repo.monsters))
	}
}

func TestGenerateMonsterFromPrompt_UnknownFaradhavenClass(t *testing.T) {
	repo := &inMemoryMonsterRepo{}
	llm := &testLLMClient{response: validMonsterJSON}
	svc := NewMonsterGenerationService(repo, nil, nil, llm)

	_, err := svc.GenerateMonsterFromPrompt(context.Background(), GenerateMonsterFromPromptRequest{
		UserID:              uuid.New(),
		Description:         "flavor",
		ChallengeRating:     "1",
		FaradhavenClassName: "NotARealClassNameXYZ",
	})
	if err == nil {
		t.Fatal("expected error for unknown class, got nil")
	}
	if len(repo.monsters) != 0 {
		t.Errorf("repo should have no monsters, got %d", len(repo.monsters))
	}
}

func TestPreviewMonsterFromPrompt_DoesNotPersist(t *testing.T) {
	repo := &inMemoryMonsterRepo{}
	llm := &testLLMClient{response: validMonsterJSON}
	svc := NewMonsterGenerationService(repo, nil, nil, llm)
	svc.schemaPath = "../docs/schemas/monster_schema.json"

	monster, err := svc.PreviewMonsterFromPrompt(context.Background(), GenerateMonsterFromPromptRequest{
		UserID:          uuid.New(),
		Description:     "preview this",
		ChallengeRating: "2",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if monster == nil {
		t.Fatal("expected preview monster")
	}
	if len(repo.monsters) != 0 {
		t.Fatalf("preview should not persist monsters, got %d", len(repo.monsters))
	}
}

func TestGenerateMonsterFromPrompt_RetrysAfterInvalidJSON(t *testing.T) {
	repo := &inMemoryMonsterRepo{}
	llm := &testLLMClient{
		responses: []string{
			`{"name":"bad"`,
			validMonsterJSON,
		},
	}
	svc := NewMonsterGenerationService(repo, nil, nil, llm)
	svc.schemaPath = "../docs/schemas/monster_schema.json"

	_, err := svc.GenerateMonsterFromPrompt(context.Background(), GenerateMonsterFromPromptRequest{
		UserID:          uuid.New(),
		Description:     "retry",
		ChallengeRating: "1",
	})
	if err != nil {
		t.Fatalf("expected retry success, got error: %v", err)
	}
	if llm.calls < 2 {
		t.Fatalf("expected at least 2 llm calls, got %d", llm.calls)
	}
}
