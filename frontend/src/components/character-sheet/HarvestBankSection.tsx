import { useState } from 'react';
import {
  GraduationCap,
  Sword,
  BookOpen,
  Plus,
  Trash2,
  Loader2,
} from 'lucide-react';
import { Separator } from '@/components/ui/separator';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

import { NormalizedCharacterSheet } from '@/types/game';
import type { ApiHarvestMinion } from '@/types/game';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { getMinions, storeHarvest, deleteMinion } from '@/lib/api/minion';
import { toast } from 'sonner';
import { useAuth } from '@/context/AuthContext';
import { Textarea } from '../ui/textarea';

interface HarvestBankSectionProps {
  sheet: NormalizedCharacterSheet;
}

const ABILITY_TYPE_ICONS: Record<string, React.ReactNode> = {
  skill: <GraduationCap className="h-3 w-3 text-blue-400" />,
  attack: <Sword className="h-3 w-3 text-red-400" />,
  action: <Sword className="h-3 w-3 text-red-400" />,
  recipe: <BookOpen className="h-3 w-3 text-purple-400" />,
};

const ABILITY_TYPE_COLORS: Record<string, string> = {
  skill: 'border-blue-500/30 bg-blue-500/10',
  attack: 'border-red-500/30 bg-red-500/10',
  action: 'border-red-500/30 bg-red-500/10',
  recipe: 'border-purple-500/30 bg-purple-500/10',
};

export function HarvestBankSection({ sheet }: HarvestBankSectionProps) {
  const { token } = useAuth();
  const queryClient = useQueryClient();
  const [dialogOpen, setDialogOpen] = useState(false);
  const [formType, setFormType] = useState<string>('skill');
  const [formName, setFormName] = useState('');
  const [formDescription, setFormDescription] = useState('');
  const [formSource, setFormSource] = useState('');

  const { data: harvests = [] } = useQuery({
    queryKey: ['minions', sheet.character.id, 'harvest'],
    queryFn: () => getMinions(sheet.character.id, 'harvest', token || undefined) as Promise<ApiHarvestMinion[]>,
  });

  const storeMutation = useMutation({
    mutationFn: () =>
      storeHarvest(
        sheet.character.id,
        {
          ability_name: formName,
          ability_description: formDescription,
          ability_type: formType,
          source_creature_type: formSource,
        },
        token || undefined
      ),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['minions', sheet.character.id, 'harvest'] });
      queryClient.invalidateQueries({ queryKey: ['character-sheet', sheet.character.id] });
      toast.success('Harvest stored');
      setDialogOpen(false);
      setFormName('');
      setFormDescription('');
      setFormSource('');
      setFormType('skill');
    },
    onError: (err: Error) => toast.error(err.message || 'Failed to store harvest'),
  });

  const removeMutation = useMutation({
    mutationFn: (minionId: string) =>
      deleteMinion(sheet.character.id, minionId, false, token || undefined),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['minions', sheet.character.id, 'harvest'] });
      queryClient.invalidateQueries({ queryKey: ['character-sheet', sheet.character.id] });
      toast.success('Harvest removed');
    },
    onError: (err: Error) => toast.error(err.message || 'Failed to remove harvest'),
  });

  // Group harvests by type
  const grouped: Record<string, ApiHarvestMinion[]> = {};
  for (const h of harvests) {
    const type = h.ability_type || 'skill';
    if (!grouped[type]) grouped[type] = [];
    grouped[type].push(h);
  }

  const typeOrder = ['skill', 'attack', 'action', 'recipe'];

  return (
    <>
      <Separator />
      <div className="flex items-center justify-between">
        <span className="flex items-center gap-2 text-sm font-tome-subheading text-primary">
          <GraduationCap className="h-4 w-4" />
          Harvest Bank
        </span>
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button size="sm" variant="outline" className="h-7 text-xs">
              <Plus className="h-3 w-3 mr-1" />
              Add Harvest
            </Button>
          </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Store Harvested Ability</DialogTitle>
              </DialogHeader>
              <div className="space-y-4 pt-2">
                <div className="space-y-2">
                  <Label>Type</Label>
                  <Select value={formType} onValueChange={setFormType}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="skill">Skill Proficiency</SelectItem>
                      <SelectItem value="attack">Attack / Action</SelectItem>
                      <SelectItem value="recipe">Spell Recipe</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label>Ability Name</Label>
                  <Input
                    value={formName}
                    onChange={(e) => setFormName(e.target.value)}
                    placeholder="e.g., Keen Hearing, Claw Attack, Fire Breath Recipe"
                  />
                </div>
                <div className="space-y-2">
                  <Label>Description</Label>
                  <Textarea
                    value={formDescription}
                    onChange={(e) => setFormDescription(e.target.value)}
                    placeholder="Describe the harvested ability..."
                    rows={3}
                  />
                </div>
                <div className="space-y-2">
                  <Label>Source Creature</Label>
                  <Input
                    value={formSource}
                    onChange={(e) => setFormSource(e.target.value)}
                    placeholder="e.g., Dire Wolf, Goblin Shaman"
                  />
                </div>
                <Button
                  className="w-full"
                  onClick={() => storeMutation.mutate()}
                  disabled={!formName || !formType || storeMutation.isPending}
                >
                  {storeMutation.isPending ? (
                    <Loader2 className="h-4 w-4 animate-spin mr-2" />
                  ) : (
                    <Plus className="h-4 w-4 mr-2" />
                  )}
                  Store Harvest
                </Button>
              </div>
            </DialogContent>
          </Dialog>
        </div>
      <div>
        {harvests.length === 0 ? (
          <p className="text-xs text-muted-foreground italic font-tome-marginalia text-center py-2">
            No harvested abilities yet. Perform Visceral Psychometry to harvest from creatures.
          </p>
        ) : (
          <div className="space-y-3">
            {typeOrder.map((type) => {
              const items = grouped[type];
              if (!items || items.length === 0) return null;
              const typeLabel = type === 'action' ? 'Attacks' : type.charAt(0).toUpperCase() + type.slice(1) + 's';
              return (
                <div key={type}>
                  <div className="flex items-center gap-1 mb-1">
                    {ABILITY_TYPE_ICONS[type]}
                    <span className="text-micro font-tome-marginalia text-muted-foreground uppercase">
                      {typeLabel}
                    </span>
                  </div>
                  <div className="space-y-1">
                    {items.map((h) => (
                      <div
                        key={h.id}
                        className={`flex items-center justify-between px-2 py-1.5 rounded border text-xs ${ABILITY_TYPE_COLORS[type] || 'border-border'}`}
                      >
                        <div className="flex-1 min-w-0">
                          <span className="font-medium">{h.ability_name || h.name}</span>
                          {h.source_creature_type && (
                            <Badge variant="secondary" size="sm" className="ml-2">
                              {h.source_creature_type}
                            </Badge>
                          )}
                        </div>
                        <Button
                          size="icon"
                          variant="ghost"
                          className="h-5 w-5 shrink-0"
                          onClick={() => removeMutation.mutate(h.id)}
                          disabled={removeMutation.isPending}
                        >
                          {removeMutation.isPending && removeMutation.variables === h.id ? (
                            <Loader2 className="h-3 w-3 animate-spin" />
                          ) : (
                            <Trash2 className="h-3 w-3" />
                          )}
                        </Button>
                      </div>
                    ))}
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </div>
    </>
  );
}


