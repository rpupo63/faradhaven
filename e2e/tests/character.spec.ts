import { test, expect } from "@playwright/test";
import { API_URL } from "../playwright.config";
import { registerUser, seedAuth, authHeaders, createCharacter, type TestUser } from "../helpers";

test.describe("character lifecycle", () => {
  let user: TestUser;

  test.beforeEach(async ({ page, request }) => {
    user = await registerUser(request);
    await seedAuth(page, user);
  });

  test("created character appears in the roster and its sheet renders", async ({ page, request }) => {
    const character = await createCharacter(request, user, "Sir Testalot");

    await page.goto("/characters");
    await expect(page.getByText("Sir Testalot").first()).toBeVisible({ timeout: 15_000 });

    await page.goto(`/character/${character.id}/sheet`);
    await expect(page.getByRole("heading", { name: "Sir Testalot" })).toBeVisible({
      timeout: 15_000,
    });
  });

  test("HP damage, temp HP, and rests round-trip through the API", async ({ request }) => {
    const character = await createCharacter(request, user, "HP Dummy");

    const damage = await request.patch(`${API_URL}/character/${character.id}/hp`, {
      headers: authHeaders(user),
      data: { delta: -3 },
    });
    expect(damage.status()).toBe(200);
    const afterDamage = await damage.json();
    expect(afterDamage.current_hp).toBe(afterDamage.max_hp - 3);

    const temp = await request.put(`${API_URL}/character/${character.id}/temp-hp`, {
      headers: authHeaders(user),
      data: { temp_hp: 5 },
    });
    expect(temp.status()).toBe(200);
    expect((await temp.json()).temp_hp).toBe(5);

    const short = await request.post(`${API_URL}/character/${character.id}/rest/short`, {
      headers: authHeaders(user),
      data: {},
    });
    expect(short.status()).toBe(200);

    const long = await request.post(`${API_URL}/character/${character.id}/rest/long`, {
      headers: authHeaders(user),
      data: {},
    });
    expect(long.status()).toBe(200);
  });

  test("level-up preview is available for a fresh character", async ({ request }) => {
    const character = await createCharacter(request, user, "Level Candidate");

    const preview = await request.get(
      `${API_URL}/character/${character.id}/level-up/preview`,
      { headers: authHeaders(user) },
    );
    expect(preview.status()).toBe(200);
  });

  test("party create, add member, and delete", async ({ request }) => {
    const character = await createCharacter(request, user, "Party Animal");

    const created = await request.post(`${API_URL}/parties`, {
      headers: authHeaders(user),
      data: { name: "E2E Fellowship" },
    });
    expect(created.status()).toBe(201);
    const party = await created.json();

    const added = await request.post(`${API_URL}/parties/${party.id}/members`, {
      headers: authHeaders(user),
      data: { character_id: character.id },
    });
    expect(added.ok()).toBeTruthy();

    const fetched = await request.get(`${API_URL}/parties/${party.id}`, {
      headers: authHeaders(user),
    });
    expect(fetched.status()).toBe(200);

    const removed = await request.delete(
      `${API_URL}/parties/${party.id}/members/${character.id}`,
      { headers: authHeaders(user) },
    );
    expect(removed.ok()).toBeTruthy();

    const deleted = await request.delete(`${API_URL}/parties/${party.id}`, {
      headers: authHeaders(user),
    });
    expect(deleted.ok()).toBeTruthy();
  });
});
