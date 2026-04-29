import { useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Layers } from 'lucide-react';
import { getSpeedDial } from '@/lib/api/character';
import type { ApiSavedSpell } from '@/types/game/api';
import type { NormalizedCharacterSheet } from '@/types/game';

interface BlueprintSlotsSummaryProps {
  characterId: string;
  token: string;
  sheet: NormalizedCharacterSheet;
  /** Shown above the slot list */
  heading: string;
}

/**
 * Read-only summary of Speed Dial / Blueprint slots (same API) so picked spells show on the class card.
 */
export function BlueprintSlotsSummary({
  characterId,
  token,
  sheet,
  heading,
}: BlueprintSlotsSummaryProps) {
  const maxSlots = sheet.class_resources?.find((r) => r.key === 'speed_dial_slots')?.value ?? 0;

  const { data: saved = [], isLoading } = useQuery({
    queryKey: ['speed-dial', characterId],
    queryFn: () => getSpeedDial(characterId, token),
    enabled: !!characterId && !!token && maxSlots > 0,
    staleTime: 15_000,
  });

  const bySlot = useMemo(() => {
    const m = new Map<number, ApiSavedSpell>();
    for (const s of saved) {
      m.set(s.slot_index, s);
    }
    return m;
  }, [saved]);

  if (maxSlots <= 0) {
    return (
      <div className="rounded-md border border-border/60 bg-muted/20 px-3 py-2 text-micro text-muted-foreground leading-snug">
        <span className="flex items-center gap-2 font-tome-subheading text-primary text-sm">
          <Layers className="h-4 w-4 shrink-0" />
          {heading}
        </span>
        <p className="mt-1.5 italic">No blueprint slots at this level.</p>
      </div>
    );
  }

  return (
    <div className="rounded-md border border-border/70 bg-card/40 px-3 py-2.5 space-y-2">
      <span className="flex items-center gap-2 text-sm font-tome-subheading text-primary">
        <Layers className="h-4 w-4 shrink-0" />
        {heading}
      </span>
      {isLoading ? (
        <p className="text-micro text-muted-foreground">Loading slots…</p>
      ) : (
        <ul className="space-y-1.5">
          {Array.from({ length: maxSlots }).map((_, i) => {
            const entry = bySlot.get(i);
            return (
              <li
                key={i}
                className="flex flex-col gap-0.5 rounded border border-border/50 bg-background/40 px-2 py-1.5"
              >
                <span className="text-micro font-tome-marginalia text-muted-foreground">Slot {i + 1}</span>
                {entry ? (
                  <span className="text-sm font-tome-marginalia text-foreground truncate" title={entry.name}>
                    {entry.name}
                  </span>
                ) : (
                  <span className="text-micro text-muted-foreground italic">Empty — set in Spell Forge</span>
                )}
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
}
