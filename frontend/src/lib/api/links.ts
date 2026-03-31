import { getBaseUrl, handleResponse } from './base';
import { ApiCharacterLink } from '@/types/game';

export async function getCharacterLinks(characterId: string, token?: string): Promise<ApiCharacterLink[]> {
    const base = getBaseUrl();
    const url = `${base}/api/character/${characterId}/links`;
    const headers: Record<string, string> = { 'Content-Type': 'application/json' };
    if (token) headers['Authorization'] = `Bearer ${token}`;
    const res = await fetch(url, { headers });
    return handleResponse<ApiCharacterLink[]>(res, 'Failed to get links');
}
