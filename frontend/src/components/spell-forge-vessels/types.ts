import type { FC } from 'react';
import type { ApiComponent } from '@/types/game';

export interface SpellForgeVesselProps {
  /** Ordered spell chain (same as forge sequence). */
  sequence: ApiComponent[];
  /** Id of the component most recently added (for entry animation); cleared by parent after timeout. */
  lastAddedComponentId: string | null;
  /** Visually stable key when sequence changes (optional). */
  sequenceRevision?: number;
  className?: string;
}

export type SpellForgeVesselComponent = FC<SpellForgeVesselProps>;
