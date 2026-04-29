const key = (characterId: string) => `faradhaven-death-saves:${characterId}`;

export interface DeathSavePersisted {
  successes: number;
  failures: number;
  stable: boolean;
  dead: boolean;
}

export const initialDeathSaveState: DeathSavePersisted = {
  successes: 0,
  failures: 0,
  stable: false,
  dead: false,
};

export function loadDeathSaves(characterId: string): DeathSavePersisted {
  try {
    const raw = localStorage.getItem(key(characterId));
    if (!raw) return { ...initialDeathSaveState };
    const parsed = JSON.parse(raw) as Partial<DeathSavePersisted>;
    return {
      ...initialDeathSaveState,
      ...parsed,
      successes: Math.min(3, Math.max(0, Number(parsed.successes) || 0)),
      failures: Math.min(3, Math.max(0, Number(parsed.failures) || 0)),
      stable: Boolean(parsed.stable),
      dead: Boolean(parsed.dead),
    };
  } catch {
    return { ...initialDeathSaveState };
  }
}

export function saveDeathSaves(characterId: string, state: DeathSavePersisted) {
  localStorage.setItem(key(characterId), JSON.stringify(state));
}

export function clearDeathSaves(characterId: string) {
  localStorage.removeItem(key(characterId));
}
