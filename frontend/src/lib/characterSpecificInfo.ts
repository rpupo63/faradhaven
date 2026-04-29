import type { NormalizedCharacterSheet } from '@/types/game';

/** Race line includes Changeling (handles “Changeling …” naming). */
export function isChangelingSheet(sheet: NormalizedCharacterSheet): boolean {
  return sheet.character.raceName.includes('Changeling');
}

/** localStorage key for the Changeling “current persona” ribbon on the sheet. */
export function getChangelingPersonaStorageKey(characterId: string): string {
  return `persona_${characterId}`;
}

/**
 * Whether the Character Specific Info card should render any content.
 * Money is treated as tracked when present on the sheet (API/normalized sheets include it).
 */
export function shouldShowCharacterSpecificCard(sheet: NormalizedCharacterSheet): boolean {
  const party = sheet.character.partyName?.trim();
  if (party) return true;
  if (isChangelingSheet(sheet)) return true;
  const money = sheet.money as number | undefined;
  if (money !== undefined && money !== null) return true;
  return false;
}
