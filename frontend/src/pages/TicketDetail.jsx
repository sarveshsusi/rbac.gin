import { useAuth } from "../auth/useAuth";
import {
  startTicket,
  closeTicket,
  submitFeedback,
} from "../api/ticket.api";
import { useState } from "react";

export default function TicketDetail({ ticket, onUpdate }) {
  const { user } = useAuth();
  const [proof, setProof] = useState(null);
  const [uploading, setUploading] = useState(false);

  /* =====================
     SUPPORT: START
  ====================== */
  const handleStart = async () => {
    try {
      await startTicket(ticket.id);
      onUpdate?.();
    } catch (err) {
      alert("Failed to start ticket");
    }
  };

  /* =====================
     SUPPORT: CLOSE
  ====================== */
  const handleSupportClose = async () => {
    if (!proof) return alert("Proof image is required");

    setUploading(true);
    try {
      const formData = new FormData();
      formData.append("proof", proof);
      await closeTicket(ticket.id, formData);
      onUpdate?.();
      alert("Ticket closed successfully");
    } catch (err) {
      alert("Failed: " + (err.response?.data?.error || err.message));
    } finally {
      setUploading(false);
    }
  };

  /* =====================
     ADMIN: CLOSE (Force)
  ====================== */
  const handleAdminClose = async () => {
    // Admin closing might not require proof or different endpoint, but using same for now implies support flow.
    // Ideally Admin uses a different endpoint or we skip proof. 
    // For now hiding Admin Close or making it just an alert that Support should close.
    alert("Please assign to Support to close with proof.");
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
      {user.role === "support" && (
        <div className="flex gap-2 items-center mt-2">
          {ticket.status === "Assigned" && (
            <button
              onClick={handleStart}
              className="bg-blue-600 text-white px-4 py-1.5 rounded text-sm"
            >
              Start Work
            </button>
          )}

          {ticket.status === "In Progress" && (
            <div className="flex flex-col gap-2">
              <input
                type="file"
                accept="image/*"
                onChange={(e) => setProof(e.target.files[0])}
                className="text-xs"
              />
              <button
                onClick={handleSupportClose}
                disabled={uploading || !proof}
                className="bg-green-600 text-white px-4 py-1.5 rounded text-sm disabled:opacity-50"
              >
                {uploading ? "Closing..." : "Close Ticket (with Proof)"}
              </button>
            </div>
          )}
        </div>
      )}

      {/* =====================
          ADMIN ACTION
      ====================== */}
      {user.role === "admin" &&
        ticket.status === "In Progress" && ( // Example condition
          <button
            onClick={handleAdminClose}
            className="bg-gray-500 text-white px-4 py-1.5 rounded text-sm"
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
