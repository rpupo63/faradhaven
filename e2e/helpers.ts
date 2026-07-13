import { expect, type APIRequestContext, type Page } from "@playwright/test";
import { API_URL } from "./playwright.config";

export interface TestUser {
  name: string;
  email: string;
  password: string;
  token: string;
  userId: string;
}

let counter = 0;

/** Register a fresh user through the real API and return its credentials. */
export async function registerUser(request: APIRequestContext): Promise<TestUser> {
  counter += 1;
  const email = `e2e-${Date.now()}-${counter}@example.com`;
  const password = "hunter22!";
  const res = await request.post(`${API_URL}/auth/register`, {
    data: { name: "E2E Keeper", email, password },
  });
  expect(res.status(), await res.text()).toBe(201);
  const body = await res.json();
  return { name: "E2E Keeper", email, password, token: body.token, userId: body.user_id };
}

/** Pre-authenticate the browser by seeding the token the frontend keeps in localStorage. */
export async function seedAuth(page: Page, user: TestUser): Promise<void> {
  await page.addInitScript(
    ({ token, userId }) => {
      localStorage.setItem("faradhaven-auth-token", token);
      localStorage.setItem("faradhaven-user-id", userId);
    },
    { token: user.token, userId: user.userId },
  );
}

export function authHeaders(user: TestUser): Record<string, string> {
  return { Authorization: `Bearer ${user.token}` };
}

/** Create a level-1 character from seeded class/race reference data. */
export async function createCharacter(
  request: APIRequestContext,
  user: TestUser,
  name: string,
): Promise<{ id: string; classId: string; raceId: string }> {
  const [classesRes, racesRes] = await Promise.all([
    request.get(`${API_URL}/classes`),
    request.get(`${API_URL}/races`),
  ]);
  const classes = await classesRes.json();
  const races = await racesRes.json();
  expect(classes.length).toBeGreaterThan(0);
  expect(races.length).toBeGreaterThan(0);

  const res = await request.post(`${API_URL}/character`, {
    headers: authHeaders(user),
    data: {
      user_id: user.userId,
      name,
      race_id: races[0].id,
      class_id: classes[0].id,
      level: 1,
      strength: 14,
      dexterity: 12,
      constitution: 13,
      intelligence: 10,
      wisdom: 11,
      charisma: 10,
      money: 100,
    },
  });
  expect(res.status(), await res.text()).toBe(201);
  const body = await res.json();
  return { id: body.id, classId: classes[0].id, raceId: races[0].id };
}
