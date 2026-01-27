import { useEffect, useState } from "react";

export default function useSLATimer(targetAt) {
  const [timeLeft, setTimeLeft] = useState("");

  useEffect(() => {
    if (!targetAt) return;

    const update = () => {
      const diff = new Date(targetAt) - new Date();

      if (diff <= 0) {
        setTimeLeft("SLA BREACHED");
        return;
      }

      const h = Math.floor(diff / 3600000);
      const m = Math.floor((diff % 3600000) / 60000);

      setTimeLeft(`${h}h ${m}m`);
    };

    update();
    const interval = setInterval(update, 60000);
    return () => clearInterval(interval);
  }, [targetAt]);

  return timeLeft;
}
