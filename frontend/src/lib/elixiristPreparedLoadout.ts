import type { ApiSpell } from '@/types/game/api';

/** Level 20 capstone: all known formulas count as prepared simultaneously. */
export function elixiristIgnoresPreparedCap(level: number): boolean {
  return level >= 20;
}

/**
 * Split character spells into the **active prepared loadout** (first `preparedCap` entries)
 * vs **overflow** when the player has more spells attached than their Prepared Formulas limit.
 *
 * Order follows the API list (typically creation order). At the table, players should keep
 * only their chosen prepared set on the character; this split is a UI hint when they have not.
 */
export function partitionElixiristPreparedSpells(
  spells: ApiSpell[],
  preparedCap: number,
  characterLevel: number,
): { active: ApiSpell[]; overflow: ApiSpell[] } {
  if (elixiristIgnoresPreparedCap(characterLevel)) {
    return { active: spells, overflow: [] };
  }
  if (preparedCap <= 0) {
    return { active: spells, overflow: [] };
  }
  return {
    active: spells.slice(0, preparedCap),
    overflow: spells.slice(preparedCap),
  };
}
