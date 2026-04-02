import { getBaseUrl, handleResponse, apiFetch } from './base';
import { ApiClassResource } from '@/types/game/api';

export async function spendResource(
  characterId: string,
  key: string,
  amount: number,
  token: string
): Promise<ApiClassResource> {
  const base = getBaseUrl();
  const url = `${base}/api/characters/${characterId}/resources/${key}/spend`;
  const res = await apiFetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify({ amount }),
  });
  return handleResponse<ApiClassResource>(res, 'Failed to spend resource');
}

export async function gainResource(
  characterId: string,
  key: string,
  amount: number,
  token: string
): Promise<ApiClassResource> {
  const base = getBaseUrl();
  const url = `${base}/api/characters/${characterId}/resources/${key}/gain`;
  const res = await apiFetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify({ amount }),
  });
  return handleResponse<ApiClassResource>(res, 'Failed to gain resource');
}
