import { useAuth } from "../auth/useAuth";
import {
  resolveTicket,
  closeTicket,
  submitFeedback,
} from "../api/ticket.api";

export default function TicketDetail({ ticket, onUpdate }) {
  const { user } = useAuth();

  /* =====================
     SUPPORT: RESOLVE
  ====================== */
  const handleResolve = async () => {
    await resolveTicket(ticket.id);
    onUpdate?.();
  };

  /* =====================
     ADMIN: CLOSE
  ====================== */
  const handleClose = async () => {
    await closeTicket(ticket.id);
    onUpdate?.();
  };

  /* =====================
     CUSTOMER: FEEDBACK
  ====================== */
  const handleFeedback = async (rating) => {
    await submitFeedback(ticket.id, {
      rating,
      comment: "",
    });
    onUpdate?.();
  };

  return (
    <div className="bg-white p-5 rounded-xl shadow space-y-3">
      <div>
        <h3 className="font-semibold text-lg">{ticket.title}</h3>
        <p className="text-sm text-gray-600">{ticket.description}</p>
      </div>

      <div className="text-sm text-gray-500">
        Status: <span className="font-medium">{ticket.status}</span>
      </div>

      {/* =====================
          SUPPORT ACTION
      ====================== */}
      {user.role === "support" &&
        ticket.status === "assigned_to_support" && (
          <button
            onClick={handleResolve}
            className="bg-green-600 text-white px-4 py-1.5 rounded text-sm"
          >
            Resolve Ticket
          </button>
        )}

      {/* =====================
          ADMIN ACTION
      ====================== */}
      {user.role === "admin" &&
        ticket.status === "resolved_by_support" && (
          <button
            onClick={handleClose}
            className="bg-red-600 text-white px-4 py-1.5 rounded text-sm"
          >
            Close Ticket
          </button>
        )}

      {/* =====================
          CUSTOMER FEEDBACK
      ====================== */}
      {user.role === "customer" &&
        ticket.status === "closed_by_admin" && (
          <div className="flex gap-2">
            {[1, 2, 3, 4, 5].map((star) => (
              <button
                key={star}
                onClick={() => handleFeedback(star)}
                className="px-2 py-1 border rounded text-sm hover:bg-yellow-100"
              >
                ⭐ {star}
              </button>
            ))}
          </div>
        )}
    </div>
  );
}
