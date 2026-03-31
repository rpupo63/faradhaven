import { Star, Hash, Box } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { NormalizedCharacterSheet } from '@/types/game';
import { cn } from '@/lib/utils';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { callTraitAbility, restoreTraitAbility } from '@/lib/api/abilities';
import { spendResource, gainResource } from '@/lib/api/resources';

interface RacialResourceTrackerProps {
  sheet: NormalizedCharacterSheet;
  characterId: string;
  token: string;
}

function TraitResourceRow({
  traitName,
  traitId,
  current,
  max,
  resetCondition,
  characterId,
  token,
}: {
  traitName: string;
  traitId: string;
  current: number;
  max: number;
  resetCondition?: string;
  characterId: string;
  token: string;
}) {
  const queryClient = useQueryClient();

  const spendMutation = useMutation({
    mutationFn: () => callTraitAbility(characterId, traitId, token),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] }),
    onError: (error) => console.error('Failed to use trait:', error),
  });

  const gainMutation = useMutation({
    mutationFn: () => restoreTraitAbility(characterId, traitId, token),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] }),
    onError: (error) => console.error('Failed to restore trait:', error),
  });

  const isPending = spendMutation.isPending || gainMutation.isPending;

  return (
    <div className="flex flex-col items-center gap-1 min-w-[100px]">
      <p className="text-[10px] font-tome-marginalia text-muted-foreground uppercase text-center leading-tight">
        {traitName}
      </p>
      
      {max === 1 ? (
        <button
          onClick={() => current > 0 ? spendMutation.mutate() : gainMutation.mutate()}
          disabled={isPending}
          className={cn(
            'px-2 py-0.5 rounded text-[10px] font-display transition-colors cursor-pointer border',
            current > 0
              ? 'bg-primary/20 text-primary border-primary/50 hover:bg-primary/30'
              : 'bg-muted text-muted-foreground border-border hover:bg-muted/80'
          )}
        >
          {current > 0 ? 'AVAILABLE' : 'SPENT'}
        </button>
      ) : (
        <div className="flex flex-wrap gap-1 items-center justify-center">
          {Array.from({ length: max }).map((_, i) => {
            const isAvailable = i < current;
            return (
              <button
                key={i}
                onClick={() => isAvailable ? spendMutation.mutate() : gainMutation.mutate()}
                disabled={isPending}
                className={cn(
                  'h-4 w-4 rounded-full transition-colors cursor-pointer',
                  isAvailable
                    ? 'bg-primary hover:bg-primary/80'
                    : 'border border-border hover:border-primary/50'
                )}
                title={isAvailable ? 'Click to use' : 'Click to restore'}
              />
            );
          })}
        </div>
      )}

      {resetCondition && (
        <Badge variant="outline" className="text-[8px] px-1 py-0 h-3 font-tome-marginalia opacity-70">
          {resetCondition}
        </Badge>
      )}
    </div>
  );
}

/** Custom tracker for Rock Gnome Clockwork Devices (Max 3) */
function ClockworkDevicesTracker({
  characterId,
  token,
  currentValue = 0,
}: {
  characterId: string;
  token: string;
  currentValue?: number;
}) {
  const queryClient = useQueryClient();
  const resourceKey = 'resource_clockwork_devices';

  const spendMutation = useMutation({
    mutationFn: () => spendResource(characterId, resourceKey, 1, token),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] }),
  });

  const gainMutation = useMutation({
    mutationFn: () => gainResource(characterId, resourceKey, 1, token),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] }),
  });

  const isPending = spendMutation.isPending || gainMutation.isPending;

  return (
    <div className="flex flex-col items-center gap-1 min-w-[100px]">
      <p className="text-[10px] font-tome-marginalia text-muted-foreground uppercase text-center leading-tight">
        Clockwork Devices
      </p>
      <div className="flex items-center gap-2">
        <Button
          size="sm"
          variant="outline"
          className="h-5 w-5 p-0"
          onClick={() => spendMutation.mutate()}
          disabled={isPending || currentValue <= 0}
        >
          -
        </Button>
        <span className="text-sm font-display text-primary">{currentValue} / 3</span>
        <Button
          size="sm"
          variant="outline"
          className="h-5 w-5 p-0"
          onClick={() => gainMutation.mutate()}
          disabled={isPending || currentValue >= 3}
        >
          +
        </Button>
      </div>
    </div>
  );
}

export function RacialResourceTracker({ sheet, characterId, token }: RacialResourceTrackerProps) {
  const traitMaxUses = sheet.trait_max_uses ?? {};
  const traitUseStates = sheet.trait_use_states ?? {};
  
  // Also check for the custom Rock Gnome clockwork device resource
  const isRockGnome = (sheet.character.raceName === 'Gnome' || sheet.character.lineageName === 'Rock Gnome') && 
                     sheet.race_traits?.some(t => t.name === 'Gnomish Lineage');
  
  // We'll store this in class_resources or create a dummy one if not exists
  const clockworkResource = sheet.class_resources?.find(r => r.key === 'resource_clockwork_devices');

  const hasTrackedTraits = Object.keys(traitMaxUses).length > 0;
  const showTracker = hasTrackedTraits || isRockGnome;

  if (!showTracker) return null;

  return (
    <Card className="arcane-border">
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-tome-marginalia uppercase tracking-wider flex items-center gap-2">
          <Star className="h-4 w-4 text-primary" />
          Racial Resources
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex flex-wrap gap-4 justify-center">
          {Object.entries(traitMaxUses).map(([traitId, max]) => {
            const trait = sheet.race_traits?.find(t => t.id === traitId);
            if (!trait) return null;
            
            // Special case: Rock Gnomes don't use the Gnomish Lineage PB resource (they use devices)
            if (trait.name === 'Gnomish Lineage' && isRockGnome) return null;
            
            return (
              <TraitResourceRow
                key={traitId}
                traitId={traitId}
                traitName={trait.name}
                current={traitUseStates[traitId] ?? 0}
                max={max}
                resetCondition={trait.reset_condition}
                characterId={characterId}
                token={token}
              />
            );
          })}
          
          {isRockGnome && (
            <ClockworkDevicesTracker
              characterId={characterId}
              token={token}
              currentValue={clockworkResource?.current_value ?? 0}
            />
          )}
        </div>
      </CardContent>
    </Card>
  );
}
