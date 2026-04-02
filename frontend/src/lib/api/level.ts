import { getBaseUrl, handleResponse, apiFetch } from './base';
import type {
  ApiCharacter,
  ApiClassLevel,
  ApiArchetype,
  ApiWeapon
} from '@/types/game';

// === Level-Up/Level-Down Types ===

export interface LevelUpRequest {
  skill_selections?: string[];
  asi_allocation?: Record<string, number>; // { strength: 1, dexterity: 1 }
  spells_learned?: string[];
  hp_roll_result?: number; // nil = use average, otherwise the rolled value
  archetype_id?: string; // required when reaching archetype level
  primary_weapon_id?: string; // selected primary weapon for signature items (e.g., Piston Core)
  mp_change?: number;
  br_change?: number;
}

export interface LevelUpResponse {
  character: ApiCharacter;
  new_level: number;
  class_level: ApiClassLevel;
  history_id?: string;
}

export interface WeaponSelectionInfo {
  description: string;
  modifier_type: string;
  eligible_weapons: ApiWeapon[];
}

export interface LevelUpPreview {
  current_level: number;
  next_level: number;
  class_level: ApiClassLevel;
  asi_points_available: number;
  new_spells_allowed: number;
  hit_die: number;
  con_mod: number;
  hp_gain_average: number;
  current_max_hp: number;
  requires_archetype_choice: boolean;
  available_archetypes?: ApiArchetype[];
  requires_weapon_selection: boolean;
  weapon_selection_info?: WeaponSelectionInfo;
}

export interface LevelUpHistoryEntry {
  id: string;
  character_id: string;
  user_id: string;
  level: number;
  skill_selections: string[] | null;
  asi_allocation: Record<string, number> | null;
  spells_learned: string[] | null;
  features_gained: string[] | null;
  created_at: string;
}

/**
 * Gets a preview of what the next level-up will provide
 */
export async function getLevelUpPreview(
  characterId: string,
  token?: string
): Promise<LevelUpPreview> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/level-up/preview`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<LevelUpPreview>(res, 'Failed to get preview');
}

/**
 * Levels up a character with the given choices
 */
export async function levelUpCharacter(
  characterId: string,
  request: LevelUpRequest,
  token?: string
): Promise<LevelUpResponse> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/level-up`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify(request),
  });
  return handleResponse<LevelUpResponse>(res, 'Level-up failed');
}

/**
 * Levels down a character, reverting to previous state
 */
export async function levelDownCharacter(
  characterId: string,
  token?: string
): Promise<LevelUpResponse> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/level-down`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { method: 'POST', headers });
  return handleResponse<LevelUpResponse>(res, 'Level-down failed');
}

/**
 * Gets the level-up history for a character
 */
export async function getLevelHistory(
  characterId: string,
  token?: string
): Promise<LevelUpHistoryEntry[]> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/level-history`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<LevelUpHistoryEntry[]>(res, 'Failed to get history');
}
