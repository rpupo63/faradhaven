import { test, expect } from "@playwright/test";
import { registerUser, seedAuth } from "../helpers";

test.describe("authentication", () => {
  test("visiting a protected page while logged out shows the login page", async ({ page }) => {
    await page.goto("/characters");
    await expect(page.getByRole("button", { name: "Sign In" })).toBeVisible();
  });

  test("register through the UI reaches the app", async ({ page }) => {
    const email = `ui-signup-${Date.now()}@example.com`;

    await page.goto("/login");
    // Toggle from login to registration mode.
    await page.getByText(/create|register|sign up/i).first().click();
    await expect(page.getByRole("button", { name: "Create Account" })).toBeVisible();

    await page.locator("#name").fill("UI Signup Keeper");
    await page.locator("#email").fill(email);
    await page.locator("#password").fill("hunter22!");
    await page.locator("#confirmPassword").fill("hunter22!");
    await page.getByRole("button", { name: "Create Account" }).click();

    await expect(page).not.toHaveURL(/\/login/, { timeout: 15_000 });
  });

  test("login with valid credentials reaches the app", async ({ page, request }) => {
    const user = await registerUser(request);

    await page.goto("/login");
    await page.locator("#email").fill(user.email);
    await page.locator("#password").fill(user.password);
    await page.getByRole("button", { name: "Sign In" }).click();

    await expect(page).not.toHaveURL(/\/login/, { timeout: 15_000 });
  });

  test("a seeded token grants access without the login form", async ({ page, request }) => {
    const user = await registerUser(request);
    await seedAuth(page, user);

    await page.goto("/characters");
    await expect(page).toHaveURL(/\/characters/);
    await expect(page.getByRole("heading", { name: "Your Characters" })).toBeVisible();
  });
});
