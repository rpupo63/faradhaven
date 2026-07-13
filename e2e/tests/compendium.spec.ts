import { test, expect } from "@playwright/test";
import { registerUser, seedAuth } from "../helpers";

// Reference-data pages must render the deterministic seed data
// (GENERATE_MODELS=true reseeds on server start).

test.describe("compendium pages render seeded data", () => {
  test.beforeEach(async ({ page, request }) => {
    const user = await registerUser(request);
    await seedAuth(page, user);
  });

  test("class compendium lists seeded classes", async ({ page }) => {
    await page.goto("/game-rules/classes");
    await expect(page.getByRole("heading", { name: "Class Compendium" })).toBeVisible();
    // The seed system guarantees these classes exist.
    await expect(page.getByText("The Mutagen").first()).toBeVisible({ timeout: 15_000 });
  });

  test("race compendium lists seeded races", async ({ page }) => {
    await page.goto("/game-rules/races");
    await expect(page.getByRole("heading", { name: "Race Compendium" })).toBeVisible();
    // At least one race card renders (h3 titles inside the grid).
    await expect(page.locator("h3.font-tome-heading").first()).toBeVisible({ timeout: 15_000 });
  });

  test("bestiary tab renders", async ({ page }) => {
    await page.goto("/game-rules/bestiary");
    await expect(page.getByRole("heading", { name: "Bestiary" })).toBeVisible();
  });

  test("shop renders the seeded market", async ({ page }) => {
    // Portrait images can hang the load event (S3 disabled), so don't wait for it.
    await page.goto("/shop", { waitUntil: "domcontentloaded" });
    await expect(page.getByRole("heading", { name: "Faradhaven Market" })).toBeVisible({
      timeout: 15_000,
    });
  });
});
