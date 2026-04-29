import { useCallback, useEffect, useState } from 'react';
import { rollDice } from '@/lib/dice';
import {
  type DeathSavePersisted,
  clearDeathSaves,
  initialDeathSaveState,
  loadDeathSaves,
  saveDeathSaves,
} from '@/lib/deathSavesStorage';

export type DeathSaveRollOutcome = {
  rolls: number[];
  best: number;
  success: boolean;
};

/**
 * Tracks death saving throws while at 0 HP (Faradhaven: 5d20 keep highest; 11+ success, 10− failure).
 * Persisted per character in localStorage; resets when current HP rises above 0.
 */
export function useDeathSaves(characterId: string, currentHP: number) {
  const [state, setState] = useState<DeathSavePersisted>(() => loadDeathSaves(characterId));

  useEffect(() => {
    setState(loadDeathSaves(characterId));
  }, [characterId]);

  useEffect(() => {
    if (currentHP > 0) {
      setState({ ...initialDeathSaveState });
      clearDeathSaves(characterId);
    }
  }, [currentHP, characterId]);

  const dying = currentHP <= 0;
  const canRoll = dying && !state.stable && !state.dead;

  const applyNextState = useCallback(
    (updater: (prev: DeathSavePersisted) => DeathSavePersisted) => {
      setState(prev => {
        const next = updater(prev);
        saveDeathSaves(characterId, next);
        return next;
      });
    },
    [characterId]
  );

  const rollDeathSave = useCallback(async (): Promise<DeathSaveRollOutcome | null> => {
    if (!canRoll) return null;
    const rolls = await rollDice(5, 20, 'Death save');
    if (!rolls || rolls.length === 0) return null;

    const best = Math.max(...rolls);
    const success = best >= 11;

    applyNextState(prev => {
      if (prev.stable || prev.dead) return prev;
      if (success) {
        const successes = prev.successes + 1;
        if (successes >= 3) {
          return {
            successes: 0,
            failures: 0,
            stable: true,
            dead: false,
          };
        }
        return { ...prev, successes, failures: prev.failures };
      }
      const failures = prev.failures + 1;
      if (failures >= 3) {
        return { ...prev, failures: 3, dead: true };
      }
      return { ...prev, failures, successes: prev.successes };
    });

    return { rolls, best, success };
  }, [applyNextState, canRoll]);

  return {
    ...state,
    dying,
    canRoll,
    rollDeathSave,
  };
}
