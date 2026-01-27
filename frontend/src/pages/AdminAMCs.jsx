import { useEffect, useState } from "react";
import { getAdminAMCs } from "../api/amc.api";

export default function AdminAMCs() {
  const [amcs, setAmcs] = useState([]);

  useEffect(() => {
    getAdminAMCs().then(res => setAmcs(res.data));
  }, []);

  return (
    <div>
      <h2 className="text-xl font-semibold mb-4">AMC Contracts</h2>

      {amcs.map(a => (
        <div key={a.id} className="card">
          <p>Start: {a.start_date}</p>
          <p>End: {a.end_date}</p>
          <p>Status: {a.status}</p>
        </div>
      ))}
    </div>
  );
}
