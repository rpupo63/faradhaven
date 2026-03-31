import type { ApiCharacterSheet, CastSpellResponse } from '@/types/game';

export const RESOURCE_KEY_LABELS: Record<string, string> = {
  max_blood_ichor: 'Ichor',
  max_stability: 'Stab',
  spell_points: 'SP',
  shadow_points: 'SP',
};

/**
 * Builds a class-appropriate toast message after a successful cast.
 */
export function buildCastToast(
  sheet: ApiCharacterSheet | undefined,
  spellName: string,
  data: CastSpellResponse,
): string {
  const className = sheet?.class?.name;

  if (data.resource_key && data.current_resource_value !== undefined) {
    const label = RESOURCE_KEY_LABELS[data.resource_key] || data.resource_key;
    return `Cast ${spellName}! (${data.current_resource_value} ${label} remaining)`;
  }

  if (className === 'The Mutagen' && data.madness_cast_count !== undefined) {
    const profBonus = sheet?.class_level?.proficiency_bonus || 2;
    const conMod = Math.floor(((sheet?.character?.constitution ?? 10) - 10) / 2);
    const threshold = profBonus + conMod;
    if (data.madness_cast_count > threshold) {
      return `Cast ${spellName}! Madness casts: ${data.madness_cast_count}/${threshold} -- Madness save required!`;
    }
    return `Cast ${spellName}! (Casts: ${data.madness_cast_count}/${threshold})`;
  }

  if (className === 'The Ironwright' || className === 'The Lorewright') {
    return `Cast ${spellName}! Components consumed.`;
  }

  if (className === 'The Powder Mage') {
    return `Cast ${spellName}!`;
  }

  return `Cast ${spellName}!`;
}
