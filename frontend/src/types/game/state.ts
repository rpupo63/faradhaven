import { Character } from './character';
import { CraftedSpell, BestiaryEntry, CreatureSize, CreatureType, DamageType } from './base';

export interface GameState {
  characters: Character[];
  activeCharacterId: string | null;
  craftedSpells: CraftedSpell[];
  bestiary: BestiaryEntry[];
}

export const CREATURE_SIZES: CreatureSize[] = ['Tiny', 'Small', 'Medium', 'Large', 'Huge', 'Gargantuan'];
export const CREATURE_TYPES: CreatureType[] = ['Aberration', 'Beast', 'Celestial', 'Construct', 'Dragon', 'Elemental', 'Fey', 'Fiend', 'Giant', 'Humanoid', 'Monstrosity', 'Ooze', 'Plant', 'Undead'];
export const DAMAGE_TYPES: DamageType[] = ['Slashing', 'Piercing', 'Bludgeoning', 'Fire', 'Cold', 'Lightning', 'Thunder', 'Poison', 'Acid', 'Necrotic', 'Radiant', 'Force', 'Psychic'];
