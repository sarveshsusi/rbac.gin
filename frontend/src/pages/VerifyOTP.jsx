import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../auth/useAuth";

export default function VerifyOTP() {
  const { verify2FA } = useAuth();
  const navigate = useNavigate();

  const [otp, setOtp] = useState("");
  const [err, setErr] = useState("");
  const [loading, setLoading] = useState(false);

  const twoFAToken = sessionStorage.getItem("two_fa_token");

  // ❌ No token → block access
  useEffect(() => {
    if (!twoFAToken) {
      navigate("/login", { replace: true });
    }
  }, [twoFAToken, navigate]);

  const submit = async (e) => {
    e.preventDefault();
    if (loading) return;

    setErr("");
    setLoading(true);

    try {
      await verify2FA(otp, twoFAToken);
      sessionStorage.removeItem("two_fa_token");
      navigate("/", { replace: true });
    } catch (err) {
      setErr(err.response?.data?.error || "Invalid or expired OTP");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4">
      <div className="w-full max-w-md bg-white rounded-lg shadow-lg p-8">
        <h1 className="text-3xl font-bold text-center mb-6">
          Verify OTP
        </h1>

        <p className="text-center text-gray-600 mb-4">
          Enter the 6-digit code sent to your email
        </p>

        {err && (
          <div className="mb-4 text-sm text-red-600 bg-red-50 border px-4 py-2 rounded">
            {err}
          </div>
        )}

        <form onSubmit={submit} className="space-y-5">
          <input
            type="text"
            maxLength={6}
            placeholder="Enter OTP"
            value={otp}
            onChange={(e) => setOtp(e.target.value)}
            className="w-full px-4 py-3 border rounded-lg text-center text-lg tracking-widest"
            required
          />

          <button
            disabled={loading}
            className="w-full py-3 bg-blue-600 text-white rounded-lg"
          >
            {loading ? "Verifying..." : "Verify OTP"}
          </button>
        </form>
      </div>
    </div>
  );
}
