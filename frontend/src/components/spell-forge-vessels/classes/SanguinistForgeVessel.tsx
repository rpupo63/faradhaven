import { AnimatePresence, motion } from 'framer-motion';
import { ComponentRpgGlyph } from '@/components/arcanum';
import { cn } from '@/lib/utils';
import type { SpellForgeVesselProps } from '@/components/spell-forge-vessels/types';

/** Sanguinist: potion vial with swirling essence and floating glyphs. */
export function SanguinistForgeVessel({ sequence, lastAddedComponentId, className }: SpellForgeVesselProps) {
  const fillPct = Math.min(92, 28 + sequence.length * 11);

  return (
    <div
      className={cn(
        'relative flex min-h-[180px] flex-col items-center justify-end rounded-xl border border-red-950/50 bg-gradient-to-b from-red-950/95 via-rose-950/90 to-black/90 p-4',
        className,
      )}
    >
      <motion.div
        className="absolute inset-x-8 top-10 h-24 rounded-[40%] border-2 border-red-900/60 bg-red-950/40 backdrop-blur-[2px]"
        style={{
          boxShadow: 'inset 0 -12px 24px rgba(127,29,29,0.45)',
        }}
        animate={{
          boxShadow: [
            'inset 0 -12px 24px rgba(127,29,29,0.45)',
            'inset 0 -18px 32px rgba(220,38,38,0.35)',
            'inset 0 -12px 24px rgba(127,29,29,0.45)',
          ],
        }}
        transition={{ duration: 4, repeat: Infinity, ease: 'easeInOut' }}
      >
        <motion.div
          className="absolute bottom-0 left-0 right-0 rounded-b-[inherit] bg-gradient-to-t from-red-700/80 via-rose-900/70 to-transparent"
          animate={{ height: `${fillPct}%` }}
          transition={{ type: 'spring', stiffness: 200, damping: 22 }}
        />
      </motion.div>
      <div className="relative z-10 mb-8 flex h-[72px] w-[52px] flex-col items-center justify-start rounded-t-full rounded-b-3xl border-2 border-red-800/70 bg-gradient-to-b from-red-950/20 to-transparent shadow-[inset_0_0_20px_rgba(0,0,0,0.6)]">
        <div className="mt-2 h-3 w-8 rounded-full border border-red-900/80 bg-red-950/90" />
        <div className="relative mt-1 flex flex-1 flex-wrap items-start justify-center gap-1 px-1 pb-2 pt-3">
          <AnimatePresence>
            {sequence.map((comp, i) => (
              <motion.span
                key={`${comp.id}-${i}`}
                initial={{ scale: 0, y: -20, opacity: 0 }}
                animate={{ scale: 1, y: 0, opacity: 1 }}
                exit={{ opacity: 0, scale: 0 }}
                transition={{ type: 'spring', stiffness: 400, damping: 24 }}
                className={cn(
                  'inline-flex h-7 w-7 items-center justify-center rounded-full border border-rose-400/40 bg-black/50 shadow-[0_0_12px_rgba(225,29,72,0.45)]',
                  lastAddedComponentId === comp.id && 'ring-2 ring-rose-400/80',
                )}
              >
                <ComponentRpgGlyph
                  component={comp}
                  iconClassName="text-xs text-rose-100"
                  fallbackClassName="font-mono text-[10px] font-bold text-rose-100"
                />
              </motion.span>
            ))}
          </AnimatePresence>
        </div>
      </div>
      <p className="relative z-10 mt-auto text-center text-[10px] font-tome-marginalia uppercase tracking-wider text-rose-300/70">
        Ichor crucible · brew in glass
      </p>
    </div>
  );
}
