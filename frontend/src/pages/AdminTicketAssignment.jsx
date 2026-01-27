import { useEffect, useState } from "react";
import {
  getAdminTickets,
  assignTicket,
  closeTicket,
} from "../api/ticket.api";
import api from "../api/axios";
import { SLABadge } from "../components/SLABadge";

export default function AdminTicketAssignment() {
  const [tickets, setTickets] = useState([]);
  const [engineers, setEngineers] = useState([]);

  useEffect(() => {
    Promise.all([
      getAdminTickets(),
      api.get("/admin/support-engineers"),
    ]).then(([t, e]) => {
      setTickets(t.data);
      setEngineers(e.data);
    });
  }, []);

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold">Admin Ticket Review</h1>

      {tickets.map((ticket) => (
        <TicketCard
          key={ticket.id}
          ticket={ticket}
          engineers={engineers}
          onUpdate={() =>
            setTickets((prev) => prev.filter((t) => t.id !== ticket.id))
          }
        />
      ))}
    </div>
  );
}

function TicketCard({ ticket, engineers, onUpdate }) {
  const [engineer, setEngineer] = useState("");

  const assign = async () => {
    await assignTicket(ticket.id, engineer);
    onUpdate();
  };

  const close = async () => {
    await closeTicket(ticket.id);
    onUpdate();
  };

  return (
    <div className="bg-white border rounded-xl p-4 shadow-sm space-y-3">
      <div className="flex justify-between">
        <h2 className="font-medium">{ticket.title}</h2>
        <SLABadge targetAt={ticket.target_at} />
      </div>

      <p className="text-sm text-gray-600">{ticket.description}</p>
      <p className="text-xs">Status: {ticket.status}</p>

      {ticket.status === "customer_created" && (
        <div className="flex gap-3">
          <select
            value={engineer}
            onChange={(e) => setEngineer(e.target.value)}
            className="border rounded px-3 py-1 text-sm"
          >
            <option value="">Select engineer</option>
            {engineers.map((e) => (
              <option key={e.id} value={e.id}>
                {e.user.name}
              </option>
            ))}
          </select>

          <button
            onClick={assign}
            className="bg-blue-600 text-white px-4 py-1.5 rounded text-sm"
          >
            Assign
          </button>
        </div>
      )}

      {ticket.status === "resolved_by_support" && (
        <button
          onClick={close}
          className="bg-red-600 text-white px-4 py-1.5 rounded text-sm"
        >
          Close Ticket
        </button>
      )}
    </div>
  );
}
