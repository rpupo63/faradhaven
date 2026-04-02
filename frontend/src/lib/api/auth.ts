import { getBaseUrl, handleResponse, apiFetch } from './base';

// === Auth Types ===
export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user_id: string;
}

export interface ApiUser {
  id: string;
  name: string;
  email: string;
  active_character_id?: string | null;
  dice_theme: string;
  dice_theme_color: string;
  dice_font_color: string;
  created_at: string;
  updated_at: string;
}

export interface DicePrefs {
  dice_theme?: string;
  dice_theme_color?: string;
  dice_font_color?: string;
}

/**
 * Update user dice appearance preferences
 */
export async function updateUserDicePrefs(
  userId: string,
  prefs: DicePrefs,
  token: string
): Promise<ApiUser> {
  const base = getBaseUrl();
  const res = await apiFetch(`${base}/api/user/${userId}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(prefs),
  });
  return handleResponse<ApiUser>(res, 'Failed to update dice preferences');
}

/**
 * Login with email and password
 */
export async function login(request: LoginRequest): Promise<AuthResponse> {
  const base = getBaseUrl();
  const res = await apiFetch(`${base}/api/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(request),
  });
  return handleResponse<AuthResponse>(res, 'Login failed');
}

/**
 * Register a new user
 */
export async function register(request: RegisterRequest): Promise<AuthResponse> {
  const base = getBaseUrl();
  const res = await apiFetch(`${base}/api/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(request),
  });
  return handleResponse<AuthResponse>(res, 'Registration failed');
}

/**
 * Get user by ID
 */
export async function getUserById(userId: string, token: string): Promise<ApiUser> {
  const base = getBaseUrl();
  const res = await apiFetch(`${base}/api/user/${userId}`, {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
  });
  return handleResponse<ApiUser>(res, 'Failed to get user');
}

/**
 * Set the user's active character for shop purchases
 */
export async function setActiveCharacter(
  userId: string,
  characterId: string | null,
  token: string
): Promise<{ message: string; active_character_id: string | null }> {
  const base = getBaseUrl();
  const res = await apiFetch(`${base}/api/user/${userId}/active-character`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ character_id: characterId }),
  });
  return handleResponse<{ message: string; active_character_id: string | null }>(
    res,
    'Failed to set active character'
  );
}

/**
 * Logout the user
 */
export async function logoutApi(): Promise<void> {
  const base = getBaseUrl();
  await apiFetch(`${base}/api/auth/logout`, {
    method: 'POST',
  }).catch(() => {
    // Ignore errors on logout
  });
}

