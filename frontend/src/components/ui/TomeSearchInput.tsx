import * as React from "react";
import { Search } from "lucide-react";
import { cn } from "@/lib/utils";

/** Search bar – simple line with sketched magnifying glass in corner (scholar's desk) */
const TomeSearchInput = React.forwardRef<
  HTMLInputElement,
  React.ComponentProps<"input">
>(({ className, ...props }, ref) => (
  <div className="tome-search-wrap flex items-center gap-2 px-2 py-1.5 focus-within:border-primary transition-colors">
    <Search
      className="w-4 h-4 text-muted-foreground shrink-0"
      strokeWidth={1.5}
      aria-hidden
    />
    <input
      ref={ref}
      type="search"
      className={cn(
        "flex-1 min-w-0 bg-transparent border-0 px-0 py-1 text-foreground font-tome-body placeholder:text-muted-foreground focus:outline-none md:text-sm",
        className,
      )}
      {...props}
    />
  </div>
));
TomeSearchInput.displayName = "TomeSearchInput";

export { TomeSearchInput };
