import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog';
import { categoryColors, elementIcons } from './ElementTile';
import { ComponentRpgGlyph } from './ComponentRpgGlyph';
import type { ApiComponent, ComponentCategory } from '@/types/game';
import { cn } from '@/lib/utils';

// Category descriptions for the detail view
const categoryDescriptions: Record<ComponentCategory, string> = {
  Forma:
    'The physical manifestation and geometric delivery of the magic. Exactly one Forma is required for a valid spell.',
  Scopus:
    'The anchor point or entity the magic is allowed to interact with. Defines where the spell originates or who it affects.',
  Essentia:
    'The fundamental matter, energy, or abstract concept being manipulated. At least one Essentia is required for a valid spell.',
  Actio:
    'The kinetic "verb" of the spell. Defines what the magic physically does to the Essentia or Scopus.',
  Magnitudo:
    'The mathematical dials that adjust the spell\'s parameters, such as power, size, or duration.',
  Logica:
    'Connectors that establish order and causality between parts of a spell. They do not add elements or shapes by themselves; reuse other components before and after them.',
};

// Category display names
const categoryNames: Record<ComponentCategory, string> = {
  Forma: 'Forma (Shape)',
  Scopus: 'Scopus (Targeting)',
  Essentia: 'Essentia (Domain)',
  Actio: 'Actio (Action)',
  Magnitudo: 'Magnitudo (Scale)',
  Logica: 'Logica (Flow)',
};

interface ElementDetailProps {
  component: ApiComponent | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function ElementDetail({ component, open, onOpenChange }: ElementDetailProps) {
  if (!component) return null;

  const colors = categoryColors[component.category] ?? categoryColors.Essentia;
  const elementIcon = component.element ? elementIcons[component.element] : null;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <div className="flex items-start gap-4">
            {/* Large symbol tile */}
            <div
              className={cn(
                'flex-shrink-0 w-20 h-24 flex flex-col items-center justify-center',
                'border-2 rounded-sm',
                colors.bg,
                colors.border
              )}
            >
              {elementIcon && (
                <span className="text-sm mb-1">{elementIcon}</span>
              )}
              <ComponentRpgGlyph
                component={component}
                iconClassName={cn('text-3xl', colors.text)}
                fallbackClassName={cn('text-3xl font-bold font-mono', colors.text)}
              />
              <span className="text-xs text-muted-foreground mt-0.5">
                {component.name}
              </span>
            </div>

            <div className="flex-1 min-w-0">
              <DialogTitle className="text-2xl font-tome-heading">
                {component.name}
              </DialogTitle>
              <div
                className={cn(
                  'inline-flex items-center gap-1.5 px-2 py-0.5 rounded text-xs mt-1',
                  colors.bg,
                  colors.border,
                  'border'
                )}
              >
                <span className={cn('font-medium', colors.text)}>
                  {categoryNames[component.category]}
                </span>
              </div>
            </div>
          </div>
        </DialogHeader>

        <div className="space-y-4 mt-4">
          {/* Description */}
          <div>
            <h4 className="text-sm font-semibold text-muted-foreground mb-1">
              Description
            </h4>
            <DialogDescription className="text-foreground text-base">
              {component.description}
            </DialogDescription>
          </div>

          {/* Element Type */}
          {component.element && (
            <div>
              <h4 className="text-sm font-semibold text-muted-foreground mb-1">
                Element Type
              </h4>
              <div className="flex items-center gap-2">
                <span className="text-lg">{elementIcon}</span>
                <span className="capitalize text-foreground">
                  {component.element}
                </span>
              </div>
            </div>
          )}

          {/* Category Info */}
          <div className="pt-3 border-t border-border">
            <h4 className="text-sm font-semibold text-muted-foreground mb-1">
              About {categoryNames[component.category]}
            </h4>
            <p className="text-sm text-muted-foreground font-tome-marginalia">
              {categoryDescriptions[component.category]}
            </p>
          </div>

          {/* Symbol explanation */}
          <div className="pt-3 border-t border-border">
            <h4 className="text-sm font-semibold text-muted-foreground mb-1">
              Alchemical Symbol
            </h4>
            <p className="text-sm text-muted-foreground font-tome-marginalia">
              <span className="font-mono font-bold text-foreground">
                {component.symbol}
              </span>{' '}
              - Used in spell notation and alchemical formulae to represent{' '}
              <span className="text-foreground">{component.name}</span>.
            </p>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
