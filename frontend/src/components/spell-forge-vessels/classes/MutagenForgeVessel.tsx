import { AnimatePresence, motion } from 'framer-motion';
import { ComponentRpgGlyph } from '@/components/arcanum';
import { cn } from '@/lib/utils';
import type { SpellForgeVesselProps } from '@/components/spell-forge-vessels/types';

/** Mutagen: glyphs drop onto a wolf forearm silhouette. */
export function MutagenForgeVessel({ sequence, lastAddedComponentId, className }: SpellForgeVesselProps) {
  return (
    <div className={cn('relative min-h-[180px] overflow-hidden rounded-xl border border-amber-900/40 bg-gradient-to-b from-stone-950/90 to-stone-900/80 p-3', className)}>
      <div className="pointer-events-none absolute inset-0 opacity-[0.12]" aria-hidden>
        <div className="h-full w-full bg-[radial-gradient(ellipse_at_50%_120%,rgba(139,92,46,0.35),transparent_55%)]" />
      </div>
      <svg className="relative mx-auto h-[120px] w-full max-w-[280px]" viewBox="0 0 240 120" fill="none" aria-hidden>
        <ellipse cx="120" cy="95" rx="85" ry="22" fill="currentColor" className="text-stone-800/80" />
        <path
          d="M45 88 Q120 35 195 88 Q195 102 120 108 Q45 102 45 88"
          fill="currentColor"
          className="text-stone-700/90"
        />
        <path d="M95 55 Q120 42 145 55 L138 95 L102 95 Z" fill="currentColor" className="text-stone-600/70" />
      </svg>
      <div className="absolute bottom-6 left-1/2 flex max-w-[240px] -translate-x-1/2 flex-wrap items-end justify-center gap-1.5 px-2">
        <AnimatePresence>
          {sequence.map((comp, i) => (
            <motion.span
              key={`${comp.id}-${i}`}
              layout
              initial={{ y: -48, opacity: 0, rotate: -12 }}
              animate={{ y: 0, opacity: 1, rotate: 0 }}
              exit={{ opacity: 0, scale: 0.5 }}
              transition={{ type: 'spring', stiffness: 380, damping: 22 }}
              className={cn(
                'inline-flex h-10 w-10 items-center justify-center rounded-full border-2 border-amber-700/50 bg-amber-950/90 shadow-lg shadow-black/40',
                lastAddedComponentId === comp.id && 'ring-2 ring-amber-400/70',
              )}
            >
              <ComponentRpgGlyph
                component={comp}
                iconClassName="text-base text-amber-100"
                fallbackClassName="font-mono text-xs font-bold text-amber-100"
              />
            </motion.span>
          ))}
        </AnimatePresence>
      </div>
      <p className="relative mt-1 text-center text-[10px] font-tome-marginalia uppercase tracking-wider text-amber-200/70">
        Feral weave · arm attunement
      </p>
    </div>
  );
}
