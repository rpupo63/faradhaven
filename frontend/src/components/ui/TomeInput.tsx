import * as React from "react";
import { cn } from "@/lib/utils";

/** Underlined fill-in-the-blank style; text in handwritten font – Bureaucrat's Desk */
const TomeInput = React.forwardRef<HTMLInputElement, React.ComponentProps<"input">>(
  ({ className, type, ...props }, ref) => (
    <input
      type={type}
      className={cn(
        "tome-input-underline flex h-10 w-full px-1 py-2 text-base placeholder:text-muted-foreground focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm transition-shadow",
        className,
      )}
      ref={ref}
      {...props}
    />
  ),
);
TomeInput.displayName = "TomeInput";

export { TomeInput };
