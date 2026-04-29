import { AnimatePresence, motion } from 'framer-motion';
import { ComponentRpgGlyph } from '@/components/arcanum';
import { cn } from '@/lib/utils';
import type { SpellForgeVesselProps } from '@/components/spell-forge-vessels/types';

/** Syllogist: premise → link → conclusion chain. */
export function SyllogistForgeVessel({ sequence, lastAddedComponentId, className }: SpellForgeVesselProps) {
  return (
    <div
      className={cn(
        'relative flex min-h-[170px] flex-col items-center justify-center rounded-xl border border-indigo-950/50 bg-gradient-to-br from-indigo-950/95 via-slate-950 to-black p-4',
        className,
      )}
    >
      <div className="flex w-full max-w-[300px] flex-wrap items-center justify-center gap-2 py-2">
        <AnimatePresence mode="popLayout">
          {sequence.map((comp, i) => (
            <span key={`${comp.id}-${i}`} className="flex items-center gap-2">
              {i > 0 && (
                <motion.span
                  className="font-tome-marginalia text-lg text-indigo-400/80"
                  initial={{ opacity: 0, scale: 0.5 }}
                  animate={{ opacity: 1, scale: 1 }}
                >
                  ∴
                </motion.span>
              )}
              <motion.span
                layout
                initial={{ scale: 0.6, opacity: 0 }}
                animate={{ scale: 1, opacity: 1 }}
                exit={{ opacity: 0, scale: 0.5 }}
                transition={{ type: 'spring', stiffness: 400, damping: 26 }}
                className={cn(
                  'inline-flex h-11 w-11 items-center justify-center rounded-lg border border-indigo-500/45 bg-indigo-950/90 shadow-[0_0_12px_rgba(99,102,241,0.25)]',
                  lastAddedComponentId === comp.id && 'ring-2 ring-indigo-400/75',
                )}
              >
                <ComponentRpgGlyph
                  component={comp}
                  iconClassName="text-lg text-indigo-100"
                  fallbackClassName="font-mono text-sm font-bold text-indigo-100"
                />
              </motion.span>
            </span>
          ))}
        </AnimatePresence>
      </div>
      <svg className="pointer-events-none absolute bottom-6 left-1/2 w-[80%] max-w-[240px] -translate-x-1/2 opacity-30" viewBox="0 0 200 24" aria-hidden>
        <path d="M 8 16 Q 100 2 192 16" stroke="currentColor" strokeWidth="1" fill="none" className="text-indigo-500" />
      </svg>
      <p className="relative mt-6 text-center text-[10px] font-tome-marginalia uppercase tracking-wider text-indigo-300/65">
        Inference chain · ordered claims
      </p>
    </div>
  );
}
