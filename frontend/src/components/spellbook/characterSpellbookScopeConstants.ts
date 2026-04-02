import type { CharacterSpellbookScope } from '@/lib/api';

export const SCOPE_LABELS: Record<CharacterSpellbookScope, string> = {
  mine: 'My spells',
  castable: 'Can cast',
  mine_or_castable: 'Mine or can cast',
  all: 'All spells',
};

export const EMPTY_BY_SCOPE: Record<CharacterSpellbookScope, string> = {
  mine: 'You have not created any spells yet.',
  castable: 'No spells are available with your current components.',
  mine_or_castable: 'No spells match this filter yet.',
  all: 'No spells in the archive.',
};
