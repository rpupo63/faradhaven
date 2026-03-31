import { useMemo, useState, useEffect } from 'react';
import { ApiSpell as BaseApiSpell } from '@/types/game';
import { Sparkles, ChevronDown, Sword, ShieldCheck, Heart, UserSquare, BookOpen } from 'lucide-react';
import { cn } from '@/lib/utils';
import { motion, AnimatePresence } from 'framer-motion';
import { LoadingQuill } from './LoadingQuill'; // NEW: Add LoadingQuill import

// The API now returns character info, so we extend the base type.
type ApiSpell = BaseApiSpell & {
  character_name?: string;
  character_class?: string;
};

const POINTS = [1, 2, 3, 4, 5, 6, 7, 8, 9];
const PAGINATION_LIMIT = 10; // Number of spells to display per page

interface SpellListPaginationProps {
  initialSpells: ApiSpell[] | undefined;
  isLoading: boolean;
  error: Error | null;
  emptyMessage: string;
  headerTitle: string;
  showFilters?: boolean; // New prop to control filter visibility
  renderSpellItem: (spell: ApiSpell) => React.ReactNode; // NEW: Custom render prop for each spell item
  
  // NEW: Pagination and Filter Props (controlled by parent)
  pointsFilter: number | 'all';
  setPointsFilter: (points: number | 'all') => void;
  currentPage: number;
  setCurrentPage: (page: number) => void;
  totalPages: number;
}

// FilterTab component (now outside SpellListPagination)
const FilterTab: React.FC<{
  points: number | 'all';
  vertical?: boolean;
  setPointsFilter: (points: number | 'all') => void;
  pointsFilter: number | 'all';
  setCurrentPage: (page: number) => void;
}> = ({ points, vertical = false, setPointsFilter, pointsFilter, setCurrentPage }) => (
  <button
    type="button"
    onClick={() => {
      setPointsFilter(points);
      // Ensure page resets to 1 when filter changes
      if (pointsFilter !== points) {
        setCurrentPage(1);
      }
    }}
    className={cn(
      'ledger-tab py-2 px-2 text-center text-xs font-tome-marginalia transition-colors',
      pointsFilter === points ? 'bg-primary/10 text-primary border-primary/50 z-10' : 'text-muted-foreground hover:text-foreground',
      vertical ? 'w-full' : 'flex-1 min-w-[3rem] border-b-0 border-r rounded-t-md rounded-bl-none rounded-br-none'
    )}
    style={vertical ? { writingMode: 'vertical-rl', textOrientation: 'mixed' } : {}}
    aria-pressed={pointsFilter === points}
  >
    {points === 'all' ? 'All' : `${points} Pts`}
  </button>
);

export function SpellListPagination({
  initialSpells,
  isLoading,
  error,
  emptyMessage,
  headerTitle,
  showFilters = true,
  renderSpellItem,
  pointsFilter,        // NEW: Destructure from props
  setPointsFilter,     // NEW: Destructure from props
  currentPage,         // NEW: Destructure from props
  setCurrentPage,      // NEW: Destructure from props
  totalPages,          // NEW: Destructure from props
}: SpellListPaginationProps) {
  // The 'initialSpells' prop now represents the currently displayed page of filtered/paginated spells.
  const displayedSpells = useMemo(() => initialSpells || [], [initialSpells]);

  // Re-group current page's spells by points for display structure
  const spellsByPointsOnPage = useMemo(() => {
    const byPoints: Record<number, ApiSpell[]> = {};
    displayedSpells.forEach((s) => {
      if (!byPoints[s.slot_level]) byPoints[s.slot_level] = [];
      byPoints[s.slot_level].push(s);
    });
    return Object.entries(byPoints)
      .map(([pts, list]) => ({ points: Number(pts), spells: list }))
      .sort((a, b) => a.points - b.points);
  }, [displayedSpells]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <LoadingQuill label="Loading spells..." />
      </div>
    );
  }

  if (error) {
    return (
      <div className="arcane-border rounded-lg p-8 text-center">
        <BookOpen className="w-12 h-12 mx-auto mb-4 text-destructive" />
        <h3 className="font-tome-heading text-lg text-destructive">Failed to Load Spells</h3>
        <p className="text-sm text-muted-foreground font-tome-marginalia">
          {error instanceof Error ? error.message : 'An error occurred'}
        </p>
      </div>
    );
  }

  if (!initialSpells || initialSpells.length === 0) {
    return (
      <div className="arcane-border rounded-lg p-8 text-center">
        <Sparkles className="w-12 h-12 mx-auto mb-4 text-muted-foreground" />
        <h3 className="font-tome-heading text-lg text-muted-foreground mb-2">{emptyMessage}</h3>
        <p className="text-sm text-muted-foreground font-tome-marginalia">
          No spells match your criteria.
        </p>
      </div>
    );
  }

  return (
    <div className="flex flex-col sm:flex-row gap-0 w-full">
      {showFilters && (
        <>
          {/* Mobile Filter Bar (Horizontal) */}
          <div className="flex sm:hidden overflow-x-auto gap-0.5 mb-0 px-1 scrollbar-hide shrink-0">
            <FilterTab points="all" setPointsFilter={setPointsFilter} pointsFilter={pointsFilter} setCurrentPage={setCurrentPage} />
            {POINTS.map((pts) => (
              <FilterTab key={pts} points={pts} setPointsFilter={setPointsFilter} pointsFilter={pointsFilter} setCurrentPage={setCurrentPage} />
            ))}
          </div>

          {/* Desktop Filter Sidebar (Vertical) */}
          <aside className="hidden sm:flex flex-col gap-1 w-14 shrink-0 pr-0">
            <FilterTab points="all" vertical setPointsFilter={setPointsFilter} pointsFilter={pointsFilter} setCurrentPage={setCurrentPage} />
            {POINTS.map((pts) => (
              <FilterTab key={pts} points={pts} vertical setPointsFilter={setPointsFilter} pointsFilter={pointsFilter} setCurrentPage={setCurrentPage} />
            ))}
          </aside>
        </>
      )}

      {/* Spellbook Content */}
      <div className="flex-1 min-w-0 border border-faded-gold/40 rounded-r-lg rounded-bl-lg sm:rounded-bl-none overflow-hidden bg-card/60">
        <div className="p-4 border-b border-faded-gold/20">
          <div className="flex items-center gap-2">
            <BookOpen className="w-5 h-5 text-primary" />
            <h3 className="font-tome-heading text-lg text-primary">
              {headerTitle}
            </h3>
            <span className="text-muted-foreground text-sm font-tome-marginalia">
              ({displayedSpells.length} {displayedSpells.length === 1 ? 'spell' : 'spells'})
            </span>
          </div>
        </div>

        {spellsByPointsOnPage.length === 0 ? (
          <div className="p-8 text-center">
            <p className="text-muted-foreground font-tome-marginalia">
              No spells with this point cost
            </p>
          </div>
        ) : (
          <div className="divide-y divide-ink/10">
            {spellsByPointsOnPage.map(({ points, spells: pointSpells }) => (
              <section key={points} className="p-4">
                <h2 className="tome-section-heading text-lg mb-3">
                  {points} Spell Points
                </h2>
                <div className="space-y-2">
                  {pointSpells.map((spell) => (
                    <div key={spell.id}>
                      {renderSpellItem(spell)}
                    </div>
                  ))}
                </div>
              </section>
            ))}
            {totalPages > 1 && (
              <div className="flex items-center justify-between p-4 border-t border-ink/10">
                <button
                  onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                  disabled={currentPage === 1}
                  className="px-4 py-2 rounded-md bg-muted-foreground/20 text-muted-foreground hover:bg-muted-foreground/30 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Previous
                </button>
                <span className="font-tome-marginalia text-sm text-muted-foreground">
                  Page {currentPage} of {totalPages}
                </span>
                <button
                  onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
                  disabled={currentPage === totalPages}
                  className="px-4 py-2 rounded-md bg-muted-foreground/20 text-muted-foreground hover:bg-muted-foreground/30 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  Next
                </button>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}