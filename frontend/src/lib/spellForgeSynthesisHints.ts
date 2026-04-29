import type { SpellSynthesis } from '@/types/game';

/** Keys used with `overrides` / `fieldOverrides` — manual edits block synthesis auto-fill for that bucket. */
export type SpellForgeManualField =
  | 'type'
  | 'range'
  | 'duration'
  | 'concentration'
  | 'damageDice'
  | 'damageType';

/** Shown next to the info icon: how Faradhaven derives suggestions (matches server synthesis rules). */
export const SPELL_SYNTHESIS_HINT_BODY =
  'Suggestions come from your component chain: Forma sets range (feet), Essentia maps to damage type, Actio sets spell type, Magnitudo feeds the Laws of Equivalency for dice (with spell level from chain length), and some shapes (Zone, Wall, Aura) or Bind set duration and concentration. Draft damage type and range sent to synthesis refine what the server returns.';

export function formatSynthesisTypeHint(s: SpellSynthesis | null): string | null {
  const t = s?.suggested_type;
  return t ? String(t) : null;
}

export function formatSynthesisRangeHint(s: SpellSynthesis | null): string | null {
  if (s?.suggested_range == null) return null;
  return `${s.suggested_range} ft`;
}

export function formatSynthesisDurationHint(s: SpellSynthesis | null): string | null {
  const d = s?.suggested_duration?.trim();
  return d ? d : null;
}

export function formatSynthesisDiceHint(s: SpellSynthesis | null): string | null {
  if (s?.suggested_damage_dice_count == null || s?.suggested_damage_die_size == null) return null;
  const dice = `${s.suggested_damage_dice_count}d${s.suggested_damage_die_size}`;
  const dt = s?.suggested_damage_type?.trim();
  return dt ? `${dice} (${dt})` : dice;
}

export function formatSynthesisDamageTypeHint(s: SpellSynthesis | null): string | null {
  const t = s?.suggested_damage_type?.trim();
  return t || null;
}

export function formatSynthesisConcentrationHint(s: SpellSynthesis | null): string | null {
  if (s == null) return null;
  return s.suggested_concentration ? 'Yes' : 'No';
}
