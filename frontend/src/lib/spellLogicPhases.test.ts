import { describe, expect, it } from 'vitest';
import type { ApiComponent } from '@/types/game';
import {
  deriveSpellPhaseWindows,
  toggleComponentWithinPhase,
} from './spellLogicPhases';

function comp(
  id: string,
  category: ApiComponent['category'],
  name = id,
): ApiComponent {
  return {
    id,
    name,
    category,
    tier: 1,
    symbol: id.toUpperCase(),
    description: `${name} description`,
  };
}

describe('deriveSpellPhaseWindows', () => {
  it('creates phase windows split by logica including empty phases', () => {
    const flame = comp('flame', 'Essentia');
    const then = comp('then', 'Logica', 'Then');
    const cold = comp('cold', 'Essentia');
    const therefore = comp('therefore', 'Logica', 'Therefore');
    const sequence = [flame, then, cold, therefore];

    const windows = deriveSpellPhaseWindows(sequence);

    expect(windows).toHaveLength(3);
    expect(windows[0]).toMatchObject({ phaseIndex: 0, start: 0, end: 1, indices: [0] });
    expect(windows[1]).toMatchObject({ phaseIndex: 1, start: 2, end: 3, indices: [2] });
    expect(windows[2]).toMatchObject({ phaseIndex: 2, start: 4, end: 4, indices: [] });
  });
});

describe('toggleComponentWithinPhase', () => {
  it('removes most recent matching component from active phase on reclick', () => {
    const flame = comp('flame', 'Essentia');
    const then = comp('then', 'Logica', 'Then');
    const cold = comp('cold', 'Essentia');
    const sequence = [flame, flame, then, cold, flame];

    const { nextSequence } = toggleComponentWithinPhase(sequence, flame, 0);

    expect(nextSequence.map((c) => c.id)).toEqual(['flame', 'then', 'cold', 'flame']);
  });

  it('adds a component only inside the requested phase window', () => {
    const flame = comp('flame', 'Essentia');
    const then = comp('then', 'Logica', 'Then');
    const cold = comp('cold', 'Essentia');
    const spark = comp('spark', 'Actio');
    const sequence = [flame, then, cold];

    const { nextSequence } = toggleComponentWithinPhase(sequence, spark, 0);

    expect(nextSequence.map((c) => c.id)).toEqual(['flame', 'spark', 'then', 'cold']);
  });

  it('supports single-phase toggle behavior when no logica exists', () => {
    const flame = comp('flame', 'Essentia');
    const gust = comp('gust', 'Actio');

    const added = toggleComponentWithinPhase([flame], gust, 0).nextSequence;
    expect(added.map((c) => c.id)).toEqual(['flame', 'gust']);

    const removed = toggleComponentWithinPhase(added, gust, 0).nextSequence;
    expect(removed.map((c) => c.id)).toEqual(['flame']);
  });
});
