export type ElementType = 'fire' | 'ice' | 'heal' | 'shield' | 'push' | 'lightning' | 'dark' | 'nature';
export type UUID = string;
export type CreatureSize = 'Tiny' | 'Small' | 'Medium' | 'Large' | 'Huge' | 'Gargantuan';
export type CreatureType = 'Aberration' | 'Beast' | 'Celestial' | 'Construct' | 'Dragon' | 'Elemental' | 'Fey' | 'Fiend' | 'Giant' | 'Humanoid' | 'Monstrosity' | 'Ooze' | 'Plant' | 'Undead';
export type DamageType = 'Slashing' | 'Piercing' | 'Bludgeoning' | 'Fire' | 'Cold' | 'Lightning' | 'Thunder' | 'Poison' | 'Acid' | 'Necrotic' | 'Radiant' | 'Force' | 'Psychic';

export interface SpellElement {
  id: ElementType;
  name: string;
  description: string;
  baseCost: number;
}

export interface Attack {
  id: string;
  name: string;
  attackBonus: number;
  damageType: DamageType;
  damageDice: string; // e.g., "2d6+4"
  reach?: string; // e.g., "5 ft." or "30/120 ft."
  description?: string;
}

export interface BestiaryEntry {
  id: string;
  name: string;
  imageUrl?: string;
  size: CreatureSize;
  type: CreatureType;
  alignment: string;
  armorClass: number;
  hitPoints: number;
  hitDice: string;
  speed: string;
  // Ability scores
  strength: number;
  dexterity: number;
  constitution: number;
  intelligence: number;
  wisdom: number;
  charisma: number;
  // Combat
  challengeRating: string;
  attacks: Attack[];
  // Traits
  abilities?: string;
  actions?: string;
  legendaryActions?: string;
  description: string;
  createdAt: number;
}

export interface CraftedSpell {
  id: string;
  name: string;
  description: string;
  elements: ElementType[];
  level: number;
  createdAt: number;
}
