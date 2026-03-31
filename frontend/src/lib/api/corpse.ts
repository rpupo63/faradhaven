import { getBaseUrl, handleResponse } from './base';
import type { ApiCorpse, CreatureSize, CreatureType, HarvestResult } from '@/types/game';

export async function getCorpses(
  mapId?: string,
  token?: string
): Promise<ApiCorpse[]> {
  const base = getBaseUrl();
  const params = mapId ? `?map_id=${encodeURIComponent(mapId)}` : '';
  const url = `${base}/api/corpses${params}`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await fetch(url, { headers });
  return handleResponse<ApiCorpse[]>(res, 'Failed to load corpses');
}

export async function createCorpse(
  data: {
    map_id?: string;
    name: string;
    creature_type: CreatureType;
    creature_size?: CreatureSize;
    challenge_rating?: number;
    grid_x?: number;
    grid_y?: number;
    available_components?: string[];
    component_yield?: number;
    source_beast_id?: string;
    expires_in_minutes?: number;
  },
  token: string
): Promise<ApiCorpse> {
  const base = getBaseUrl();
  const url = `${base}/api/corpses`;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`,
  };
  const res = await fetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify(data),
  });
  return handleResponse<ApiCorpse>(res, 'Failed to create corpse');
}

export async function harvestCorpse(
  corpseId: string,
  characterId: string,
  token: string
): Promise<HarvestResult> {
  const base = getBaseUrl();
  const url = `${base}/api/corpses/${corpseId}/harvest`;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`,
  };
  const res = await fetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify({ character_id: characterId }),
  });
  return handleResponse<HarvestResult>(res, 'Failed to harvest corpse');
}

export async function consumeCorpse(
  corpseId: string,
  characterId: string,
  token: string
): Promise<{ creature_type: string; challenge_rating: number; trauma_dc: number; message: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/corpses/${corpseId}/consume`;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`,
  };
  const res = await fetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify({ character_id: characterId }),
  });
  return handleResponse(res, 'Failed to consume corpse');
}

export async function deleteCorpse(
  corpseId: string,
  token: string
): Promise<{ status: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/corpses/${corpseId}`;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`,
  };
  const res = await fetch(url, { method: 'DELETE', headers });
  return handleResponse<{ status: string }>(res, 'Failed to delete corpse');
}
