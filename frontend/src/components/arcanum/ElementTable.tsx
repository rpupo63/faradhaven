import { useMemo, useState } from 'react';
import { ElementTile, categoryColors } from './ElementTile';
import { ElementDetail } from './ElementDetail';
import type { ApiComponent, ComponentCategory, ApiCharacterComponent } from '@/types/game';
import { cn } from '@/lib/utils';

// Category display names and groupings
const categoryMeta: Record<
  ComponentCategory,
  { name: string; shortName: string; group: 'required' | 'optional' }
> = {
  Forma: { name: 'Forma (Shape)', shortName: 'Forma', group: 'required' },
  Scopus: { name: 'Scopus (Targeting)', shortName: 'Scopus', group: 'optional' },
  Essentia: { name: 'Essentia (Domain)', shortName: 'Essentia', group: 'required' },
  Actio: { name: 'Actio (Action)', shortName: 'Actio', group: 'optional' },
  Magnitudo: { name: 'Magnitudo (Scale)', shortName: 'Magnitudo', group: 'optional' },
  Logica: { name: 'Logica (Flow)', shortName: 'Logica', group: 'optional' },
};

// Category order for display
const categoryOrder: ComponentCategory[] = [
  'Forma',
  'Scopus',
  'Essentia',
  'Actio',
  'Magnitudo',
  'Logica',
];

interface ElementTableProps {
  components: ApiComponent[];
  onComponentClick?: (component: ApiComponent) => void;
  /** If provided, components not in this set will be shown as locked/blurred */
  availableComponentIds?: Set<string>;
  /** If provided, components in this set will be shown as selected/highlighted */
  selectedComponentIds?: Set<string>;
  /** If true, clicking a component will not open the detail modal */
  disableDetailPopup?: boolean;
  /** Character's inventory of components to display counts */
  characterComponents?: ApiCharacterComponent[];
}

export function ElementTable({
  components,
  onComponentClick,
  availableComponentIds,
  selectedComponentIds,
  disableDetailPopup = false,
  characterComponents,
}: ElementTableProps) {
  const [selectedComponent, setSelectedComponent] = useState<ApiComponent | null>(null);
  const [detailOpen, setDetailOpen] = useState(false);

  // Group components by category
  const componentsByCategory = useMemo(() => {
    const grouped: Record<ComponentCategory, ApiComponent[]> = {} as Record<ComponentCategory, ApiComponent[]>;
    for (const cat of categoryOrder) {
      grouped[cat] = [];
    }
    for (const comp of components) {
      if (grouped[comp.category]) {
        grouped[comp.category].push(comp);
      }
    }
    // Sort each category by name
    for (const cat of categoryOrder) {
      grouped[cat].sort((a, b) => a.name.localeCompare(b.name));
    }
    return grouped;
  }, [components]);

  // Create a map of component counts if characterComponents is provided
  const componentCounts = useMemo(() => {
    if (!characterComponents) return null;
    const counts = new Map<string, number>();
    for (const cc of characterComponents) {
      counts.set(cc.component_id, cc.count);
    }
    return counts;
  }, [characterComponents]);

  const handleTileClick = (component: ApiComponent) => {
    if (!disableDetailPopup) {
      setSelectedComponent(component);
      setDetailOpen(true);
    }
    onComponentClick?.(component);
  };

  return (
    <div className="arcanum-table w-full">
      {/* Title Section */}
      <div className="text-center mb-8">
        <h1 className="text-4xl md:text-5xl mb-2">Arcanum Elementis</h1>
        <p className="text-muted-foreground font-tome-marginalia text-lg">
          Pillars and sequential logic for spell crafting
        </p>
      </div>

      {/* Pillar Sections */}
      <div className="space-y-8">
        {categoryOrder.map((cat) => {
          const meta = categoryMeta[cat];
          const comps = componentsByCategory[cat];
          if (!comps || comps.length === 0) return null;

          return (
            <section key={cat}>
              <div className="flex items-center gap-3 mb-4 border-b border-faded-gold/30 pb-2">
                <span
                  className={cn(
                    'w-4 h-4 rounded-sm',
                    categoryColors[cat].bg,
                    categoryColors[cat].border,
                    'border'
                  )}
                />
                <h2 className="text-xl md:text-2xl">
                  {meta.name}
                </h2>
                {meta.group === 'required' && (
                  <span className="text-red-500 text-sm font-bold ml-1" title="Required">*</span>
                )}
              </div>
              <p className="text-center text-sm text-muted-foreground mb-4 font-tome-marginalia">
                {cat === 'Forma' && 'The physical manifestation and geometric delivery of the magic. (One required)'}
                {cat === 'Scopus' && 'The anchor point or entity the magic is allowed to interact with.'}
                {cat === 'Essentia' && 'The fundamental matter, energy, or abstract concept being manipulated. (At least one required)'}
                {cat === 'Actio' && 'What the magic physically does to the Essentia or Scopus.'}
                {cat === 'Magnitudo' && "The mathematical dials that adjust the spell's parameters."}
                {cat === 'Logica' &&
                  'Optional — for multi-phase spells only. If / Then / Therefore split the chain into phases (default spells stay a single phase). Order matters — e.g. water, Then, cold.'}
              </p>
              <div className="flex flex-wrap justify-center gap-2 md:gap-3">
                {comps.map((comp) => (
                  <ElementTile
                    key={comp.id}
                    component={comp}
                    onClick={handleTileClick}
                    size="md"
                    showElement={cat === 'Essentia'}
                    locked={availableComponentIds ? !availableComponentIds.has(comp.id) : false}
                    selected={selectedComponentIds?.has(comp.id)}
                    count={componentCounts && !availableComponentIds?.has(comp.id) ? (componentCounts.get(comp.id) ?? 0) : undefined}
                  />
                ))}
              </div>
            </section>
          );
        })}
      </div>

      {/* Legend */}
      <section className="mt-8 p-4 arcane-border rounded-lg">
        <h3 className="text-lg font-semibold mb-3 text-center">Category Legend</h3>
        <div className="flex flex-wrap justify-center gap-3">
          {categoryOrder.map((cat) => (
            <div
              key={cat}
              className={cn(
                'flex items-center gap-1.5 px-2 py-1 rounded text-xs',
                categoryColors[cat].bg,
                categoryColors[cat].border,
                'border'
              )}
            >
              <span className={cn('font-medium', categoryColors[cat].text)}>
                {categoryMeta[cat].shortName}
              </span>
              {categoryMeta[cat].group === 'required' && (
                <span className="text-red-500 font-bold">*</span>
              )}
            </div>
          ))}
        </div>
      </section>

      {/* Detail Modal */}
      <ElementDetail
        component={selectedComponent}
        open={detailOpen}
        onOpenChange={setDetailOpen}
      />
    </div>
  );
}

export { categoryMeta, categoryOrder };
