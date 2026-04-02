import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/lib/utils";

const badgeVariants = cva(
  "inline-flex items-center rounded-full border font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 whitespace-nowrap overflow-hidden text-ellipsis max-w-full shrink",
  {
    variants: {
      variant: {
        default: "border-transparent bg-primary text-primary-foreground hover:bg-primary/80",
        secondary: "border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80",
        destructive: "border-transparent bg-destructive text-destructive-foreground hover:bg-destructive/80",
        outline: "border border-foreground text-foreground bg-background",
        "element-heal": "border-transparent bg-element-heal text-background hover:bg-element-heal/80",
        "element-lightning": "border-transparent bg-element-lightning text-background hover:bg-element-lightning/80",
        warning: "border-destructive/30 bg-destructive/10 text-destructive",
        accent: "border-transparent bg-accent text-accent-foreground hover:bg-accent/80",
        success: "border-transparent bg-green-500 text-white hover:bg-green-500/80",
      },
      size: {
        default: "px-2.5 py-0.5 text-xs",
        sm: "h-4 px-2 py-0 text-micro",
        tiny: "h-3.5 px-1.5 py-0 text-[0.625rem] leading-none", // Slightly smaller than micro
        lg: "px-3 py-1 text-base sm:text-lg", // Responsive sizing
        xl: "px-4 py-2 text-lg sm:text-xl", // Responsive sizing
        "xl-compact": "px-3 py-1 text-lg sm:text-xl",
        "h5-sm": "h-5 px-2 py-0 text-micro",
        "min-w-sm": "min-w-[3rem] justify-center",
        "min-w-xs": "min-w-[1.2rem] flex justify-center",
      },
      font: {
        tomeMarginalia: "font-tome-marginalia",
        mono: "font-mono",
      },
      theme: {
        "primary-light-outline": "bg-primary/10 border-primary/30 text-primary",
        "spell-type-default": "bg-slate-500/20 text-foreground",
        "spell-type-attack": "bg-red-500/20 text-foreground",
        "spell-type-healing": "bg-pink-500/20 text-foreground",
        "spell-type-save": "bg-amber-500/20 text-foreground",
        "element-lightning-outline": "bg-element-lightning/10 text-element-lightning border-element-lightning/30",
        "element-nature-outline": "bg-element-nature/10 text-element-nature border-element-nature/30",
        "background-subtle": "bg-background/50",
        "muted-subtle-outline": "bg-muted/50 text-muted-foreground border-muted-foreground/30",
        "element-ice-outline": "bg-element-ice/10 text-element-ice border-element-ice/30",
        "element-dark-outline": "bg-element-dark/10 text-element-dark border-element-dark/30",
        "destructive-outline": "border-destructive text-destructive",
        "warning": "bg-amber-500/10 border-amber-500/30 text-amber-500",
        "info": "bg-sky-500/10 border-sky-500/30 text-sky-500",
      }
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  },
);

export interface BadgeProps extends React.HTMLAttributes<HTMLDivElement>, VariantProps<typeof badgeVariants> {}

const Badge = React.forwardRef<HTMLDivElement, BadgeProps>(
  ({ className, variant, size, font, theme, ...props }, ref) => {
    return (
      <div 
        ref={ref} 
        className={cn(badgeVariants({ variant, size, font, theme }), className)} 
        {...props} 
      />
    );
  }
);
Badge.displayName = "Badge";

export { Badge, badgeVariants };