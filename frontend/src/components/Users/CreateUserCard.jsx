import { useState } from "react";
import { Info, UserPlus, Loader2, CheckCircle } from "lucide-react";
import api from "../../api/axios";

const ROLE_MAP = {
  Admin: "admin",
  Support: "support",
  Customer: "customer",
};

export default function CreateUserCard() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [role, setRole] = useState("Admin");

  // 🔹 Customer-only fields
  const [company, setCompany] = useState("");
  const [phone, setPhone] = useState("");
  const [address, setAddress] = useState("");

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const isCustomer = role === "Customer";

  const submit = async () => {
    setError("");
    setSuccess("");

    if (!name.trim()) return setError("Name is required");
    if (!email.trim()) return setError("Email is required");

    if (isCustomer && !company.trim()) {
      return setError("Company name is required for customers");
    }

    try {
      setLoading(true);

      await api.post("/admin/users", {
        name: name.trim(),
        email: email.trim().toLowerCase(),
        role: ROLE_MAP[role],

        // 👇 Only sent for customers
        ...(isCustomer && {
          company: company.trim(),
          phone: phone.trim(),
          address: address.trim(),
        }),
      });

      setSuccess("User created. Password setup email sent.");

      // reset
      setName("");
      setEmail("");
      setRole("Admin");
      setCompany("");
      setPhone("");
      setAddress("");
    } catch (err) {
      setError(
        err.response?.data?.error ||
          "Failed to create user. Please try again."
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div
      className="
        w-full bg-white border border-slate-200
        shadow-[0_10px_40px_rgba(15,23,42,0.08)]
        rounded-none px-4 py-4
        sm:rounded-2xl sm:px-6 sm:py-6 sm:mx-4
        lg:mx-auto lg:max-w-5xl lg:px-8 lg:py-8
      "
    >
      {/* HEADER */}
      <div>
        <h2 className="text-lg sm:text-xl font-semibold text-slate-900">
          Create New User
        </h2>
        <p className="mt-1 text-sm text-slate-500">
          Add a new team member. They will receive an email to set their password.
        </p>
      </div>

      {/* BASIC FORM */}
      <div className="mt-6 grid grid-cols-1 gap-5 sm:grid-cols-2">
        <div>
          <label className="block text-sm font-medium text-slate-700">
            Full Name
          </label>
          <input
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Jane Smith"
            className="mt-1 w-full rounded-xl border border-slate-300
                       px-4 py-3 text-sm
                       focus:ring-2 focus:ring-blue-500 focus:outline-none"
          />
        </div>

        <div>
          <label className="block text-sm font-medium text-slate-700">
            Email Address
          </label>
          <input
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="jane@example.com"
            className="mt-1 w-full rounded-xl border border-slate-300
                       px-4 py-3 text-sm
                       focus:ring-2 focus:ring-blue-500 focus:outline-none"
          />
        </div>
      </div>

      {/* ROLE + INFO */}
      <div className="mt-6 grid grid-cols-1 gap-4 lg:grid-cols-2">
        <div>
          <label className="block text-sm font-medium text-slate-700">
            Role
          </label>
          <select
            value={role}
            onChange={(e) => setRole(e.target.value)}
            className="mt-1 w-full rounded-xl border border-slate-300
                       px-4 py-3 text-sm
                       focus:ring-2 focus:ring-blue-500 focus:outline-none"
          >
            <option>Admin</option>
            <option>Support</option>
            <option>Customer</option>
          </select>
        </div>

        <div className="flex items-start gap-3 rounded-xl bg-blue-50
                        px-4 py-3 text-sm text-blue-700">
          <Info className="mt-0.5 h-4 w-4 shrink-0" />
          <p>
            The user will receive a secure email link to set their password.
            The link expires after 24 hours.
          </p>
        </div>
      </div>

      {/* 👇 CUSTOMER FIELDS (SMOOTH + NO LAYOUT BREAK) */}
      <div
        className={`
          overflow-hidden transition-all duration-500 ease-in-out
          ${isCustomer ? "max-h-[500px] opacity-100 mt-6" : "max-h-0 opacity-0"}
        `}
      >
        <div className="grid grid-cols-1 gap-5 sm:grid-cols-2">
          <div>
            <label className="block text-sm font-medium text-slate-700">
              Company Name
            </label>
            <input
              value={company}
              onChange={(e) => setCompany(e.target.value)}
              placeholder="Acme Corp"
              className="mt-1 w-full rounded-xl border border-slate-300
                         px-4 py-3 text-sm
                         focus:ring-2 focus:ring-blue-500 focus:outline-none"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-slate-700">
              Phone
            </label>
            <input
              value={phone}
              onChange={(e) => setPhone(e.target.value)}
              placeholder="+91 9876543210"
              className="mt-1 w-full rounded-xl border border-slate-300
                         px-4 py-3 text-sm
                         focus:ring-2 focus:ring-blue-500 focus:outline-none"
            />
          </div>

          <div className="sm:col-span-2">
            <label className="block text-sm font-medium text-slate-700">
              Address
            </label>
            <textarea
              value={address}
              onChange={(e) => setAddress(e.target.value)}
              placeholder="Company address"
              rows={3}
              className="mt-1 w-full rounded-xl border border-slate-300
                         px-4 py-3 text-sm
                         focus:ring-2 focus:ring-blue-500 focus:outline-none"
            />
          </div>
        </div>
      </div>

      {/* FEEDBACK */}
      {error && <p className="mt-4 text-sm text-red-600">{error}</p>}

      {success && (
        <div className="mt-4 flex items-center gap-2 text-sm text-green-600">
          <CheckCircle size={16} />
          {success}
        </div>
      )}

      {/* ACTION */}
      <div className="mt-8 flex justify-end">
        <button
          onClick={submit}
          disabled={loading}
          className="
            inline-flex w-full sm:w-auto items-center justify-center gap-2
            rounded-xl bg-blue-600 px-6 py-3
            text-sm font-semibold text-white
            shadow-lg shadow-blue-600/30
            hover:bg-blue-700 transition
            disabled:opacity-60 disabled:cursor-not-allowed
          "
        >
          {loading ? (
            <Loader2 className="animate-spin" size={16} />
          ) : (
            <UserPlus size={16} />
          )}
          Create User
        </button>
      </div>
    </div>
  );
}
