import { useEffect, useRef, useState } from "react";
import { useAuth } from "../auth/useAuth";
import usePageTitle from "../components/hooks/usePageTitle";
import { getAdminDashboard } from "../api/dashboard.api";

const CARD_COUNT = 4;
const CARD_HEIGHT = 96; // 👈 keeps skeleton + card same height

export default function Dashboard() {
  usePageTitle("Dashboard • EMERD");

  const { user } = useAuth();
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);
  const [activeIndex, setActiveIndex] = useState(0);
  const scrollRef = useRef(null);

  /* ================= FETCH ================= */
  useEffect(() => {
    if (!user || user.role !== "admin") return;

    let mounted = true;

    (async () => {
      try {
        const res = await getAdminDashboard();
        if (mounted) setStats(res.data);
      } finally {
        if (mounted) setLoading(false);
      }
    })();

    return () => (mounted = false);
  }, [user]);

  if (!user || user.role !== "admin") return null;

  const cards = stats
    ? [
        { label: "Users", value: stats.users },
        { label: "Total Tickets", value: stats.tickets },
        { label: "Pending Tickets", value: stats.pending_tickets },
        { label: "Closed Tickets", value: stats.closed_tickets },
      ]
    : Array(CARD_COUNT).fill(null);

  const onScroll = () => {
    if (!scrollRef.current) return;
    const { scrollLeft, clientWidth } = scrollRef.current;
    setActiveIndex(Math.round(scrollLeft / (clientWidth * 0.85)));
  };

  return (
    <div className="space-y-6 px-3 pb-safe">
      <h1 className="text-2xl font-bold">Dashboard</h1>

      {/* ================= MOBILE CAROUSEL ================= */}
      <div className="sm:hidden">
        <div
          ref={scrollRef}
          onScroll={onScroll}
          className="
            flex gap-4 overflow-x-auto
            snap-x snap-mandatory
            scrollbar-hide
            touch-pan-x
            scroll-smooth
            pb-2
          "
          style={{ WebkitOverflowScrolling: "touch" }}
        >
          {cards.map((item, i) => (
            <div
              key={i}
              className="snap-center shrink-0"
              style={{ width: "85%" }} // 👈 peek next card
            >
              {loading ? (
                <StatSkeleton />
              ) : (
                <Stat label={item.label} value={item.value} />
              )}
            </div>
          ))}
        </div>

        {/* DOTS */}
        <div className="flex justify-center gap-2 mt-2">
          {cards.map((_, i) => (
            <span
              key={i}
              className={`h-2 w-2 rounded-full transition ${
                activeIndex === i ? "bg-blue-600" : "bg-gray-300"
              }`}
            />
          ))}
        </div>
      </div>

      {/* ================= DESKTOP GRID ================= */}
      <div className="hidden sm:grid grid-cols-2 lg:grid-cols-4 gap-4">
        {cards.map((item, i) =>
          loading ? (
            <StatSkeleton key={i} />
          ) : (
            <Stat key={item.label} label={item.label} value={item.value} />
          )
        )}
      </div>
    </div>
  );
}

/* ================= CARD ================= */
function Stat({ label, value }) {
  return (
    <div
      className="bg-white border border-slate-200 rounded-xl p-4 shadow-sm"
      style={{ minHeight: CARD_HEIGHT }}
    >
      <p className="text-sm text-gray-500">{label}</p>
      <p className="text-2xl font-semibold mt-1">{value}</p>
    </div>
  );
}

/* ================= SKELETON ================= */
function StatSkeleton() {
  return (
    <div
      className="bg-white border border-slate-200 rounded-xl p-4 shadow-sm animate-pulse"
      style={{ minHeight: CARD_HEIGHT }}
    >
      <div className="h-3 w-24 bg-gray-200 rounded mb-3" />
      <div className="h-8 w-16 bg-gray-300 rounded" />
    </div>
  );
}
