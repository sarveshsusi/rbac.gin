import { GalleryVerticalEnd } from "lucide-react"
import { LoginForm } from "../components/Login-from"

export default function LoginPage() {
  return (
    <div className="grid min-h-screen lg:grid-cols-2">
      <div className="flex flex-col gap-4 p-6 md:p-10">
        <div className="flex items-center gap-2 font-medium">
          <div className="flex h-6 w-6 items-center justify-center rounded-md bg-blue-600 text-white">
            <GalleryVerticalEnd size={16} />
          </div>
          ABC Corp
        </div>

        <div className="flex flex-1 items-center justify-center">
          <div className="w-full max-w-xs">
            <LoginForm />
          </div>
        </div>
      </div>

      <div className="relative hidden lg:block">
  <img
    src="https://images.unsplash.com/photo-1522202176988-66273c2fd55f"
    alt="Login illustration"
    className="absolute inset-0 h-full w-full object-cover"
  />
</div>
    </div>
  )
}
