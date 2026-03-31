import { getBaseUrl, handleResponse } from './base';
import { gainResource } from './resources';

export interface UseTraitAbilityResponse {
  trait_id: string;
  current_uses?: number;
  max_uses?: number;
}

export interface UseFeatureAbilityResponse {
  feature_id: string;
  message: string;
}

export async function callTraitAbility(
  characterId: string,
  traitId: string,
  token: string
): Promise<UseTraitAbilityResponse> {
  const base = getBaseUrl();
  const url = `${base}/api/characters/${characterId}/traits/${traitId}/use`;
  const res = await fetch(url, {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${token}` },
  });
  return handleResponse<UseTraitAbilityResponse>(res, 'Failed to use trait ability');
}

export async function restoreTraitAbility(
  characterId: string,
  traitId: string,
  token: string
): Promise<void> {
  await gainResource(characterId, 'trait_' + traitId, 1, token);
}

export async function callFeatureAbility(
  characterId: string,
  featureId: string,
  token: string
): Promise<UseFeatureAbilityResponse> {
  const base = getBaseUrl();
  const url = `${base}/api/characters/${characterId}/features/${featureId}/use`;
  const res = await fetch(url, {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${token}` },
  });
  return handleResponse<UseFeatureAbilityResponse>(res, 'Failed to use feature ability');
}
