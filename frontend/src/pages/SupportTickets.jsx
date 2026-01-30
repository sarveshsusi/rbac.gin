import { useEffect, useState } from "react";
import { getSupportTickets, startTicket, closeTicket } from "../api/ticket.api";
import { SLABadge } from "../components/SLABadge";

export default function SupportTickets() {
  const [tickets, setTickets] = useState([]);
  const [uploading, setUploading] = useState(false);

  useEffect(() => {
    getSupportTickets().then((res) => setTickets(res.data));
  }, []);

  const handleStart = async (id) => {
    try {
      await startTicket(id);
      const updated = tickets.map(t => t.id === id ? { ...t, status: "In Progress" } : t);
      setTickets(updated);
    } catch (err) {
      alert("Failed to start ticket");
    }
  };

  const handleClose = async (id, file) => {
    if (!file) return alert("Proof image is required");

    setUploading(true);
    try {
      const formData = new FormData();
      formData.append("proof", file);

      await closeTicket(id, formData);

      const updated = tickets.map(t => t.id === id ? { ...t, status: "Closed" } : t);
      setTickets(updated);
      alert("Ticket closed successfully");
    } catch (err) {
      alert("Failed to close ticket: " + (err.response?.data?.error || err.message));
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold">My Assigned Tickets</h1>

      {tickets.map((t) => (
        <TicketRow
          key={t.id}
          ticket={t}
          onStart={() => handleStart(t.id)}
          onClose={(file) => handleClose(t.id, file)}
          uploading={uploading}
        />
      ))}
    </div>
  );
}

function TicketRow({ ticket, onStart, onClose, uploading }) {
  const [proof, setProof] = useState(null);

  return (
    <div className="bg-white p-4 rounded-xl border shadow-sm space-y-2">
      <div className="flex justify-between">
        <h2 className="font-medium">{ticket.title}</h2>
        <SLABadge targetAt={ticket.target_at} />
      </div>

      <p className="text-sm text-gray-600">{ticket.description}</p>

      <div className="flex justify-between items-center mt-4">
        <span className={`text-xs px-2 py-1 rounded ${getStatusColor(ticket.status)}`}>
          {ticket.status}
        </span>

        {ticket.status === "Assigned" && (
          <button
            onClick={onStart}
            className="bg-blue-600 text-white px-4 py-1.5 rounded text-sm hover:bg-blue-700"
          >
            Start Work
          </button>
        )}

        {ticket.status === "In Progress" && (
          <div className="flex gap-2 items-center">
            <input
              type="file"
              accept="image/*"
              onChange={(e) => setProof(e.target.files[0])}
              className="text-xs"
            />
            <button
              onClick={() => onClose(proof)}
              disabled={uploading || !proof}
              className="bg-green-600 text-white px-4 py-1.5 rounded text-sm hover:bg-green-700 disabled:opacity-50"
            >
              {uploading ? "Closing..." : "Close Ticket"}
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

function getStatusColor(status) {
  switch (status) {
    case "Open": return "bg-yellow-100 text-yellow-800";
    case "Assigned": return "bg-blue-100 text-blue-800";
    case "In Progress": return "bg-purple-100 text-purple-800";
    case "Closed": return "bg-green-100 text-green-800";
    default: return "bg-gray-100 text-gray-800";
  }
}
