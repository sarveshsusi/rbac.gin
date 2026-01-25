import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Lock, Eye, EyeOff } from "lucide-react";
import { useAuth } from "../auth/useAuth";
import api, { setAccessToken } from "../api/axios";

export default function Login() {
  const { login } = useAuth();
  const navigate = useNavigate();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [otp, setOtp] = useState("");

  const [showPassword, setShowPassword] = useState(false);
  const [err, setErr] = useState("");
  const [loading, setLoading] = useState(false);

  // 🔐 2FA STATE
  const [twoFARequired, setTwoFARequired] = useState(false);
  const [twoFAToken, setTwoFAToken] = useState("");

  /* ======================
     LOGIN
  ====================== */
  const submit = async (e) => {
    e.preventDefault();
    if (loading) return;

    setErr("");
    setLoading(true);

    try {
      const res = await login(email.trim(), password);

      // 🔐 2FA REQUIRED
      if (res?.two_fa_required) {
        setTwoFARequired(true);
        setTwoFAToken(res.two_fa_token);
        return;
      }

      // ✅ NORMAL LOGIN
      navigate("/", { replace: true });
    } catch (err) {
      if (
        err.response?.status === 403 &&
        err.response?.data?.error === "password_reset_required"
      ) {
        setErr("You must reset your password before logging in.");
        return;
      }

      setErr("Invalid email or password");
    } finally {
      setLoading(false);
    }
  };

  /* ======================
     VERIFY OTP
  ====================== */
  const verifyOTP = async (e) => {
    e.preventDefault();
    if (loading) return;

    setErr("");
    setLoading(true);

    try {
      const res = await api.post(
        "/auth/verify-2fa",
        { code: otp },
        {
          headers: {
            "X-2FA-Token": twoFAToken,
          },
        }
      );

      setAccessToken(res.data.access_token);
      navigate("/", { replace: true });
    } catch (err) {
      setErr(err.response?.data?.error || "Invalid or expired OTP");
    } finally {
      setLoading(false);
    }
  };

  /* ======================
     UI
  ====================== */
  return (
    <div className="min-h-screen flex flex-col bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 shadow-sm">
        <div className="max-w-7xl mx-auto px-3 py-3 flex justify-between items-center">
          <div className="flex items-center gap-2">
            <div className="w-6 h-6 bg-blue-600 rounded-sm rotate-45 flex items-center justify-center">
              <div className="w-3 h-3 bg-white rounded-sm -rotate-45" />
            </div>
            <span className="text-xl font-semibold text-gray-900">Zoho</span>
          </div>

          <Link
            to="/help"
            className="px-4 py-2 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 transition"
          >
            Help
          </Link>
        </div>
      </header>

      {/* Main */}
      <main className="flex-1 flex items-center justify-center px-4 py-12">
        <div className="w-full max-w-md bg-white rounded-lg shadow-lg p-8">
          {/* Title */}
          <div className="text-center mb-8">
            <div className="flex justify-center mb-4">
              <div className="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
                <Lock className="w-6 h-6 text-blue-600" />
              </div>
            </div>

            <h1 className="text-3xl font-bold text-gray-900 mb-2">
              {twoFARequired ? "Enter OTP" : "Sign In"}
            </h1>

            <p className="text-gray-600">
              {twoFARequired
                ? "We sent a 6-digit code to your email"
                : "to access your account"}
            </p>
          </div>

          {/* Error */}
          {err && (
            <div className="mb-4 text-sm text-red-600 bg-red-50 border border-red-200 rounded-lg px-4 py-2">
              {err}
            </div>
          )}

          {/* FORM */}
          {!twoFARequired ? (
            <form onSubmit={submit} className="space-y-5">
              {/* Email */}
              <div>
                <label className="block text-sm font-semibold text-gray-900 mb-2">
                  Email Address
                </label>
                <input
                  type="email"
                  required
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="w-full px-4 py-3 border rounded-lg focus:ring-2 focus:ring-blue-500"
                />
              </div>

              {/* Password */}
              <div>
                <label className="block text-sm font-semibold text-gray-900 mb-2">
                  Password
                </label>

                <div className="relative">
                  <input
                    type={showPassword ? "text" : "password"}
                    required
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="w-full px-4 py-3 pr-12 border rounded-lg focus:ring-2 focus:ring-blue-500"
                  />

                  <button
                    type="button"
                    onClick={() => setShowPassword((v) => !v)}
                    className="absolute inset-y-0 right-3 flex items-center"
                  >
                    {showPassword ? (
                      <EyeOff className="w-5 h-5" />
                    ) : (
                      <Eye className="w-5 h-5" />
                    )}
                  </button>
                </div>
              </div>

              <button
                type="submit"
                disabled={loading}
                className="w-full py-3 bg-blue-600 text-white rounded-lg"
              >
                {loading ? "Signing in..." : "Sign In"}
              </button>
            </form>
          ) : (
            <form onSubmit={verifyOTP} className="space-y-5">
              <input
                type="text"
                maxLength={6}
                placeholder="Enter 6-digit OTP"
                value={otp}
                onChange={(e) => setOtp(e.target.value)}
                className="w-full px-4 py-3 border rounded-lg text-center tracking-widest text-lg"
                required
              />

              <button
                type="submit"
                disabled={loading}
                className="w-full py-3 bg-blue-600 text-white rounded-lg"
              >
                {loading ? "Verifying..." : "Verify OTP"}
              </button>
            </form>
          )}
        </div>
      </main>
    </div>
  );
}
