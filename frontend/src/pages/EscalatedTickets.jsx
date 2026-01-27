import { useEffect, useState } from "react";
import api from "../api/axios";

export default function EscalatedTickets() {
  const [tickets, setTickets] = useState([]);

  useEffect(() => {
    api.get("/admin/escalations").then(res => setTickets(res.data));
  }, []);

  return (
    <div>
      <h2 className="text-xl font-semibold mb-4 text-red-600">
        Escalated Tickets
      </h2>

      {tickets.map(t => (
        <div key={t.id} className="card border-l-4 border-red-500">
          <p>{t.title}</p>
          <p>Status: {t.status}</p>
        </div>
      ))}
    </div>
  );
}
