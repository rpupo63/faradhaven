import type { NormalizedCharacterSheet } from '@/types/game';

/** Props passed to class-specific blocks under Class Resources (registry-driven). */
export type ClassResourceExtraSectionProps = {
  characterId: string;
  token: string;
  sheet: NormalizedCharacterSheet;
  onGoToSpellForge?: () => void;
};
