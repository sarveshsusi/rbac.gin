import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Eye, EyeOff } from "lucide-react";
import { useAuth } from "../auth/useAuth";

export default function Login() {
  const { login } = useAuth();
  const navigate = useNavigate();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [err, setErr] = useState("");
  const [loading, setLoading] = useState(false);

const submit = async (e) => {
  e.preventDefault();
  if (loading) return;

  setErr("");
  setLoading(true);

  try {
    const result = await login(email.trim(), password);
    console.log(result);

    // 🔐 2FA FLOW
    if (result.status === "2FA") {
      sessionStorage.setItem("two_fa_token", result.token);
      navigate("/verify-otp", { replace: true });
      return; // ⛔ STOP HERE
    }

    // ✅ NORMAL LOGIN
    navigate("/", { replace: true });
  } catch (err) {
    if (
      err.response?.status === 403 &&
      err.response?.data?.error === "password_reset_required"
    ) {
      setErr("You must reset your password before logging in.");
    } else {
      setErr("Invalid email or password");
    }
  } finally {
    setLoading(false);
  }
};


  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4">
      <div className="w-full max-w-md bg-white rounded-lg shadow-lg p-8">
        <h1 className="text-3xl font-bold text-center mb-6">Sign In</h1>

        {err && (
          <div className="mb-4 text-sm text-red-600 bg-red-50 border px-4 py-2 rounded">
            {err}
          </div>
        )}

        <form onSubmit={submit} className="space-y-5">
          <input
            type="email"
            placeholder="Email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full px-4 py-3 border rounded-lg"
          />

          <div className="relative">
            <input
              type={showPassword ? "text" : "password"}
              placeholder="Password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-3 pr-12 border rounded-lg"
            />

            <button
              type="button"
              onClick={() => setShowPassword((v) => !v)}
              className="absolute right-3 top-3"
            >
              {showPassword ? <EyeOff /> : <Eye />}
            </button>
          </div>

          <button
            disabled={loading}
            className="w-full py-3 bg-blue-600 text-white rounded-lg"
          >
            {loading ? "Signing in..." : "Sign In"}
          </button>
        </form>
      </div>
    </div>
  );
}
