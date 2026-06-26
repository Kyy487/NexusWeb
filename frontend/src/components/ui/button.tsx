import { cloneElement, isValidElement } from "react";

import { cn } from "@/lib/utils";

type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "primary" | "secondary" | "ghost" | "outline";
  size?: "sm" | "md" | "lg";
  asChild?: boolean;
};

export function Button({ className, variant = "primary", size = "md", asChild, children, ...props }: ButtonProps) {
  const classes = cn(
    "inline-flex items-center justify-center rounded-2xl font-medium transition-all duration-200 disabled:cursor-not-allowed disabled:opacity-50",
    variant === "primary" && "bg-blue-600 text-white shadow-lg shadow-blue-600/20 hover:-translate-y-0.5 hover:bg-blue-700",
    variant === "secondary" && "bg-slate-900 text-white shadow-lg shadow-slate-900/10 hover:-translate-y-0.5 hover:bg-slate-800",
    variant === "outline" && "border border-slate-200 bg-white text-slate-900 hover:-translate-y-0.5 hover:border-slate-300",
    variant === "ghost" && "bg-transparent text-slate-700 hover:bg-slate-100",
    size === "sm" && "h-9 px-4 text-sm",
    size === "md" && "h-11 px-5 text-sm",
    size === "lg" && "h-12 px-6 text-base",
    className
  );

  if (asChild && isValidElement(children)) {
    return cloneElement(children, {
      className: cn(classes, (children.props as { className?: string }).className),
    });
  }

  return (
    <button className={classes} {...props}>
      {children}
    </button>
  );
}