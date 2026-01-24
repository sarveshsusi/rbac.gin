import { cn } from "../../lib/utils"

export function Button({
  className,
  variant = "default",
  ...props
}) {
  return (
    <button
      className={cn(
        "inline-flex w-full items-center justify-center rounded-md px-4 py-2 text-sm font-medium transition disabled:opacity-60 disabled:cursor-not-allowed",
        variant === "default" &&
          "bg-blue-600 text-white hover:bg-blue-700",
        variant === "outline" &&
          "border border-gray-300 hover:bg-gray-100",
        className
      )}
      {...props}
    />
  )
}
