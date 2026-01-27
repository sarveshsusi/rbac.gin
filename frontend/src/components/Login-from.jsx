import { useState } from "react"
import { useNavigate } from "react-router-dom"
import { Eye, EyeOff } from "lucide-react"

import { cn } from "../lib/utils"
import { Button } from "../components/ui/button"
import {
  Field,
  FieldGroup,
  FieldLabel,
  FieldSeparator,
} from "./ui/field"
import { Input } from "./ui/input"
import { useAuth } from "../auth/useAuth"

export function LoginForm({ className, ...props }) {
  const { login } = useAuth()
  const navigate = useNavigate()

  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [rememberDevice, setRememberDevice] = useState(false) // 👈 NEW
  const [showPassword, setShowPassword] = useState(false)
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)

  const submit = async (e) => {
    e.preventDefault()
    if (loading) return

    setError("")
    setLoading(true)

    try {
      // 👇 PASS rememberDevice
      const result = await login(
        email.trim(),
        password,
        rememberDevice
      )
    


      // 🔐 2FA REQUIRED → redirect
      if (result.status === "2FA") {
        sessionStorage.setItem("two_fa_token", result.token)
        navigate("/verify-otp", { replace: true })
        return
      }

      // ✅ NORMAL LOGIN
      navigate("/", { replace: true })
    } catch (err) {
      if (
        err.response?.status === 403 &&
        err.response?.data?.error === "password_reset_required"
      ) {
        setError("You must reset your password before logging in.")
      } else {
        setError("Invalid email or password")
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <form
      onSubmit={submit}
      className={cn("flex flex-col gap-6", className)}
      autoComplete="off"
      {...props}
    >
      <FieldGroup>
        <div className="text-center">
          <h1 className="text-2xl font-bold">Login to your account</h1>
          <p className="text-sm text-gray-600 mt-1">
            Enter your email below to login
          </p>
        </div>

        {error && (
          <div className="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-600">
            {error}
          </div>
        )}

        <Field>
          <FieldLabel>Email</FieldLabel>
          <Input
            type="email"
            required
            autoComplete="username"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </Field>

        <Field>
          <div className="flex items-center">
            <FieldLabel>Password</FieldLabel>
            <a
              href="/forgot-password"
              className="ml-auto text-sm hover:underline"
            >
              Forgot?
            </a>
          </div>

          <div className="relative">
            <Input
              type={showPassword ? "text" : "password"}
              required
              autoComplete="current-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="pr-10"
            />

            <button
              type="button"
              aria-label="Toggle password visibility"
              onClick={() => setShowPassword((v) => !v)}
              className="absolute right-3 top-2.5 text-gray-500"
            >
              {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
            </button>
          </div>
        </Field>

        {/* 👇 REMEMBER DEVICE CHECKBOX */}
        <label className="flex items-center gap-2 text-sm text-gray-600">
          <input
            type="checkbox"
            checked={rememberDevice}
            onChange={(e) => setRememberDevice(e.target.checked)}
            className="accent-blue-600"
          />
          Remember this device
        </label>

        <Button type="submit" disabled={loading}>
          {loading ? "Signing in…" : "Login"}
        </Button>

        <FieldSeparator>Powered by Emerald</FieldSeparator>
      </FieldGroup>
    </form>
  )
}
