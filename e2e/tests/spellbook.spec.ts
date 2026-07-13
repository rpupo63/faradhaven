import { test, expect } from "@playwright/test";
import { API_URL } from "../playwright.config";
import { registerUser, seedAuth, authHeaders, type TestUser } from "../helpers";

// Spell creation runs against the mock LLM client (no GOOGLE_API_KEY), so the
// AI review path is exercised without any external calls.

test.describe("spellbook", () => {
  let user: TestUser;

  test.beforeEach(async ({ page, request }) => {
    user = await registerUser(request);
    await seedAuth(page, user);
  });

  test("spell composed of seeded components appears in the spellbook", async ({ page, request }) => {
    const componentsRes = await request.get(`${API_URL}/components`);
    const components: Array<{ id: string; category: string }> = await componentsRes.json();
    const byCategory = (category: string) => {
      const found = components.find((c) => c.category === category);
      expect(found, `seeded component with category ${category}`).toBeTruthy();
      return found!.id;
    };

    const created = await request.post(`${API_URL}/spell`, {
      headers: authHeaders(user),
      data: {
        user_id: user.userId,
        name: "E2E Testing Bolt",
        description: "A deterministic bolt of end-to-end energy.",
        type: "Attack",
        concentration: false,
        add_modifier: false,
        damage_dice_count: 1,
        damage_die_size: 8,
        component_ids: [byCategory("Forma"), byCategory("Scopus"), byCategory("Essentia")],
      },
    });
    expect(created.status(), await created.text()).toBe(201);
    const spell = await created.json();

    // The user's paginated spell list contains it.
    const list = await request.get(`${API_URL}/user/${user.userId}/spells`, {
      headers: authHeaders(user),
    });
    expect(list.status()).toBe(200);
    const paginated = await list.json();
    expect(
      (paginated.spells ?? []).some((s: { name: string }) => s.name === "E2E Testing Bolt"),
    ).toBe(true);

    // The arcanum spellbook page renders.
    await page.goto("/spellbook", { waitUntil: "domcontentloaded" });
    await expect(page.getByText("Arcanum").first()).toBeVisible({ timeout: 15_000 });

    // Delete round-trip.
    const deleted = await request.delete(`${API_URL}/spell/${spell.id}`, {
      headers: authHeaders(user),
    });
    expect(deleted.ok()).toBeTruthy();
  });
});
