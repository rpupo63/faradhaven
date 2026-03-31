import React, { useState, useEffect } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Label } from '@/components/ui/label';
import { Loader2 } from 'lucide-react';
import { toast } from 'sonner';
import { generateLoot, type LootSource, type LootTier, type LootResult } from '@/lib/api/loot';

interface LootModalProps {
  isOpen: boolean;
  onClose: () => void;
  characterId: string;
  token: string;
}

const SOURCE_OPTIONS: { value: LootSource; label: string }[] = [
  { value: 'common_enemy', label: 'Common Enemy' },
  { value: 'boss_enemy', label: 'Boss Enemy' },
  { value: 'room', label: 'Room' },
];

const TIER_OPTIONS: { value: LootTier; label: string }[] = [
  { value: 'low', label: 'Low' },
  { value: 'medium', label: 'Medium' },
  { value: 'high', label: 'High' },
];

function formatCopper(cp: number): string {
  if (cp >= 100) {
    const gp = Math.floor(cp / 100);
    const remainder = cp % 100;
    return remainder > 0 ? `${gp} gp ${remainder} cp` : `${gp} gp`;
  }
  if (cp >= 10) {
    const sp = Math.floor(cp / 10);
    const remainder = cp % 10;
    return remainder > 0 ? `${sp} sp ${remainder} cp` : `${sp} sp`;
  }
  return `${cp} cp`;
}

export const LootModal: React.FC<LootModalProps> = ({ isOpen, onClose, characterId, token }) => {
  const queryClient = useQueryClient();
  const [source, setSource] = useState<LootSource>('common_enemy');
  const [tier, setTier] = useState<LootTier>('medium');
  const [result, setResult] = useState<LootResult | null>(null);

  const lootMutation = useMutation<LootResult, Error>({
    mutationFn: () => generateLoot(characterId, source, tier, token),
    onSuccess: (data) => {
      setResult(data);
      queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] });
      toast.success(data.message);
    },
    onError: (error) => {
      toast.error(error.message || 'Failed to generate loot');
    },
  });

  useEffect(() => {
    if (!isOpen) {
      setTimeout(() => setResult(null), 0);
    }
  }, [isOpen]);

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Generate Loot</DialogTitle>
          <DialogDescription>Roll for rewards from encounters and exploration.</DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label>Source</Label>
            <Select value={source} onValueChange={(v) => setSource(v as LootSource)}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {SOURCE_OPTIONS.map((o) => (
                  <SelectItem key={o.value} value={o.value}>
                    {o.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div className="grid gap-2">
            <Label>Tier</Label>
            <Select value={tier} onValueChange={(v) => setTier(v as LootTier)}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {TIER_OPTIONS.map((o) => (
                  <SelectItem key={o.value} value={o.value}>
                    {o.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {result && (
            <div className="rounded-md border p-3 space-y-2 text-sm">
              <p className="font-medium">
                Gold earned: {formatCopper(result.gold_earned)}
              </p>
              <p className="text-muted-foreground">
                Item drops rolled: {result.items_rolled ?? result.items.length} · Weapon drops rolled:{' '}
                {result.weapons_rolled ?? result.weapons.length}
              </p>
              {result.drops && result.drops.length > 0 ? (
                <div>
                  <p className="font-medium">Discovered (in order)</p>
                  <ol className="list-decimal list-inside space-y-1">
                    {result.drops.map((drop, i) => (
                      <li key={`${drop.kind}-${drop.name}-${i}`}>
                        <span className="capitalize text-muted-foreground">{drop.kind}:</span> {drop.name}{' '}
                        <span className="text-muted-foreground">({drop.rarity})</span>
                      </li>
                    ))}
                  </ol>
                </div>
              ) : (
                <>
                  {result.items.length > 0 && (
                    <div>
                      <p className="font-medium">Items:</p>
                      <ul className="list-disc list-inside">
                        {result.items.map((item, i) => (
                          <li key={`${item.id}-${i}`}>
                            {item.name}{' '}
                            <span className="text-muted-foreground">({item.rarity})</span>
                          </li>
                        ))}
                      </ul>
                    </div>
                  )}
                  {result.weapons.length > 0 && (
                    <div>
                      <p className="font-medium">Weapons:</p>
                      <ul className="list-disc list-inside">
                        {result.weapons.map((w, i) => (
                          <li key={`${w.id}-${i}`}>
                            {w.name}{' '}
                            <span className="text-muted-foreground">({w.rarity})</span>
                          </li>
                        ))}
                      </ul>
                    </div>
                  )}
                </>
              )}
              {(!result.drops || result.drops.length === 0) &&
                result.items.length === 0 &&
                result.weapons.length === 0 && (
                  <p className="text-muted-foreground">No items or weapons found — just gold this time.</p>
                )}
            </div>
          )}
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            Close
          </Button>
          <Button onClick={() => lootMutation.mutate()} disabled={lootMutation.isPending}>
            {lootMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Roll Loot
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
