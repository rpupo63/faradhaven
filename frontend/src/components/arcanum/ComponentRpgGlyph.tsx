import { RaIcon } from '@/components/ui/RaIcon';
import type { ApiComponent } from '@/types/game';
import { cn } from '@/lib/utils';

interface ComponentRpgGlyphProps {
  component: ApiComponent;
  /** Applied to RaIcon when an RPG Awesome icon is present */
  iconClassName?: string;
  /** Applied to the alchemical symbol fallback */
  fallbackClassName?: string;
}

/** RPG Awesome glyph when seeded; otherwise the short alchemical symbol. */
export function ComponentRpgGlyph({ component, iconClassName, fallbackClassName }: ComponentRpgGlyphProps) {
  const icon = component.rpg_awesome_icon?.trim();
  if (icon) {
    return <RaIcon name={icon} className={iconClassName} />;
  }
  return <span className={fallbackClassName}>{component.symbol}</span>;
}

/** Inline chip: icon or symbol plus optional label */
export function ComponentRpgChip({
  component,
  showName = false,
  className,
  iconClassName,
  nameClassName,
}: {
  component: ApiComponent;
  showName?: boolean;
  className?: string;
  iconClassName?: string;
  nameClassName?: string;
}) {
  return (
    <span className={cn('inline-flex items-center gap-1', className)}>
      <ComponentRpgGlyph
        component={component}
        iconClassName={cn('leading-none', iconClassName)}
        fallbackClassName={cn('font-mono font-bold leading-none', iconClassName)}
      />
      {showName && <span className={nameClassName}>{component.name}</span>}
    </span>
  );
}
