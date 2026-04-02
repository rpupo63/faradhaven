import type { ApiSpell, ApiComponent, ApiSavedSpell } from '@/types/game';
import { spellChainHasLogica } from '@/lib/spellLogicPhases';

function multisetSignature(ids: string[]): string {
  return [...ids].sort().join('\0');
}

function spellComponentIds(spell: ApiSpell): string[] {
  return spell.components?.map((c) => c.id) ?? [];
}

/**
 * Finds a prepared spell that matches the forge sequence:
 * - No Logica in the sequence: multiset equality of component IDs.
 * - Logica present: exact ordered ID match (same length, same ID at each index).
 */
export function findMatchingPreparedSpell(
  sequence: ApiComponent[],
  prepared: ApiSpell[],
): ApiSpell | null {
  if (sequence.length === 0) return null;

  const seqIds = sequence.map((c) => c.id);
  const useOrdered = spellChainHasLogica(sequence);

  for (const spell of prepared) {
    const spellIds = spellComponentIds(spell);
    if (spellIds.length !== seqIds.length) continue;

    if (useOrdered) {
      if (spellIds.every((id, i) => id === seqIds[i])) return spell;
    } else {
      if (multisetSignature(spellIds) === multisetSignature(seqIds)) return spell;
    }
  }
  return null;
}

/** True if a grimoire spell matches a Speed Dial snapshot (same rules as forge matching). */
export function spellMatchesSavedSpeedDial(spell: ApiSpell, saved: ApiSavedSpell): boolean {
  const seq = spell.components?.map((c) => c.id) ?? [];
  const ids = saved.component_ids;
  if (seq.length !== ids.length) return false;
  const comps = spell.components ?? [];
  const useOrdered = spellChainHasLogica(comps);
  if (useOrdered) {
    return seq.every((id, i) => id === ids[i]);
  }
  return multisetSignature(seq) === multisetSignature(ids);
}

export function spellMatchesAnySavedSpeedDial(spell: ApiSpell, savedList: ApiSavedSpell[]): boolean {
  return savedList.some((s) => spellMatchesSavedSpeedDial(spell, s));
}
