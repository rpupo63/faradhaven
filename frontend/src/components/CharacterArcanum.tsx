import { useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import { getComponents } from '@/lib/api';
import { ElementTable } from '@/components/arcanum';
import { Loader2, Atom } from 'lucide-react';
import type { ApiComponent } from '@/types/game';

interface CharacterArcanumProps {
  /** Components the character has access to (from class/race) */
  availableComponents?: ApiComponent[];
}

/**
 * Arcanum page variant for character sheets.
 * Shows the full periodic table but with components the character
 * doesn't have access to blurred/locked.
 */
export function CharacterArcanum({ availableComponents }: CharacterArcanumProps) {
  const { data: allComponents, isLoading, error } = useQuery({
    queryKey: ['components'],
    queryFn: getComponents,
    staleTime: 60_000 * 10, // 10 minutes
  });

  // Create a Set of available component IDs for efficient lookup
  const availableComponentIds = useMemo(() => {
    if (!availableComponents) return undefined;
    return new Set(availableComponents.map((c) => c.id));
  }, [availableComponents]);

  if (isLoading) {
    return (
      <div className="min-h-[40vh] flex flex-col items-center justify-center gap-4">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
        <p className="text-muted-foreground font-tome-marginalia">
          Consulting the arcane manuscripts...
        </p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-[40vh] flex flex-col items-center justify-center gap-4">
        <Atom className="h-16 w-16 text-muted-foreground" />
        <h2 className="text-xl font-tome-heading text-primary">
          The Arcanum is Sealed
        </h2>
        <p className="text-muted-foreground font-tome-marginalia text-center max-w-md">
          Unable to access the component registry.
        </p>
      </div>
    );
  }

  if (!allComponents || allComponents.length === 0) {
    return (
      <div className="min-h-[40vh] flex flex-col items-center justify-center gap-4">
        <Atom className="h-16 w-16 text-muted-foreground" />
        <h2 className="text-xl font-tome-heading text-primary">
          Empty Arcanum
        </h2>
        <p className="text-muted-foreground font-tome-marginalia text-center max-w-md">
          No spell components have been catalogued yet.
        </p>
      </div>
    );
  }

  const availableCount = availableComponentIds?.size ?? allComponents.length;
  const totalCount = allComponents.length;

  return (
    <div className="w-full">
      {/* Character component summary */}
      {availableComponentIds && (
        <div className="mb-6 p-4 arcane-border rounded-lg bg-card/60 text-center">
          <p className="text-sm text-muted-foreground font-tome-marginalia">
            Your character has access to{' '}
            <span className="text-primary font-semibold">{availableCount}</span> of{' '}
            <span className="text-muted-foreground">{totalCount}</span> components.
            <br />
            <span className="text-xs">Locked components appear blurred.</span>
          </p>
        </div>
      )}
      
      <ElementTable
        components={allComponents}
        availableComponentIds={availableComponentIds}
      />
    </div>
  );
}
