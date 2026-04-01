import { cn } from '@/lib/utils';
import type { ApiComponent, ComponentCategory } from '@/types/game';

// Category color mapping for periodic table visual distinction
const categoryColors: Record<ComponentCategory, { bg: string; border: string; text: string }> = {
  Forma: {
    bg: 'bg-yellow-100 dark:bg-yellow-900/30',
    border: 'border-yellow-600/60',
    text: 'text-yellow-900 dark:text-yellow-200',
  },
  Scopus: {
    bg: 'bg-sky-100 dark:bg-sky-900/30',
    border: 'border-sky-500/60',
    text: 'text-sky-900 dark:text-sky-200',
  },
  Essentia: {
    bg: 'bg-amber-100 dark:bg-amber-900/30',
    border: 'border-amber-500/60',
    text: 'text-amber-900 dark:text-amber-200',
  },
  Actio: {
    bg: 'bg-red-100 dark:bg-red-900/30',
    border: 'border-red-500/60',
    text: 'text-red-900 dark:text-red-200',
  },
  Magnitudo: {
    bg: 'bg-violet-100 dark:bg-violet-900/30',
    border: 'border-violet-500/60',
    text: 'text-violet-900 dark:text-violet-200',
  },
  Logica: {
    bg: 'bg-slate-200/80 dark:bg-slate-800/50',
    border: 'border-slate-500/50',
    text: 'text-slate-800 dark:text-slate-100',
  },
};

// Element type icons (for primordial elements)
const elementIcons: Record<string, string> = {
  fire: '🔥',
  water: '💧',
  earth: '🪨',
  air: '💨',
  lightning: '⚡',
  metal: '⚙️',
  nature: '🌿',
  shadow: '🌑',
  light: '✨',
  arcane: '🔮',
  thunder: '🔊',
  acid: '🧪',
  psychic: '🧠',
  radiant: '☀️',
  force: '💫',
};

interface ElementTileProps {
  component: ApiComponent;
  onClick?: (component: ApiComponent) => void;
  size?: 'sm' | 'md' | 'lg';
  showElement?: boolean;
  locked?: boolean;
  selected?: boolean;
  count?: number;
}

export function ElementTile({
  component,
  onClick,
  size = 'md',
  showElement = true,
  locked = false,
  selected = false,
  count,
}: ElementTileProps) {
  const colors = categoryColors[component.category] || categoryColors.Essentia;
  const elementIcon = component.element ? elementIcons[component.element] : null;

  const sizeClasses = {
    sm: 'w-12 h-14 md:w-14 md:h-16',
    md: 'w-16 h-20 md:w-20 md:h-24',
    lg: 'w-20 h-24 md:w-24 md:h-28',
  };

  const symbolSizes = {
    sm: 'text-lg md:text-xl',
    md: 'text-2xl md:text-3xl',
    lg: 'text-3xl md:text-4xl',
  };

  const nameSizes = {
    sm: 'text-nano md:text-tiny',
    md: 'text-tiny md:text-xs',
    lg: 'text-xs md:text-sm',
  };

  return (
    <button
      onClick={() => !locked && onClick?.(component)}
      disabled={locked}
      className={cn(
        'element-tile group relative flex flex-col items-center justify-center',
        'border-2 rounded-sm transition-all duration-200',
        'focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-1',
        colors.bg,
        colors.border,
        sizeClasses[size],
        locked
          ? 'opacity-40 blur-[1px] cursor-not-allowed grayscale'
          : 'hover:scale-105 hover:shadow-lg hover:z-10',
        selected && !locked && 'ring-2 ring-primary ring-offset-2 scale-105 shadow-lg shadow-primary/30'
      )}
      title={locked ? 'Locked - not available to this character' : `${component.name} - ${component.description}`}
    >
      {/* Element icon indicator */}
      {showElement && elementIcon && (
        <span className="absolute top-0.5 right-0.5 text-micro md:text-xs opacity-70">
          {elementIcon}
        </span>
      )}

      {/* Count Badge */}
      {count !== undefined && (
        <span 
          className={cn(
            "absolute -top-1.5 -left-1.5 z-20 flex h-5 min-w-[1.25rem] items-center justify-center rounded-full border border-background bg-primary px-1 text-micro font-bold text-primary-foreground shadow-sm",
            size === 'sm' && "h-4 min-w-[1rem] text-tiny -top-1 -left-1"
          )}
          title={`Quantity: ${count >= 999 ? 'Infinite' : count}`}
        >
          {count >= 999 ? '∞' : count}
        </span>
      )}

      {/* Symbol - large centered text */}
      <span
        className={cn(
          'font-bold font-mono leading-none',
          colors.text,
          symbolSizes[size]
        )}
      >
        {component.symbol}
      </span>

      {/* Name - small text below */}
      <span
        className={cn(
          'font-medium leading-tight text-center px-0.5 truncate w-full',
          'text-muted-foreground group-hover:text-foreground',
          nameSizes[size]
        )}
      >
        {component.name}
      </span>

      {/* Hover glow effect */}
      {!locked && (
        <div
          className={cn(
            'absolute inset-0 opacity-0 group-hover:opacity-100 transition-opacity',
            'bg-gradient-to-t from-transparent via-transparent to-white/20 dark:to-white/10',
            'pointer-events-none rounded-sm'
          )}
        />
      )}
    </button>
  );
}

export { categoryColors, elementIcons };
