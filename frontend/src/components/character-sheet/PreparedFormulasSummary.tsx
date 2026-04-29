import { useQuery } from '@tanstack/react-query';
import { FlaskConical, Wand2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { getCharacterSpells } from '@/lib/api/character';
import type { NormalizedCharacterSheet } from '@/types/game';
import {
  elixiristIgnoresPreparedCap,
  partitionElixiristPreparedSpells,
} from '@/lib/elixiristPreparedLoadout';

const SUMMARY_LIMIT = 80;

interface PreparedFormulasSummaryProps {
  characterId: string;
  token: string;
  sheet: NormalizedCharacterSheet;
  onGoToSpellForge?: () => void;
}

function SpellRows({ spells }: { spells: { id: string; name: string }[] }) {
  return (
    <ul className="space-y-1 max-h-[min(12rem,40vh)] overflow-y-auto pr-1">
      {spells.map((s) => (
        <li
          key={s.id}
          className="text-sm font-tome-marginalia text-foreground/90 leading-snug truncate"
          title={s.name}
        >
          • {s.name}
        </li>
      ))}
    </ul>
  );
}

/**
 * Elixirist: shows **active prepared formulas** (at-will for 1 Emax each) vs class cap.
 * When more spells exist on the character than the cap, lists overflow separately.
 */
export function PreparedFormulasSummary({
  characterId,
  token,
  sheet,
  onGoToSpellForge,
}: PreparedFormulasSummaryProps) {
  const level = sheet.character.level;
  const capIgnored = elixiristIgnoresPreparedCap(level);

  const capResource = sheet.class_resources?.find((r) => r.key === 'prepared_formulas');
  const cap = capResource?.value ?? 0;

  const { data, isLoading } = useQuery({
    queryKey: ['character-spells', characterId, 'sheet-summary'],
    queryFn: () =>
      getCharacterSpells(characterId, token, { page: 1, limit: SUMMARY_LIMIT }),
    enabled: !!characterId && !!token,
    staleTime: 15_000,
  });

  const spells = data?.spells ?? [];
  const totalCount = data?.total_count ?? spells.length;
  const { active, overflow } = partitionElixiristPreparedSpells(spells, cap, level);
  const overCap = !capIgnored && cap > 0 && totalCount > cap;

  return (
    <div className="rounded-md border border-border/70 bg-card/40 px-3 py-2.5 space-y-2">
      <div className="flex items-start gap-2">
        <FlaskConical className="h-4 w-4 shrink-0 text-primary mt-0.5" />
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2 flex-wrap">
            <span className="text-sm font-tome-subheading text-primary">Active prepared formulas</span>
            <span className="text-micro font-tome-marginalia text-muted-foreground ml-auto tabular-nums shrink-0">
              {capIgnored ? (
                <>{Math.min(totalCount, spells.length)} ready</>
              ) : (
                <>
                  {Math.min(active.length, cap)} / {cap}
                </>
              )}
            </span>
          </div>
          <p className="text-micro font-tome-marginalia text-muted-foreground mt-1 leading-snug">
            At-will for <span className="text-foreground/85">1 Emax</span> per use while prepared. Remix after a{' '}
            <span className="text-foreground/85">Short</span> or <span className="text-foreground/85">Long Rest</span>.
          </p>
        </div>
      </div>
      {isLoading ? (
        <p className="text-micro text-muted-foreground">Loading formulas…</p>
      ) : totalCount === 0 ? (
        <p className="text-micro text-muted-foreground italic leading-snug">
          No formulas on this character yet. Forge them in the Spellbook, then choose your prepared set after a rest.
        </p>
      ) : (
        <>
          <div>
            <p className="text-micro font-tome-marginalia text-muted-foreground uppercase mb-1">
              {capIgnored ? 'All formulas prepared' : 'Prepared loadout'}
            </p>
            <SpellRows spells={capIgnored ? spells : active} />
          </div>
          {!capIgnored && overflow.length > 0 ? (
            <div className="rounded border border-amber-500/30 bg-amber-500/5 px-2 py-2 space-y-1">
              <p className="text-micro font-tome-marginalia text-amber-700 dark:text-amber-400 leading-snug">
                Not in loadout ({overflow.length}) — trim or swap after rest so only {cap} formulas match your Prepared
                Formulas limit.
              </p>
              <ul className="space-y-0.5 max-h-24 overflow-y-auto opacity-90">
                {overflow.map((s) => (
                  <li key={s.id} className="text-micro truncate text-muted-foreground" title={s.name}>
                    ○ {s.name}
                  </li>
                ))}
              </ul>
            </div>
          ) : null}
        </>
      )}
      {totalCount > SUMMARY_LIMIT && (
        <p className="text-micro text-muted-foreground">
          Showing first {SUMMARY_LIMIT} of {totalCount} — see Spellbook for the full list.
        </p>
      )}
      {overCap && overflow.length === 0 && (
        <p className="text-micro text-amber-600 dark:text-amber-500 leading-snug">
          More spells on this character than your Prepared Formulas limit — adjust in the Spellbook so your active set
          matches your rest picks.
        </p>
      )}
      {onGoToSpellForge && (
        <Button type="button" variant="outline" size="sm" className="w-full gap-2 h-8" onClick={onGoToSpellForge}>
          <Wand2 className="h-3.5 w-3.5" />
          Spellbook
        </Button>
      )}
    </div>
  );
}
