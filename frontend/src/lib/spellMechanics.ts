/** Matches backend models.ValidateSpellDuration (backend/models/spell_constraints.go). Keep patterns in sync. */
export function isValidSpellDuration(s: string): boolean {
  const t = s.trim();
  if (!t) return false;
  if (/^\d+\s+(min|minute|minutes|hour|hours|day|days|week|weeks|month|months|year|years)$/i.test(t)) {
    return true;
  }
  if (/^\d+\s+rounds?$/i.test(t)) return true;
  return /^(concentration|instantaneous|instant|until dispelled|until triggered|special|permanent)$/i.test(t);
}

export const SPELL_TYPES = ['Attack', 'Save', 'Effect', 'Healing', 'Utility'] as const;
export type SpellMechanicType = (typeof SPELL_TYPES)[number];

export function formatSpellRangeFeet(feet: number | null | undefined): string {
  if (feet == null) return '';
  if (feet === 0) return 'Self';
  return `${feet} ft`;
}

/** Allowed die faces for spell damage (matches backend models.IsAllowedSpellDieSize). */
export const STANDARD_SPELL_DIE_SIZES = [4, 6, 8, 10, 12, 20, 100] as const;

export type StandardSpellDieSize = (typeof STANDARD_SPELL_DIE_SIZES)[number];

export function formatSpellDamageDice(count: number, size: number): string {
  return `${count}d${size}`;
}

/** True if both omitted, or both set with count >= 1 and allowed die size. */
export function isValidSpellDamageDicePair(
  count: number | undefined,
  size: number | undefined
): boolean {
  if (count === undefined && size === undefined) return true;
  if (count === undefined || size === undefined) return false;
  if (count < 1) return false;
  return (STANDARD_SPELL_DIE_SIZES as readonly number[]).includes(size);
}

export const SPELL_SAVE_ATTRIBUTES = ['STR', 'DEX', 'CON', 'INT', 'WIS', 'CHA'] as const;

/** Matches backend `Spell.concentration` default (true). Use when JSON omits the field. */
export function spellRequiresConcentration(concentration: boolean | undefined | null): boolean {
  return concentration ?? true;
}
