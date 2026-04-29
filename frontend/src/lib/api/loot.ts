import { getBaseUrl, handleResponse, apiFetch } from './base';
import type {
  LootSource,
  LootResult,
  LootRoomTheme,
  LootLocation,
  LootOptionsResponse,
  LootPreviewResponse,
  LootAssignmentPayload,
} from '@/types/game/api';

export type { LootSource, LootResult, LootRoomTheme, LootLocation };

export async function generateLootPreview(
  characterId: string,
  source: LootSource,
  roomTheme: LootRoomTheme,
  location: LootLocation | undefined,
  lootLevel: number,
  token: string
): Promise<LootPreviewResponse> {
  const base = getBaseUrl();
  const res = await apiFetch(`${base}/api/character/${characterId}/loot`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({
      character_id: characterId,
      source,
      room_theme: roomTheme,
      location,
      loot_level: lootLevel,
    }),
  });
  return handleResponse<LootPreviewResponse>(res, 'Failed to generate loot preview');
}

export async function confirmLootPickup(
  characterId: string,
  sessionId: string,
  assignments: LootAssignmentPayload[],
  token: string
): Promise<LootResult> {
  const base = getBaseUrl();
  const res = await apiFetch(`${base}/api/character/${characterId}/loot/confirm`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({
      session_id: sessionId,
      assignments,
    }),
  });
  return handleResponse<LootResult>(res, 'Failed to confirm loot pickup');
}

export async function getLootOptions(): Promise<LootOptionsResponse> {
  const base = getBaseUrl();
  const res = await apiFetch(`${base}/api/loot/options`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });
  return handleResponse<LootOptionsResponse>(res, 'Failed to load loot options');
}
