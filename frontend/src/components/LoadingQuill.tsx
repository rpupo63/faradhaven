import { cn } from '@/lib/utils';

interface LoadingQuillProps {
  className?: string;
  /** Optional short label, e.g. "Loading..." */
  label?: string;
}

/** Loading state: quill writing script / inkwell – scholar's desk */
export function LoadingQuill({ className, label = 'Loading...' }: LoadingQuillProps) {
  return (
    <div className={cn('flex flex-col items-center justify-center gap-4 py-8', className)} role="status" aria-live="polite">
      <div className="relative">
        <img src="/wolf.gif" alt="Loading..." className="w-72 h-72 object-contain" />
      </div>
      <p className="font-tome-marginalia text-charcoal text-sm animate-pulse">{label}</p>
    </div>
  );
}
