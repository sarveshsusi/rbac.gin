import { useEffect, useState } from "react";
import { getCustomerTickets } from "../api/ticket.api";
import { Link } from "react-router-dom";
import { SLABadge } from "../components/SLABadge";

export default function CustomerTickets() {
  const [tickets, setTickets] = useState([]);

  useEffect(() => {
    getCustomerTickets().then((res) => setTickets(res.data));
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-semibold">My Tickets</h1>

        <Link
          to="/tickets/new"
          className="bg-blue-600 text-white px-4 py-2 rounded"
        >
          + Create Ticket
        </Link>
      </div>

      {tickets.map((t) => (
        <div key={t.id} className="bg-white p-4 rounded-xl shadow">
          <div className="flex justify-between">
            <h3 className="font-medium">{t.title}</h3>
            <SLABadge targetAt={t.target_at} />
          </div>

          <p className="text-sm text-gray-600">{t.description}</p>
          <p className="text-xs mt-2 text-gray-500">Status: {t.status}</p>
        </div>
      ))}
    </div>
  );
}
