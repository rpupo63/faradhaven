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

export type SpellPhaseWindow = {
  phaseIndex: number;
  /** Start index in flat sequence (inclusive). */
  start: number;
  /** End index in flat sequence (exclusive). */
  end: number;
  /** Non-Logica component indices inside this phase window. */
  indices: number[];
};

/**
 * Derives contiguous phase windows from a flat sequence.
 * Phases are split by Logica components and can be empty.
 */
export function deriveSpellPhaseWindows(components: ApiComponent[]): SpellPhaseWindow[] {
  const windows: SpellPhaseWindow[] = [{ phaseIndex: 0, start: 0, end: components.length, indices: [] }];
  let phaseIndex = 0;

  for (let i = 0; i < components.length; i++) {
    const comp = components[i];
    if (comp.category === 'Logica') {
      windows[phaseIndex].end = i;
      phaseIndex += 1;
      windows.push({ phaseIndex, start: i + 1, end: components.length, indices: [] });
      continue;
    }
    windows[phaseIndex].indices.push(i);
  }

  if (windows.length > 0) {
    windows[windows.length - 1].end = components.length;
  }
  return windows;
}

export function clampPhaseIndex(
  requestedPhaseIndex: number,
  windows: SpellPhaseWindow[],
): number {
  if (windows.length === 0) return 0;
  return Math.max(0, Math.min(requestedPhaseIndex, windows.length - 1));
}

/**
 * Toggle behavior for component clicks in an active phase:
 * - remove one matching component if present in that phase (most recent occurrence),
 * - otherwise insert at phase end (before next Logica link).
 */
export function toggleComponentWithinPhase(
  sequence: ApiComponent[],
  component: ApiComponent,
  requestedPhaseIndex: number,
): { nextSequence: ApiComponent[]; phaseWindows: SpellPhaseWindow[]; resolvedPhaseIndex: number } {
  const phaseWindows = deriveSpellPhaseWindows(sequence);
  const resolvedPhaseIndex = clampPhaseIndex(requestedPhaseIndex, phaseWindows);
  const activeWindow = phaseWindows[resolvedPhaseIndex];
  if (!activeWindow) {
    return { nextSequence: [...sequence, component], phaseWindows, resolvedPhaseIndex: 0 };
  }

  const matchingIndices = activeWindow.indices.filter((idx) => sequence[idx]?.id === component.id);
  if (matchingIndices.length > 0) {
    const removeAt = matchingIndices[matchingIndices.length - 1];
    return {
      nextSequence: sequence.filter((_, idx) => idx !== removeAt),
      phaseWindows,
      resolvedPhaseIndex,
    };
  }

  const insertAt = activeWindow.end;
  return {
    nextSequence: [
      ...sequence.slice(0, insertAt),
      component,
      ...sequence.slice(insertAt),
    ],
    phaseWindows,
    resolvedPhaseIndex,
  };
}

export function validateFormaScopusPerPhase(components: ApiComponent[]): string[] {
  const windows = deriveSpellPhaseWindows(components);
  const errors: string[] = [];

  for (const w of windows) {
    const phaseComponents = w.indices.map((idx) => components[idx]).filter(Boolean);
    if (phaseComponents.length === 0) continue;

    const formaCount = phaseComponents.filter((c) => c.category === 'Forma').length;
    const scopusCount = phaseComponents.filter((c) => c.category === 'Scopus').length;
    const phaseNum = w.phaseIndex + 1;

    if (formaCount === 0) errors.push(`Phase ${phaseNum} requires exactly 1 Forma (shape) component`);
    else if (formaCount > 1) errors.push(`Phase ${phaseNum} allows only 1 Forma (shape) component`);

    if (scopusCount === 0) errors.push(`Phase ${phaseNum} requires exactly 1 Scopus (targeting) component`);
    else if (scopusCount > 1) errors.push(`Phase ${phaseNum} allows only 1 Scopus (targeting) component`);
  }

  return errors;
}
