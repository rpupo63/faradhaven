import { getBaseUrl, handleResponse, apiFetch } from './base';
import type {
  ApiCharacterSheet,
  ApiClass,
  ApiClassWithLevels,
  ApiCharacter,
  ApiRace,
  ApiRaceWithTraits,
  ApiWeapon,
  ApiComponent,
  ApiSpell,
  CreateCharacterRequest,
  ApiCreationOptions,
  ApiItem,
  ApiEffect,
  PaginatedCharactersResponse,
  CastSpellResponse,
  ApiSavedSpell,
  ApiStoreOwner
} from '@/types/game';

// NEW: PaginationOptions Interface
export interface PaginationOptions {
  page?: number;
  limit?: number;
  slot_level?: number | 'all';
}

// NEW: PaginatedSpellsResponse Interface
export interface PaginatedSpellsResponse {
  spells: ApiSpell[];
  total_count: number;
  page: number;
  limit: number;
}

/**
 * Fetches all races (for dropdowns and race list).
 */
export async function getRaces(token?: string): Promise<ApiRace[]> {
  const base = getBaseUrl();
  const url = `${base}/api/races`;
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<ApiRace[]>(res, 'Failed to load races');
}

/**
 * Fetches a race with traits (for D&D Beyond-style compendium detail view).
 */
export async function getRaceWithTraits(
  raceId: string,
  token?: string
): Promise<ApiRaceWithTraits> {
  const base = getBaseUrl();
  const url = `${base}/api/races/${raceId}`;
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<ApiRaceWithTraits>(res, 'Failed to load race');
}

/**
 * Fetches all classes (for dropdowns and class list).
 */
export async function getClasses(token?: string): Promise<ApiClass[]> {
  const base = getBaseUrl();
  const url = `${base}/api/classes`;
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<ApiClass[]>(res, 'Failed to load classes');
}

/**
 * Fetches a class with all levels 1-20 (for D&D Beyond-style compendium).
 */
export async function getClassWithLevels(
  classId: string,
  token?: string
): Promise<ApiClassWithLevels> {
  const base = getBaseUrl();
  const url = `${base}/api/classes/${classId}`;
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<ApiClassWithLevels>(res, 'Failed to load class');
}

/**
 * Fetches all items (for compendium/items list).
 */
export async function getItems(token?: string): Promise<ApiItem[]> {
  const base = getBaseUrl();
  const url = `${base}/api/items`;
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<ApiItem[]>(res, 'Failed to load items');
}

/**
 * Fetches all status effects and conditions.
 */
export async function getEffects(token?: string): Promise<ApiEffect[]> {
  const base = getBaseUrl();
  const url = `${base}/api/effects`;
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<ApiEffect[]>(res, 'Failed to load effects');
}

/**
 * Fetches all weapons (for compendium/weapons list).
 */
export async function getWeapons(token?: string): Promise<ApiWeapon[]> {
  const base = getBaseUrl();
  const url = `${base}/api/weapons`;
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<ApiWeapon[]>(res, 'Failed to load weapons');
}

/**
 * Fetches all store owners with catalog rules (shop vendors).
 */
export async function getStoreOwners(token?: string): Promise<ApiStoreOwner[]> {
  const base = getBaseUrl();
  const url = `${base}/api/store-owners`;
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<ApiStoreOwner[]>(res, 'Failed to load store owners');
}

/**
 * Fetches a single weapon by ID.
 */
export async function getWeaponById(weaponId: string, token?: string): Promise<ApiWeapon> {
  const base = getBaseUrl();
  const url = `${base}/api/weapons/${weaponId}`;
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<ApiWeapon>(res, `Failed to load weapon with ID ${weaponId}`);
}

/**
 * Fetches all spell components (for periodic table/arcanum view).
 */
export async function getComponents(): Promise<ApiComponent[]> {
  const base = getBaseUrl();
  const url = `${base}/api/components`;
  const res = await apiFetch(url);
  return handleResponse<ApiComponent[]>(res, 'Failed to load components');
}

/**
 * Fetches all spells from all users (aggregate spellbook).
 */
export async function getSpells(token?: string, options?: PaginationOptions): Promise<PaginatedSpellsResponse> {
  const base = getBaseUrl();
  const url = new URL(`${base}/api/spells`);
  if (options) {
    if (options.page) url.searchParams.append('page', options.page.toString());
    if (options.limit) url.searchParams.append('limit', options.limit.toString());
    if (options.slot_level && options.slot_level !== 'all') {
      url.searchParams.append('level', options.slot_level.toString());
    }
  }
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url.toString(), { headers });
  return handleResponse<PaginatedSpellsResponse>(res, 'Failed to load spells');
}

export interface CreateSpellRequest {
  user_id: string;
  character_id?: string; // Optional: character who prepared this spell
  name: string;
  description: string;
  component_ids: string[];
  slot_level: number;
  
  // Spell Mechanics
  type?: string;
  /** Feet; omit or use 0 for self-centered. */
  range?: number;
  duration?: string;
  concentration?: boolean;
  save_attr?: string;
  damage_dice_count?: number;
  damage_die_size?: number;
  damage_type?: string;
  add_modifier?: boolean;
}

/**
 * Creates a new spell for the user.
 */
export async function createSpell(
  request: CreateSpellRequest,
  token?: string
): Promise<ApiSpell> {
  const base = getBaseUrl();
  const url = `${base}/api/spell`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify(request),
  });
  return handleResponse<ApiSpell>(res, 'Failed to create spell');
}

/** Same shape as returned by GET /api/spell/:id/opinion and POST preview (SpellAIService). */
export interface SpellAIOpinion {
  description_opinion: string;
  damage_opinion: string;
  effect_opinion: string;
  overall_verdict: string;
  recommended_name?: string | null;
  recommended_description?: string | null;
  recommended_type?: string | null;
  recommended_range?: string | null;
  recommended_duration?: string | null;
  recommended_damage_dice_count?: number | null;
  recommended_damage_die_size?: number | null;
  recommended_damage_type?: string | null;
  recommended_save_attr?: string | null;
}

/**
 * Runs the same AI spell review as GM tools / spell opinion, using a draft payload (no spell row created).
 */
export async function previewSpellAIOpinion(
  request: CreateSpellRequest,
  token?: string
): Promise<SpellAIOpinion> {
  const base = getBaseUrl();
  const url = `${base}/api/spell/preview-ai-opinion`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify(request),
  });
  return handleResponse<SpellAIOpinion>(res, 'Failed to generate spell description');
}

/**
 * Calls the synthesis preview endpoint to get auto-derived spell properties.
 */
export async function synthesizeSpell(
  componentIds: string[],
  draft?: {
    damageType?: string;
    range?: number;
  },
  token?: string
): Promise<import('@/types/game').SpellSynthesis> {
  const base = getBaseUrl();
  const url = `${base}/api/spell/synthesize`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify({
      component_ids: componentIds,
      damage_type: draft?.damageType?.trim() ? draft.damageType.trim() : undefined,
      range: typeof draft?.range === 'number' ? draft.range : undefined,
    }),
  });
  return handleResponse(res, 'Failed to synthesize spell');
}

export interface UpdateSpellRequest extends CreateSpellRequest {
  id: string;
  checked?: boolean;
}

/**
 * Updates an existing spell.
 */
export async function updateSpell(
  request: UpdateSpellRequest,
  token?: string
): Promise<ApiSpell> {
  const base = getBaseUrl();
  const url = `${base}/api/spell/${request.id}`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'PUT',
    headers,
    body: JSON.stringify(request),
  });
  return handleResponse<ApiSpell>(res, 'Failed to update spell');
}

/**
 * GM: Fetches all spells that have not been reviewed yet.
 */
export async function getUncheckedSpells(token: string): Promise<ApiSpell[]> {
  const base = getBaseUrl();
  const url = `${base}/api/gm/spells/unchecked`;
  const res = await apiFetch(url, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return handleResponse<ApiSpell[]>(res, 'Failed to load unchecked spells');
}

/**
 * GM: Fetches all spells that have already been reviewed (checked).
 */
export async function getCheckedSpells(token: string): Promise<ApiSpell[]> {
  const base = getBaseUrl();
  const url = `${base}/api/gm/spells/checked`;
  const res = await apiFetch(url, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return handleResponse<ApiSpell[]>(res, 'Failed to load checked spells');
}

/**
 * GM: Update any spell (including marking it as checked).
 */
export async function gmUpdateSpell(
  spellId: string,
  updates: Partial<Omit<UpdateSpellRequest, 'id'>>,
  token: string
): Promise<ApiSpell> {
  const base = getBaseUrl();
  const url = `${base}/api/spell/${spellId}`;
  const res = await apiFetch(url, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
    body: JSON.stringify(updates),
  });
  return handleResponse<ApiSpell>(res, 'Failed to update spell');
}

/**
 * Fetches all options needed for character creation (races, classes, etc.)
 */
export async function getCreationOptions(token?: string): Promise<ApiCreationOptions> {
  const base = getBaseUrl();
  const url = `${base}/api/characters/options`;
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<ApiCreationOptions>(res, 'Failed to load options');
}

/**
 * Creates a new character.
 */
export async function createCharacter(
  request: CreateCharacterRequest,
  token?: string
): Promise<ApiCharacter> {
  const base = getBaseUrl();
  const url = `${base}/api/character`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify(request),
  });
  return handleResponse<ApiCharacter>(res, 'Failed to create character');
}

/**
 * Fetches all characters across all users (for DM token linking).
 */
export async function getAllCharacters(token?: string, options?: PaginationOptions): Promise<PaginatedCharactersResponse> {
  const base = getBaseUrl();
  const url = new URL(`${base}/api/characters`);
  if (options) {
    if (options.page) url.searchParams.append('page', options.page.toString());
    if (options.limit) url.searchParams.append('limit', options.limit.toString());
  }
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url.toString(), { headers });
  return handleResponse<PaginatedCharactersResponse>(res, 'Failed to load characters');
}

/**
 * Fetches all characters for a specific user.
 */
export async function getCharactersByUserId(
  userId: string,
  token?: string
): Promise<PaginatedCharactersResponse> {
  const base = getBaseUrl();
  const url = `${base}/api/user/${userId}/characters`;
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<PaginatedCharactersResponse>(res, 'Failed to load characters');
}

/**
 * Fetches the fully calculated character sheet (Class + ClassLevel joined).
 * Requires auth: pass token for protected backend.
 */
export async function getCharacterSheet(
  characterId: string,
  token?: string
): Promise<ApiCharacterSheet> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/sheet`;
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<ApiCharacterSheet>(res, 'Failed to load sheet');
}

/**
 * Resets CurrentSpellPoints to MaxSpellPoints (short/long rest).
 */
export async function restSpellPoints(
  characterId: string,
  token?: string
): Promise<{ current_spell_points: number; max_spell_points: number }> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/rest`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { method: 'POST', headers });
  return handleResponse<{ current_spell_points: number; max_spell_points: number }>(res, 'Failed to rest');
}

/**
 * Casts a spell, deducting spell points and updating class-specific resources.
 */
export async function castSpell(
  characterId: string,
  spellLevel: number,
  token?: string,
  spellId?: string
): Promise<CastSpellResponse> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/cast`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify({ spell_level: spellLevel, spell_id: spellId }),
  });
  return handleResponse(res, 'Failed to cast spell');
}

export interface SaveSpeedDialRequest {
  name: string;
  component_ids: string[];
}

/** Loads Speed Dial / blueprint slots for a character. */
export async function getSpeedDial(characterId: string, token?: string): Promise<ApiSavedSpell[]> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/speed-dial`;
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<ApiSavedSpell[]>(res, 'Failed to load speed dial');
}

/** Saves the current component sequence to a slot (Powder Mage Speed Dial, Piston blueprint, etc.). */
export async function saveSpeedDialSlot(
  characterId: string,
  slotIndex: number,
  body: SaveSpeedDialRequest,
  token?: string,
): Promise<ApiSavedSpell> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/speed-dial/${slotIndex}`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'PUT',
    headers,
    body: JSON.stringify(body),
  });
  return handleResponse<ApiSavedSpell>(res, 'Failed to save speed dial slot');
}

export async function clearSpeedDialSlot(
  characterId: string,
  slotIndex: number,
  token?: string,
): Promise<{ message: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/speed-dial/${slotIndex}`;
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { method: 'DELETE', headers });
  return handleResponse<{ message: string }>(res, 'Failed to clear speed dial slot');
}

/**
 * Manually consumes one unit of a specific component.
 */
export async function consumeComponent(
  characterId: string,
  componentId: string,
  token?: string
): Promise<{ message: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/component/${componentId}/consume`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'POST',
    headers,
  });
  return handleResponse<{ message: string }>(res, 'Failed to consume component');
}

/**
 * Manually increments one unit of a specific component.
 */
export async function gainComponent(
  characterId: string,
  componentId: string,
  token?: string
): Promise<{ message: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/component/${componentId}/gain`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'POST',
    headers,
  });
  return handleResponse<{ message: string }>(res, 'Failed to gain component');
}

/**
 * Fetches a character's backstory.
 */
export async function getBackstory(
  characterId: string,
  token?: string
): Promise<{ backstory: string; backstory_hex_color?: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/backstory`;
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, { headers });
  return handleResponse<{ backstory: string; backstory_hex_color?: string }>(res, 'Failed to load backstory');
}

/**
 * Updates a character's backstory.
 */
export async function updateBackstory(
  characterId: string,
  backstory: string,
  token?: string,
  backstory_hex_color?: string
): Promise<{ backstory: string; message: string; backstory_hex_color?: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/backstory`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'PUT',
    headers,
    body: JSON.stringify({ backstory, backstory_hex_color }),
  });
  return handleResponse<{ backstory: string; message: string; backstory_hex_color?: string }>(res, 'Failed to update backstory');
}

/**
 * Fetches all prepared spells for a character.
 */
export async function getCharacterSpells(
  characterId: string,
  token?: string,
  options?: PaginationOptions
): Promise<PaginatedSpellsResponse> {
  const base = getBaseUrl();
  const url = new URL(`${base}/api/character/${characterId}/spells`);
  if (options) {
    if (options.page) url.searchParams.append('page', options.page.toString());
    if (options.limit) url.searchParams.append('limit', options.limit.toString());
    if (options.slot_level && options.slot_level !== 'all') {
      url.searchParams.append('level', options.slot_level.toString());
    }
  }
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url.toString(), { headers });
  return handleResponse<PaginatedSpellsResponse>(res, 'Failed to load character spells');
}

/**
 * Deletes a character.
 */
export async function deleteCharacter(
  characterId: string,
  token?: string
): Promise<{ message: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'DELETE',
    headers,
  });
  return handleResponse<{ message: string }>(res, 'Failed to delete character');
}

/**
 * Updates a character's basic fields.
 */
export async function updateCharacter(
  characterId: string,
  updates: Partial<ApiCharacter> & { money?: number },
  token?: string
): Promise<ApiCharacter> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'PUT',
    headers,
    body: JSON.stringify(updates),
  });
  return handleResponse<ApiCharacter>(res, 'Failed to update character');
}

/**
 * Purchases an item or weapon for a character.
 */
export async function purchaseItem(
  characterId: string,
  itemId: string,
  itemType: 'item' | 'weapon',
  token?: string,
  storeOwnerId?: string | null
): Promise<{ message: string; money: number; cost_deducted: number }> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/purchase`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const body: Record<string, unknown> = {
    item_id: itemId,
    item_type: itemType,
  };
  if (storeOwnerId) {
    body.store_owner_id = storeOwnerId;
  }
  const res = await apiFetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify(body),
  });
  return handleResponse<{ message: string; money: number; cost_deducted: number }>(res, 'Failed to purchase item');
}

/**
 * Sells an owned item or weapon from a character inventory.
 */
export async function sellItem(
  characterId: string,
  payload: {
    item_id: string;
    item_type: 'item' | 'weapon';
    character_weapon_id?: string;
  },
  token?: string
): Promise<{ message: string; money: number; value_received: number }> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/sell`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const res = await apiFetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify(payload),
  });
  return handleResponse<{ message: string; money: number; value_received: number }>(
    res,
    'Failed to sell item'
  );
}

/**
 * Permanently discards an owned item or weapon without receiving money.
 */
export async function tossItem(
  characterId: string,
  payload: {
    item_id: string;
    item_type: 'item' | 'weapon';
    character_weapon_id?: string;
  },
  token?: string
): Promise<{ message: string; money: number }> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/toss`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const res = await apiFetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify(payload),
  });
  return handleResponse<{ message: string; money: number }>(res, 'Failed to toss item');
}

/**
 * Uploads a profile picture for a character.
 */
export async function uploadCharacterImage(
  characterId: string,
  file: File,
  token?: string
): Promise<{ image_url: string; message: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/image`;
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  
  const formData = new FormData();
  formData.append('image', file);

  const res = await apiFetch(url, {
    method: 'POST',
    headers, // Do NOT set Content-Type for FormData, browser sets it with boundary
    body: formData,
  });
  return handleResponse<{ image_url: string; message: string }>(res, 'Failed to upload image');
}

/**
 * Manually updates notoriety (delta).
 */
export async function updateNotoriety(
  characterId: string,
  delta: number,
  token?: string
): Promise<{ notoriety: number; message: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/notoriety`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify({ delta }),
  });
  return handleResponse<{ notoriety: number; message: string }>(res, 'Failed to update notoriety');
}

/**
 * Updates Sanguine MP and BR points.
 */
export async function updateSanguineNotoriety(
  characterId: string,
  { mpChange, brChange }: { mpChange: number; brChange: number },
  token?: string
): Promise<{ sanguine_mp: number; sanguine_br: number; sanguine_notoriety: number; message: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/sanguine-notoriety`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify({ mp_change: mpChange, br_change: brChange }),
  });
  return handleResponse(res, 'Failed to update Sanguine notoriety');
}

export interface UpdateEquipmentRequest {
  item_id: string;
  is_weapon: boolean;
  equip: boolean;
  slot: 'armor' | 'shield' | 'weapon';
}

/**
 * Equips or unequips an item or weapon for a character.
 */
export async function updateEquipment(
  characterId: string,
  request: UpdateEquipmentRequest,
  token?: string,
): Promise<ApiCharacterSheet> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/equipment`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const res = await apiFetch(url, {
    method: 'POST',
    headers,
    body: JSON.stringify(request),
  });
  return handleResponse<ApiCharacterSheet>(res, 'Failed to update equipment');
}

/**
 * Sets or clears a character's party affiliation.
 * Pass null for partyId to clear the party.
 */
export async function setCharacterParty(
  characterId: string,
  partyId: string | null,
  token?: string
): Promise<{ message: string; character_id: string; party_id: string | null }> {
  const base = getBaseUrl();
  const url = `${base}/api/characters/${characterId}/party`;
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url, {
    method: 'PUT',
    headers,
    body: JSON.stringify({ party_id: partyId }),
  });
  return handleResponse<{ message: string; character_id: string; party_id: string | null }>(
    res,
    'Failed to set character party'
  );
}

/**
 * Forages random spell components not in the class's own component pool.
 * Valid for Sanguinist, Ironwright, and Lorewright classes.
 */
export async function forageComponents(
  characterId: string,
  token: string
): Promise<{ components: Array<{ id: string; name: string }>; yield: number; message: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/character/${characterId}/forage-components`;
  const headers = { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` };
  const res = await apiFetch(url, { method: 'POST', headers });
  return handleResponse(res, 'Failed to forage components');
}

export type CharacterSpellbookScope = 'mine' | 'castable' | 'mine_or_castable' | 'all';

export interface GetCharacterSpellbookParams {
  scope: CharacterSpellbookScope;
  page?: number;
  /** Default on server is 20 for this endpoint */
  limit?: number;
  /** Spell level (points); omit or use only when filtering by tier */
  level?: number;
}

/**
 * Fetches a paginated spell list for a character spellbook by scope.
 */
export async function getCharacterSpellbook(
  characterId: string,
  token: string | undefined,
  params: GetCharacterSpellbookParams
): Promise<PaginatedSpellsResponse> {
  const base = getBaseUrl();
  const url = new URL(`${base}/api/character/${characterId}/spellbook`);
  url.searchParams.set('scope', params.scope);
  if (params.page != null) url.searchParams.set('page', String(params.page));
  if (params.limit != null) url.searchParams.set('limit', String(params.limit));
  if (params.level != null && params.level >= 1) {
    url.searchParams.set('level', String(params.level));
  }
  const headers: Record<string, string> = {};
  if (token) headers['Authorization'] = `Bearer ${token}`;
  const res = await apiFetch(url.toString(), { headers });
  return handleResponse<PaginatedSpellsResponse>(res, 'Failed to load character spellbook');
}

/**
 * Fetches an AI-generated review of a spell.
 */
export async function getSpellAIReview(
  spellId: string,
  token: string
): Promise<{ description_opinion: string; damage_opinion: string; effect_opinion: string; overall_verdict: string }> {
  const base = getBaseUrl();
  const url = `${base}/api/spell/${spellId}/opinion`;
  const res = await apiFetch(url, {
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
  });
  return handleResponse(res, 'Failed to get AI opinion');
}
