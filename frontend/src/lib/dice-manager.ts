/**
 * DiceManager — singleton that owns the @3d-dice/dice-box instance.
 *
 * All async roll functions in dice.ts go through here so that:
 *   1. The 3D animation is the authoritative source of die values.
 *   2. Cancelling (clear) while dice are mid-flight resolves null so callers
 *      can bail out gracefully without updating state.
 */

export interface DiceResult {
  rolls: number[];
  total: number;
  label?: string;
  modifier?: number;
}

export type ResultListener = (result: DiceResult & { notation: string }) => void;

interface PendingRoll {
  resolve: (r: DiceResult | null) => void;
  notation: string;
  label?: string;
  modifier?: number;
}

class DiceManagerClass {
  private box: { roll: (n: unknown) => Promise<void>; clear: () => void } | null = null;
  private pending: PendingRoll | null = null;
  private onResultListener: ResultListener | null = null;

  /** Resolves when dice-box finishes init(). roll() awaits this. */
  private readyResolve: (() => void) | null = null;
  private readyPromise: Promise<void> = new Promise(r => { this.readyResolve = r; });

  /** Called by DiceAnimation once the dice-box instance is ready. */
  register(box: { roll: (n: unknown) => Promise<void>; clear: () => void }) {
    this.box = box;
    this.readyResolve?.();
    this.readyResolve = null;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (box as any).onRollComplete = (groups: Array<{ rolls: Array<{ value: number }> }>) => {
      const rolls = groups.flatMap(g => g.rolls.map(d => d.value));
      const rawTotal = rolls.reduce((a, b) => a + b, 0);
      const notation = this.pending?.notation ?? '';
      const label = this.pending?.label;
      const modifier = this.pending?.modifier ?? 0;
      const total = rawTotal + modifier;

      if (this.pending) {
        const { resolve } = this.pending;
        this.pending = null;
        resolve({ rolls, total, label, modifier });
      }

      this.onResultListener?.({ rolls, total, notation, label, modifier });
    };
  }

  /** DiceAnimation subscribes here to update its overlay after dice settle. */
  setOnResult(fn: ResultListener | null) {
    this.onResultListener = fn;
  }

  private onRollStartListener: (() => void) | null = null;

  /** DiceAnimation subscribes here to play sound the moment dice appear. */
  setOnRollStart(fn: (() => void) | null) {
    this.onRollStartListener = fn;
  }

  isReady(): boolean {
    return this.box !== null;
  }

  /**
   * Roll dice via the 3D engine. Resolves once dice physically settle.
   * Returns null if the roll is cancelled (via clear()) before settling.
   */
  async roll(notation: string, label?: string, modifier?: number): Promise<DiceResult | null> {
    // Clear any existing dice from the tray first.
    this.box?.clear();

    // Cancel any in-flight roll promise — new roll wins.
    if (this.pending) {
      this.pending.resolve(null);
      this.pending = null;
    }

    // Wait for dice-box to finish initializing (WASM + asset loading).
    // Cap at 6 seconds — fall back to random if the engine never comes up.
    if (!this.box) {
      await Promise.race([
        this.readyPromise,
        new Promise<void>(r => setTimeout(r, 6000)),
      ]);
    }

    if (!this.box) {
      return this.fallback(notation, label, modifier);
    }

    return new Promise<DiceResult | null>((resolve) => {
      this.pending = { resolve, notation, label, modifier };
      this.onRollStartListener?.();
      this.box!.roll(notation).catch(() => {
        if (this.pending?.notation === notation) {
          const { label, modifier } = this.pending;
          this.pending = null;
          resolve(this.fallback(notation, label, modifier));
        }
      });
    });
  }

  /**
   * Synchronous random roll — no animation.
   * Use for batch operations where per-die animation would be disruptive
   * (e.g. rolling initiative for every token, auto-rolling 6 ability scores).
   */
  rollSync(notation: string, label?: string, modifier?: number): DiceResult {
    const fallback = this.fallback(notation, label, modifier);
    return fallback;
  }

  /** Stop animation and cancel any pending roll promise (resolves null). */
  clear() {
    if (this.pending) {
      this.pending.resolve(null);
      this.pending = null;
    }
    this.box?.clear();
  }

  private fallback(notation: string, label?: string, modifier?: number): DiceResult {
    const match = notation.toLowerCase().match(/^(\d*)d(\d+)$/);
    if (!match) return { rolls: [1], total: 1 + (modifier ?? 0), label, modifier };
    const count = match[1] ? parseInt(match[1], 10) : 1;
    const sides = parseInt(match[2], 10);
    const rolls = Array.from({ length: count }, () => Math.floor(Math.random() * sides) + 1);
    const rawTotal = rolls.reduce((a, b) => a + b, 0);
    const total = rawTotal + (modifier ?? 0);
    return { rolls, total, label, modifier };
  }
}

export const DiceManager = new DiceManagerClass();
