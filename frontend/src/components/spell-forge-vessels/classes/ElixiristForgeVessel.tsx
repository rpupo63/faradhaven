import { AnimatePresence, motion } from 'framer-motion';
import { ComponentRpgGlyph } from '@/components/arcanum';
import { cn } from '@/lib/utils';
import type { SpellForgeVesselProps } from '@/components/spell-forge-vessels/types';

/** Elixirist: cauldron with bubbles and rim glyphs. */
export function ElixiristForgeVessel({ sequence, lastAddedComponentId, className }: SpellForgeVesselProps) {
  return (
    <div
      className={cn(
        'relative min-h-[190px] overflow-hidden rounded-xl border border-emerald-900/50 bg-gradient-to-b from-emerald-950/95 via-teal-950/85 to-black/90 p-4',
        className,
      )}
    >
      <motion.div
        className="pointer-events-none absolute bottom-[52px] left-1/2 h-16 w-[85%] max-w-[220px] -translate-x-1/2 rounded-[50%] bg-emerald-500/15 blur-xl"
        animate={{ opacity: [0.35, 0.55, 0.35], scale: [1, 1.05, 1] }}
        transition={{ duration: 3.5, repeat: Infinity }}
      />
      <svg className="relative mx-auto block h-[130px] w-full max-w-[260px]" viewBox="0 0 200 120" fill="none" aria-hidden>
        <ellipse cx="100" cy="102" rx="72" ry="14" fill="currentColor" className="text-emerald-950/90" />
        <path
          d="M38 55 Q100 22 162 55 L148 95 Q100 108 52 95 Z"
          fill="currentColor"
          className="text-emerald-900/95"
          stroke="currentColor"
          strokeWidth="2"
          strokeOpacity="0.5"
        />
        <ellipse cx="100" cy="52" rx="58" ry="12" fill="currentColor" className="text-teal-900/70" />
      </svg>
      {[0, 1, 2].map((i) => (
        <motion.span
          key={i}
          className="pointer-events-none absolute bottom-[72px] h-2 w-2 rounded-full bg-emerald-300/40"
          style={{ left: `calc(50% + ${(i - 1) * 22}px - 4px)` }}
          animate={{ y: [0, -40, -70], opacity: [0.9, 0.5, 0], scale: [1, 0.8, 0.4] }}
          transition={{ duration: 2.2 + i * 0.4, repeat: Infinity, ease: 'easeOut', delay: i * 0.5 }}
        />
      ))}
      <div className="absolute bottom-[56px] left-1/2 flex max-w-[200px] -translate-x-1/2 flex-wrap items-center justify-center gap-1.5 px-2">
        <AnimatePresence>
          {sequence.map((comp, i) => (
            <motion.span
              key={`${comp.id}-${i}`}
              initial={{ scale: 0, rotate: 180 }}
              animate={{ scale: 1, rotate: 0 }}
              exit={{ opacity: 0 }}
              transition={{ type: 'spring', stiffness: 260, damping: 18 }}
              className={cn(
                'inline-flex h-9 w-9 items-center justify-center rounded-full border border-emerald-400/50 bg-black/60 shadow-[0_0_14px_rgba(52,211,153,0.35)]',
                lastAddedComponentId === comp.id && 'ring-2 ring-emerald-400/70',
              )}
            >
              <ComponentRpgGlyph
                component={comp}
                iconClassName="text-sm text-emerald-100"
                fallbackClassName="font-mono text-xs font-bold text-emerald-100"
              />
            </motion.span>
          ))}
        </AnimatePresence>
      </div>
      <p className="relative mt-8 text-center text-[10px] font-tome-marginalia uppercase tracking-wider text-emerald-300/75">
        Still crucible · vapors rise
      </p>
    </div>
  );
}
