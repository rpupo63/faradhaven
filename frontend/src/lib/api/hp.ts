import { getBaseUrl, handleResponse, apiFetch } from './base';

// === HP Management Types ===

export interface HPResponse {
  current_hp: number;
  max_hp: number;
  temp_hp: number;
}

export interface UseHitDiceResponse {
  current_hp: number;
  max_hp: number;
  hp_healed: number;
  dice_used: number;
  dice_results: number[];
  hit_dice_remaining: number;
}

export interface RestResponse {
  current_hp: number;
  max_hp: number;
  temp_hp: number;
  current_spell_points: number;
  max_spell_points: number;
  hit_dice_remaining: number;
  hit_dice_total: number;
  // Class-specific resources
  current_stability?: number;
  max_stability?: number;
  current_blood_ichor?: number;
  max_blood_ichor?: number;
  madness_cast_count?: number;
}

/**
 * Updates character HP by delta (positive = heal, negative = damage)
 */
export async function updateHP(
  characterId: string,
  delta: number,
  source?: string,
  token?: string
): Promise<HPResponse> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/hp`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'PATCH',
    headers,
    body: JSON.stringify({ delta, source }),
  });
  return handleResponse<HPResponse>(res, 'Failed to update HP');
}

/**
 * Sets temporary HP (replaces existing, per 5e rules)
 */
export async function setTempHP(
  characterId: string,
  tempHP: number,
  token?: string
): Promise<HPResponse> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/temp-hp`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'PUT',
    headers,
    body: JSON.stringify({ temp_hp: tempHP }),
  });
  return handleResponse<HPResponse>(res, 'Failed to set temp HP');
}

/**
 * Uses hit dice during short rest. Rolls should be array of (d{hitDie} + conMod) results.
 */
export async function spendHitDice(
  characterId: string,
  rolls: number[],
  token?: string
): Promise<UseHitDiceResponse> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/hit-dice`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify({ rolls }),
  });
  return handleResponse<UseHitDiceResponse>(res, 'Failed to use hit dice');
}

/**
 * Performs a short rest: restores spell points (hit dice use is separate)
 */
export async function shortRest(
  characterId: string,
  token?: string
): Promise<RestResponse> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/rest/short`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { method: 'POST', headers });
  return handleResponse<RestResponse>(res, 'Failed to short rest');
}

/**
 * Performs a long rest: restores HP to max, restores half hit dice, restores spell points
 */
export async function longRest(
  characterId: string,
  token?: string
): Promise<RestResponse> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/rest/long`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { method: 'POST', headers });
  return handleResponse<RestResponse>(res, 'Failed to long rest');
}
