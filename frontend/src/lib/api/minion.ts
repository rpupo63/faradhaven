import { getBaseUrl, handleResponse, apiFetch } from './base';
import type {
  ApiMinion,
  ApiConstructTemplate,
  CreateMinionRequest,
  ApiHarvestMinion, // Import the new type
} from '@/types/game';

/**
 * Fetches all minions for a character, optionally filtered by type.
 */
export async function getMinions(
  characterId: string,
  type?: string,
  token?: string
): Promise<(ApiMinion | ApiHarvestMinion)[]> { // Updated return type
  const base = getBaseUrl();
  const params = type ? `?type=${encodeURIComponent(type)}` : '';
  const url = `${base}/api/characters/${characterId}/minions${params}`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<(ApiMinion | ApiHarvestMinion)[]>(res, 'Failed to load minions');
}

/**
 * Fetches available construct templates (Sentry, Striker, Titan).
 */
export async function getMinionTemplates(
  characterId: string,
  token?: string
): Promise<Record<string, ApiConstructTemplate>> {
  const base = getBaseUrl();
  const url = `${base}/api/characters/${characterId}/minions/templates`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<Record<string, ApiConstructTemplate>>(res, 'Failed to load construct templates');
}

/**
 * Creates a new construct for a character.
 */
export async function createConstruct(
  characterId: string,
  request: CreateMinionRequest,
  token?: string
): Promise<ApiMinion> {
  const base = getBaseUrl();
  const url = `${base}/api/characters/${characterId}/minions`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify(request),
  });
  return handleResponse<ApiMinion>(res, 'Failed to create construct');
}

/**
 * Updates a minion's HP by a delta (positive = heal, negative = damage).
 */
export async function updateMinionHP(
  characterId: string,
  minionId: string,
  delta: number,
  token?: string
): Promise<ApiMinion> {
  const base = getBaseUrl();
  const url = `${base}/api/characters/${characterId}/minions/${minionId}/hp`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'PATCH',
    headers,
    body: JSON.stringify({ delta }),
  });
  return handleResponse<ApiMinion>(res, 'Failed to update minion HP');
}

/**
 * Stores a harvested ability (skill, attack, or recipe) as a "harvest" type minion.
 */
export async function storeHarvest(
  characterId: string,
  harvest: {
    ability_name: string;
    ability_description: string;
    ability_type: string; // "skill", "attack", "recipe", "action", "trait", "language"
    source_creature_type: string;
    slot_index?: number;
  },
  token?: string
): Promise<ApiMinion> {
  const request: CreateMinionRequest = {
    minion_type: 'harvest',
    custom_name: harvest.ability_name,
  };
  // The backend CreateMinionRequest supports these echo-specific fields
  const body = {
    ...request,
    ability_name: harvest.ability_name,
    ability_description: harvest.ability_description,
    ability_type: harvest.ability_type,
    source_creature_type: harvest.source_creature_type,
    slot_index: harvest.slot_index ?? 0,
  };
  const base = getBaseUrl();
  const url = `${base}/api/characters/${characterId}/minions`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify(body),
  });
  return handleResponse<ApiMinion>(res, 'Failed to store harvest');
}

/**
 * Deletes a minion. For constructs, optionally reclaim components.
 */
export async function deleteMinion(
  characterId: string,
  minionId: string,
  reclaim?: boolean,
  token?: string
): Promise<{ status: string; components_reclaimed?: number }> {
  const base = getBaseUrl();
  const params = reclaim ? '?reclaim=true' : '';
  const url = `${base}/api/characters/${characterId}/minions/${minionId}${params}`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'DELETE',
    headers,
  });
  return handleResponse<{ status: string; components_reclaimed?: number }>(res, 'Failed to delete minion');
}
