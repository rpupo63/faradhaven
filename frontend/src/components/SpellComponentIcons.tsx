import { cn } from '@/lib/utils';

/** Hand-stamped component icons: V (mouth), S (hand), M (pouch) – Master Ledger */
export function SpellComponentIcons({
  verbal = false,
  somatic = false,
  material = false,
  className,
}: {
  verbal?: boolean;
  somatic?: boolean;
  material?: boolean;
  className?: string;
}) {
  return (
    <span className={cn('inline-flex items-center gap-0.5 text-charcoal/80', className)} aria-label="Spell components">
      {/* V – tiny mouth */}
      <span
        className={cn(
          'w-5 h-5 flex items-center justify-center rounded border border-ink/30 text-micro font-tome-marginalia font-bold',
          verbal ? 'bg-ink/10 text-ink' : 'bg-parchment/50 text-muted-foreground'
        )}
        title="Verbal"
      >
        V
      </span>
      {/* S – hand gesture */}
      <span
        className={cn(
          'w-5 h-5 flex items-center justify-center rounded border border-ink/30 text-micro font-tome-marginalia font-bold',
          somatic ? 'bg-ink/10 text-ink' : 'bg-parchment/50 text-muted-foreground'
        )}
        title="Somatic"
      >
        S
      </span>
      {/* M – pouch */}
      <span
        className={cn(
          'w-5 h-5 flex items-center justify-center rounded border border-ink/30 text-micro font-tome-marginalia font-bold',
          material ? 'bg-ink/10 text-ink' : 'bg-parchment/50 text-muted-foreground'
        )}
        title="Material"
      >
        M
      </span>
    </span>
  );
}
