import { motion } from 'framer-motion';
import { ComponentRpgGlyph } from '@/components/arcanum';
import { cn } from '@/lib/utils';
import type { SpellForgeVesselProps } from '@/components/spell-forge-vessels/types';

/** Ironwright: nested gears; complexity scales with component count. */
export function IronwrightForgeVessel({ sequence, lastAddedComponentId, className }: SpellForgeVesselProps) {
  const n = sequence.length;
  const gearTurn = Math.min(28, 8 + n * 5);

  return (
    <div
      className={cn(
        'relative flex min-h-[180px] flex-col items-center justify-center overflow-hidden rounded-xl border border-slate-500/40 bg-gradient-to-br from-slate-900 via-slate-950 to-zinc-950 p-4',
        className,
      )}
    >
      <div className="relative flex h-[140px] w-[140px] items-center justify-center">
        <motion.svg
          className="absolute text-slate-600/90"
          width={140}
          height={140}
          viewBox="0 0 100 100"
          animate={{ rotate: gearTurn }}
          transition={{ type: 'spring', stiffness: 120, damping: 14 }}
        >
          <circle cx="50" cy="50" r="46" fill="none" stroke="currentColor" strokeWidth="3" />
          {[0, 45, 90, 135, 180, 225, 270, 315].map((deg) => (
            <rect
              key={deg}
              x="46"
              y="4"
              width="8"
              height="14"
              fill="currentColor"
              transform={`rotate(${deg} 50 50)`}
            />
          ))}
        </motion.svg>
        <motion.svg
          className="absolute text-slate-400/85"
          width={92}
          height={92}
          viewBox="0 0 100 100"
          animate={{ rotate: -gearTurn * 1.3 }}
          transition={{ type: 'spring', stiffness: 100, damping: 12 }}
        >
          <circle cx="50" cy="50" r="38" fill="none" stroke="currentColor" strokeWidth="2.5" />
          {[0, 60, 120, 180, 240, 300].map((deg) => (
            <rect key={deg} x="47" y="10" width="6" height="12" fill="currentColor" transform={`rotate(${deg} 50 50)`} />
          ))}
        </motion.svg>
        <motion.svg
          className="absolute text-amber-700/90"
          width={56}
          height={56}
          viewBox="0 0 100 100"
          animate={{ rotate: gearTurn * 2 }}
          transition={{ type: 'spring', stiffness: 160, damping: 16 }}
        >
          <circle cx="50" cy="50" r="22" fill="none" stroke="currentColor" strokeWidth="2" className="text-amber-700/90" />
          {[0, 90, 180, 270].map((deg) => (
            <rect key={deg} x="46" y="22" width="8" height="10" fill="currentColor" transform={`rotate(${deg} 50 50)`} />
          ))}
        </motion.svg>
        <div className="relative z-10 flex h-14 w-14 flex-wrap items-center justify-center gap-0.5 rounded-full border border-amber-600/50 bg-slate-950/95 p-1 shadow-inner">
          {sequence.slice(-5).map((comp, i) => (
            <motion.span
              key={`${comp.id}-${sequence.length - 5 + i}`}
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              className={cn(
                'inline-flex h-7 w-7 items-center justify-center rounded-sm bg-amber-950/80',
                lastAddedComponentId === comp.id && 'ring-1 ring-amber-400',
              )}
            >
              <ComponentRpgGlyph
                component={comp}
                iconClassName="text-xs text-amber-100"
                fallbackClassName="font-mono text-[10px] font-bold text-amber-100"
              />
            </motion.span>
          ))}
        </div>
      </div>
      <p className="mt-2 text-center text-[10px] font-tome-marginalia uppercase tracking-wider text-slate-400">
        Clockwork assembly · {n} piece{n === 1 ? '' : 's'}
      </p>
    </div>
  );
}
