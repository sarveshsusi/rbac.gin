import { useEffect, useState } from "react";
import {
  getAdminTickets,
  assignTicket,
  closeTicket, // Although Admin doesn't close usually, keeping if needed or removing
} from "../api/ticket.api";
import api from "../api/axios";
import { SLABadge } from "../components/SLABadge";

import { Link } from "react-router-dom";

export default function AdminTicketAssignment() {
  const [tickets, setTickets] = useState([]);
  const [engineers, setEngineers] = useState([]);

  useEffect(() => {
    Promise.all([
      getAdminTickets(),
      api.get("/admin/support-engineers"), // Check if this endpoint exists
    ]).then(([t, e]) => {
      setTickets(t.data);
      setEngineers(e.data);
    });
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-semibold">Admin Ticket Review</h1>
        <Link to="/admin/tickets/new" className="bg-blue-600 text-white px-4 py-2 rounded text-sm hover:bg-blue-700">
          + Create Ticket (On Behalf)
        </Link>
      </div>

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
  const [priority, setPriority] = useState("Standard");
  const [supportMode, setSupportMode] = useState("Remote");
  const [serviceType, setServiceType] = useState("Service");

  const assign = async () => {
    if (!engineer) return alert("Please select an engineer");

    try {
      await assignTicket(ticket.id, {
        engineer_id: engineer,
        priority: priority,
        support_mode: supportMode,
        service_call_type: serviceType
      });
      alert("Ticket assigned successfully");
      onUpdate();
    } catch (err) {
      alert("Failed to assign ticket: " + (err.response?.data?.error || err.message));
    }
  };

  return (
    <div className="bg-white border rounded-xl p-4 shadow-sm space-y-3">
      <div className="flex justify-between">
        <h2 className="font-medium">{ticket.title}</h2>
        <SLABadge targetAt={ticket.target_at} />
      </div>

      <p className="text-sm text-gray-600">{ticket.description}</p>
      <div className="flex gap-4 text-xs text-gray-500">
        <span>Product: {ticket.product_id}</span>
        <span>Status: {ticket.status}</span>
      </div>

      {ticket.status === "Open" && (
        <div className="grid grid-cols-2 gap-3 bg-gray-50 p-3 rounded">

          {/* CONFIGURATION */}
          <div className="space-y-2">
            <label className="block text-xs font-semibold">Priority</label>
            <select value={priority} onChange={e => setPriority(e.target.value)} className="w-full text-sm border rounded p-1">
              <option value="Low">Low</option>
              <option value="Standard">Standard</option>
              <option value="Critical">Critical</option>
            </select>

            <label className="block text-xs font-semibold">Support Mode</label>
            <select value={supportMode} onChange={e => setSupportMode(e.target.value)} className="w-full text-sm border rounded p-1">
              <option value="On-site">On-site</option>
              <option value="Remote">Remote</option>
              <option value="Phone">Phone</option>
            </select>
          </div>

          <div className="space-y-2">
            <label className="block text-xs font-semibold">Service Type</label>
            <select value={serviceType} onChange={e => setServiceType(e.target.value)} className="w-full text-sm border rounded p-1">
              <option value="Warranty">Warranty</option>
              <option value="Service">Service</option>
              <option value="AMC">AMC</option>
            </select>

            <label className="block text-xs font-semibold">Assign Engineer</label>
            <select
              value={engineer}
              onChange={(e) => setEngineer(e.target.value)}
              className="w-full text-sm border rounded p-1"
            >
              <option value="">Select engineer</option>
              {engineers.map((e) => (
                <option key={e.id} value={e.id}>
                  {e.user.name}
                </option>
              ))}
            </select>
          </div>

          <div className="col-span-2 mt-2">
            <button
              onClick={assign}
              className="w-full bg-blue-600 text-white px-4 py-2 rounded text-sm font-medium hover:bg-blue-700 transition"
            >
              Confirm Assignment
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
