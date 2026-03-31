import { Monster, MonsterAttack, MonsterAction } from "@/types/monster";
import { makeApiRequest } from "./base";

// Request body for creating a monster
export interface CreateMonsterRequest {
  description: string;
  challenge_rating: string;
  user_id: string; // UUID string
  /** When set, backend themes the stat block from seed/faradhaven_classes (must match class name). */
  faradhaven_class_name?: string;
}

// Request body for updating a monster (all fields optional for partial update)
export interface UpdateMonsterRequest {
  name?: string;
  size?: string;
  type?: string;
  alignment?: string;
  challenge_rating?: string;
  proficiency_bonus?: number;
  armor_class?: number;
  hit_points?: number;
  hit_dice?: string;
  speed?: string;
  strength?: number;
  dexterity?: number;
  constitution?: number;
  intelligence?: number;
  wisdom?: number;
  charisma?: number;
  saving_throws?: string[];
  skills?: string[];
  damage_immunities?: string[];
  damage_resistances?: string[];
  condition_immunities?: string[];
  senses?: string;
  languages?: string;
  notes?: string;
  image_url?: string;
  special_traits?: string[];
  legendary_actions?: string[];
  reactions?: string[];
  lair_actions?: string[];
  environments?: string[];
  source?: string;
  visual_description?: string;
  attacks?: MonsterAttack[]; // Use specific type
  actions?: MonsterAction[]; // Use specific type
}

// Create a new monster (LLM generation)
export async function createMonster(data: CreateMonsterRequest, token?: string): Promise<Monster> {
  return makeApiRequest<Monster>('POST', `/api/monsters`, data, token);
}

// Get all monsters for a user
export async function getMonstersByUser(userId: string, token?: string): Promise<Monster[]> {
  return makeApiRequest<Monster[]>('GET', `/api/user/${userId}/monsters`, null, token);
}

// Get a single monster by ID
export async function getMonsterById(monsterId: string, token?: string): Promise<Monster> {
  return makeApiRequest<Monster>('GET', `/api/monsters/${monsterId}`, null, token);
}

// Update an existing monster
export async function updateMonster(monsterId: string, data: UpdateMonsterRequest, token?: string): Promise<Monster> {
  return makeApiRequest<Monster>('PUT', `/api/monsters/${monsterId}`, data, token);
}

// Delete a monster
export async function deleteMonster(monsterId: string, token?: string): Promise<{ message: string }> {
  return makeApiRequest<{ message: string }>('DELETE', `/api/monsters/${monsterId}`, null, token);
}
