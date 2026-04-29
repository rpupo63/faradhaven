import { useMemo, useState } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Layers, Wand2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { cn } from '@/lib/utils';
import type { NormalizedCharacterSheet } from '@/types/game';
import type { ApiSavedSpell } from '@/types/game/api';
import { clearSpeedDialSlot, getSpeedDial } from '@/lib/api/character';
import { toast } from 'sonner';

function resolveComponentName(componentId: string, sheet: NormalizedCharacterSheet): string {
  const pool = sheet.available_components?.find((c) => c.id === componentId);
  if (pool?.name) return pool.name;
  const inv = sheet.components?.find((cc) => cc.component_id === componentId)?.component;
  if (inv?.name) return inv.name;
  return 'Unknown';
}

interface PowderMageSpeedDialMenuProps {
  sheet: NormalizedCharacterSheet;
  characterId: string;
  token: string;
  onGoToSpellForge?: () => void;
  className?: string;
  /** Piston Brawler shares the same saved-slot API (blueprints). */
  variant?: 'speed_dial' | 'blueprint';
}

export function PowderMageSpeedDialMenu({
  sheet,
  characterId,
  token,
  onGoToSpellForge,
  className,
  variant = 'speed_dial',
}: PowderMageSpeedDialMenuProps) {
  const isBlueprint = variant === 'blueprint';
  const [open, setOpen] = useState(false);
  const queryClient = useQueryClient();

  const maxSlots = sheet.class_resources?.find((r) => r.key === 'speed_dial_slots')?.value ?? 0;
  const sdResource = sheet.class_resources?.find((r) => r.key === 'speed_dial_slots');
  const slotsRemaining =
    sdResource?.is_trackable && sdResource.current_value != null
      ? sdResource.current_value
      : maxSlots;

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

  const clearMutation = useMutation({
    mutationFn: (slotIndex: number) => clearSpeedDialSlot(characterId, slotIndex, token),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['speed-dial', characterId] });
      toast.success(isBlueprint ? 'Blueprint slot cleared' : 'Speed Dial slot cleared');
    },
    onError: (e: Error) => toast.error(e.message || 'Failed to clear slot'),
  });

  const goForge = () => {
    setOpen(false);
    onGoToSpellForge?.();
  };

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          type="button"
          variant="secondary"
          className={cn('w-full gap-2 font-tome-marginalia', className)}
          onClick={(e) => e.stopPropagation()}
        >
          <Layers className="h-4 w-4 shrink-0" />
          <span>{isBlueprint ? 'Blueprint slots' : 'Speed Dial'}</span>
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-[min(22rem,calc(100vw-2rem))] space-y-3" align="start">
        <div>
          <p className="text-sm font-tome-subheading text-primary">
            {isBlueprint ? 'Blueprint pre-mix' : 'Pre-mix (Speed Dial)'}
          </p>
          <p className="text-micro font-tome-marginalia text-muted-foreground mt-1 leading-snug">
            {isBlueprint ? (
              <>
                Pre-calculate combinations during rests (level 13+). Save component strings in Spell Forge.
                Slot uses restore on <span className="text-foreground/90">short or long rest</span>.
              </>
            ) : (
              <>
                Save up to 3 components per slot in Spell Forge. Slot uses restore on{' '}
                <span className="text-foreground/90">short or long rest</span> (see Resources). You can
                reassign your strings anytime in the forge.
              </>
            )}
          </p>
        </div>

        {maxSlots <= 0 ? (
          <div className="space-y-2">
            <p className="text-sm text-muted-foreground italic">
              {isBlueprint
                ? 'Blueprint slots unlock at higher level (see class progression).'
                : 'Speed Dial (Pre-mix) unlocks at level 3.'}
            </p>
            <Button type="button" variant="outline" className="w-full gap-2" onClick={goForge}>
              <Wand2 className="h-4 w-4" />
              Open Spell Forge
            </Button>
          </div>
        ) : (
          <>
            <p className="text-micro font-tome-marginalia text-muted-foreground">
              Uses remaining:{' '}
              <span className="font-display text-foreground tabular-nums">
                {slotsRemaining} / {maxSlots}
              </span>
            </p>

            {isLoading ? (
              <p className="text-sm text-muted-foreground">Loading…</p>
            ) : (
              <ul className="space-y-2 max-h-[min(50vh,18rem)] overflow-y-auto pr-1">
                {Array.from({ length: maxSlots }).map((_, i) => {
                  const entry = bySlot.get(i);
                  return (
                    <li
                      key={i}
                      className="rounded-md border border-border/70 bg-card/60 px-2.5 py-2 text-sm"
                    >
                      <div className="flex items-start justify-between gap-2">
                        <span className="text-micro font-tome-marginalia text-muted-foreground shrink-0">
                          Slot {i + 1}
                        </span>
                        {entry && (
                          <Button
                            type="button"
                            variant="ghost"
                            size="sm"
                            className="h-7 text-micro px-2 shrink-0"
                            disabled={clearMutation.isPending}
                            onClick={() => clearMutation.mutate(i)}
                          >
                            Clear
                          </Button>
                        )}
                      </div>
                      {entry ? (
                        <>
                          <p className="font-tome-subheading text-primary mt-0.5">{entry.name}</p>
                          <p className="text-micro text-muted-foreground mt-1 leading-snug">
                            {entry.component_ids
                              .map((id) => resolveComponentName(id, sheet))
                              .join(' → ')}
                          </p>
                        </>
                      ) : (
                        <p className="text-muted-foreground italic mt-1">No pre-mix saved yet.</p>
                      )}
                    </li>
                  );
                })}
              </ul>
            )}

            <Button type="button" variant="default" className="w-full gap-2" onClick={goForge}>
              <Wand2 className="h-4 w-4" />
              {saved.length > 0 ? 'Change in Spell Forge' : 'Add in Spell Forge'}
            </Button>
          </>
        )}
      </PopoverContent>
    </Popover>
  );
}
