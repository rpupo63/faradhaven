import { cn } from '@/lib/utils';

export type LogicConnectorVariant = 'if' | 'then' | 'therefore';

const labels: Record<LogicConnectorVariant, string> = {
  if: 'If',
  then: 'Then',
  therefore: 'Therefore',
};

interface LogicConnectorProps {
  variant: LogicConnectorVariant;
  className?: string;
}

/** Compact connector label for spell sequence UI (If / Then / Therefore). */
export function LogicConnector({ variant, className }: LogicConnectorProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center px-2 py-0.5 rounded-md border border-slate-500/40',
        'bg-slate-500/10 text-slate-700 dark:text-slate-200',
        'font-tome-marginalia text-xs font-semibold uppercase tracking-wide',
        className
      )}
    >
      {labels[variant]}
    </span>
  );
}

function variantFromComponentName(name: string): LogicConnectorVariant | null {
  switch (name) {
    case 'If':
      return 'if';
    case 'Then':
      return 'then';
    case 'Therefore':
      return 'therefore';
    default:
      return null;
  }
}

/** Maps backend Logica component names to connector styling. */
export function logicVariantFromName(name: string): LogicConnectorVariant | null {
  return variantFromComponentName(name);
}
