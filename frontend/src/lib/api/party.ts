import { getBaseUrl, handleResponse } from './base';
import { ApiParty, ApiCharacter, ApiBeast } from '@/types/game/api';

/**
 * Creates a new party, optionally joining a character to it atomically.
 */
export async function createParty(
  name: string,
  token?: string,
  characterId?: string
): Promise<ApiParty> {
  const base = getBaseUrl();
  const url = `${base}/api/parties`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const body: Record<string, unknown> = { name };
  if (characterId) body.character_id = characterId;
  const res = await fetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify(body),
  });
  return handleResponse<ApiParty>(res, 'Failed to create party');
}

/**
 * Creates a new party and adds a character to it in a single request.
 */
export async function createPartyAndJoinWithCharacter(
  name: string,
  characterId: string,
  token?: string
): Promise<ApiParty> {
  return createParty(name, token, characterId);
}

/**
 * Fetches a party by its ID, including members and identified beasts.
 */
export async function getParty(
  partyId: string,
  token?: string
): Promise<ApiParty> {
  const base = getBaseUrl();
  const url = `${base}/api/parties/${partyId}`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await fetch(url, { headers });
  return handleResponse<ApiParty>(res, 'Failed to get party');
}

/**
 * Fetches all parties owned by a specific user.
 */
export async function getPartiesByUserId(
  userId: string,
  token?: string
): Promise<ApiParty[]> {
  const base = getBaseUrl();
  const url = `${base}/api/user/${userId}/parties`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await fetch(url, { headers });
  return handleResponse<ApiParty[]>(res, 'Failed to get user parties');
}

/**
 * Updates an existing party.
 */
export async function updateParty(
  partyId: string,
  name: string,
  token?: string
): Promise<ApiParty> {
  const base = getBaseUrl();
  const url = `${base}/api/parties/${partyId}`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await fetch(url, {
    method: 'PUT',
    headers,
    body: JSON.stringify({ name }),
  });
  return handleResponse<ApiParty>(res, 'Failed to update party');
}

/**
 * Deletes a party.
 */
export async function deleteParty(
  partyId: string,
  token?: string
): Promise<{ message: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/parties/${partyId}`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await fetch(url, {
    method: 'DELETE',
    headers,
  });
  return handleResponse<{ message: string }>(res, 'Failed to delete party');
}

/**
 * Adds a character to a party.
 */
export async function addCharacterToParty(
  partyId: string,
  characterId: string,
  token?: string
): Promise<{ message: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/parties/${partyId}/members`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await fetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify({ character_id: characterId }),
  });
  return handleResponse<{ message: string }>(res, 'Failed to add character to party');
}

/**
 * Removes a character from a party.
 */
export async function removeCharacterFromParty(
  partyId: string,
  characterId: string,
  token?: string
): Promise<{ message: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/parties/${partyId}/members/${characterId}`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await fetch(url, {
    method: 'DELETE',
    headers,
  });
  return handleResponse<{ message: string }>(res, 'Failed to remove character from party');
}

/**
 * Adds a beast to a party's identified beasts.
 */
export async function addBeastToParty(
  partyId: string,
  beastData: Omit<ApiBeast, 'id'>,
  token?: string
): Promise<ApiBeast> {
  const base = getBaseUrl();
  const url = `${base}/api/parties/${partyId}/beasts`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await fetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify(beastData),
  });
  return handleResponse<ApiBeast>(res, 'Failed to add beast to party');
}

/**
 * Adds a beast to a party's identified beasts.
 */
export async function addIdentifiedBeast(
  partyId: string,
  beastId: string,
  token?: string
): Promise<{ message: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/parties/${partyId}/identified-beasts`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await fetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify({ beast_id: beastId }),
  });
  return handleResponse<{ message: string }>(res, 'Failed to add identified beast');
}

/**
 * Removes a beast from a party's identified beasts.
 */
export async function removeIdentifiedBeast(
  partyId: string,
  beastId: string,
  token?: string
): Promise<{ message: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/parties/${partyId}/identified-beasts/${beastId}`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await fetch(url, {
    method: 'DELETE',
    headers,
  });
  return handleResponse<{ message: string }>(res, 'Failed to remove identified beast');
}