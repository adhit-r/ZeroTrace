import * as React from "react"

import { cn } from "@/lib/utils"

export interface InputProps
  extends React.InputHTMLAttributes<HTMLInputElement> {}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, type, ...props }, ref) => {
    return (
      <input
        type={type}
        className={cn(
          "flex h-11 w-full border-3 border-black bg-white px-3 py-2 text-sm font-semibold shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] file:border-0 file:bg-transparent file:text-sm file:font-bold placeholder:text-black/60 focus-visible:outline-none focus-visible:shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] focus-visible:-translate-x-0.5 focus-visible:-translate-y-0.5 disabled:cursor-not-allowed disabled:opacity-50 transition-all duration-150",
          className
        )}
        ref={ref}
        {...props}
      />
    )
  }
)
Input.displayName = "Input"

export { Input }