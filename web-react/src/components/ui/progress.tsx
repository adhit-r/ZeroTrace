import * as React from "react"

import { cn } from "@/lib/utils"

const Progress = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement> & { value?: number }
>(({ className, value = 0, ...props }, ref) => (
  <div
    ref={ref}
    className={cn(
      "relative h-8 w-full border-3 border-black bg-white shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] overflow-hidden",
      className
    )}
    {...props}
  >
    <div
      className="h-full bg-black transition-all duration-300 relative"
      style={{ width: `${Math.max(0, Math.min(100, value))}%` }}
    >
      <div
        className="absolute inset-0 bg-repeat"
        style={{
          backgroundImage: `repeating-linear-gradient(
            45deg,
            transparent,
            transparent 2px,
            rgba(255, 255, 255, 0.1) 2px,
            rgba(255, 255, 255, 0.1) 4px
          )`,
        }}
      />
    </div>
  </div>
))
Progress.displayName = "Progress"

export { Progress }

