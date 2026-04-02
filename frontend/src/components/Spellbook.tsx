import { useQuery } from '@tanstack/react-query';
import { getSpells } from '@/lib/api';
import { useAuth } from '@/context/AuthContext';
import { SpellListPagination } from '@/components/SpellListPagination';
import { SpellItem } from '@/components/SpellItem';
import {
  REALM_SPELLBOOK_PAGE_SIZE,
  useSpellbookListState,
  type SpellbookListSpell,
} from '@/components/spellbook';

/** Realm-wide spell catalog (not per-character). Character pools use `CharacterSpellbook` + `resolveSpellPoolComponents`. */
export function Spellbook() {
  const { token } = useAuth();
  const { pointsFilter, setPointsFilter, currentPage, setCurrentPage } = useSpellbookListState();

  const { data: paginatedSpells, isLoading, error } = useQuery({
    queryKey: ['spells', { page: currentPage, limit: REALM_SPELLBOOK_PAGE_SIZE, pointsFilter }],
    queryFn: () =>
      getSpells(token ?? undefined, {
        page: currentPage,
        limit: REALM_SPELLBOOK_PAGE_SIZE,
        slot_level: pointsFilter,
      }),
  });

  const spells = paginatedSpells?.spells;
  const totalCount = paginatedSpells?.total_count || 0;
  const totalPages = Math.ceil(totalCount / REALM_SPELLBOOK_PAGE_SIZE);

  const renderSpellItem = (spell: SpellbookListSpell) => <SpellItem spell={spell} />;

  return (
    <SpellListPagination
      initialSpells={spells}
      isLoading={isLoading}
      error={error}
      emptyMessage="No spells have been crafted in the realm."
      headerTitle="All Spells"
      showFilters={true}
      pointsFilter={pointsFilter}
      setPointsFilter={setPointsFilter}
      currentPage={currentPage}
      setCurrentPage={setCurrentPage}
      totalPages={totalPages}
      renderSpellItem={renderSpellItem}
    />
  );
}
