import { beforeEach, describe, expect, it, vi } from 'vitest';
import {
  abilityModsFromScores,
  buildSkillBonusMap,
  collectFeaturesThroughLevel,
  computeCharacterSheetFromApi,
  computeInitiativeModifier,
  ensureSanguinistBiteWeapon,
  getSavingThrowsFromClass,
  getSkillFocusFromClass,
  normalizeApiSheet,
  normalizeComputedSheet,
  readAbilityScoresForSheet,
  virtualSanguinistBiteWeapon,
} from './classLevelData';
import type { ApiCharacterSheet, ApiClass, ApiClassLevel, ApiLevelFeature } from '@/types/game';

describe('classLevelData shared helpers', () => {
  it('readAbilityScoresForSheet + abilityModsFromScores match manual mods', () => {
    const c = { strength: 14, dexterity: 12, constitution: 10, intelligence: 8, wisdom: 16, charisma: 9 };
    const scores = readAbilityScoresForSheet(c);
    expect(scores.str).toBe(14);
    const { strMod, dexMod, abilityMods } = abilityModsFromScores(scores);
    expect(strMod).toBe(2);
    expect(dexMod).toBe(1);
    expect(abilityMods.strength).toBe(2);
  });

  it('buildSkillBonusMap adds proficiency only for listed skills', () => {
    const { abilityMods } = abilityModsFromScores(
      readAbilityScoresForSheet({ dexterity: 14, intelligence: 10, wisdom: 10, charisma: 10, strength: 10, constitution: 10 })
    );
    const prof = 3;
    const map = buildSkillBonusMap(abilityMods, ['acrobatics'], prof);
    expect(map.acrobatics).toBe(abilityMods.dexterity + prof);
    expect(map.arcana).toBe(abilityMods.intelligence);
  });

  it('computeInitiativeModifier adds half prof for Jack of All Trades', () => {
    const features: ApiLevelFeature[] = [{ name: 'Jack of All Trades' } as ApiLevelFeature];
    expect(computeInitiativeModifier(2, 4, features)).toBe(4);
    expect(computeInitiativeModifier(2, 4, [])).toBe(2);
  });

  it('collectFeaturesThroughLevel respects level and archetype', () => {
    const levels: ApiClassLevel[] = [
      {
        id: 'l1',
        class_id: 'c',
        level: 1,
        hp_gain: 0,
        proficiency_bonus: 2,
        max_spell_points: 50,
        ability_score_improvement: 0,
        level_features: [{ name: 'A', archetype_id: undefined } as ApiLevelFeature],
      },
      {
        id: 'l2',
        class_id: 'c',
        level: 2,
        hp_gain: 0,
        proficiency_bonus: 2,
        max_spell_points: 52,
        ability_score_improvement: 0,
        level_features: [
          { name: 'B', archetype_id: 'arch-1' } as ApiLevelFeature,
          { name: 'C', archetype_id: undefined } as ApiLevelFeature,
        ],
      },
    ];
    const out = collectFeaturesThroughLevel(levels, 2, 'arch-1', undefined);
    const names = out.map((f) => f.name).sort();
    expect(names).toEqual(['A', 'B', 'C']);
    const outL1 = collectFeaturesThroughLevel(levels, 1, 'arch-1', undefined);
    expect(outL1.map((f) => f.name)).toEqual(['A']);
  });

  it('ensureSanguinistBiteWeapon injects once', () => {
    const bite = virtualSanguinistBiteWeapon();
    expect(bite.weapon.name).toBe('Bite');
    const first = ensureSanguinistBiteWeapon([], 'The Sanguinist');
    expect(first).toHaveLength(1);
    expect(ensureSanguinistBiteWeapon(first, 'The Sanguinist')).toHaveLength(1);
    expect(ensureSanguinistBiteWeapon([], 'The Mutagen')).toHaveLength(0);
  });
});

describe('normalizeApiSheet vs normalizeComputedSheet parity', () => {
  beforeEach(() => {
    vi.stubGlobal('localStorage', { getItem: () => null, setItem: vi.fn(), removeItem: vi.fn() });
  });

  it('produces matching skill modifiers for same stats and class level', () => {
    const apiClass: ApiClass = {
      id: 'class-1',
      name: 'The Mutagen',
      hit_die: 8,
      primary_ability: 'intelligence',
      proficiencies: '',
      skill_choice: [],
      skill_focus: [],
      tools: [],
      saving_throws: ['Intelligence', 'Wisdom'],
      levels: [],
    };
    const apiClassLevel: ApiClassLevel = {
      id: 'cl-1',
      class_id: 'class-1',
      level: 3,
      hp_gain: 4,
      proficiency_bonus: 2,
      max_spell_points: 54,
      ability_score_improvement: 0,
      level_features: [],
    };
    const api: ApiCharacterSheet = {
      character: {
        id: 'ch-1',
        name: 'Test',
        level: 3,
        spellbook: [],
        strength: 10,
        dexterity: 14,
        constitution: 10,
        intelligence: 16,
        wisdom: 10,
        charisma: 10,
        race: { name: 'Human' },
      } as ApiCharacterSheet['character'],
      class: apiClass,
      class_level: apiClassLevel,
      max_hp: 20,
      current_hp: 20,
      temp_hp: 0,
      ac: 12,
      save_dc: 15,
      max_spell_points: 54,
      current_spell_points: 54,
      hit_dice_total: 3,
      hit_dice_remaining: 3,
      hit_die: 8,
      money: 0,
      melee_attack_bonus: 2,
      ranged_attack_bonus: 4,
      saving_throw_proficiencies: ['intelligence', 'wisdom'],
    };

    const normalizedApi = normalizeApiSheet(api);

    const computed = computeCharacterSheetFromApi(
      {
        id: 'ch-1',
        name: 'Test',
        race: 'Human',
        class: 'The Mutagen',
        level: 3,
        spellbook: [],
        strength: 10,
        dexterity: 14,
        constitution: 10,
        intelligence: 16,
        wisdom: 10,
        charisma: 10,
        current_spell_points: 54,
        createdAt: 0,
        skill_proficiencies: [],
        saving_throw_proficiencies: [],
      },
      apiClass,
      apiClassLevel
    );
    const normalizedLocal = normalizeComputedSheet(
      computed,
      getSkillFocusFromClass(apiClass),
      getSavingThrowsFromClass(apiClass)
    );

    expect(normalizedApi.modifiers.skills).toEqual(normalizedLocal.modifiers.skills);
    expect(normalizedApi.modifiers.saving_throws).toEqual(normalizedLocal.modifiers.saving_throws);
    expect(normalizedApi.modifiers.proficiency).toBe(normalizedLocal.modifiers.proficiency);
    expect(normalizedApi.modifiers.spell_attack).toBe(normalizedLocal.modifiers.spell_attack);
  });
});
