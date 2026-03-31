/**
 * Dice rolling utilities.
 *
 * All public roll functions are async and use DiceManager so that
 * the 3D physics simulation is authoritative — the face that lands up
 * determines the result.  If the user closes the triggering context
 * (dialog, panel) before dice settle, DiceManager.clear() is called
 * and every awaited function returns null.  Callers should early-return
 * on null.
 *
 * For batch operations that must not trigger an animation (e.g. rolling
 * initiative for every map token), use DiceManager.rollSync() directly.
 */

import { DiceManager } from './dice-manager';

// ---------------------------------------------------------------------------
// Clear event — still dispatched by components on close so DiceAnimation
// can stop the in-progress 3D animation.
// ---------------------------------------------------------------------------

export const DICE_CLEAR_EVENT = 'faradhaven-dice-clear';

export function dispatchClearDice() {
  if (typeof window !== 'undefined') {
    window.dispatchEvent(new CustomEvent(DICE_CLEAR_EVENT));
  }
}

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

export type RollType = 'spell_attack' | 'melee_attack' | 'ranged_attack' | 'ability_check' | 'saving_throw';

export interface RollResult {
  type: RollType;
  label: string;
  total: number;
  natural: number;
  modifier: number;
  timestamp: number;
  advantage?: boolean;
  disadvantage?: boolean;
}

export interface DamageRollResult {
  damageType: string;
  dice: string;
  rolls: number[];
  modifier: number;
  total: number;
  isCrit?: boolean;
}

// ---------------------------------------------------------------------------
// Notation parser  ("2d6+3" → { count:2, sides:6, modifier:3 })
// ---------------------------------------------------------------------------

export function parseDiceNotation(notation: string): {
  count: number;
  sides: number;
  modifier: number;
} {
  const match = notation.toLowerCase().match(/^(\d*)d(\d+)(?:\s*([+-]\s*\d+))?$/);
  if (!match) return { count: 1, sides: 20, modifier: 0 };
  const count = match[1] ? parseInt(match[1], 10) : 1;
  const sides = parseInt(match[2], 10);
  const modStr = match[3] ? match[3].replace(/\s/g, '') : undefined;
  const modifier = modStr ? parseInt(modStr, 10) : 0;
  return { count, sides, modifier };
}

// ---------------------------------------------------------------------------
// Core async roll functions
// All return null if the roll was cancelled mid-flight.
// ---------------------------------------------------------------------------

export async function rollD20(modifier: number, label?: string): Promise<{
  total: number;
  natural: number;
  modifier: number;
} | null> {
  const result = await DiceManager.roll('1d20', label, modifier);
  if (result === null) return null;
  const natural = result.rolls[0];
  return { total: natural + modifier, natural, modifier };
}

export async function rollD20Advantage(modifier: number, label?: string): Promise<{
  total: number;
  natural: number;
  modifier: number;
  rolls: [number, number];
} | null> {
  const result = await DiceManager.roll('2d20', label, modifier);
  if (result === null) return null;
  const [r1, r2] = result.rolls as [number, number];
  const natural = Math.max(r1, r2);
  return { total: natural + modifier, natural, modifier, rolls: [r1, r2] };
}

export async function rollD20Disadvantage(modifier: number, label?: string): Promise<{
  total: number;
  natural: number;
  modifier: number;
  rolls: [number, number];
} | null> {
  const result = await DiceManager.roll('2d20', label, modifier);
  if (result === null) return null;
  const [r1, r2] = result.rolls as [number, number];
  const natural = Math.min(r1, r2);
  return { total: natural + modifier, natural, modifier, rolls: [r1, r2] };
}

/** Roll an arbitrary notation string like "2d6+3". */
export async function rollDiceNotation(notation: string, label?: string): Promise<{
  total: number;
  rolls: number[];
  modifier: number;
} | null> {
  const parsed = parseDiceNotation(notation);
  const result = await DiceManager.roll(`${parsed.count}d${parsed.sides}`, label, parsed.modifier);
  if (result === null) return null;
  return { total: result.total + parsed.modifier, rolls: result.rolls, modifier: parsed.modifier };
}

/** Roll a single die of the given number of sides. */
export async function rollD(sides: number, label?: string): Promise<number | null> {
  const result = await DiceManager.roll(`1d${sides}`, label);
  if (result === null) return null;
  return result.rolls[0];
}

/** Roll multiple dice of the same type. Returns array of individual values. */
export async function rollDice(count: number, sides: number, label?: string): Promise<number[] | null> {
  const result = await DiceManager.roll(`${count}d${sides}`, label);
  if (result === null) return null;
  return result.rolls;
}

export async function rollHitDie(hitDie: number, label?: string): Promise<number | null> {
  return rollD(hitDie, label);
}

export async function rollHitDice(count: number, hitDie: number, label?: string): Promise<number[] | null> {
  return rollDice(count, hitDie, label);
}

/**
 * Roll 4d6 and drop the lowest.
 *
 * @param silent  Skip the animation (use for batch operations like auto-rolling
 *                all six ability scores at once in character creation).
 */
export async function roll4d6DropLowest(silent = false): Promise<{
  total: number;
  rolls: number[];
  dropped: number;
} | null> {
  const raw = silent ? DiceManager.rollSync('4d6') : await DiceManager.roll('4d6');
  if (raw === null) return null;
  const sorted = [...raw.rolls].sort((a, b) => a - b);
  const dropped = sorted[0];
  const kept = sorted.slice(1);
  return { total: kept.reduce((a, b) => a + b, 0), rolls: kept, dropped };
}

/**
 * Exploding dice (re-roll on max).  Synchronous — no 3D animation since
 * the number of dice is unknowable upfront.
 */
export function rollExploding(sides: number): { total: number; rolls: number[] } {
  const rolls: number[] = [];
  let total = 0;
  let current: number;
  do {
    current = DiceManager.rollSync(`1d${sides}`).rolls[0];
    rolls.push(current);
    total += current;
  } while (current === sides);
  return { total, rolls };
}
