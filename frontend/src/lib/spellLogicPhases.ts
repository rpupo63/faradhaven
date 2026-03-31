import type { ApiComponent } from '@/types/game';

/** One slot in the spell chain with its index in the flat ordered array (for reorder/remove). */
export type IndexedSpellComponent = { comp: ApiComponent; index: number };

/**
 * A segment of the spell sequence: either a **phase** (non-Logica components that fire together
 * before the next link) or a **logic link** (If / Then / Therefore) separating phases.
 */
export type SpellSequenceSegment =
  | { kind: 'phase'; phaseNumber: number; items: IndexedSpellComponent[] }
  | { kind: 'logic'; item: IndexedSpellComponent };

/**
 * Splits the flat forge/spell order into phases separated by Logica components.
 * Order is preserved; duplicate Essentia names in different phases stay distinct (e.g. Aqua then cold Aqua).
 */
export function splitSpellSequenceByLogica(components: ApiComponent[]): SpellSequenceSegment[] {
  const segments: SpellSequenceSegment[] = [];
  let phaseBuf: IndexedSpellComponent[] = [];
  let phaseNumber = 0;

  const flushPhase = () => {
    if (phaseBuf.length === 0) return;
    phaseNumber += 1;
    segments.push({ kind: 'phase', phaseNumber, items: phaseBuf });
    phaseBuf = [];
  };

  components.forEach((comp, index) => {
    const slot: IndexedSpellComponent = { comp, index };
    if (comp.category === 'Logica') {
      flushPhase();
      segments.push({ kind: 'logic', item: slot });
    } else {
      phaseBuf.push(slot);
    }
  });
  flushPhase();
  return segments;
}

export function spellChainHasLogica(components: ApiComponent[]): boolean {
  return components.some((c) => c.category === 'Logica');
}
