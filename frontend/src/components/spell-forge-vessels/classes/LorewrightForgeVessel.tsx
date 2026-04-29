import { AnimatePresence, motion } from 'framer-motion';
import { ComponentRpgGlyph } from '@/components/arcanum';
import { cn } from '@/lib/utils';
import type { SpellForgeVesselProps } from '@/components/spell-forge-vessels/types';

/** Lorewright: open tome; glyphs line up as inscribed runes. */
export function LorewrightForgeVessel({ sequence, lastAddedComponentId, className }: SpellForgeVesselProps) {
  return (
    <div
      className={cn(
        'relative min-h-[180px] overflow-hidden rounded-xl border border-amber-900/35 bg-gradient-to-b from-amber-950/40 via-stone-900/90 to-stone-950 p-4',
        className,
      )}
    >
      <div className="relative mx-auto max-w-[260px]">
        <svg className="w-full" viewBox="0 0 240 130" fill="none" aria-hidden>
          <path
            d="M20 12 H120 L118 118 H22 Q20 12 20 12 Z"
            fill="currentColor"
            className="text-amber-950/85"
            stroke="currentColor"
            strokeWidth="1.5"
            strokeOpacity="0.4"
          />
          <path
            d="M120 12 H220 L218 118 H120 Q118 12 120 12 Z"
            fill="currentColor"
            className="text-amber-900/75"
            stroke="currentColor"
            strokeWidth="1.5"
            strokeOpacity="0.35"
          />
          <line x1="120" y1="14" x2="120" y2="116" stroke="currentColor" strokeOpacity="0.25" strokeWidth="1" />
        </svg>
        <div className="absolute left-[14%] right-[14%] top-[22%] flex min-h-[72px] flex-wrap content-start items-start justify-start gap-1.5">
          <AnimatePresence>
            {sequence.map((comp, i) => (
              <motion.span
                key={`${comp.id}-${i}`}
                initial={{ opacity: 0, x: -8 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0 }}
                transition={{ delay: i * 0.03 }}
                className={cn(
                  'inline-flex h-8 w-8 items-center justify-center rounded-sm border border-amber-700/40 bg-amber-950/70',
                  lastAddedComponentId === comp.id && 'ring-2 ring-amber-400/70',
                )}
              >
                <ComponentRpgGlyph
                  component={comp}
                  iconClassName="text-xs text-amber-50"
                  fallbackClassName="font-mono text-[10px] font-bold text-amber-50"
                />
              </motion.span>
            ))}
          </AnimatePresence>
        </div>
      </div>
      <p className="mt-2 text-center text-[10px] font-tome-marginalia uppercase tracking-wider text-amber-200/60">
        Codex thread · marginal binding
      </p>
    </div>
  );
}
