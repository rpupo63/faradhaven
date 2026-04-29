import { AnimatePresence, motion } from 'framer-motion';
import { ComponentRpgGlyph } from '@/components/arcanum';
import { cn } from '@/lib/utils';
import type { SpellForgeVesselProps } from '@/components/spell-forge-vessels/types';

/** Piston Brawler: pressure chamber; impact flash when a component lands. */
export function PistonBrawlerForgeVessel({ sequence, lastAddedComponentId, className }: SpellForgeVesselProps) {
  return (
    <div
      className={cn(
        'relative min-h-[180px] overflow-hidden rounded-xl border border-stone-600/50 bg-gradient-to-b from-zinc-900 to-stone-950 p-4',
        className,
      )}
    >
      <div className="relative mx-auto flex h-[150px] max-w-[240px] flex-col items-center justify-end">
        <motion.div
          className="relative z-0 h-8 w-32 rounded-t-md border-2 border-b-0 border-stone-500 bg-stone-800/90"
          animate={{ y: [0, 6, 0] }}
          transition={{ duration: 0.6, ease: 'easeInOut' }}
        />
        <div className="relative z-10 -mt-1 flex h-16 w-40 items-start justify-center rounded-b-lg border-2 border-t-0 border-stone-500 bg-stone-900/95 pt-1 shadow-[inset_0_4px_12px_rgba(0,0,0,0.5)]">
          <AnimatePresence>
            {sequence.map((comp, i) => (
              <motion.span
                key={`${comp.id}-${i}`}
                initial={{ y: -30, scale: 0.5 }}
                animate={{ y: 0, scale: 1 }}
                exit={{ opacity: 0 }}
                transition={{ type: 'spring', stiffness: 500, damping: 24 }}
                className={cn(
                  'm-0.5 inline-flex h-8 w-8 items-center justify-center rounded border border-amber-600/50 bg-stone-950/90',
                  lastAddedComponentId === comp.id && 'ring-2 ring-amber-500/80',
                )}
              >
                <ComponentRpgGlyph
                  component={comp}
                  iconClassName="text-sm text-amber-100"
                  fallbackClassName="font-mono text-xs font-bold text-amber-100"
                />
              </motion.span>
            ))}
          </AnimatePresence>
        </div>
        {lastAddedComponentId && (
          <motion.div
            className="pointer-events-none absolute bottom-4 left-1/2 h-20 w-20 -translate-x-1/2 rounded-full bg-amber-400/25 blur-2xl"
            initial={{ scale: 0.2, opacity: 0.9 }}
            animate={{ scale: 1.6, opacity: 0 }}
            transition={{ duration: 0.45 }}
            key={lastAddedComponentId}
          />
        )}
        <p className="mt-2 text-center text-[10px] font-tome-marginalia uppercase tracking-wider text-stone-400">
          Pressure stack · stability spent
        </p>
      </div>
    </div>
  );
}
