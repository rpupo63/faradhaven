import type { CreatureSize, CreatureType, DamageType } from './base';
import { AbilityId } from './schema';
import { ApiClassResourceDefinition } from './class-resources';

export type SpellMechanicType = 'Attack' | 'Save' | 'Effect' | 'Healing' | 'Utility';

/** D&D-style ability codes for saving throws (matches Go models.SaveAttribute). */
export type SpellSaveAttribute = 'STR' | 'DEX' | 'CON' | 'INT' | 'WIS' | 'CHA';

/** Backend API Spell model (matches Go models.Spell) */
export interface ApiSpell {
  id: string;
  user_id: string;
  character_id?: string; // Character who prepared this spell (nil = user-level spell)
  name: string;
  description: string;
  level: number;
  components?: ApiComponent[];
  /** Legacy alias; API returns tier as `level`. Use `getSpellComponentCount()` for display/cost. */
  slot_level?: number;

  // Spell Mechanics
  type?: SpellMechanicType;
  /** Distance in feet; 0 = self-centered. */
  range?: number;
  duration?: string;
  concentration?: boolean;
  save_attr?: SpellSaveAttribute;
  /** Number of dice (e.g. 2 for 2d6); use with damage_die_size. */
  damage_dice_count?: number;
  /** Die faces: 4, 6, 8, 10, 12, 20, or 100. */
  damage_die_size?: number;
  damage_type?: string;
  suggested_damage_dice_count?: number;
  suggested_damage_die_size?: number;
  add_modifier?: boolean;
  checked?: boolean;

  // AI Review Opinions
  ai_description_opinion?: string;
  ai_damage_opinion?: string;
  ai_effect_opinion?: string;
  ai_overall_verdict?: string;
  ai_raw_output?: string;

  // AI Recommended Edits
  ai_recommended_name?: string;
  ai_recommended_description?: string;
  ai_recommended_type?: SpellMechanicType;
  /** Feet; aligns with [range]. */
  ai_recommended_range?: number;
  ai_recommended_duration?: string;
  ai_recommended_damage_dice_count?: number;
  ai_recommended_damage_die_size?: number;
  ai_recommended_damage_type?: string;
  ai_recommended_save_attr?: SpellSaveAttribute;

  created_at?: string;
  updated_at?: string;
}

// NEW: ApiBeast, ApiAttack, ApiBeastSkill interfaces
export interface ApiBeast {
  id: string;
  user_id: string;
  name: string;
  image_url?: string;
  size: string; // CreatureSize enum
  type: string; // CreatureType enum
  alignment: string;
  armor_class: number;
  hit_points: number;
  hit_dice: string;
  speed: string;
  strength: number;
  dexterity: number;
  constitution: number;
  intelligence: number;
  wisdom: number;
  charisma: number;
  challenge_rating: string;
  abilities: string[]; // pq.StringArray
  actions: string[]; // pq.StringArray
  legendary_actions: string[]; // pq.StringArray
  description: string;
  attacks?: ApiAttack[]; // Relationship
  skills?: ApiBeastSkill[]; // Relationship
  created_at?: string;
  updated_at?: string;
}

export interface ApiAttack {
  id: string;
  beast_id: string;
  name: string;
  attack_bonus: number;
  damage_type: string; // DamageType enum
  damage_dice: string;
  reach?: string;
  description?: string;
}

export interface ApiBeastSkill {
  id: string;
  beast_id: string;
  name: string;
  modifier: number;
}
// End NEW: ApiBeast, ApiAttack, ApiBeastSkill interfaces

/** Backend API types (match Go models) */
export interface ApiTraitOption {
  id: string;
  trait_id: string;
  name: string;
  description: string;
  ability_score_bonuses?: Record<string, number>;
}

export interface ApiTrait {
  id: string;
  race_id: string;
  name: string;
  description: string;
  level_req: number;
  action_type: string;
  uses_per_rest?: string;
  reset_condition?: string;
  range_value?: string;
  aoe?: string;
  save_ability?: string;
  options?: ApiTraitOption[];
}

export interface CastSpellResponse {
  current_spell_points: number;
  madness_cast_count?: number;
  current_resource_value?: number;
  resource_key?: string;
}

/** Speed Dial / blueprint slot (saved_spells table) */
export interface ApiSavedSpell {
  id: string;
  character_id: string;
  name: string;
  component_ids: string[];
  slot_index: number;
  total_cost?: number;
  damage_type?: string | null;
  description?: string | null;
  uses_remaining?: number | null;
  max_uses?: number | null;
  created_at?: string;
  updated_at?: string;
}

/** Store owner catalog rule (matches models.StoreOwnerCatalogRule JSON) */
export interface ApiStoreOwnerCatalogRule {
  id: string;
  store_owner_id: string;
  item_id?: string | null;
  weapon_id?: string | null;
  category?: string | null;
  allowed_rarities?: string[];
  created_at?: string;
  updated_at?: string;
}

/** Vendor NPC for shop (matches models.StoreOwner JSON) */
export interface ApiStoreOwner {
  id: string;
  name: string;
  location: string;
  personality?: string;
  willingness_to_purchase: number;
  exchange_rate: number;
  categories_obtained?: string[];
  catalog_rules?: ApiStoreOwnerCatalogRule[];
  created_at?: string;
  updated_at?: string;
}

export interface ApiRace {
  id: string;
  name: string;
  description?: string;
  creature_type?: string;
  size?: string;
  base_speed?: number;
  photo_url?: string;
  ability_score_bonuses?: Record<string, number>;
  languages?: string[];
  bonus_language_count?: number;
  traits?: ApiTrait[];
  /** Join-table spell pool (same as backend Race.Components). */
  components?: ApiComponent[];
  created_at?: string;
  updated_at?: string;
}

// NEW: ApiParty interface
export interface ApiParty {
  id: string;
  name: string;
  owner_id: string;
  created_at?: string;
  updated_at?: string;
  members?: ApiCharacter[]; // Optional, to be preloaded
  identified_beasts?: ApiBeast[]; // Optional, to be preloaded
}

export interface ApiRaceWithTraits extends ApiRace {
  traits?: ApiTrait[];
}

export type ComponentCategory =
  | 'Forma'
  | 'Scopus'
  | 'Essentia'
  | 'Actio'
  | 'Magnitudo'
  | 'Logica';

export interface ApiComponent {
  id: string;
  name: string;
  symbol: string; // 2-3 letter alchemical symbol for periodic table display
  /** RPG Awesome icon class suffix (e.g. "fire" for .ra-fire); empty = use symbol only */
  rpg_awesome_icon?: string;
  category: ComponentCategory;
  description: string;
  element?: string;
  tier: number; // Tier level of the component (e.g., 1 for basic, 2 for advanced)
  created_at?: string;
  updated_at?: string;
}

export interface SpellSynthesis {
  level: number;
  suggested_type: SpellMechanicType;
  suggested_range?: number;
  suggested_damage_type?: string;
  suggested_damage_dice_count?: number;
  suggested_damage_die_size?: number;
  suggested_duration?: string;
  suggested_concentration: boolean;
  validation_errors: string[];
}

export interface ApiArchetype {
  id: string;
  class_id: string;
  name: string;
  description: string;
  sort_order: number;
  created_at?: string;
  updated_at?: string;
}

export interface ApiEquipmentOption {
  id: string;
  choice_id: string;
  description: string;
  items: string[];
}

export interface ApiEquipmentChoice {
  id: string;
  class_id: string;
  instruction: string;
  sort_order: number;
  options: ApiEquipmentOption[];
}

export interface ApiClass {
  id: string;
  name: string;
  description?: string;
  hit_die: number;
  primary_ability: string;
  photo_url?: string;
  proficiencies?: string;
  skill_focus?: string[];
  skill_choice?: string[];
  skill_choice_count?: number;
  tools?: string[];
  saving_throws?: string[];
  starting_equipment?: string[];
  equipment_choices?: ApiEquipmentChoice[];
  archetype_level?: number; // Level at which players choose their archetype
  archetypes?: ApiArchetype[]; // Available archetypes for this class
  components?: ApiComponent[];
  levels?: ApiClassLevel[]; // Level 1 included in list view for feature previews
  created_at?: string;
  updated_at?: string;

  // Resource metadata
  resource_type?: string;
  resource_name?: string;
  resource_restore_type?: string;
  resource_definitions?: ApiClassResourceDefinition[]; // New property
}

export interface ApiClassWithLevels extends ApiClass {
  levels: ApiClassLevel[];
  madness_table?: Record<number, string>;
}

export interface ApiLevelFeature {
  id: string;
  class_level_id: string;
  archetype_id?: string; // nil = all archetypes get this feature
  name: string;
  description: string;
  sort_order: number;
  action_type?: string;       // "Action", "Bonus Action", "Reaction", or ""
  uses_per_rest?: string;     // "1", "Proficiency Bonus", or ""
  reset_condition?: string;   // "Long Rest", "Short Rest", or ""
  resource_costs?: Array<{ key: string; amount: number }>;
  created_at?: string;
  updated_at?: string;
}

export interface ApiClassLevel {
  id: string;
  class_id: string;
  level: number;
  hp_gain: number;
  proficiency_bonus: number;
  max_spell_points: number;
  ability_score_improvement?: number; // 0 = none, 2 = +2 to one or +1 to two (levels 4,8,12,16,19)
  cantrips_known?: number;
  spells_known?: number;
  level_features?: ApiLevelFeature[]; // structured features with name and description
  extra_attack_count?: number;
  sneak_attack_dice?: number;
  rage_damage_bonus?: number;
  martial_arts_die?: number;
  unarmored_movement?: number;
  superiority_dice?: number;
  superiority_die?: number;
  bardic_inspiration?: number;
  max_spell_level?: number; // Piston Brawler/Casters: spell level cap
  created_at?: string;
  updated_at?: string;

  // Faradhaven class resources are now in class_resources (via ClassLevelResource)
  resources?: Record<string, number>; // New property
}

export interface ApiCharacterComponent {
  character_id: string;
  component_id: string;
  count: number;
  component?: ApiComponent;
}

export interface PaginatedResponse<T> {
  total_count: number;
  page: number;
  limit: number;
  [key: string]: T[] | number; // This allows "characters": ApiCharacter[] or "spells": ApiSpell[]
}

export interface PaginatedCharactersResponse {
  characters: ApiCharacter[];
  total_count: number;
  page: number;
  limit: number;
}

export interface ApiCharacter {
  id: string;
  user_id: string;
  name: string;
  race_id: string;
  lineage_id?: string;
  class_id: string;
  archetype_id?: string; // nil until archetype is chosen
  party_id?: string; // NEW: Foreign key to Party
  level: number;
  spellbook: string[];
  strength: number;
  dexterity: number;
  constitution: number;
  intelligence: number;
  wisdom: number;
  charisma: number;
  current_spell_points: number;
  notoriety: number;
  sanguine_mp?: number;
  sanguine_br?: number;
  sanguine_notoriety?: number;
  backstory?: string; // Rich text HTML backstory
  notes?: string; // Miscellaneous notes
  skill_proficiencies?: string[];
  languages?: string[];
  saving_throw_proficiencies?: AbilityId[];
  inventory?: string[];
  components?: ApiCharacterComponent[];
  image_url?: string;
  trauma?: number;
  madness_cast_count?: number;
  current_madness_dc?: number;
  equipped_armor_id?: string;
  equipped_shield_id?: string;
  dice_theme?: string;
  dice_theme_color?: string;
  dice_font_color?: string;
  created_at?: string;
  updated_at?: string;
  race?: ApiRace;
  class?: ApiClass;
  archetype?: ApiArchetype;
  party?: ApiParty; // NEW: Relationship to Party
}

export interface ApiItem {
  id: string;
  name: string;
  description: string;
  category: string;
  rarity: string;
  cost: string;
  weight: string;
  effects: string;
  is_consumable: boolean;
  armor_type?: string;
  base_ac?: number;
  strength_requirement?: number;
  stealth_disadvantage?: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface ApiEffect {
  id: string;
  name: string;
  description: string;
  category: string;
  mechanics: string;
  created_at?: string;
  updated_at?: string;
}

export interface ApiWeaponDamage {
  id: string;
  weapon_id: string;
  damage_dice: string;
  damage_type: DamageType;
  damage_category: 'Base' | 'Bonus';
}

export interface ApiWeapon {
  id: string;
  name: string;
  description: string;
  category: string;
  rarity: string;
  range_type: string;
  cost: string;
  weight: string;
  attack_modifier: string;
  properties: string[];
  range_normal: number;
  range_long: number;
  versatile_damage_dice?: string;
  secondary_effect?: string;

  // Explicit Two-Handed Properties (if different from base and not just Versatile)
  two_handed_damage_dice?: string;
  two_handed_properties?: string[];

  // Transformation Properties
  transformed_to_weapon_id?: string; // UUIDs are strings
  transformation_description?: string;

  // For self-modifying transformations
  transformed_self_damage_dice?: string;
  transformed_self_properties?: string[];

  damages?: ApiWeaponDamage[];
  created_at?: string;
  updated_at?: string;
}

export interface ApiBonusDamageInfo {
  dice: string;
  damage_type: string;
}

export interface ApiWeaponModifier {
  modifier_type: string;
  is_active: boolean;
  bonus_damage?: ApiBonusDamageInfo[];
  metadata?: Record<string, unknown>;
}

export interface ApiCharacterWeapon {
  character_weapon_id: string;
  weapon: ApiWeapon;
  is_primary: boolean;
  is_equipped: boolean;
  custom_name?: string;
  active_modifiers?: ApiWeaponModifier[];

  // Frontend-managed state for weapon modes
  is_two_handed?: boolean;
  is_transformed?: boolean;
  transformed_weapon_data?: ApiWeapon;
}

export interface ApiCharacterSheet {
  character: ApiCharacter;
  class: ApiClass;
  class_level: ApiClassLevel;
  max_hp: number;
  current_hp: number;
  temp_hp: number;
  ac: number;
  save_dc: number;
  max_spell_points: number;
  current_spell_points: number;
  saving_throw_proficiencies?: string[];
  available_components?: ApiComponent[];
  hit_dice_total: number;
  hit_dice_remaining: number;
  hit_die: number;
  money: number;
  melee_attack_bonus: number;
  ranged_attack_bonus: number;
  race_traits?: ApiTrait[];
  lineage?: ApiRace; // Reusing ApiRace for lineage (same fields)
  inventory_weapons?: ApiCharacterWeapon[]; // Weapons with modifiers
  inventory_items?: ApiItem[];
  components?: ApiCharacterComponent[];

  // --- Generic Class Resources ---
  class_resources?: ApiClassResource[];

  // Madness table (Mutagen / Lorewright)
  madness_table?: Record<number, string>;
  // Lorewright harvested abilities (stored on Character as JSONB)
  harvested_abilities?: HarvestedAbilities;
  // Trait use states: traitID → current uses remaining
  trait_use_states?: Record<string, number>;
  trait_max_uses?: Record<string, number>;
}

export interface HarvestedAbilities {
  skills: string[];
  attacks: string[];
  recipes: string[];
}

/** Generic class resource from the API (replaces individual class-specific resource fields) */
export interface ApiClassResource {
  key: string;
  display_name: string;
  category: 'pool' | 'die_size' | 'limit' | 'slot_count' | 'modifier' | 'state';
  description?: string;
  value: number;
  current_value?: number;
  max_value?: number;
  is_trackable: boolean;
  display_order: number;
}

// === Minion / Construct Types ===

export interface ApiConstructAction {
  name: string;
  description: string;
  damage_dice?: string;
  damage_type?: string;
  range?: string;
  action_type: string; // "action", "bonus_action", "reaction"
}

export interface ApiConstructTemplate {
  key: string;
  name: string;
  description: string;
  base_ac: number;
  base_hp: number;
  base_speed: number;
  size: string; // "Small", "Medium", "Large", "Tiny"
  actions: ApiConstructAction[];
  required_components?: string[];
}

export interface ApiMinion {
  id: string;
  character_id: string;
  name: string;
  minion_type: string; // "construct", "echo", "renfield", "summon", "drone", "harvest"
  armor_class: number;
  max_hp: number;
  current_hp: number;
  speed: number;
  component_cost: number;
  template_key?: string;
  is_active: boolean;
  expires_at?: string;
  metadata?: Record<string, unknown>;
  created_at?: string;
  updated_at?: string;
}

export interface ApiHarvestMinion extends ApiMinion {
  ability_name: string;
  ability_description: string;
  ability_type: string; // "skill", "attack", "recipe", "action", "trait", "language"
  source_creature_type: string;
  slot_index?: number;
}

export interface CreateMinionRequest {
  minion_type: string;
  template_key?: string;
  custom_name?: string;
  component_id?: string; // For drone: which component to spend
}

export interface UpdateMinionHPRequest {
  delta: number;
}

export interface ApiCreationOptions {
  races: ApiRace[];
  classes: ApiClass[];
  points_max: number;
}

export interface CreateCharacterRequest {
  user_id: string;
  name: string;
  race_id: string;
  lineage_id?: string | null;
  class_id: string;
  level: number;
  spellbook: string[];
  strength: number;
  dexterity: number;
  constitution: number;
  intelligence: number;
  wisdom: number;
  charisma: number;
  current_spell_points: number;
  money: number;
  skill_proficiencies: string[];
  languages: string[];
  equipment_choices: string[];
  primary_weapon_id?: string;
}

// === Level-Up Preview Types ===

export interface ApiWeaponSelectionInfo {
  description: string;
  modifier_type: string;
  eligible_weapons: ApiWeapon[];
}

export interface ApiLevelUpPreview {
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
  weapon_selection_info?: ApiWeaponSelectionInfo;
}

// === Corpse Types ===

export interface ApiCorpse {
  id: string;
  map_id?: string;
  name: string;
  creature_type: CreatureType;
  creature_size: CreatureSize;
  challenge_rating: number;
  grid_x?: number;
  grid_y?: number;
  available_components: string[];
  component_yield: number;
  has_been_harvested: boolean;
  has_been_consumed: boolean;
  died_at: string;
  expires_at?: string;
  source_beast_id?: string;
  created_at: string;
}

export interface HarvestResult {
  components_yielded: number;
  component_ids: string[];
  component_names: string[];
  message: string;
}

// === Initiative Types ===

export interface ApiInitiativeEntry {
  token_id: string;
  order: number;
}

export interface ApiCharacterLink {
    id: string;
    source_character_id: string;
    target_character_id: string;
    link_type: string;
    is_active: boolean;
    notes?: string;
}

// === Loot Types ===
export type LootSource = 'common_enemy' | 'boss_enemy' | 'room';
export type LootTier = 'low' | 'medium' | 'high';

export interface LootDrop {
  kind: 'item' | 'weapon';
  name: string;
  rarity: string;
}

export interface LootResult {
  items: ApiItem[];
  weapons: ApiWeapon[];
  gold_earned: number;
  total_money: number;
  message: string;
  items_rolled: number;
  weapons_rolled: number;
  drops: LootDrop[];
}

