import useSLATimer from "./hooks/useSLATimer";

export function SLABadge({ targetAt }) {
  const timeLeft = useSLATimer(targetAt);
  const breached = timeLeft === "SLA BREACHED";

  return (
    <span
      className={`px-2 py-1 rounded text-xs font-medium ${
        breached
          ? "bg-red-100 text-red-700"
          : "bg-green-100 text-green-700"
      }`}
    >
      {timeLeft}
    </span>
  );
}
