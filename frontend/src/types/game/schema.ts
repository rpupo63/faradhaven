export type AbilityId = 'strength' | 'dexterity' | 'constitution' | 'intelligence' | 'wisdom' | 'charisma';

export interface Dnd5eSkill {
  id: string;
  name: string;
  ability: AbilityId;
}

export const DND5E_SKILLS: Dnd5eSkill[] = [
  { id: 'acrobatics', name: 'Acrobatics', ability: 'dexterity' },
  { id: 'animal_handling', name: 'Animal Handling', ability: 'wisdom' },
  { id: 'arcana', name: 'Arcana', ability: 'intelligence' },
  { id: 'athletics', name: 'Athletics', ability: 'strength' },
  { id: 'deception', name: 'Deception', ability: 'charisma' },
  { id: 'history', name: 'History', ability: 'intelligence' },
  { id: 'insight', name: 'Insight', ability: 'wisdom' },
  { id: 'intimidation', name: 'Intimidation', ability: 'charisma' },
  { id: 'investigation', name: 'Investigation', ability: 'intelligence' },
  { id: 'medicine', name: 'Medicine', ability: 'wisdom' },
  { id: 'nature', name: 'Nature', ability: 'intelligence' },
  { id: 'perception', name: 'Perception', ability: 'wisdom' },
  { id: 'performance', name: 'Performance', ability: 'charisma' },
  { id: 'persuasion', name: 'Persuasion', ability: 'charisma' },
  { id: 'religion', name: 'Religion', ability: 'intelligence' },
  { id: 'sleight_of_hand', name: 'Sleight of Hand', ability: 'dexterity' },
  { id: 'stealth', name: 'Stealth', ability: 'dexterity' },
  { id: 'survival', name: 'Survival', ability: 'wisdom' },
];

export const DND5E_SAVING_THROWS: { id: AbilityId; name: string }[] = [
  { id: 'strength', name: 'Strength' },
  { id: 'dexterity', name: 'Dexterity' },
  { id: 'constitution', name: 'Constitution' },
  { id: 'intelligence', name: 'Intelligence' },
  { id: 'wisdom', name: 'Wisdom' },
  { id: 'charisma', name: 'Charisma' },
];

import { SpellElement } from './base';

export const SPELL_ELEMENTS: SpellElement[] = [
  { id: 'fire', name: 'Fire', description: 'Destructive flames that burn enemies', baseCost: 2 },
  { id: 'ice', name: 'Ice', description: 'Freezing cold that slows and damages', baseCost: 2 },
  { id: 'heal', name: 'Heal', description: 'Restorative energy that mends wounds', baseCost: 3 },
  { id: 'shield', name: 'Shield', description: 'Protective barrier against attacks', baseCost: 2 },
  { id: 'push', name: 'Force', description: 'Kinetic energy that pushes targets', baseCost: 1 },
  { id: 'lightning', name: 'Lightning', description: 'Electric strikes with chain potential', baseCost: 3 },
  { id: 'dark', name: 'Shadow', description: 'Dark energy that drains life', baseCost: 3 },
  { id: 'nature', name: 'Nature', description: 'Living magic from the earth', baseCost: 2 },
];
