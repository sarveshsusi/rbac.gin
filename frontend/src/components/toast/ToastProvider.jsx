import { useState, useCallback } from "react";


import { ToastContext } from "./toastContext";

export function ToastProvider({ children }) {
  const [toasts, setToasts] = useState([]);

  const showToast = useCallback((message, type = "success") => {
    const id = crypto.randomUUID();

    setToasts((prev) => [...prev, { id, message, type }]);

    setTimeout(() => {
      setToasts((prev) => prev.filter((t) => t.id !== id));
    }, 3500);
  }, []);

  return (
    <ToastContext.Provider value={{ showToast }}>
      {children}

      {/* TOAST UI */}
      <div className="fixed top-4 right-4 z-50 space-y-3 w-[90%] sm:w-80">
        {toasts.map((t) => (
          <div
            key={t.id}
            className={`
              animate-slide-in rounded-xl px-4 py-3 shadow-lg text-sm font-medium
              ${
                t.type === "success"
                  ? "bg-green-600 text-white"
                  : t.type === "error"
                  ? "bg-red-600 text-white"
                  : "bg-blue-600 text-white"
              }
            `}
          >
            {t.message}
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  );
}


