import type { ApiCharacterSheet, ApiSpell, ApiSavedSpell } from '@/types/game';
import { spellMatchesAnySavedSpeedDial } from '@/lib/powderMageSpellMatch';

/** True if the character can spend a Speed Dial–style cast (trackable slot uses, or unlocked non-trackable slots). */
export function hasSpeedDialResourceRemaining(sheet: ApiCharacterSheet | undefined): boolean {
  if (!sheet?.class_resources) return false;
  const sd = sheet.class_resources.find((r) => r.key === 'speed_dial_slots');
  if (!sd || sd.value <= 0) return false;
  if (sd.is_trackable) {
    const cur = sd.current_value ?? sd.value;
    return cur > 0;
  }
  return true;
}

/** Spell matches a saved Speed Dial snapshot and Speed Dial can be used (class + resources). */
export function spellQualifiesForSpeedDialCast(
  spell: ApiSpell,
  speedDialSpells: ApiSavedSpell[] | undefined,
  sheet: ApiCharacterSheet | undefined,
): boolean {
  if (sheet?.class?.name !== 'The Powder Mage') return false;
  if (!speedDialSpells?.length) return false;
  if (!hasSpeedDialResourceRemaining(sheet)) return false;
  return spellMatchesAnySavedSpeedDial(spell, speedDialSpells);
}
