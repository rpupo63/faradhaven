import { getBaseUrl, handleResponse, apiFetch } from './base';
import type {
  GameMap,
  CreateMapRequest,
  UpdateMapRequest,
  CreateTokenRequest,
  UpdateTokenRequest,
  MapToken,
  MapElement,
  MapElementType,
  CreateMapElementRequest,
  UpdateMapElementRequest
} from '@/types/map';

export async function createMap(
  data: CreateMapRequest,
  token: string
): Promise<GameMap> {
  const base = getBaseUrl();
  const url = `${base}/api/map`;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  };
  const res = await apiFetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify(data),
  });
  return handleResponse<GameMap>(res, 'Failed to create map');
}

export async function getMap(
  mapId: string,
  token?: string
): Promise<GameMap> {
  const base = getBaseUrl();
  const url = `${base}/api/map/${mapId}`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<GameMap>(res, 'Failed to load map');
}

export async function getMapByRoom(
  roomCode: string,
  token?: string
): Promise<GameMap> {
  const base = getBaseUrl();
  const url = `${base}/api/map/room/${roomCode}`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<GameMap>(res, 'Failed to load map');
}

export async function getUserMaps(
  userId: string,
  token: string
): Promise<GameMap[]> {
  const base = getBaseUrl();
  const url = `${base}/api/user/${userId}/maps`;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  };
  const res = await apiFetch(url, { headers });
  return handleResponse<GameMap[]>(res, 'Failed to load maps');
}

/**
 * Uploads a battle map background image to S3 (maps/backgrounds/) and sets background_url on the map.
 */
export async function uploadMapBackgroundImage(
  mapId: string,
  file: File,
  token: string
): Promise<{ background_url: string; message: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/map/${mapId}/background`;
  const headers: Record<string, string> = {
    Authorization: `Bearer ${token}`,
  };
  const formData = new FormData();
  formData.append('image', file);
  const res = await apiFetch(url, {
    method: 'POST',
    headers,
    body: formData,
  });
  return handleResponse<{ background_url: string; message: string }>(
    res,
    'Failed to upload map background'
  );
}

export async function updateMap(
  mapId: string,
  data: UpdateMapRequest,
  token: string
): Promise<GameMap> {
  const base = getBaseUrl();
  const url = `${base}/api/map/${mapId}`;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  };
  const res = await apiFetch(url, {
    method: 'PUT',
    headers,
    body: JSON.stringify(data),
  });
  return handleResponse<GameMap>(res, 'Failed to update map');
}

export async function deleteMap(
  mapId: string,
  token: string
): Promise<{ message: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/map/${mapId}`;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  };
  const res = await apiFetch(url, { method: 'DELETE', headers });
  return handleResponse<{ message: string }>(res, 'Failed to delete map');
}

export async function addToken(
  mapId: string,
  data: CreateTokenRequest,
  token: string
): Promise<MapToken> {
  const base = getBaseUrl();
  const url = `${base}/api/map/${mapId}/token`;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  };
  const res = await apiFetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify(data),
  });
  return handleResponse<MapToken>(res, 'Failed to add token');
}

export async function updateToken(
  mapId: string,
  tokenId: string,
  data: UpdateTokenRequest,
  token: string
): Promise<MapToken> {
  const base = getBaseUrl();
  const url = `${base}/api/map/${mapId}/token/${tokenId}`;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  };
  const res = await apiFetch(url, {
    method: 'PUT',
    headers,
    body: JSON.stringify(data),
  });
  return handleResponse<MapToken>(res, 'Failed to update token');
}

export async function deleteToken(
  mapId: string,
  tokenId: string,
  token: string
): Promise<{ message: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/map/${mapId}/token/${tokenId}`;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  };
  const res = await apiFetch(url, { method: 'DELETE', headers });
  return handleResponse<{ message: string }>(res, 'Failed to delete token');
}

export async function createMapElement(
  mapId: string,
  data: CreateMapElementRequest,
  token: string
): Promise<MapElement> {
  const base = getBaseUrl();
  const url = `${base}/api/map/${mapId}/elements`;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  };
  const res = await apiFetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify(data),
  });
  return handleResponse<MapElement>(res, 'Failed to create map element');
}

export async function updateMapElement(
  mapId: string,
  elementId: string,
  data: UpdateMapElementRequest,
  token: string
): Promise<MapElement> {
  const base = getBaseUrl();
  const url = `${base}/api/map/${mapId}/elements/${elementId}`;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  };
  const res = await apiFetch(url, {
    method: 'PUT',
    headers,
    body: JSON.stringify(data),
  });
  return handleResponse<MapElement>(res, 'Failed to update map element');
}

export async function deleteMapElement(
  mapId: string,
  elementId: string,
  token: string
): Promise<{ message: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/map/${mapId}/elements/${elementId}`;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  };
  const res = await apiFetch(url, { method: 'DELETE', headers });
  return handleResponse<{ message: string }>(res, 'Failed to delete map element');
}

export async function getInitiative(
  mapId: string,
  token?: string
): Promise<MapToken[]> {
  const base = getBaseUrl();
  const url = `${base}/api/map/${mapId}/initiative`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<MapToken[]>(res, 'Failed to load initiative order');
}

export async function setInitiative(
  mapId: string,
  entries: { token_id: string; order: number }[],
  token: string
): Promise<MapToken[]> {
  const base = getBaseUrl();
  const url = `${base}/api/map/${mapId}/initiative`;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`,
  };
  const res = await apiFetch(url, {
    method: 'PUT',
    headers,
    body: JSON.stringify({ entries }),
  });
  return handleResponse<MapToken[]>(res, 'Failed to set initiative order');
}

export async function clearInitiative(
  mapId: string,
  token: string
): Promise<{ status: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/map/${mapId}/initiative`;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`,
  };
  const res = await apiFetch(url, { method: 'DELETE', headers });
  return handleResponse<{ status: string }>(res, 'Failed to clear initiative');
}
