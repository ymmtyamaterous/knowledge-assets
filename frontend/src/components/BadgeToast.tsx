"use client";

import { useEffect, useState } from "react";
import type { UserBadge } from "@/types/badge";

type BadgeToastProps = {
  badges: UserBadge[];
  onDismiss: () => void;
};

export default function BadgeToast({ badges, onDismiss }: BadgeToastProps) {
  const [visible, setVisible] = useState(true);

  useEffect(() => {
    if (badges.length === 0) return;
    const timer = setTimeout(() => {
      setVisible(false);
      setTimeout(onDismiss, 300);
    }, 4000);
    return () => clearTimeout(timer);
  }, [badges, onDismiss]);

  if (badges.length === 0 || !visible) return null;

  return (
    <div
      className={`fixed bottom-6 right-6 z-50 max-w-xs space-y-2 transition-all duration-300 ${
        visible ? "translate-y-0 opacity-100" : "translate-y-4 opacity-0"
      }`}
    >
      {badges.map((ub) => (
        <div
          key={ub.id}
          className="flex items-start gap-3 rounded-xl bg-gradient-to-r from-pink-500 to-rose-500 p-4 text-white shadow-lg"
        >
          <span className="text-2xl">🏅</span>
          <div>
            <p className="text-xs font-semibold opacity-80">バッジを取得しました！</p>
            <p className="font-bold">{ub.badge.name}</p>
            <p className="mt-0.5 text-xs opacity-80">{ub.badge.description}</p>
          </div>
        </div>
      ))}
    </div>
  );
}
