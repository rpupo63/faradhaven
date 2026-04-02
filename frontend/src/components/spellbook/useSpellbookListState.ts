import { useCallback, useState } from 'react';

/**
 * Points filter + pagination for {@link SpellListPagination}.
 * Changing the points filter resets to page 1. Reset page yourself when other
 * query keys change (e.g. character spell scope).
 */
export function useSpellbookListState() {
  const [pointsFilter, setPointsFilterState] = useState<number | 'all'>('all');
  const [currentPage, setCurrentPage] = useState(1);

  const setPointsFilter = useCallback((points: number | 'all') => {
    setPointsFilterState(points);
    setCurrentPage(1);
  }, []);

  return { pointsFilter, setPointsFilter, currentPage, setCurrentPage };
}
