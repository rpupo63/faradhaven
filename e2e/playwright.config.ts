import { defineConfig, devices } from "@playwright/test";

// Playwright starts both servers. The backend boots with GENERATE_MODELS=true,
// which migrates AND runs the full deterministic seed (classes, races, items,
// components, ...) before serving — that is why the backend webServer timeout
// is generous. Without GOOGLE_API_KEY the spell AI uses the built-in mock
// client and without AWS creds S3 uploads are disabled, so no secrets needed.
//
// The vite proxy targets :8080, so the backend port is fixed.
const BACKEND_PORT = 8080;
const FRONTEND_PORT = 4173;

export const API_URL = `http://localhost:${BACKEND_PORT}/api`;

const DATABASE_URL =
  process.env.E2E_DATABASE_URL ??
  "postgresql://postgres:postgres@localhost:5432/faradhaven_e2e?sslmode=disable";

export default defineConfig({
  testDir: "./tests",
  // Specs share one backend database, so run them sequentially.
  fullyParallel: false,
  workers: 1,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  reporter: process.env.CI ? [["list"], ["html", { open: "never" }]] : "list",
  use: {
    baseURL: `http://localhost:${FRONTEND_PORT}`,
    trace: "on-first-retry",
    screenshot: "only-on-failure",
  },
  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
  ],
  webServer: [
    {
      command: `cd ../backend && go build -o /tmp/faradhaven-e2e . && exec /tmp/faradhaven-e2e`,
      url: `http://localhost:${BACKEND_PORT}/healthcheck`,
      reuseExistingServer: !process.env.CI,
      timeout: 240_000, // migrate + full reseed before serving
      env: {
        PORT: String(BACKEND_PORT),
        DATABASE_URL,
        GENERATE_MODELS: "true",
        ACCEPTED_ORIGINS: `http://localhost:${FRONTEND_PORT}`,
      },
    },
    {
      command: `cd ../frontend && npm run dev -- --port ${FRONTEND_PORT} --strictPort`,
      url: `http://localhost:${FRONTEND_PORT}`,
      reuseExistingServer: !process.env.CI,
      timeout: 120_000,
    },
  ],
});
