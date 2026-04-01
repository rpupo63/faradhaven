import { useState } from 'react';
import { Swords, ChevronRight, RotateCcw } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useAuth } from '@/context/AuthContext';
import { getInitiative, clearInitiative } from '@/lib/api/map';
import { getMinions } from '@/lib/api';
import type { MapToken } from '@/types/map';
import type { ApiMinion } from '@/types/game';

interface InitiativeTrackerProps {
  mapId: string;
  isDM: boolean;
}

interface InitiativeSlot {
  type: 'token' | 'minion';
  token?: MapToken;
  minion?: ApiMinion;
  label: string;
  characterId?: string;
}

export function InitiativeTracker({ mapId, isDM }: InitiativeTrackerProps) {
  const { token } = useAuth();
  const queryClient = useQueryClient();
  const [currentTurnIndex, setCurrentTurnIndex] = useState(0);

  // Fetch initiative order
  const { data: initiativeTokens = [] } = useQuery({
    queryKey: ['initiative', mapId],
    queryFn: () => getInitiative(mapId, token || undefined),
    enabled: !!mapId && !!token,
  });

  // Collect character IDs from PC tokens to fetch their minions
  const pcCharacterIds = initiativeTokens
    .filter((t) => t.character_id)
    .map((t) => t.character_id!);

  // Fetch minions for each PC (active constructs + drones)
  const { data: allMinions = [] } = useQuery({
    queryKey: ['initiative-minions', pcCharacterIds],
    queryFn: async () => {
      const results: ApiMinion[] = [];
      for (const charId of pcCharacterIds) {
        const minions = await getMinions(charId, undefined, token || undefined);
        results.push(...minions.filter((m) => m.is_active && (m.minion_type === 'construct' || m.minion_type === 'drone')));
      }
      return results;
    },
    enabled: pcCharacterIds.length > 0 && !!token,
  });

  // Build interleaved initiative list: PC token → their minions → next token
  const slots: InitiativeSlot[] = [];
  for (const t of initiativeTokens) {
    slots.push({
      type: 'token',
      token: t,
      label: t.name,
      characterId: t.character_id || undefined,
    });
    // Insert minions belonging to this character right after
    if (t.character_id) {
      const charMinions = allMinions.filter((m) => m.character_id === t.character_id);
      for (const m of charMinions) {
        slots.push({
          type: 'minion',
          minion: m,
          label: `  ${m.name}`,
          characterId: t.character_id,
        });
      }
    }
  }

  const clearMutation = useMutation({
    mutationFn: () => clearInitiative(mapId, token!),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['initiative', mapId] });
      setCurrentTurnIndex(0);
    },
  });

  const handleNextTurn = () => {
    if (slots.length === 0) return;
    setCurrentTurnIndex((prev) => (prev + 1) % slots.length);
  };

  if (initiativeTokens.length === 0) {
    return null;
  }

  return (
    <Card className="arcane-border bg-card">
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center justify-between text-base font-tome-subheading text-primary">
          <div className="flex items-center gap-2">
            <Swords className="h-4 w-4" />
            Initiative
          </div>
          <div className="flex items-center gap-1">
            <Button
              size="icon"
              variant="ghost"
              className="h-6 w-6"
              onClick={handleNextTurn}
              title="Next turn"
            >
              <ChevronRight className="h-4 w-4" />
            </Button>
            {isDM && (
              <Button
                size="icon"
                variant="ghost"
                className="h-6 w-6 text-red-400"
                onClick={() => clearMutation.mutate()}
                disabled={clearMutation.isPending}
                title="Clear initiative"
              >
                <RotateCcw className="h-3 w-3" />
              </Button>
            )}
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-1">
        {slots.map((slot, idx) => {
          const isActive = idx === currentTurnIndex;
          return (
            <div
              key={slot.type === 'token' ? `t-${slot.token?.id}` : `m-${slot.minion?.id}`}
              className={`flex items-center justify-between px-2 py-1 rounded text-sm ${
                isActive
                  ? 'bg-primary/20 border border-primary/50 text-primary font-semibold'
                  : 'text-muted-foreground'
              }`}
            >
              <span className={slot.type === 'minion' ? 'pl-4 text-xs italic' : ''}>
                {slot.label}
              </span>
              <div className="flex items-center gap-1">
                {slot.type === 'token' && slot.token?.initiative_order != null && (
                  <Badge variant="outline" size="sm">
                    {slot.token.initiative_order}
                  </Badge>
                )}
                {slot.type === 'minion' && slot.minion && (
                  <span className="text-micro text-muted-foreground">
                    HP {slot.minion.current_hp}/{slot.minion.max_hp}
                  </span>
                )}
                {isActive && <ChevronRight className="h-3 w-3 text-primary" />}
              </div>
            </div>
          );
        })}
      </CardContent>
    </Card>
  );
}
