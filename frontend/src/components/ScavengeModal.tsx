import { useState } from 'react';
import { useMutation } from '@tanstack/react-query';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Loader2, Package } from 'lucide-react';
import { RaIcon } from '@/components/ui/RaIcon';
import { forageComponents } from '@/lib/api';

interface ScavengeModalProps {
  isOpen: boolean;
  onClose: () => void;
  characterId: string;
  token: string;
  className: string;
  yieldDie?: number;
  characterLevel: number;
  onSuccess: () => void;
}

function getClassConfig(className: string, yieldDie: number | undefined, level: number) {
  switch (className) {
    case 'The Ironwright':
      return {
        title: 'Scavenge Components',
        icon: <RaIcon name="regeneration" className="text-base" />,
        description: `As an action, harvest exotic components from your environment. You roll 1d${yieldDie ?? 4} to determine how many distinct components you receive.`,
        actionLabel: 'Scavenge',
      };
    case 'The Sanguinist':
      return {
        title: 'Sanguine Extraction',
        icon: <RaIcon name="water-drop" className="text-base text-red-400" />,
        description: `After using your Bite attack, extract unstable components from your prey. You gain ${level >= 11 ? 3 : level >= 5 ? 2 : 1} component${(level >= 11 ? 3 : level >= 5 ? 2 : 1) > 1 ? 's' : ''} at your current level.`,
        actionLabel: 'Extract',
      };
    case 'The Lorewright':
      return {
        title: 'Component Scavenging',
        icon: <Package className="h-5 w-5" />,
        description: `Spend 10 minutes harvesting a magical creature to extract its essence as a usable component. You gain ${level >= 7 ? '2 components' : '1 component'} at your current level.`,
        actionLabel: 'Scavenge',
      };
    default:
      return {
        title: 'Scavenge Components',
        icon: <RaIcon name="regeneration" className="text-base" />,
        description: 'Gather components from the environment.',
        actionLabel: 'Scavenge',
      };
  }
}

export function ScavengeModal({
  isOpen,
  onClose,
  characterId,
  token,
  className,
  yieldDie,
  characterLevel,
  onSuccess,
}: ScavengeModalProps) {
  const [result, setResult] = useState<{ yield: number; components: Array<{ id: string; name: string }> } | null>(null);

  const config = getClassConfig(className, yieldDie, characterLevel);

  const scavengeMutation = useMutation({
    mutationFn: () => forageComponents(characterId, token),
    onSuccess: (data) => {
      setResult(data);
      onSuccess();
    },
  });

  const handleClose = () => {
    setResult(null);
    scavengeMutation.reset();
    onClose();
  };

  const handleScavengeAgain = () => {
    setResult(null);
    scavengeMutation.reset();
    scavengeMutation.mutate();
  };

  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[440px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            {config.icon}
            {config.title}
          </DialogTitle>
          <DialogDescription>{config.description}</DialogDescription>
        </DialogHeader>

        <div className="py-4 space-y-4">
          {!result && !scavengeMutation.isPending && (
            <p className="text-sm text-muted-foreground text-center">
              {className === 'The Ironwright'
                ? `Your Scavenge Yield Die is 1d${yieldDie ?? 4}.`
                : className === 'The Sanguinist'
                ? 'Components are distinct — no repeats from your current inventory.'
                : 'You will receive new components not already in your inventory.'}
            </p>
          )}

          {scavengeMutation.isPending && (
            <div className="flex flex-col items-center gap-3 py-4">
              <Loader2 className="h-8 w-8 animate-spin text-primary" />
              <p className="text-sm text-muted-foreground">Rolling dice and selecting components...</p>
            </div>
          )}

          {scavengeMutation.isError && (
            <div className="rounded-md border border-destructive/50 bg-destructive/10 px-4 py-3 text-sm text-destructive">
              {(scavengeMutation.error as Error)?.message || 'Failed to scavenge components.'}
            </div>
          )}

          {result && (
            <div className="space-y-3">
              <div className="rounded-md border border-primary/30 bg-primary/5 px-4 py-3 text-center">
                <p className="text-2xl font-bold text-primary">{result.yield}</p>
                <p className="text-xs text-muted-foreground mt-1">
                  component{result.yield !== 1 ? 's' : ''} gained
                </p>
              </div>
              <div className="space-y-2">
                {result.components.map((comp) => (
                  <div
                    key={comp.id}
                    className="flex items-center justify-between rounded-md border border-border bg-card px-3 py-2"
                  >
                    <span className="text-sm font-medium">{comp.name}</span>
                    <Badge variant="secondary" className="text-xs">+1</Badge>
                  </div>
                ))}
              </div>
              <p className="text-xs text-muted-foreground text-center">
                Components added to your inventory. Use them in the Spell Forge to craft spells.
              </p>
            </div>
          )}
        </div>

        <DialogFooter className="gap-2">
          <Button variant="outline" onClick={handleClose}>
            Close
          </Button>
          {result ? (
            <Button onClick={handleScavengeAgain} disabled={scavengeMutation.isPending}>
              {config.actionLabel} Again
            </Button>
          ) : (
            <Button
              onClick={() => scavengeMutation.mutate()}
              disabled={scavengeMutation.isPending}
            >
              {scavengeMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {config.actionLabel}
            </Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
