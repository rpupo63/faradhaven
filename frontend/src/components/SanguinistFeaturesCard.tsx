import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import { updateSanguineNotoriety, forageComponents } from '@/lib/api';
import type { NormalizedCharacterSheet } from '@/types/game';
import { ArrowDown, ArrowUp, Loader2 } from 'lucide-react';
import { toast } from 'sonner';

interface SanguinistFeaturesCardProps {
  sheet: NormalizedCharacterSheet;
  token: string;
}

export function SanguinistFeaturesCard({ sheet, token }: SanguinistFeaturesCardProps) {
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: ({ mpChange, brChange }: { mpChange: number; brChange: number }) =>
      updateSanguineNotoriety(sheet.character.id, { mpChange, brChange }, token),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['character-sheet', sheet.character.id] });
    },
    onError: (error) => {
      console.error("Failed to update notoriety:", error);
      alert("Failed to update notoriety. See console for details.");
    }
  });

  const forageMutation = useMutation({
    mutationFn: () => forageComponents(sheet.character.id, token),
    onSuccess: (data) => {
      const names = data.components.map((c) => c.name).join(', ');
      toast.success(`Extracted ${data.yield} component(s): ${names}`);
      queryClient.invalidateQueries({ queryKey: ['character-sheet', sheet.character.id] });
    },
    onError: (error: Error) => {
      toast.error(error.message || 'Failed to extract components');
    },
  });

  const handleUpdate = (mpChange: number, brChange: number) => {
    mutation.mutate({ mpChange, brChange });
  };

  const sanguineMp = sheet.character.sanguine_mp ?? 0;
  const sanguineBr = sheet.character.sanguine_br ?? 0;
  const notoriety = sanguineMp - sanguineBr;

  const ichorResource = sheet.class_resources?.find(r => r.key === 'max_blood_ichor');
  const ichorBelowHalf =
    ichorResource &&
    ichorResource.max_value != null &&
    ichorResource.max_value > 0 &&
    (ichorResource.current_value ?? 0) < ichorResource.max_value * 0.5;

  return (
    <div className="space-y-4">
      <div className="flex justify-around items-center text-center">
        <div>
          <h3 className="font-bold text-lg text-blue-400">Medical Prodigy (MP)</h3>
          <p className="text-2xl font-bold">{sanguineMp}</p>
          <div className="flex gap-2 mt-2">
            <Button size="sm" onClick={() => handleUpdate(1, 0)} disabled={mutation.isPending}>
              <ArrowUp className="h-4 w-4 mr-1" /> +1 MP
            </Button>
            <Button size="sm" variant="outline" onClick={() => handleUpdate(-1, 0)} disabled={mutation.isPending || sanguineMp <= 0}>
              <ArrowDown className="h-4 w-4 mr-1" /> -1 MP
            </Button>
          </div>
        </div>
        <div>
          <h3 className="font-bold text-lg text-red-400">Blood Rage (BR)</h3>
          <p className="text-2xl font-bold">{sanguineBr}</p>
          <div className="flex gap-2 mt-2">
            <Button size="sm" onClick={() => handleUpdate(0, 1)} disabled={mutation.isPending}>
              <ArrowUp className="h-4 w-4 mr-1" /> +1 BR
            </Button>
            <Button size="sm" variant="outline" onClick={() => handleUpdate(0, -1)} disabled={mutation.isPending || sanguineBr <= 0}>
              <ArrowDown className="h-4 w-4 mr-1" /> -1 BR
            </Button>
          </div>
        </div>
      </div>
      <div className="text-center pt-4 border-t">
        <h3 className="font-bold text-lg">Net Notoriety</h3>
        <p className="text-3xl font-bold">{notoriety}</p>
        <p className="text-sm text-muted-foreground">
          {notoriety >= 3 ? "Overloaded Healer" : notoriety <= -3 ? "Starving Predator" : "Balanced"}
        </p>
      </div>
      <div className="text-xs text-muted-foreground p-3 bg-secondary rounded-md">
        <p className="font-bold">Notoriety Effects:</p>
        <ul className="list-disc list-inside">
          <li><span className="font-semibold">Overloaded Healer (MP is 3+ higher than BR):</span> Damage-dealing class features may trigger 'Sanguine Backfires'.</li>
          <li><span className="font-semibold">Starving Predator (BR is 3+ higher than MP):</span> Blood Graft becomes 'Ravenous'.</li>
        </ul>
      </div>
      <div className="pt-2 border-t">
        {ichorBelowHalf && (
          <div className="mb-2 rounded border border-yellow-500/50 bg-yellow-500/10 px-3 py-2 text-xs text-yellow-400">
            ⚠ Exhaustion risk: DC 13 Con save if you spend &gt;50% Ichor per turn
          </div>
        )}
        <Button
          className="w-full"
          variant="outline"
          onClick={() => forageMutation.mutate()}
          disabled={forageMutation.isPending}
        >
          {forageMutation.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
          Extract Components
        </Button>
        <p className="text-xs text-muted-foreground mt-1 text-center">
          Sanguine Extraction — gains exotic components (1/2/3 at L1/L5/L11)
        </p>
      </div>
    </div>
  );
}