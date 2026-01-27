import { useEffect, useState } from "react";
import { getCustomerAMCs } from "../api/amc.api";

export default function CustomerAMCs() {
  const [amcs, setAmcs] = useState([]);

  useEffect(() => {
    getCustomerAMCs().then(res => setAmcs(res.data));
  }, []);

  return (
    <div>
      <h2 className="text-xl font-semibold mb-4">My AMC</h2>

      {amcs.map(a => (
        <div key={a.id} className="card">
          <p>Valid till: {a.end_date}</p>
        </div>
      ))}
    </div>
  );
}
