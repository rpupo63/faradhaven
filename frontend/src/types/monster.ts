// frontend/src/types/monster.ts

import type { CreatureSize, CreatureType, UUID } from "./game";

// Matches backend models.MonsterAttack
export interface MonsterAttack {
  id: UUID;
  name: string;
  description: string;
  is_legendary?: boolean;
  weapon_id?: UUID;
  // weapon?: Weapon; // If we need to embed weapon details
  created_at: string;
  updated_at: string;
}

// Matches backend models.MonsterAction
export interface MonsterAction {
  id: UUID;
  name: string;
  description: string;
  is_legendary?: boolean;
  created_at: string;
  updated_at: string;
}

// Matches backend models.Monster
export interface Monster {
  id: UUID;
  user_id: UUID;
  name: string;
  size: CreatureSize;
  type: CreatureType;
  alignment: string;
  challenge_rating: string;
  proficiency_bonus: number;
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
  saving_throws: string[];
  skills: string[];
  damage_immunities: string[];
  damage_resitances: string[];
  condition_immunities: string[];
  senses: string;
  languages: string;
  notes: string;
  image_url?: string;
  attacks: MonsterAttack[];
  actions: MonsterAction[];
  special_traits: string[];
  legendary_actions: string[];
  reactions: string[];
  lair_actions: string[];
  environments: string[];
  source: string;
  generation_mode?: string;
  generation_template?: string;
  generation_class_name?: string;
  generation_context?: {
    role?: string;
    environment?: string;
    temperament?: string;
    encounter_goal?: string;
    party_level?: number;
    template_id?: string;
    class_theme_intensity?: "light" | "strong";
  };
  visual_description: string;
  parsed_cr?: number; // Added from custom MarshalJSON in backend
  created_at: string;
  updated_at: string;
}