import { getBaseUrl, handleResponse, apiFetch } from './base';
import type { LootSource, LootTier, LootResult } from '@/types/game/api';

export type { LootSource, LootTier, LootResult };

export async function generateLoot(
  characterId: string,
  source: LootSource,
  tier: LootTier,
  token: string
): Promise<LootResult> {
  const base = getBaseUrl();
  const res = await apiFetch(`${base}/api/character/${characterId}/loot`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ character_id: characterId, source, tier }),
  });
  return handleResponse<LootResult>(res, 'Failed to generate loot');
}
