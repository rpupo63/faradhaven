import { useMemo, useState, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import { getSpells } from '@/lib/api';
import { useAuth } from '@/context/AuthContext';
import { ApiSpell as BaseApiSpell } from '@/types/game';
import { LoadingQuill } from './LoadingQuill';
import { SpellListPagination } from './SpellListPagination'; // NEW: Import SpellListPagination
import { SpellItem } from './SpellItem'; // NEW: Import SpellItem

// The API now returns character info, so we extend the base type.
type ApiSpell = BaseApiSpell & {
  character_name?: string;
  character_class?: string;
};

const PAGINATION_LIMIT = 10; // Number of spells to display per page

export function Spellbook() {
  const { token } = useAuth();
  const [pointsFilter, setPointsFilter] = useState<number | 'all'>('all'); // NEW: pointsFilter state
  const [currentPage, setCurrentPage] = useState(1); // NEW: currentPage state

  const { data: paginatedSpells, isLoading, error } = useQuery({
    queryKey: ['spells', { page: currentPage, limit: PAGINATION_LIMIT, pointsFilter }],
    queryFn: () => getSpells(token ?? undefined, { page: currentPage, limit: PAGINATION_LIMIT, slot_level: pointsFilter }),
  });

  const spells = paginatedSpells?.spells;
  const totalCount = paginatedSpells?.total_count || 0;
  const totalPages = Math.ceil(totalCount / PAGINATION_LIMIT);

  // NEW: Define renderSpellItem function
  const renderSpellItem = (spell: ApiSpell) => (
    <SpellItem spell={spell} />
  );

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
      renderSpellItem={renderSpellItem} // NEW: Pass renderSpellItem prop
    />
  );
}

