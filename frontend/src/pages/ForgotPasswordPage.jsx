import { useState } from "react"
import api from "../api/axios"
import { Mail, ArrowLeft, Loader2 } from "lucide-react"
import { Link } from "react-router-dom"

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState("")
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState("")
  const [error, setError] = useState("")

  const submit = async (e) => {
    e.preventDefault()
    setError("")
    setSuccess("")
    setLoading(true)

    try {
      await api.post("/auth/forgot-password", { email })

      // 🔒 Always show same success message
      setSuccess(
        "If an account exists for this email, a password reset link has been sent."
      )
    } catch {
      setError("Something went wrong. Please try again.")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-50 px-4">
      <div className="w-full max-w-md rounded-2xl bg-white p-6 shadow-lg sm:p-8">
        {/* HEADER */}
        <div className="mb-6 text-center">
          <div className="mx-auto mb-3 flex h-12 w-12 items-center justify-center rounded-full bg-blue-100">
            <Mail className="h-6 w-6 text-blue-600" />
          </div>
          <h1 className="text-xl font-semibold text-slate-900">
            Forgot your password?
          </h1>
          <p className="mt-1 text-sm text-slate-500">
            Enter your email and we’ll send you a reset link.
          </p>
        </div>

        {/* FORM */}
        <form onSubmit={submit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-slate-700">
              Email address
            </label>
            <input
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
              className="mt-1 w-full rounded-xl border border-slate-300 px-4 py-3 text-sm
                         focus:border-blue-500 focus:ring-2 focus:ring-blue-500 focus:outline-none"
            />
          </div>

          {error && (
            <p className="text-sm text-red-600">{error}</p>
          )}

          {success && (
            <p className="text-sm text-green-600">{success}</p>
          )}

          <button
            type="submit"
            disabled={loading}
            className="flex w-full items-center justify-center gap-2 rounded-xl
                       bg-blue-600 px-4 py-3 text-sm font-semibold text-white
                       hover:bg-blue-700 disabled:opacity-60 transition"
          >
            {loading && <Loader2 className="h-4 w-4 animate-spin" />}
            Send reset link
          </button>
        </form>

        {/* FOOTER */}
        <div className="mt-6 text-center">
          <Link
            to="/login"
            className="inline-flex items-center gap-1 text-sm text-slate-600 hover:text-slate-900"
          >
            <ArrowLeft size={14} />
            Back to login
          </Link>
        </div>
      </div>
    </div>
  )
}
