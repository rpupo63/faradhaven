import { cn } from '@/lib/utils';
import { SpellElement as SpellElementType, ElementType } from '@/types/game';
import {
  Flame,
  Snowflake,
  Heart,
  Shield,
  Wind,
  Zap,
  Moon,
  Leaf,
  type LucideIcon,
} from 'lucide-react';

const ELEMENT_ICONS: Record<ElementType, LucideIcon> = {
  fire: Flame,
  ice: Snowflake,
  heal: Heart,
  shield: Shield,
  push: Wind,
  lightning: Zap,
  dark: Moon,
  nature: Leaf,
};

function ElementIcon({ elementId, className }: { elementId: ElementType; className?: string }) {
  const Icon = ELEMENT_ICONS[elementId];
  return Icon ? <Icon className={cn('shrink-0', className)} /> : null;
}

interface SpellElementProps {
  element: SpellElementType;
  selected?: boolean;
  onClick?: () => void;
  size?: 'sm' | 'md' | 'lg';
}

const elementStyles: Record<ElementType, string> = {
  fire: 'bg-element-fire/20 border-element-fire/50 text-element-fire hover:bg-element-fire/30',
  ice: 'bg-element-ice/20 border-element-ice/50 text-element-ice hover:bg-element-ice/30',
  heal: 'bg-element-heal/20 border-element-heal/50 text-element-heal hover:bg-element-heal/30',
  shield: 'bg-element-shield/20 border-element-shield/50 text-element-shield hover:bg-element-shield/30',
  push: 'bg-element-push/20 border-element-push/50 text-element-push hover:bg-element-push/30',
  lightning: 'bg-element-lightning/20 border-element-lightning/50 text-element-lightning hover:bg-element-lightning/30',
  dark: 'bg-element-dark/20 border-element-dark/50 text-foreground hover:bg-element-dark/30',
  nature: 'bg-element-nature/20 border-element-nature/50 text-element-nature hover:bg-element-nature/30',
};

const selectedStyles: Record<ElementType, string> = {
  fire: 'bg-element-fire/40 border-element-fire shadow-[0_0_15px_hsl(var(--element-fire)/0.5)]',
  ice: 'bg-element-ice/40 border-element-ice shadow-[0_0_15px_hsl(var(--element-ice)/0.5)]',
  heal: 'bg-element-heal/40 border-element-heal shadow-[0_0_15px_hsl(var(--element-heal)/0.5)]',
  shield: 'bg-element-shield/40 border-element-shield shadow-[0_0_15px_hsl(var(--element-shield)/0.5)]',
  push: 'bg-element-push/40 border-element-push shadow-[0_0_15px_hsl(var(--element-push)/0.5)]',
  lightning: 'bg-element-lightning/40 border-element-lightning shadow-[0_0_15px_hsl(var(--element-lightning)/0.5)]',
  dark: 'bg-element-dark/40 border-element-dark shadow-[0_0_15px_hsl(var(--element-dark)/0.5)]',
  nature: 'bg-element-nature/40 border-element-nature shadow-[0_0_15px_hsl(var(--element-nature)/0.5)]',
};

const sizeClasses = {
  sm: 'p-2 text-lg',
  md: 'p-3 text-2xl',
  lg: 'p-4 text-3xl',
};

const iconSizeClasses = {
  sm: 'w-5 h-5',
  md: 'w-6 h-6',
  lg: 'w-8 h-8',
};

export function SpellElementCard({ element, selected, onClick, size = 'md' }: SpellElementProps) {
  return (
    <button
      onClick={onClick}
      className={cn(
        'spell-element border rounded-lg transition-all duration-300 flex flex-col items-center gap-1',
        sizeClasses[size],
        elementStyles[element.id],
        selected && selectedStyles[element.id],
        onClick && 'cursor-pointer',
        !onClick && 'cursor-default'
      )}
    >
      <span className={cn('animate-float flex items-center justify-center', iconSizeClasses[size])}>
        <ElementIcon elementId={element.id} className="w-full h-full" />
      </span>
      <span className="text-xs font-display tracking-wider">{element.name}</span>
      {size !== 'sm' && (
        <span className="text-micro text-muted-foreground">Cost: {element.baseCost}</span>
      )}
    </button>
  );
}

export function SpellElementBadge({
  element,
  selected,
}: {
  element: SpellElementType;
  selected?: boolean;
}) {
  return (
    <span
      className={cn(
        'inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs border',
        elementStyles[element.id],
        selected && selectedStyles[element.id]
      )}
    >
      <ElementIcon elementId={element.id} className="w-3.5 h-3.5" />
      <span className="font-display">{element.name}</span>
    </span>
  );
}
