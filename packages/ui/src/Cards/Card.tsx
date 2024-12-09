import { cn } from "@authentication-service-go/utils"

interface CardProps {
  children?: React.ReactNode
  className?: string
}

export function Card({ children, className }: CardProps) {
  return (
    <div
      className={cn(
        "border-2 rounded-xl",
        "shadow-lg shadow-slate-300",
        "p-4 m-4",
        className
      )}
    >
      {children}
    </div>
  )
}
