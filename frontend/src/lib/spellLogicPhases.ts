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

/** Non-Logica components merged by id for crucible display; `indices` are positions in the flat sequence. */
export type SpellComponentBucket = {
  comp: ApiComponent;
  indices: number[];
};

/**
 * Groups a flat sequence into buckets (first-seen id order). Used when no Logica — cast order is cosmetic.
 */
export function bucketSequenceByComponentId(components: ApiComponent[]): SpellComponentBucket[] {
  const order: string[] = [];
  const map = new Map<string, SpellComponentBucket>();
  for (let i = 0; i < components.length; i++) {
    const comp = components[i];
    const id = comp.id;
    let b = map.get(id);
    if (!b) {
      b = { comp, indices: [] };
      map.set(id, b);
      order.push(id);
    }
    b.indices.push(i);
  }
  return order.map((id) => map.get(id)!);
}

/** Buckets one multi-phase segment (non-Logica items only). */
export function bucketIndexedPhase(items: IndexedSpellComponent[]): SpellComponentBucket[] {
  const order: string[] = [];
  const map = new Map<string, SpellComponentBucket>();
  for (const { comp, index } of items) {
    const id = comp.id;
    let b = map.get(id);
    if (!b) {
      b = { comp, indices: [] };
      map.set(id, b);
      order.push(id);
    }
    b.indices.push(index);
  }
  return order.map((id) => map.get(id)!);
}
