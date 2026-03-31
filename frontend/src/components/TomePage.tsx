import { ReactNode } from 'react';
import { cn } from '@/lib/utils';

interface TomePageProps {
  children: ReactNode;
  /** Two-column layout for long lore text (RPG book style) */
  columns?: boolean;
  /** Add deckled-edge effect to the content area */
  deckled?: boolean;
  /** Prose with drop-cap support – wrap main text in .tome-prose and use .drop-cap on first paragraph */
  prose?: boolean;
  className?: string;
}

/**
 * Wrapper for page content that should feel like a spread from a handbook:
 * optional two-column layout, deckled edges, and prose styling with drop caps.
 */
export function TomePage({
  children,
  columns = false,
  deckled = false,
  prose = false,
  className,
}: TomePageProps) {
  return (
    <div
      className={cn(
        'max-w-4xl mx-auto',
        columns && 'tome-columns',
        deckled && 'deckled-edge',
        prose && 'tome-prose prose-headings:font-tome-heading prose-headings:text-primary prose-p:font-tome-body prose-p:text-foreground',
        className
      )}
    >
      {children}
    </div>
  );
}
