import { cn } from "../../lib/utils"

export function FieldGroup({ className, ...props }) {
  return (
    <div className={cn("flex flex-col gap-4", className)} {...props} />
  )
}

export function Field({ className, ...props }) {
  return (
    <div className={cn("flex flex-col gap-1", className)} {...props} />
  )
}

export function FieldLabel({ className, ...props }) {
  return (
    <label
      className={cn("text-sm font-medium text-gray-900", className)}
      {...props}
    />
  )
}

export function FieldDescription({ className, ...props }) {
  return (
    <p className={cn("text-sm text-gray-600", className)} {...props} />
  )
}

export function FieldSeparator({ children }) {
  return (
    <div className="relative text-center text-sm text-gray-500">
      <span className="bg-white px-2">{children}</span>
      <div className="absolute inset-x-0 top-1/2 -z-10 h-px bg-gray-200" />
    </div>
  )
}
