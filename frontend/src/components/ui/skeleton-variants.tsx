import { Skeleton } from "./skeleton";
import { cn } from "@/lib/utils";

interface SkeletonProps {
  className?: string;
}

/**
 * Card skeleton for list items
 */
export function CardSkeleton({ className }: SkeletonProps) {
  return (
    <div className={cn("arcane-border rounded-xl p-4 space-y-3", className)}>
      <div className="flex items-start justify-between gap-2">
        <Skeleton className="h-6 w-32" />
        <Skeleton className="h-5 w-16 rounded-full" />
      </div>
      <Skeleton className="h-4 w-24" />
      <Skeleton className="h-16 w-full" />
      <div className="flex gap-2">
        <Skeleton className="h-5 w-14 rounded-full" />
        <Skeleton className="h-5 w-14 rounded-full" />
        <Skeleton className="h-5 w-14 rounded-full" />
      </div>
    </div>
  );
}

/**
 * Stat block skeleton for character sheet sections
 */
export function StatBlockSkeleton({ className }: SkeletonProps) {
  return (
    <div className={cn("arcane-border rounded-xl p-4 space-y-4", className)}>
      <div className="flex items-center gap-2">
        <Skeleton className="h-5 w-5 rounded-full" />
        <Skeleton className="h-5 w-28" />
      </div>
      <div className="grid grid-cols-3 gap-4">
        <div className="text-center space-y-2">
          <Skeleton className="h-8 w-12 mx-auto" />
          <Skeleton className="h-3 w-8 mx-auto" />
        </div>
        <div className="text-center space-y-2">
          <Skeleton className="h-8 w-12 mx-auto" />
          <Skeleton className="h-3 w-8 mx-auto" />
        </div>
        <div className="text-center space-y-2">
          <Skeleton className="h-8 w-12 mx-auto" />
          <Skeleton className="h-3 w-8 mx-auto" />
        </div>
      </div>
    </div>
  );
}

/**
 * Grid of card skeletons
 */
export function CardGridSkeleton({
  count = 4,
  className
}: SkeletonProps & { count?: number }) {
  return (
    <div className={cn("grid gap-6 md:grid-cols-2 lg:grid-cols-3", className)}>
      {Array.from({ length: count }).map((_, i) => (
        <CardSkeleton key={i} />
      ))}
    </div>
  );
}
