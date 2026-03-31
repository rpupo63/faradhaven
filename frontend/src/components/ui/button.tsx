import * as React from "react";
import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/lib/utils";

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap text-sm font-medium ring-offset-background transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 font-display tracking-wide",
  {
    variants: {
      variant: {
        default: "rounded-md bg-primary text-primary-foreground hover:bg-primary/90 shadow-glow hover:shadow-lg",
        destructive: "rounded-md bg-destructive text-destructive-foreground hover:bg-destructive/90",
        outline: "rounded-md border border-primary/50 bg-transparent text-primary hover:bg-primary/10 hover:border-primary",
        secondary: "rounded-md bg-secondary text-secondary-foreground hover:bg-secondary/80",
        ghost: "rounded-md hover:bg-accent hover:text-accent-foreground",
        link: "text-primary underline-offset-4 hover:underline",
        arcane: "rounded-md bg-gradient-to-r from-accent to-primary text-primary-foreground hover:opacity-90 shadow-arcane hover:shadow-lg",
        forge: "rounded-md bg-element-fire/20 border border-element-fire/50 text-element-fire hover:bg-element-fire/30 hover:border-element-fire",
        /* Tome: quill underline (ink-splatter hover feel) */
        quill: "rounded-none bg-transparent text-primary font-tome-body border-b-2 border-transparent hover:border-primary hover:shadow-seal",
        /* Tome: wax seal – circular, seal-like shape */
        seal: "rounded-full border-2 border-ink/30 bg-primary text-primary-foreground shadow-seal hover:shadow-card hover:scale-105 [border-radius:50%]",
      },
      size: {
        default: "h-10 px-4 py-2",
        xs: "h-7 rounded-md px-2 text-xs",
        sm: "h-9 rounded-md px-3",
        lg: "h-11 rounded-md px-8",
        icon: "h-10 w-10",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  }
);

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : "button";
    return (
      <Comp
        className={cn(buttonVariants({ variant, size, className }))}
        ref={ref}
        {...props}
      />
    );
  }
);
Button.displayName = "Button";

export { Button, buttonVariants };
