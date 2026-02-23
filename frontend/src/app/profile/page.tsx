"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/features/auth/AuthContext";
import { updateMe } from "@/lib/api";

export default function ProfilePage() {
  const { user, setUser, loading: authLoading } = useAuth();
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    if (authLoading) return; // 初期化中
    if (user === null) {
      router.replace("/auth/login");
      return;
    }
    setUsername(user.username ?? "");
  }, [user, authLoading, router]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess(false);
    if (!username.trim()) {
      setError("ユーザー名を入力してください。");
      return;
    }
    setLoading(true);
    try {
      const updated = await updateMe({ username });
      setUser(updated);
      setSuccess(true);
    } catch {
      setError("更新に失敗しました。もう一度お試しください。");
    } finally {
      setLoading(false);
    }
  };

  if (authLoading) {
    return (
      <main>
        <div className="h-8 w-40 animate-pulse rounded bg-slate-200" />
      </main>
    );
  }

  return (
    <main className="max-w-md">
      <h1 className="mb-6 text-2xl font-bold text-slate-800">プロフィール</h1>

      <div className="mb-6 rounded-xl border border-slate-200 bg-white p-6 shadow-sm">
        <p className="text-sm text-slate-500">メールアドレス</p>
        <p className="mt-1 font-medium text-slate-800">{user?.email}</p>
      </div>

      <form onSubmit={handleSubmit} className="rounded-xl border border-slate-200 bg-white p-6 shadow-sm space-y-4">
        <h2 className="text-lg font-semibold text-slate-800">情報の編集</h2>

        <div>
          <label htmlFor="username" className="block text-sm font-medium text-slate-700">
            ユーザー名
          </label>
          <input
            id="username"
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-pink-400 focus:outline-none focus:ring-1 focus:ring-pink-400"
          />
        </div>

        {error && <p className="text-sm text-red-500">{error}</p>}
        {success && <p className="text-sm text-green-500">更新しました！</p>}

        <button
          type="submit"
          disabled={loading}
          className="w-full rounded-lg bg-pink-400 py-2 text-sm font-semibold text-white hover:bg-pink-500 disabled:opacity-50 transition"
        >
          {loading ? "更新中..." : "変更を保存"}
        </button>
      </form>
    </main>
  );
}
