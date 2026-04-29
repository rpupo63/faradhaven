import { AnimatePresence, motion } from 'framer-motion';
import { ComponentRpgGlyph } from '@/components/arcanum';
import { cn } from '@/lib/utils';
import type { SpellForgeVesselProps } from '@/components/spell-forge-vessels/types';

/** Powder Mage: horn trail + fuse spark along a line of glyphs. */
export function PowderMageForgeVessel({ sequence, lastAddedComponentId, className }: SpellForgeVesselProps) {
  const sparks = sequence.length;

  return (
    <div
      className={cn(
        'relative min-h-[170px] overflow-hidden rounded-xl border border-amber-900/45 bg-gradient-to-br from-amber-950/98 via-stone-950 to-neutral-950 p-4',
        className,
      )}
    >
      <svg className="pointer-events-none absolute inset-x-4 top-12 h-12 w-[calc(100%-2rem)]" viewBox="0 0 280 48" fill="none" aria-hidden>
        <motion.path
          d="M 12 28 Q 140 8 268 28"
          stroke="currentColor"
          strokeWidth="2"
          strokeDasharray="6 8"
          className="text-amber-700/60"
          initial={false}
          animate={{ pathLength: sparks > 0 ? 1 : 0.15 }}
          transition={{ duration: 0.5 }}
        />
        {[...Array(Math.min(6, sparks))].map((_, i) => (
          <motion.circle
            key={i}
            cx={28 + i * 46}
            cy={26 - Math.sin(i) * 6}
            r="5"
            fill="currentColor"
            className="text-orange-500/70"
            initial={{ opacity: 0.3 }}
            animate={{ opacity: [0.3, 1, 0.3] }}
            transition={{ duration: 0.8, repeat: Infinity, delay: i * 0.12 }}
          />
        ))}
      </svg>
      <div className="relative z-10 mt-10 flex min-h-[72px] flex-wrap items-center justify-center gap-2 px-2">
        <AnimatePresence>
          {sequence.map((comp, i) => (
            <motion.span
              key={`${comp.id}-${i}`}
              layout
              initial={{ scale: 0, x: -20, rotate: -25 }}
              animate={{ scale: 1, x: 0, rotate: 0 }}
              exit={{ opacity: 0 }}
              transition={{ type: 'spring', stiffness: 380, damping: 22 }}
              className={cn(
                'inline-flex h-10 w-10 items-center justify-center rounded-lg border border-orange-500/45 bg-black/55 shadow-[0_0_16px_rgba(251,146,60,0.35)]',
                lastAddedComponentId === comp.id && 'ring-2 ring-orange-400/90',
              )}
            >
              <ComponentRpgGlyph
                component={comp}
                iconClassName="text-base text-orange-100"
                fallbackClassName="font-mono text-xs font-bold text-orange-100"
              />
            </motion.span>
          ))}
        </AnimatePresence>
      </div>
      <p className="relative mt-3 text-center text-[10px] font-tome-marginalia uppercase tracking-wider text-amber-200/65">
        Powder trace · timed weave
      </p>
    </div>
  );
}
