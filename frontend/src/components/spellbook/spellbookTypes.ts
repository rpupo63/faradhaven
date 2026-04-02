import type { ApiSpell } from '@/types/game';

/** Spell shape from list APIs that may include crafter / character metadata. */
export type SpellbookListSpell = ApiSpell & {
  character_name?: string;
  character_class?: string;
};
