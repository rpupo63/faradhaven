import { AnimatePresence, motion } from 'framer-motion';
import { ComponentRpgGlyph } from '@/components/arcanum';
import { cn } from '@/lib/utils';
import type { SpellForgeVesselProps } from '@/components/spell-forge-vessels/types';

/** Fallback when class is unknown or missing from the Faradhaven registry. */
export function DefaultForgeVessel({ sequence, lastAddedComponentId, className }: SpellForgeVesselProps) {
  return (
    <div
      className={cn(
        'flex min-h-[140px] flex-wrap items-center justify-center gap-2 rounded-lg border border-border/60 bg-muted/20 p-4',
        className,
      )}
    >
      <AnimatePresence mode="popLayout">
        {sequence.map((comp, i) => (
          <motion.span
            key={`${comp.id}-${i}`}
            layout
            initial={{ scale: 0.6, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            exit={{ scale: 0.5, opacity: 0 }}
            transition={{ type: 'spring', stiffness: 420, damping: 28 }}
            className={cn(
              'inline-flex h-11 w-11 items-center justify-center rounded-md border border-primary/30 bg-primary/10',
              lastAddedComponentId === comp.id && 'ring-2 ring-primary/60 ring-offset-2 ring-offset-background',
            )}
          >
            <ComponentRpgGlyph
              component={comp}
              iconClassName="text-lg text-primary"
              fallbackClassName="font-mono text-sm font-bold text-primary"
            />
          </motion.span>
        ))}
      </AnimatePresence>
      {sequence.length === 0 && (
        <p className="text-center text-xs italic text-muted-foreground">Components appear here as you forge.</p>
      )}
    </div>
  );
}
