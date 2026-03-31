import { getBaseUrl, handleResponse } from './base';
import type { ApiEffect } from '@/types/game';

export interface RollTableResponse {
  effect: ApiEffect;
  roll: number;
}

export interface MutagenCastResponse {
  madness_cast_count: number;
  current_dc: number;
}

export interface ActiveEffect {
  id: string;
  effect_id: string;
  name: string;
  description: string;
  mechanics: string;
  duration: string;
  source: string;
  category: string;
}

/**
 * Rolls on a mechanics table (e.g. mutagen_feral, lorewright_madness)
 */
export async function rollTable(
  characterId: string,
  tableName: 'mutagen_feral' | 'lorewright_madness',
  token?: string
): Promise<RollTableResponse> {
  const base = getBaseUrl();
  const url = `${base}/api/characters/${characterId}/mechanics/roll-table`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  
  const res = await fetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify({ table_name: tableName }),
  });
  return handleResponse<RollTableResponse>(res, 'Failed to roll table');
}

/**
 * Increments mutagen cast count and returns new count/DC
 */
export async function castMutagen(
  characterId: string,
  token?: string
): Promise<MutagenCastResponse> {
  const base = getBaseUrl();
  const url = `${base}/api/characters/${characterId}/mechanics/mutagen-cast`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  
  const res = await fetch(url, {
    method: 'POST',
    headers,
  });
  return handleResponse<MutagenCastResponse>(res, 'Failed to cast mutation');
}

/**
 * Gets all active effects for a character
 */
export async function getActiveEffects(
  characterId: string,
  token?: string
): Promise<ActiveEffect[]> {
  const base = getBaseUrl();
  const url = `${base}/api/characters/${characterId}/mechanics/effects`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  
  const res = await fetch(url, { headers });
  return handleResponse<ActiveEffect[]>(res, 'Failed to fetch active effects');
}

/**
 * Removes an active effect
 */
export async function removeEffect(
  characterId: string,
  effectId: string,
  token?: string
): Promise<void> {
  const base = getBaseUrl();
  const url = `${base}/api/characters/${characterId}/mechanics/effects/${effectId}`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  
  const res = await fetch(url, { method: 'DELETE', headers });
  return handleResponse<void>(res, 'Failed to remove effect');
}
