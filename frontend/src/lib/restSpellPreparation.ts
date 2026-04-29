/**
 * Classes whose rules tie **changing which spells/combinations you have ready** to rests
 * (seed: `backend/seed/faradhaven_classes/`).
 *
 * - The Elixirist — Potion Formulary: swap Prepared Formulas on Short or Long Rest (`prepared_formulas` cap).
 * - The Piston Brawler (13+) — Advanced Calibration: pre-calculate Blueprints during Short or Long Rest.
 *
 * Other classes may refresh pools (e.g. Powder Mage timer/Speed Dial **uses**) but do not use the same
 * “swap your prepared picks on this rest” rule in seed text.
 */
export function shouldOfferRestSpellPreparation(
  gameClassName: string | undefined,
  level: number,
): boolean {
  if (!gameClassName) return false;
  if (gameClassName === 'The Elixirist') return true;
  if (gameClassName === 'The Piston Brawler' && level >= 13) return true;
  return false;
}

export function restSpellPreparationHint(
  gameClassName: string | undefined,
  restKind: 'short' | 'long',
): string {
  const restLabel = restKind === 'short' ? 'Short Rest' : 'Long Rest';
  if (gameClassName === 'The Elixirist') {
    return `After a ${restLabel}, you may swap any or all Prepared Formulas (up to your Prepared Formulas limit). Use your spell list below to pick what you’re preparing for the day — same view as your Prepared Spells page.`;
  }
  if (gameClassName === 'The Piston Brawler') {
    return `At level 13+, you pre-calculate Blueprint combinations during Short or Long Rests. Review your spells here; save combinations to Blueprint slots from the Spellbook tab (like Speed Dial).`;
  }
  return '';
}
