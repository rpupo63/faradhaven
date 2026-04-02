import { useMemo } from 'react';
import { Sparkles, BookOpen, ChevronDown } from 'lucide-react';
import { cn } from '@/lib/utils';
import { LoadingQuill } from './LoadingQuill';
import type { SpellbookListSpell } from '@/components/spellbook/spellbookTypes';
import { getSpellComponentCount } from '@/lib/spellUtils';
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible';

const POINTS = [1, 2, 3, 4, 5, 6, 7, 8, 9];

interface SpellListPaginationProps {
  initialSpells: SpellbookListSpell[] | undefined;
  isLoading: boolean;
  error: Error | null;
  emptyMessage: string;
  headerTitle: string;
  showFilters?: boolean;
  renderSpellItem: (spell: SpellbookListSpell) => React.ReactNode;
  
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
    title={points === 'all' ? 'All spells' : `${points} components`}
    aria-label={points === 'all' ? 'Show all component counts' : `Filter spells with ${points} components`}
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
    {points === 'all' ? 'All' : points}
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

  const spellsByComponentCountOnPage = useMemo(() => {
    const byCount: Record<number, SpellbookListSpell[]> = {};
    displayedSpells.forEach((s) => {
      const n = getSpellComponentCount(s);
      const key = n > 0 ? n : 0;
      if (!byCount[key]) byCount[key] = [];
      byCount[key].push(s);
    });
    return Object.entries(byCount)
      .map(([count, list]) => ({ count: Number(count), spells: list }))
      .sort((a, b) => a.count - b.count);
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
      {/* sm:flex-1 only: in flex-col (mobile) plain flex-1 + overflow-hidden can clip tall lists and block main scroll */}
      <div className="w-full min-w-0 sm:flex-1 border border-faded-gold/40 rounded-r-lg rounded-bl-lg sm:rounded-bl-none overflow-x-clip overflow-y-visible sm:overflow-hidden bg-card/60">
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

        {spellsByComponentCountOnPage.length === 0 ? (
          <div className="p-8 text-center">
            <p className="text-muted-foreground font-tome-marginalia">
              No spells with this component count
            </p>
          </div>
        ) : (
          <div className="divide-y divide-ink/10">
            {spellsByComponentCountOnPage.map(({ count, spells: countSpells }) => (
              <section key={count} className="p-4">
                <h2 className="tome-section-heading text-lg mb-3">
                  {count === 0
                    ? 'Unknown size'
                    : `${count} ${count === 1 ? 'component' : 'components'}`}
                </h2>
                <div className="space-y-2">
                  {countSpells.map((spell) => {
                    const n = getSpellComponentCount(spell);
                    return (
                      <Collapsible key={spell.id} defaultOpen={false} className="rounded-lg border border-ink/15 bg-card/30">
                        <CollapsibleTrigger
                          className={cn(
                            'group flex w-full items-center gap-2 px-3 py-2.5 text-left font-tome-item-heading text-sm sm:text-base',
                            'hover:bg-muted/25 transition-colors data-[state=open]:border-b data-[state=open]:border-ink/10',
                          )}
                        >
                          <ChevronDown className="h-4 w-4 shrink-0 text-muted-foreground transition-transform duration-200 group-data-[state=open]:rotate-180" />
                          <span className="min-w-0 flex-1 truncate">{spell.name}</span>
                          <span className="shrink-0 text-xs text-muted-foreground font-tome-marginalia tabular-nums">
                            {n > 0 ? `${n} comp.` : '—'}
                          </span>
                        </CollapsibleTrigger>
                        <CollapsibleContent className="overflow-hidden">
                          <div className="border-t border-ink/10 p-2 sm:p-3">{renderSpellItem(spell)}</div>
                        </CollapsibleContent>
                      </Collapsible>
                    );
                  })}
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