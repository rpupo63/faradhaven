import { describe, it, expect } from 'vitest';
import { findMatchingPreparedSpell, spellMatchesSavedSpeedDial } from './powderMageSpellMatch';
import type { ApiComponent, ApiSpell, ApiSavedSpell } from '@/types/game';

const mk = (id: string, name: string, category: ApiComponent['category']): ApiComponent =>
  ({
    id,
    name,
    symbol: 'x',
    category,
    tier: 1,
    element: '',
    description: '',
  }) as ApiComponent;

describe('findMatchingPreparedSpell', () => {
  const a = mk('a', 'Ignis', 'Essentia');
  const b = mk('b', 'Nova', 'Forma');
  const ifComp = mk('if1', 'If', 'Logica');

  it('matches multiset when no Logica (order ignored)', () => {
    const prepared: ApiSpell[] = [
      {
        id: 's1',
        name: 'Blast',
        slot_level: 2,
        components: [a, b],
      } as ApiSpell,
    ];
    const match = findMatchingPreparedSpell([b, a], prepared);
    expect(match?.name).toBe('Blast');
  });

  it('requires order when Logica is in the sequence', () => {
    const prepared: ApiSpell[] = [
      {
        id: 's1',
        name: 'Ordered',
        slot_level: 2,
        components: [a, ifComp, b],
      } as ApiSpell,
    ];
    expect(findMatchingPreparedSpell([a, ifComp, b], prepared)?.name).toBe('Ordered');
    expect(findMatchingPreparedSpell([ifComp, a, b], prepared)).toBeNull();
  });

  it('returns null when no spell matches', () => {
    const prepared: ApiSpell[] = [
      {
        id: 's1',
        name: 'Only A',
        slot_level: 1,
        components: [a],
      } as ApiSpell,
    ];
    expect(findMatchingPreparedSpell([a, b], prepared)).toBeNull();
  });
});

describe('spellMatchesSavedSpeedDial', () => {
  const a = mk('a', 'Ignis', 'Essentia');
  const b = mk('b', 'Nova', 'Forma');

  it('matches multiset when no Logica', () => {
    const spell = { id: 's1', name: 'X', slot_level: 2, components: [a, b] } as ApiSpell;
    const saved: ApiSavedSpell = {
      id: 'sd1',
      character_id: 'c1',
      name: 'Dial',
      component_ids: [b.id, a.id],
      slot_index: 0,
    };
    expect(spellMatchesSavedSpeedDial(spell, saved)).toBe(true);
  });
});
