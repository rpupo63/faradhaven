import { cn } from "@/lib/utils";

interface RefreshIndicatorProps {
  /** Show indicator when true (typically isFetching && !isLoading) */
  isFetching: boolean;
  /** Additional class names */
  className?: string;
  /** Position of the indicator */
  position?: "top-right" | "top-left" | "bottom-right" | "bottom-left";
}

/**
 * Subtle pulsing dot indicator for background refetches.
 * Non-blocking - doesn't disable interactions.
 */
export function RefreshIndicator({
  isFetching,
  className,
  position = "top-right",
}: RefreshIndicatorProps) {
  if (!isFetching) return null;

  const positionClasses = {
    "top-right": "top-2 right-2",
    "top-left": "top-2 left-2",
    "bottom-right": "bottom-2 right-2",
    "bottom-left": "bottom-2 left-2",
  };

  return (
    <div
      className={cn(
        "absolute z-10",
        positionClasses[position],
        className
      )}
      role="status"
      aria-label="Refreshing data"
    >
      <div className="relative flex h-2.5 w-2.5">
        <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75" />
        <span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-primary" />
      </div>
    </div>
  );
}
