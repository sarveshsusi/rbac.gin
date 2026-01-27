import { useEffect, useState } from "react";
import { getSupportTickets, resolveTicket } from "../api/ticket.api";
import { SLABadge } from "../components/SLABadge";

export default function SupportTickets() {
  const [tickets, setTickets] = useState([]);

  useEffect(() => {
    getSupportTickets().then((res) => setTickets(res.data));
  }, []);

  const handleResolve = async (id) => {
    await resolveTicket(id);
    setTickets((prev) =>
      prev.map((t) =>
        t.id === id
          ? { ...t, status: "resolved_by_support" }
          : t
      )
    );
  };

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-semibold">My Assigned Tickets</h1>

      {tickets.map((t) => (
        <div
          key={t.id}
          className="bg-white p-4 rounded-xl border shadow-sm space-y-2"
        >
          <div className="flex justify-between">
            <h2 className="font-medium">{t.title}</h2>
            <SLABadge targetAt={t.target_at} />
          </div>

          <p className="text-sm text-gray-600">{t.description}</p>

          <div className="flex justify-between items-center">
            <span className="text-xs text-gray-500">
              Status: {t.status}
            </span>

            {t.status === "assigned_to_support" && (
              <button
                onClick={() => handleResolve(t.id)}
                className="bg-green-600 text-white px-4 py-1.5 rounded text-sm"
              >
                Mark Resolved
              </button>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}
