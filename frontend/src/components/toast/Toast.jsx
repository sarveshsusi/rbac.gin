 export function Toast({ message, type, onClose }) {
  const styles = {
    success: "bg-green-600",
    error: "bg-red-600",
    info: "bg-blue-600",
  };

  return (
    <div
      className={`
        ${styles[type]}
        text-white px-4 py-3 rounded-xl shadow-lg
        flex items-start justify-between gap-3
        animate-toast-in
      `}
    >
      <p className="text-sm leading-snug">{message}</p>

      <button
        onClick={onClose}
        className="text-white/80 hover:text-white text-lg leading-none"
      >
        ×
      </button>
    </div>
  );
}
