import { useMemo } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import type { Components } from 'react-markdown';

type LoreMarkdownProps = {
  source: string;
};

/**
 * Renders lore-notes markdown with GFM (tables, strikethrough, etc.) and tome typography.
 */
export function LoreMarkdown({ source }: LoreMarkdownProps) {
  const components = useMemo(() => {
    let firstP = true;
    return {
      p: ({ children, ...props }) => {
        const isFirst = firstP;
        if (firstP) firstP = false;
        return (
          <p className={isFirst ? 'drop-cap' : undefined} {...props}>
            {children}
          </p>
        );
      },
      table: ({ children, ...props }) => (
        <div className="my-4 overflow-x-auto rounded-lg border border-faded-gold/30 bg-card/40">
          <table className="w-full min-w-[520px] text-sm font-tome-body" {...props}>
            {children}
          </table>
        </div>
      ),
    } satisfies Components;
  }, [source]);

  return (
    <ReactMarkdown remarkPlugins={[remarkGfm]} components={components}>
      {source}
    </ReactMarkdown>
  );
}
