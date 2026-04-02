import type {
  ApiCharacterSheet,
  ApiClass,
  ApiClassLevel,
  ApiClassWithLevels,
  ApiLevelFeature,
  NormalizedCharacterSheet,
  AbilityId,
  ApiTrait,
  ApiWeapon,
  ApiWeaponDamage,
} from '@/types/game';
import { resolveSpellPoolComponents } from '@/lib/spellUtils';
import { DND5E_SKILLS, DND5E_SAVING_THROWS } from '@/types/game';
import type { Class, ClassLevel } from '@/types/game';

/** Map D&D skill name (from backend) to skill id */
const SKILL_NAME_TO_ID: Record<string, string> = Object.fromEntries(
  DND5E_SKILLS.map((s) => [s.name.toLowerCase(), s.id])
);

/** Convert backend skill names to D&D 5e skill ids */
export function skillNamesToIds(names: string[] | undefined): string[] {
  if (!names?.length) return [];
  return names
    .map((n) => SKILL_NAME_TO_ID[n.toLowerCase()] ?? n.toLowerCase().replace(/\s+/g, '_'))
    .filter(Boolean);
}

/** Get skill choice options (ids) from ApiClass */
export function getSkillChoiceFromClass(apiClass: ApiClass): string[] {
  return skillNamesToIds(apiClass.skill_choice);
}

/** Get fixed skill proficiencies (skill_focus) from ApiClass */
export function getSkillFocusFromClass(apiClass: ApiClass): string[] {
  return skillNamesToIds(apiClass.skill_focus);
}

/** Get saving throw proficiencies from ApiClass (backend uses ability names) */
export function getSavingThrowsFromClass(apiClass: ApiClass): AbilityId[] {
  const names = apiClass.saving_throws ?? [];
  const ids: AbilityId[] = [];
  for (const n of names) {
    const id = n.toLowerCase().trim() as AbilityId;
    if (['strength', 'dexterity', 'constitution', 'intelligence', 'wisdom', 'charisma'].includes(id)) {
      ids.push(id);
    }
  }
  return ids;
}

export function abilityMod(score: number): number {
  return Math.floor((score - 10) / 2);
}

// Helper to safely get ability score value by string name
function getAbilityScoreValue(character: {
  strength?: number;
  dexterity?: number;
  constitution?: number;
  intelligence?: number;
  wisdom?: number;
  charisma?: number;
}, ability: string): number {
  switch (ability.toLowerCase()) {
    case 'strength': return character.strength ?? 10;
    case 'dexterity': return character.dexterity ?? 10;
    case 'constitution': return character.constitution ?? 10;
    case 'intelligence': return character.intelligence ?? 10;
    case 'wisdom': return character.wisdom ?? 10;
    case 'charisma': return character.charisma ?? 10;
    default: return 10;
  }
}

function primaryAbilityMod(
  character: {
    strength?: number;
    dexterity?: number;
    constitution?: number;
    intelligence?: number;
    wisdom?: number;
    charisma?: number;
  },
  primaryAbility: string
): number {
  const s = (v: number) => abilityMod(v);
  return s(getAbilityScoreValue(character, primaryAbility));
}

export interface CharacterForSheet {
  id: string;
  name: string;
  race: string;
  class: string;
  level: number;
  spellbook: string[];
  archetype_id?: string;
  strength?: number;
  dexterity?: number;
  constitution?: number;
  intelligence?: number;
  wisdom?: number;
  charisma?: number;
  current_spell_points?: number;
  backstory?: string;
  notes?: string;
  skill_proficiencies?: string[];
  saving_throw_proficiencies?: AbilityId[];
  inventory?: string[];
  image_url?: string;
  createdAt: number;
}

/** Same shape as API CharacterSheet for UI consumption */
export interface CharacterSheetComputed {
  character: CharacterForSheet;
  class: Class & { proficiencies?: string; tools?: string[]; skill_focus?: string[] };
  class_level: ClassLevel & { level_features?: ApiLevelFeature[] };
  ac: number;
  save_dc: number;
  max_spell_points: number;
  current_spell_points: number;
  max_hp: number;
}

/**
 * Compute character sheet from API class data (for local characters with class_id).
 * TotalHP = BaseHP + (AvgHitDie * (Level - 1)) + (ConMod * Level); BaseHP = HitDie.
 */
export function computeCharacterSheetFromApi(
  character: CharacterForSheet,
  apiClass: ApiClass,
  apiClassLevel: ApiClassLevel
): CharacterSheetComputed {
  const conMod = abilityMod(character.constitution ?? 10);
  const dexMod = abilityMod(character.dexterity ?? 10);
  const chaMod = abilityMod(character.charisma ?? 10);
  const primaryMod = primaryAbilityMod(character, apiClass.primary_ability);

  // Collect all features up to current level
  const allFeatures: ApiLevelFeature[] = [];
  if (apiClass.levels) {
    const sortedLevels = [...apiClass.levels].sort((a, b) => a.level - b.level);
    for (const level of sortedLevels) {
      if (level.level <= character.level) {
        if (level.level_features) {
          for (const feature of level.level_features) {
            // Include if no archetype requirement, or if it matches character's archetype
            const charArchetypeId = character.archetype_id;
            if (!feature.archetype_id || feature.archetype_id === charArchetypeId) {
              allFeatures.push(feature);
            }
          }
        }
      }
    }
  } else if (apiClassLevel.level_features) {
    allFeatures.push(...apiClassLevel.level_features);
  }

  // Calculate AC including Unarmored Defense logic
  let ac = 8 + apiClassLevel.proficiency_bonus + dexMod;
  const hasUnarmoredDefense = allFeatures.some(f => f.name === 'Unarmored Defense');
  // Local characters don't have inventory tracking for shields yet in this function
  if (hasUnarmoredDefense) {
    if (apiClass.name === 'The Sanguinist') {
      ac = 10 + dexMod + chaMod;
    }
  }

  const saveDC = 8 + apiClassLevel.proficiency_bonus + primaryMod;
  const maxSP = apiClassLevel.max_spell_points;
  const currentSP = Math.min(character.current_spell_points ?? maxSP, maxSP);

  // Calculate maxHP based on the logic that was previously used for totalHP
  const avgHitDie = Math.floor((apiClass.hit_die + 1) / 2);
  const baseHP = apiClass.hit_die;
  const maxHP =
    baseHP +
    avgHitDie * (Math.max(1, character.level) - 1) +
    conMod * Math.max(1, character.level);

  return {
    character: { ...character },
    class: {
      id: apiClass.id,
      name: apiClass.name,
      hit_die: apiClass.hit_die,
      primary_ability: apiClass.primary_ability,
      proficiencies: apiClass.proficiencies,
      tools: apiClass.tools,
      skill_focus: apiClass.skill_focus,
    },
    class_level: {
      id: apiClassLevel.id,
      class_id: apiClassLevel.class_id,
      level: apiClassLevel.level,
      hp_gain: apiClassLevel.hp_gain,
      proficiency_bonus: apiClassLevel.proficiency_bonus,
      max_spell_points: apiClassLevel.max_spell_points,
      ability_score_improvement: apiClassLevel.ability_score_improvement,
      level_features: allFeatures,
    },
    ac,
    save_dc: saveDC,
    max_spell_points: maxSP,
    current_spell_points: currentSP,
    max_hp: maxHP,
  };
}

/** Normalize API sheet to NormalizedCharacterSheet */
export function normalizeApiSheet(api: ApiCharacterSheet): NormalizedCharacterSheet {
  const c = api.character;
  const cls = api.class;
  const cl = api.class_level;
  const str = getAbilityScoreValue(c, 'strength');
  const dex = getAbilityScoreValue(c, 'dexterity');
  const con = getAbilityScoreValue(c, 'constitution');
  const int = getAbilityScoreValue(c, 'intelligence');
  const wis = getAbilityScoreValue(c, 'wisdom');
  const cha = getAbilityScoreValue(c, 'charisma');
  const strMod = abilityMod(str);
  const dexMod = abilityMod(dex);
  const conMod = abilityMod(con);
  const intMod = abilityMod(int);
  const wisMod = abilityMod(wis);
  const chaMod = abilityMod(cha);
  const primaryMod = primaryAbilityMod(
    { strength: str, dexterity: dex, constitution: con, intelligence: int, wisdom: wis, charisma: cha },
    cls.primary_ability
  );
  const prof = cl.proficiency_bonus;
  const skillProfs = [...(c.skill_proficiencies ?? [])];
  
  // Kalashtar daily skill logic
  const kalashtarSkill = localStorage.getItem(`kalashtar_skill_${c.id}`);
  if (kalashtarSkill && (c.race?.name === 'Kalashtar') && !skillProfs.includes(kalashtarSkill)) {
    skillProfs.push(kalashtarSkill);
  }

  const saveProfs = api.saving_throw_proficiencies ?? c.saving_throw_proficiencies ?? [];

  const abilityMods: Record<AbilityId, number> = {
    strength: strMod,
    dexterity: dexMod,
    constitution: conMod,
    intelligence: intMod,
    wisdom: wisMod,
    charisma: chaMod,
  };

  const skills: Record<string, number> = {};
  for (const skill of DND5E_SKILLS) {
    const base = abilityMods[skill.ability];
    skills[skill.id] = base + (skillProfs.includes(skill.id) ? prof : 0);
  }

  const saving_throws: Record<AbilityId, number> = {} as Record<AbilityId, number>;
  for (const save of DND5E_SAVING_THROWS) {
    const base = abilityMods[save.id];
    saving_throws[save.id] = base + (saveProfs.includes(save.id) ? prof : 0);
  }

  // Combine race traits from sheet (backend now combines them)
  const combinedTraits = (api.race_traits ?? (c.race as { traits?: ApiTrait[] })?.traits ?? [])
    .filter(trait => (trait.level_req || 1) <= c.level);
  const speed = (c.race as { base_speed?: number })?.base_speed;

  // Collect all features up to current level, filtering by archetype
  const allFeatures: ApiLevelFeature[] = [];
  if (cls.levels) {
    const sortedLevels = [...cls.levels].sort((a, b) => a.level - b.level);
    for (const level of sortedLevels) {
      if (level.level <= c.level) {
        if (level.level_features) {
          for (const feature of level.level_features) {
            // Include if no archetype requirement, or if it matches character's archetype
            if (!feature.archetype_id || feature.archetype_id === c.archetype_id) {
              allFeatures.push(feature);
            }
          }
        }
      }
    }
  } else if (cl.level_features) {
    // Fallback to just current level features if levels list is missing
    allFeatures.push(...cl.level_features);
  }

  // AC is fully computed by the backend (armor, unarmored defense, shields all factored in)
  const finalAc = api.ac;

  // Calculate initiative (base Dex mod, + half proficiency if Jack of All Trades)
  let initiative = dexMod;
  if (allFeatures.some(f => f.name === 'Jack of All Trades')) {
    initiative += Math.floor(prof / 2);
  }

  return {
    character: {
      id: c.id,
      name: c.name,
      raceName: c.race?.name ?? 'Unknown',
      lineageName: api.lineage?.name,
      className: cls.name,
      archetypeName: c.archetype?.name,
      level: c.level,
      spellbook: c.spellbook ?? [],
      strength: str,
      dexterity: dex,
      constitution: con,
      intelligence: int,
      wisdom: wis,
      charisma: cha,
      current_spell_points: api.current_spell_points,
      sanguine_mp: c.sanguine_mp,
      sanguine_br: c.sanguine_br,
      sanguine_notoriety: c.sanguine_notoriety,
      backstory: c.backstory,
      notes: c.notes,
      inventory: c.inventory,
      image_url: c.image_url,
      equipped_armor_id: c.equipped_armor_id,
      equipped_shield_id: c.equipped_shield_id,
      partyName: c.party?.name?.trim() || undefined,
    },
    class_level: {
      proficiency_bonus: prof,
      level_features: allFeatures,
      sneak_attack_dice: cl.sneak_attack_dice,
    },
    current_hp: api.current_hp,
    max_hp: api.max_hp,
    temp_hp: api.temp_hp ?? 0,
    hit_dice_total: api.hit_dice_total ?? c.level,
    hit_dice_remaining: api.hit_dice_remaining ?? c.level,
    hit_die: api.hit_die ?? cls.hit_die,
    money: api.money ?? 0,
    bite_damage_dice: api.class_resources?.find(r => r.key === 'bite_damage_dice')?.value,
    ac: finalAc,
    save_dc: api.save_dc,
    max_spell_points: api.max_spell_points,
    current_spell_points: api.current_spell_points,
    skill_proficiencies: skillProfs,
    saving_throw_proficiencies: saveProfs as AbilityId[],
    race_traits: combinedTraits,
    lineage: api.lineage,
    speed,
    available_components: resolveSpellPoolComponents(api),
    inventory_weapons: (function() {
      const weapons = [...(api.inventory_weapons || [])];
      // Inject virtual Bite weapon for Sanguinists if not already present
      if (cls.name === 'The Sanguinist' && !weapons.some(w => w.weapon.name === 'Bite')) {
        weapons.push({
          character_weapon_id: 'virtual-bite',
          is_primary: false,
          is_equipped: true, // Should be true since it's a natural weapon
          weapon: {
            id: 'bite-id',
            name: 'Bite',
            description: 'A vampiric bite that extracts blood ichor.',
            category: 'Natural Melee',
            rarity: 'Common',
            range_type: 'Melee',
            attack_modifier: 'Charisma',
            properties: ['Finesse'],
            range_normal: 5,
            damages: [{ damage_dice: '2d8', damage_type: 'Necrotic', damage_category: 'Base' }] as ApiWeaponDamage[]
          } as ApiWeapon
        });
      }
      return weapons;
    })(),
    inventory_items: api.inventory_items,
    class: {
      name: cls.name,
      primary_ability: cls.primary_ability,
      proficiencies: cls.proficiencies,
      tools: cls.tools,
      skill_focus: cls.skill_focus,
    },
    modifiers: {
      strength: strMod,
      dexterity: dexMod,
      constitution: conMod,
      intelligence: intMod,
      wisdom: wisMod,
      charisma: chaMod,
      primary: primaryMod,
      proficiency: prof,
      initiative,
      spell_attack: prof + primaryMod,
      melee_attack: api.melee_attack_bonus ?? prof + strMod,
      ranged_attack: api.ranged_attack_bonus ?? prof + dexMod,
      skills,
      saving_throws,
    },

    // --- Generic Class Resources ---
    class_resources: api.class_resources,

    // Madness table and harvested abilities (Lorewright)
    madness_table: api.madness_table,
    harvested_abilities: api.harvested_abilities,
    trauma: c.trauma,
    trait_use_states: api.trait_use_states,
    trait_max_uses: api.trait_max_uses,
  };
}

/** Normalize local computed sheet to NormalizedCharacterSheet */
export function normalizeComputedSheet(
  computed: CharacterSheetComputed,
  skillDefaults?: string[],
  saveDefaults?: AbilityId[]
): NormalizedCharacterSheet {
  const c = computed.character;
  const cls = computed.class;
  const cl = computed.class_level;
  const str = getAbilityScoreValue(c, 'strength');
  const dex = getAbilityScoreValue(c, 'dexterity');
  const con = getAbilityScoreValue(c, 'constitution');
  const int = getAbilityScoreValue(c, 'intelligence');
  const wis = getAbilityScoreValue(c, 'wisdom');
  const cha = getAbilityScoreValue(c, 'charisma');
  const strMod = abilityMod(str);
  const dexMod = abilityMod(dex);
  const conMod = abilityMod(con);
  const intMod = abilityMod(int);
  const wisMod = abilityMod(wis);
  const chaMod = abilityMod(cha);
  const primaryMod = primaryAbilityMod(c, cls.primary_ability);
  const prof = cl.proficiency_bonus;
  const skillProfs =
    (c.skill_proficiencies?.length ?? 0) > 0
      ? c.skill_proficiencies ?? []
      : skillDefaults ?? [];
  const saveProfs =
    (c.saving_throw_proficiencies?.length ?? 0) > 0
      ? c.saving_throw_proficiencies ?? []
      : saveDefaults ?? [];

  const abilityMods: Record<AbilityId, number> = {
    strength: strMod,
    dexterity: dexMod,
    constitution: conMod,
    intelligence: intMod,
    wisdom: wisMod,
    charisma: chaMod,
  };

  const skills: Record<string, number> = {};
  for (const skill of DND5E_SKILLS) {
    const base = abilityMods[skill.ability];
    skills[skill.id] = base + (skillProfs.includes(skill.id) ? prof : 0);
  }

  const saving_throws: Record<AbilityId, number> = {} as Record<AbilityId, number>;
  for (const save of DND5E_SAVING_THROWS) {
    const base = abilityMods[save.id];
    saving_throws[save.id] = base + (saveProfs.includes(save.id) ? prof : 0);
  }

  return {
    character: {
      id: c.id,
      name: c.name,
      raceName: c.race,
      className: cls.name,
      level: c.level,
      spellbook: c.spellbook ?? [],
      strength: str,
      dexterity: dex,
      constitution: con,
      intelligence: int,
      wisdom: wis,
      charisma: cha,
      current_spell_points: computed.current_spell_points,
      backstory: c.backstory,
      notes: c.notes,
      inventory: c.inventory,
      image_url: c.image_url,
    },
    class: {
      name: cls.name,
      primary_ability: cls.primary_ability,
      proficiencies: cls.proficiencies,
      tools: cls.tools,
      skill_focus: cls.skill_focus,
    },
    class_level: { 
      proficiency_bonus: prof,
      level_features: cl.level_features,
    },
    current_hp: computed.max_hp, // Local characters don't track HP yet
    max_hp: computed.max_hp,
    temp_hp: 0,
    hit_dice_total: c.level,
    hit_dice_remaining: c.level,
    hit_die: cls.hit_die,
    money: 0, // Local characters don't track money yet
    ac: computed.ac,
    save_dc: computed.save_dc,
    max_spell_points: computed.max_spell_points,
    current_spell_points: computed.current_spell_points,
    skill_proficiencies: skillProfs,
    saving_throw_proficiencies: saveProfs,
    inventory_weapons: (function() {
      // Inject virtual Bite weapon for Sanguinists
      if (cls.name === 'The Sanguinist') {
        return [{
          character_weapon_id: 'virtual-bite',
          is_primary: false,
          is_equipped: true,
          weapon: {
            id: 'bite-id',
            name: 'Bite',
            description: 'A vampiric bite that extracts blood ichor.',
            category: 'Natural Melee',
            rarity: 'Common',
            range_type: 'Melee',
            attack_modifier: 'Charisma',
            properties: ['Finesse'],
            range_normal: 5,
            damages: [{ damage_dice: '2d8', damage_type: 'Necrotic', damage_category: 'Base' }] as ApiWeaponDamage[]
          } as ApiWeapon
        }];
      }
      return [];
    })(),
    modifiers: {
      strength: strMod,
      dexterity: dexMod,
      constitution: conMod,
      intelligence: intMod,
      wisdom: wisMod,
      charisma: chaMod,
      primary: primaryMod,
      proficiency: prof,
      initiative: dexMod,
      spell_attack: prof + primaryMod,
      melee_attack: prof + strMod,
      ranged_attack: prof + dexMod,
      skills,
      saving_throws,
    },
  };
}