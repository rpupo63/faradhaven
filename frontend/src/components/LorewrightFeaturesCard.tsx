import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Separator } from '@/components/ui/separator';
import { Skull, Loader2 } from 'lucide-react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from '@/components/ui/use-toast';
import { rollMadness, forageComponents } from '@/lib/api';
import { NormalizedCharacterSheet } from '@/types/game/character';

interface LorewrightFeaturesCardProps {
  sheet: NormalizedCharacterSheet;
  token: string;
}

export const LorewrightFeaturesCard: React.FC<LorewrightFeaturesCardProps> = ({ sheet, token }) => {
  const queryClient = useQueryClient();
  const [strainCR, setStrainCR] = useState('');

  const scavengeMutation = useMutation({
    mutationFn: () => forageComponents(sheet.character.id, token),
    onSuccess: (data) => {
      const names = data.components.map((c: { id: string; name: string }) => c.name).join(', ');
      toast({
        title: 'Components Scavenged',
        description: `Gained ${data.yield} component(s): ${names}`,
      });
      queryClient.invalidateQueries({ queryKey: ['character-sheet', sheet.character.id] });
    },
    onError: (error: Error) => {
      toast({
        title: 'Scavenge Failed',
        description: error.message || 'Could not scavenge components.',
        variant: 'destructive',
      });
    },
  });

  const madnessRollMutation = useMutation({
    mutationFn: () => rollMadness(sheet.character.id, token),
    onSuccess: (data) => {
      toast({
        title: 'Madness Roll Result',
        description: `Rolled a ${data.roll} on a d${data.madness_die}. ${data.effect ? `Effect: ${data.effect}` : 'No effect.'}`,
        variant: data.effect ? 'destructive' : 'default',
      });
      queryClient.invalidateQueries({ queryKey: ['character-sheet', sheet.character.id] });
    },
    onError: (error) => {
      toast({
        title: 'Madness Roll Failed',
        description: error.message || 'Could not roll for madness.',
        variant: 'destructive',
      });
    },
  });

  const harvestedSkills = sheet.harvested_abilities?.skills || [];
  const harvestedAttacks = sheet.harvested_abilities?.attacks || [];
  const harvestedRecipes = sheet.harvested_abilities?.recipes || [];

  const hasHarvestedAbilities =
    harvestedSkills.length > 0 ||
    harvestedAttacks.length > 0 ||
    harvestedRecipes.length > 0;

  return (
    <div className="space-y-4">
      {/* Harvest Bank Summary */}
      <div>
        <h3 className="text-md font-semibold mb-2">Harvest Bank</h3>
        {!hasHarvestedAbilities ? (
          <p className="text-muted-foreground text-sm">No abilities harvested yet.</p>
        ) : (
          <div className="space-y-2 text-sm">
            {harvestedSkills.length > 0 && (
              <p><strong>Skills:</strong> {harvestedSkills.join(', ')}</p>
            )}
            {harvestedAttacks.length > 0 && (
              <p><strong>Attacks:</strong> {harvestedAttacks.join(', ')}</p>
            )}
            {harvestedRecipes.length > 0 && (
              <p><strong>Recipes:</strong> {harvestedRecipes.join(', ')}</p>
            )}
          </div>
        )}
      </div>

      <Separator />

      {/* Trauma */}
      <div>
        <h3 className="text-md font-semibold mb-2">Trauma Level</h3>
        <div className="flex items-center gap-2">
          <Skull className="h-5 w-5 text-red-500" />
          <span className="text-lg font-bold">{sheet.trauma ?? 0}</span>
          <p className="text-sm text-muted-foreground">
            (Effects: {sheet.trauma === 1 && 'Disadvantage on Charisma checks.'}
            {sheet.trauma === 2 && 'Temporary flaw.'}
            {(sheet.trauma ?? 0) >= 3 && 'Severe Disorientation (Random movement and action).'}
            {(!sheet.trauma || sheet.trauma === 0) && 'None'}
            )
          </p>
        </div>
      </div>

      <Separator />

      {/* Psychic Strain DC Calculator */}
      <div>
        <h3 className="text-md font-semibold mb-2">Psychic Strain DC</h3>
        <p className="text-xs text-muted-foreground mb-2">
          Safe threshold: CR ≤ level/3. Above that, a Wisdom save is required.
        </p>
        <div className="flex gap-2 items-center">
          <Input
            type="number"
            placeholder="Creature CR"
            value={strainCR}
            onChange={e => setStrainCR(e.target.value)}
            className="h-8 text-sm w-32"
            step="0.25"
            min="0"
          />
          {(() => {
            const cr = parseFloat(strainCR);
            const level = sheet.character.level;
            if (isNaN(cr) || strainCR === '') return null;
            const safeThreshold = level / 3;
            if (cr <= safeThreshold) {
              return <p className="text-sm text-green-400">No save required</p>;
            }
            const excess = cr - safeThreshold;
            const dc = 10 + Math.ceil(excess);
            return <p className="text-sm text-red-400">Save required — DC {dc}</p>;
          })()}
        </div>
      </div>

      <Separator />

      {/* Madness Roll */}
      <div>
        <h3 className="text-md font-semibold mb-2">Madness Roll</h3>
        <Button
          onClick={() => madnessRollMutation.mutate()}
          disabled={madnessRollMutation.isPending}
          className="w-full"
        >
          {madnessRollMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
          Roll Madness Die
        </Button>
        <p className="text-sm text-muted-foreground mt-2">
          (Triggers when using a harvested ability or passively.)
        </p>
      </div>

      <Separator />

      {/* Component Scavenging */}
      <div>
        <h3 className="text-md font-semibold mb-2">Component Scavenging</h3>
        <Button
          onClick={() => scavengeMutation.mutate()}
          disabled={scavengeMutation.isPending}
          className="w-full"
          variant="outline"
        >
          {scavengeMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
          Scavenge Components
        </Button>
        <p className="text-sm text-muted-foreground mt-2">
          Extracts essence from magical creatures (1 component L1-6, 2 at L7+).
        </p>
      </div>
    </div>
  );
};
