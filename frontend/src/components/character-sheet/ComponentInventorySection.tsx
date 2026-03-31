import { useMemo } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Atom, Plus, Minus, Loader2 } from 'lucide-react';
import { cn } from '@/lib/utils';
import { NormalizedCharacterSheet, ApiCharacterComponent, ComponentCategory } from '@/types/game';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { consumeComponent, gainComponent } from '@/lib/api/character';
import { categoryColors } from '../arcanum/ElementTile';
import { toast } from 'sonner';

interface ComponentInventorySectionProps {
  sheet: NormalizedCharacterSheet;
  token: string;
}

const CATEGORY_ORDER: ComponentCategory[] = [
  'Forma',
  'Scopus',
  'Essentia',
  'Actio',
  'Magnitudo',
  'Logica',
];

export function ComponentInventorySection({ sheet, token }: ComponentInventorySectionProps) {
  const queryClient = useQueryClient();
  const characterId = sheet.character.id;
  const components = sheet.components || [];

  // Filter to components with count > 0
  const ownedComponents = useMemo(() => {
    return components.filter((c) => c.count > 0);
  }, [components]);

  // Group by category
  const grouped = useMemo(() => {
    const g: Record<string, ApiCharacterComponent[]> = {};
    for (const c of ownedComponents) {
      const cat = c.component?.category || 'Other';
      if (!g[cat]) g[cat] = [];
      g[cat].push(c);
    }
    // Sort each group by name
    for (const cat in g) {
      g[cat].sort((a, b) => (a.component?.name || '').localeCompare(b.component?.name || ''));
    }
    return g;
  }, [ownedComponents]);

  const consumeMutation = useMutation({
    mutationFn: (componentId: string) => consumeComponent(characterId, componentId, token),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] }),
    onError: (err: Error) => toast.error(err.message || 'Failed to consume component'),
  });

  const gainMutation = useMutation({
    mutationFn: (componentId: string) => gainComponent(characterId, componentId, token),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['character-sheet', characterId] }),
    onError: (err: Error) => toast.error(err.message || 'Failed to add component'),
  });

  if (ownedComponents.length === 0) return null;

  return (
    <Card className="arcane-border bg-card">
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center gap-2 text-base font-tome-subheading text-primary">
          <Atom className="h-4 w-4" />
          Component Inventory
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {CATEGORY_ORDER.map((cat) => {
          const items = grouped[cat];
          if (!items || items.length === 0) return null;

          return (
            <div key={cat} className="space-y-1.5">
              <div className="flex items-center gap-2 border-b border-border pb-0.5">
                <div className={cn("w-2 h-2 rounded-full", categoryColors[cat as ComponentCategory]?.bg || 'bg-muted')} />
                <span className="text-[10px] font-tome-marginalia text-muted-foreground uppercase tracking-wider">
                  {cat}
                </span>
              </div>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                {items.map((c) => (
                  <div
                    key={c.component_id}
                    className="flex items-center justify-between p-1.5 rounded border border-border bg-primary/5 group"
                  >
                    <div className="flex flex-col min-w-0">
                      <span className="text-xs font-medium truncate" title={c.component?.name}>
                        {c.component?.name}
                      </span>
                      {c.component?.element && (
                        <span className="text-[9px] text-muted-foreground italic">
                          {c.component.element}
                        </span>
                      )}
                    </div>
                    
                    <div className="flex items-center gap-1.5">
                      <Badge variant="outline" className="h-5 px-1.5 text-xs font-mono">
                        ×{c.count}
                      </Badge>
                      
                      <div className="flex items-center opacity-0 group-hover:opacity-100 transition-opacity">
                        <Button
                          size="icon"
                          variant="ghost"
                          className="h-5 w-5"
                          onClick={() => consumeMutation.mutate(c.component_id)}
                          disabled={consumeMutation.isPending || c.count <= 0}
                        >
                          <Minus className="h-3 w-3" />
                        </Button>
                        <Button
                          size="icon"
                          variant="ghost"
                          className="h-5 w-5"
                          onClick={() => gainMutation.mutate(c.component_id)}
                          disabled={gainMutation.isPending}
                        >
                          <Plus className="h-3 w-3" />
                        </Button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          );
        })}
      </CardContent>
    </Card>
  );
}
