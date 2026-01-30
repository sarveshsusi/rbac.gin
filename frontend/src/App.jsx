import { Routes, Route, Navigate } from "react-router-dom";
import { useAuth } from "./auth/useAuth";
import RequireAuth from "./auth/RequireAuth";
import RequireGuest from "./auth/RequireGuest";

import LoginPage from "./pages/Login-page";
import VerifyOTP from "./pages/VerifyOTP";
import ResetPassword from "./pages/ResetPassword";
import ForgotPasswordPage from "./pages/ForgotPasswordPage";

import Dashboard from "./pages/Dashboard";
import AdminUsers from "./pages/AdminUsers";
import AppLayout from "./layout/AppLayout";

// ✅ ADMIN
import AdminProducts from "./pages/AdminProducts";
import AdminAMCs from "./pages/AdminAMCs";
import AdminTicketAssignment from "./pages/AdminTicketAssignment";
import AdminProductAssignment from "./pages/AdminProductAssignment";
// ✅ SUPPORT
import SupportTickets from "./pages/SupportTickets";

// ✅ CUSTOMER
import CustomerTickets from "./pages/CustomerTickets";
import CustomerAMCs from "./pages/CustomerAMCs";

// ✅ SHARED
import CreateTicket from "./pages/CreateTicket";
import TicketDetail from "./pages/TicketDetail";
import AdminCreateProductAllInOne from "./pages/AdminCreateProductAllInOne";
import AdminCreateTicket from "./pages/AdminCreateTicket";

function App() {
  const { loading } = useAuth();
  if (loading) return null;

  return (
    <Routes>
      {/* =====================
          PUBLIC / GUEST
      ====================== */}
      <Route
        path="/login"
        element={
          <RequireGuest>
            <LoginPage />
          </RequireGuest>
        }
      />

      <Route
        path="/verify-otp"
        element={
          <RequireGuest>
            <VerifyOTP />
          </RequireGuest>
        }
      />

      <Route path="/reset-password" element={<ResetPassword />} />
      <Route path="/forgot-password" element={<ForgotPasswordPage />} />

      {/* =====================
          PROTECTED (AUTH)
      ====================== */}
      <Route element={<RequireAuth />}>
        <Route element={<AppLayout />}>
          {/* DASHBOARDS */}
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/customer/dashboard" element={<Dashboard />} />

          {/* ADMIN */}
          <Route path="/admin/users" element={<AdminUsers />} />
          <Route path="/admin/products" element={<AdminCreateProductAllInOne />} />
          <Route path="/admin/assign-products" element={<AdminProductAssignment />} />
          <Route path="/admin/amc" element={<AdminAMCs />} />
          <Route path="/admin/tickets" element={<AdminTicketAssignment />} />
          <Route path="/admin/tickets/new" element={<AdminCreateTicket />} />

          {/* SUPPORT */}
          <Route path="/support/tickets" element={<SupportTickets />} />

          {/* CUSTOMER */}
          <Route path="/customer/tickets" element={<CustomerTickets />} />
          <Route path="/customer/amc" element={<CustomerAMCs />} />

          {/* TICKETS (SHARED) */}
          <Route path="/tickets/new" element={<CreateTicket />} />
          <Route path="/tickets/:id" element={<TicketDetail />} />
        </Route>
      </Route>

      {/* =====================
          FALLBACK
      ====================== */}
      <Route path="*" element={<Navigate to="/dashboard" replace />} />
    </Routes>
  );
}

export default App;
