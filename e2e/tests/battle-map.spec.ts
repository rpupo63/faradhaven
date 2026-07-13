import { test, expect } from "@playwright/test";
import { API_URL } from "../playwright.config";
import { registerUser, seedAuth, authHeaders, createCharacter, type TestUser } from "../helpers";

// Map CRUD + tokens + initiative through the API (websocket collaboration and
// S3 background uploads are intentionally out of scope).

test.describe("battle maps", () => {
  let user: TestUser;

  test.beforeEach(async ({ page, request }) => {
    user = await registerUser(request);
    await seedAuth(page, user);
  });

  test("map with token and initiative round-trips; battle map page renders", async ({ page, request }) => {
    const created = await request.post(`${API_URL}/map`, {
      headers: authHeaders(user),
      data: {
        name: "E2E Skirmish Grounds",
        room_code: `E2E${Date.now() % 100000}`,
        grid_rows: 10,
        grid_cols: 10,
        tile_size: 50,
      },
    });
    expect(created.status(), await created.text()).toBe(201);
    const map = await created.json();

    // Public map read works.
    const fetched = await request.get(`${API_URL}/map/${map.id}`);
    expect(fetched.status()).toBe(200);

    // Add a token for one of our characters.
    const character = await createCharacter(request, user, "Token Bearer");
    const tokenRes = await request.post(`${API_URL}/map/${map.id}/token`, {
      headers: authHeaders(user),
      data: {
        character_id: character.id,
        name: "Token Bearer",
        token_type: "pc",
        grid_x: 2,
        grid_y: 3,
      },
    });
    expect(tokenRes.ok(), await tokenRes.text()).toBeTruthy();
    const token = await tokenRes.json();

    // Initiative set + get.
    const setInit = await request.put(`${API_URL}/map/${map.id}/initiative`, {
      headers: authHeaders(user),
      data: { entries: [{ token_id: token.id, order: 1 }] },
    });
    expect(setInit.ok(), await setInit.text()).toBeTruthy();

    const getInit = await request.get(`${API_URL}/map/${map.id}/initiative`, {
      headers: authHeaders(user),
    });
    expect(getInit.status()).toBe(200);

    // The user's maps page shows the new map.
    await page.goto("/battle-map");
    await expect(page.getByText("E2E Skirmish Grounds").first()).toBeVisible({ timeout: 15_000 });

    // Cleanup.
    const deleted = await request.delete(`${API_URL}/map/${map.id}`, {
      headers: authHeaders(user),
    });
    expect(deleted.ok()).toBeTruthy();
  });
});
