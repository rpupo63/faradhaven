package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/rpupo63/faradhaven/backend/internal/testutils"
)

// The tests below drive the real chi router (all middleware, auth, and
// handlers) against a seeded Postgres database. Without GOOGLE_API_KEY the
// spell AI falls back to the built-in mock client, and without AWS creds the
// S3 service is disabled — both are the intended offline test setup.

var (
	routerOnce sync.Once
	testRouter http.Handler
)

func apiRouter(t *testing.T) http.Handler {
	t.Helper()
	_, app := testutils.SetupSeededTestDB(t)
	routerOnce.Do(func() {
		testRouter = newRouter(app)
	})
	return testRouter
}

type apiClient struct {
	t      *testing.T
	router http.Handler
	token  string
}

func newAPIClient(t *testing.T) *apiClient {
	return &apiClient{t: t, router: apiRouter(t)}
}

func (c *apiClient) do(method, path string, body any) (*httptest.ResponseRecorder, map[string]any) {
	c.t.Helper()
	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			c.t.Fatalf("marshal body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	rr := httptest.NewRecorder()
	c.router.ServeHTTP(rr, req)

	var decoded map[string]any
	if len(rr.Body.Bytes()) > 0 {
		_ = json.Unmarshal(rr.Body.Bytes(), &decoded)
	}
	return rr, decoded
}

func (c *apiClient) doList(method, path string) (*httptest.ResponseRecorder, []map[string]any) {
	c.t.Helper()
	rr, _ := c.do(method, path, nil)
	var list []map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &list)
	return rr, list
}

var userCounter int

// registerAndLogin creates a fresh user through the public auth endpoints and
// stores its Bearer token on the client.
func (c *apiClient) registerAndLogin() (userID string) {
	c.t.Helper()
	userCounter++
	email := fmt.Sprintf("it-%d-%d@example.com", userCounter, len(c.token))
	rr, body := c.do("POST", "/api/auth/register", map[string]string{
		"name": "Integration User", "email": email, "password": "hunter22!",
	})
	if rr.Code != http.StatusCreated {
		c.t.Fatalf("register failed: %d %s", rr.Code, rr.Body.String())
	}
	c.token = body["token"].(string)
	return body["user_id"].(string)
}

func mustStatus(t *testing.T, rr *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rr.Code != want {
		t.Fatalf("got status %d, want %d — body: %s", rr.Code, want, rr.Body.String())
	}
}

func TestAuthEndpoints(t *testing.T) {
	c := newAPIClient(t)

	rr, body := c.do("POST", "/api/auth/register", map[string]string{
		"name": "Auth Tester", "email": "auth-tester@example.com", "password": "hunter22!",
	})
	mustStatus(t, rr, http.StatusCreated)
	if body["token"] == nil || body["token"] == "" {
		t.Fatalf("expected token, got %v", body)
	}

	// Duplicate email rejected.
	rr, _ = c.do("POST", "/api/auth/register", map[string]string{
		"name": "Auth Tester 2", "email": "auth-tester@example.com", "password": "hunter22!",
	})
	mustStatus(t, rr, http.StatusConflict)

	rr, body = c.do("POST", "/api/auth/login", map[string]string{
		"email": "auth-tester@example.com", "password": "hunter22!",
	})
	mustStatus(t, rr, http.StatusOK)
	if body["token"] == nil {
		t.Fatalf("expected token on login, got %v", body)
	}

	rr, _ = c.do("POST", "/api/auth/login", map[string]string{
		"email": "auth-tester@example.com", "password": "wrong",
	})
	mustStatus(t, rr, http.StatusUnauthorized)

	// Protected route without a token is rejected.
	anon := newAPIClient(t)
	rr, _ = anon.do("GET", "/api/characters", nil)
	mustStatus(t, rr, http.StatusUnauthorized)
}

func TestReferenceDataSeeded(t *testing.T) {
	c := newAPIClient(t)

	for _, path := range []string{"/api/classes", "/api/races", "/api/components", "/api/weapons"} {
		rr, list := c.doList("GET", path)
		mustStatus(t, rr, http.StatusOK)
		if len(list) == 0 {
			t.Fatalf("expected seeded data at %s, got empty list", path)
		}
	}

	rr, _ := c.do("GET", "/api/characters/options", nil)
	mustStatus(t, rr, http.StatusOK)
}

// createTestCharacter registers a user and creates a level-1 character using
// seeded class/race reference data. Returns the character ID.
func createTestCharacter(t *testing.T, c *apiClient, name string) (userID, characterID string) {
	t.Helper()
	userID = c.registerAndLogin()

	_, classes := c.doList("GET", "/api/classes")
	_, races := c.doList("GET", "/api/races")
	if len(classes) == 0 || len(races) == 0 {
		t.Fatal("seeded classes/races required")
	}

	rr, body := c.do("POST", "/api/character", map[string]any{
		"user_id":      userID,
		"name":         name,
		"race_id":      races[0]["id"],
		"class_id":     classes[0]["id"],
		"level":        1,
		"strength":     14,
		"dexterity":    12,
		"constitution": 13,
		"intelligence": 10,
		"wisdom":       11,
		"charisma":     10,
		"money":        100,
	})
	mustStatus(t, rr, http.StatusCreated)
	characterID, _ = body["id"].(string)
	if characterID == "" {
		t.Fatalf("expected character id in response: %v", body)
	}
	return userID, characterID
}

func TestCharacterHPAndRest(t *testing.T) {
	c := newAPIClient(t)
	_, characterID := createTestCharacter(t, c, "HP Tester")

	// Baseline sheet
	rr, _ := c.do("GET", "/api/character/"+characterID+"/sheet", nil)
	mustStatus(t, rr, http.StatusOK)

	// Damage
	rr, body := c.do("PATCH", "/api/character/"+characterID+"/hp", map[string]any{"delta": -3})
	mustStatus(t, rr, http.StatusOK)
	maxHP := int(body["max_hp"].(float64))
	currentHP := int(body["current_hp"].(float64))
	if currentHP != maxHP-3 {
		t.Fatalf("expected current_hp %d, got %d", maxHP-3, currentHP)
	}

	// Overheal clamps at max
	rr, body = c.do("PATCH", "/api/character/"+characterID+"/hp", map[string]any{"delta": 999})
	mustStatus(t, rr, http.StatusOK)
	if int(body["current_hp"].(float64)) != maxHP {
		t.Fatalf("expected HP clamped to max %d, got %v", maxHP, body["current_hp"])
	}

	// Temp HP
	rr, body = c.do("PUT", "/api/character/"+characterID+"/temp-hp", map[string]any{"temp_hp": 5})
	mustStatus(t, rr, http.StatusOK)
	if int(body["temp_hp"].(float64)) != 5 {
		t.Fatalf("expected temp_hp 5, got %v", body["temp_hp"])
	}

	// Rests succeed
	rr, _ = c.do("POST", "/api/character/"+characterID+"/rest/short", map[string]any{})
	mustStatus(t, rr, http.StatusOK)
	rr, _ = c.do("POST", "/api/character/"+characterID+"/rest/long", map[string]any{})
	mustStatus(t, rr, http.StatusOK)
}

func TestLevelUpFlow(t *testing.T) {
	c := newAPIClient(t)
	_, characterID := createTestCharacter(t, c, "Level Tester")

	rr, _ := c.do("GET", "/api/character/"+characterID+"/level-up/preview", nil)
	mustStatus(t, rr, http.StatusOK)

	// Level up with defaults; if the next level requires an archetype choice,
	// retry with the first archetype of the character's class.
	rr, body := c.do("POST", "/api/character/"+characterID+"/level-up", map[string]any{})
	if rr.Code == http.StatusBadRequest {
		charRr, charBody := c.do("GET", "/api/character/"+characterID, nil)
		mustStatus(t, charRr, http.StatusOK)
		classID, _ := charBody["class_id"].(string)
		_, classBody := c.do("GET", "/api/classes/"+classID, nil)
		archetypes, _ := classBody["archetypes"].([]any)
		if len(archetypes) == 0 {
			t.Fatalf("level-up failed and no archetypes to retry with: %s", rr.Body.String())
		}
		first := archetypes[0].(map[string]any)
		rr, body = c.do("POST", "/api/character/"+characterID+"/level-up", map[string]any{
			"archetype_id": first["id"],
		})
	}
	mustStatus(t, rr, http.StatusOK)
	if int(body["new_level"].(float64)) != 2 {
		t.Fatalf("expected new_level 2, got %v", body["new_level"])
	}

	// History records the level-up.
	rr, _ = c.do("GET", "/api/character/"+characterID+"/level-history", nil)
	mustStatus(t, rr, http.StatusOK)
}

func TestPartyLifecycle(t *testing.T) {
	c := newAPIClient(t)
	_, characterID := createTestCharacter(t, c, "Party Tester")

	rr, party := c.do("POST", "/api/parties", map[string]any{"name": "The Test Fellowship"})
	mustStatus(t, rr, http.StatusCreated)
	partyID := party["id"].(string)

	rr, _ = c.do("POST", "/api/parties/"+partyID+"/members", map[string]any{
		"character_id": characterID,
	})
	if rr.Code != http.StatusOK && rr.Code != http.StatusCreated && rr.Code != http.StatusNoContent {
		t.Fatalf("add member failed: %d %s", rr.Code, rr.Body.String())
	}

	rr, _ = c.do("GET", "/api/parties/"+partyID, nil)
	mustStatus(t, rr, http.StatusOK)

	rr, _ = c.do("DELETE", "/api/parties/"+partyID+"/members/"+characterID, nil)
	if rr.Code != http.StatusOK && rr.Code != http.StatusNoContent {
		t.Fatalf("remove member failed: %d %s", rr.Code, rr.Body.String())
	}

	rr, _ = c.do("DELETE", "/api/parties/"+partyID, nil)
	if rr.Code != http.StatusOK && rr.Code != http.StatusNoContent {
		t.Fatalf("delete party failed: %d %s", rr.Code, rr.Body.String())
	}
}

func TestSpellCRUDWithMockAI(t *testing.T) {
	c := newAPIClient(t)
	userID, _ := createTestCharacter(t, c, "Spell Tester")

	_, components := c.doList("GET", "/api/components")
	if len(components) == 0 {
		t.Fatal("expected seeded components")
	}

	// Valid composition: exactly 1 Forma (shape) + 1 Scopus (targeting) + ≥1 Essentia.
	byCategory := func(category string) any {
		for _, comp := range components {
			if comp["category"] == category {
				return comp["id"]
			}
		}
		t.Fatalf("no seeded component with category %q", category)
		return nil
	}

	rr, spell := c.do("POST", "/api/spell", map[string]any{
		"user_id":       userID,
		"name":          "Integration Bolt",
		"description":   "A deterministic bolt of testing energy.",
		"type":              "Attack",
		"concentration":     false,
		"add_modifier":      false,
		"damage_dice_count": 1,
		"damage_die_size":   8,
		"component_ids": []any{byCategory("Forma"), byCategory("Scopus"), byCategory("Essentia")},
	})
	mustStatus(t, rr, http.StatusCreated)
	spellID, _ := spell["id"].(string)
	if spellID == "" {
		t.Fatalf("expected spell id, got %v", spell)
	}

	rr, got := c.do("GET", "/api/spell/"+spellID, nil)
	mustStatus(t, rr, http.StatusOK)
	if got["name"] != "Integration Bolt" {
		t.Fatalf("expected spell name, got %v", got["name"])
	}

	rr, _ = c.do("DELETE", "/api/spell/"+spellID, nil)
	if rr.Code != http.StatusOK && rr.Code != http.StatusNoContent {
		t.Fatalf("delete spell failed: %d %s", rr.Code, rr.Body.String())
	}
}

func TestResourceEndpoints(t *testing.T) {
	c := newAPIClient(t)
	_, characterID := createTestCharacter(t, c, "Resource Tester")

	rr, _ := c.do("GET", "/api/characters/"+characterID+"/resources", nil)
	mustStatus(t, rr, http.StatusOK)
}
