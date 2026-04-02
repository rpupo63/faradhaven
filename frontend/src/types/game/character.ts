import { AbilityId } from './schema';
import { ApiLevelFeature, ApiTrait, ApiRace, ApiComponent, ApiCharacterWeapon, ApiItem, ApiCharacterComponent, ApiClassResource, HarvestedAbilities } from './api';

export interface Character {
  id: string;
  name: string;
  race: string;
  race_id?: string;
  class: string;
  class_id?: string;
  archetype_id?: string;
  level: number;
  spellbook: string[];
  strength?: number;
  dexterity?: number;
  constitution?: number;
  intelligence?: number;
  wisdom?: number;
  charisma?: number;
  current_spell_points?: number;
  backstory?: string;
  notes?: string;
  money?: number;
  skill_proficiencies?: string[];
  saving_throw_proficiencies?: AbilityId[];
  inventory?: string[];
  image_url?: string;
  createdAt: number;
}

export interface Class {
  id: string;
  name: string;
  hit_die: number;
  primary_ability: string;
  created_at?: string;
  updated_at?: string;
}

export interface ClassLevel {
  id: string;
  class_id: string;
  level: number;
  hp_gain: number;
  proficiency_bonus: number;
  max_spell_points: number;
  ability_score_improvement?: number;
  cantrips_known?: number;
  spells_known?: number;
  spell_slots?: Record<string, number>;
  resource_pools?: Record<string, number>;
  extra_attack_count?: number;
  sneak_attack_dice?: number;
  rage_damage_bonus?: number;
  martial_arts_die?: number;
  unarmored_movement?: number;
  superiority_dice?: number;
  superiority_die?: number;
  bardic_inspiration?: number;
  created_at?: string;
  updated_at?: string;
}

export interface CharacterSheet {
  character: Character;
  class: Class;
  class_level: ClassLevel;
  ac: number;
  save_dc: number;
  max_spell_points: number;
  current_spell_points: number;
}

/** Normalized sheet for display + dice rolls (works with API or local data) */
export interface NormalizedCharacterSheet {
  character: {
    id: string;
    name: string;
    raceName: string;
    lineageName?: string;
    className: string;
    archetypeName?: string;
    level: number;
    spellbook: string[];
    strength: number;
    dexterity: number;
    constitution: number;
    intelligence: number;
    wisdom: number;
    charisma: number;
    current_spell_points: number;
    sanguine_mp?: number;
    sanguine_br?: number;
    sanguine_notoriety?: number;
    backstory?: string;
    notes?: string;
    inventory?: string[];
    image_url?: string;
    equipped_armor_id?: string;
    equipped_shield_id?: string;
    /** Set when the API returns an affiliated party for this character. */
    partyName?: string;
  };
  class: {
    name: string;
    primary_ability: string;
    proficiencies?: string;
    tools?: string[];
    skill_focus?: string[];
  };
  class_level: {
    proficiency_bonus: number;
    level_features?: ApiLevelFeature[];
    sneak_attack_dice?: number;
  };
  current_hp: number;
  max_hp: number;
  temp_hp: number;
  ac: number;
  save_dc: number;
  max_spell_points: number;
  current_spell_points: number;
  hit_dice_total: number;
  hit_dice_remaining: number;
  hit_die: number;
  money: number;
  bite_damage_dice?: number;
  modifiers: {
    strength: number;
    dexterity: number;
    constitution: number;
    intelligence: number;
    wisdom: number;
    charisma: number;
    primary: number;
    proficiency: number;
    initiative: number;
    spell_attack: number;
    melee_attack: number;
    ranged_attack: number;
    skills: Record<string, number>;
    saving_throws: Record<AbilityId, number>;
  };
  skill_proficiencies: string[];
  saving_throw_proficiencies: AbilityId[];
  race_traits?: ApiTrait[];
  lineage?: ApiRace;
  speed?: number;
  /** Always-known spell pool from current class + race definitions (not stored per character). */
  available_components?: ApiComponent[];
  /** Expendable inventory (foraged/extracted/etc.); spell pool is `available_components`. */
  components?: ApiCharacterComponent[];
  inventory_weapons?: ApiCharacterWeapon[];
  inventory_items?: ApiItem[];

  // --- Generic Class Resources ---
  class_resources?: ApiClassResource[];

  // Madness table (Mutagen / Lorewright)
  madness_table?: Record<number, string>;
  // Lorewright harvested abilities (stored on Character as JSONB)
  harvested_abilities?: HarvestedAbilities;
  // Lorewright: psychological damage from harvesting
  trauma?: number;
  // Trait use states: traitID → current uses remaining
  trait_use_states?: Record<string, number>;
  trait_max_uses?: Record<string, number>;
}
