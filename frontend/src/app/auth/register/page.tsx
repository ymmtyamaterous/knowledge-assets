"use client";

import { useState, type FormEvent } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Eye, EyeOff } from "lucide-react";
import { useAuth } from "@/features/auth/AuthContext";

export default function RegisterPage() {
  const { register } = useAuth();
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);

    if (password.length < 8) {
      setError("パスワードは8文字以上で入力してください。");
      return;
    }
    if (password !== confirmPassword) {
      setError("パスワードが一致しません。");
      return;
    }

    setLoading(true);
    try {
      await register(email, password, username);
      router.push("/");
    } catch (err) {
      setError(err instanceof Error ? err.message : "登録に失敗しました");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="flex min-h-[70vh] items-center justify-center">
      <div className="w-full max-w-md rounded-2xl border border-pink-100 bg-white p-8 shadow-md">
        <h1 className="mb-6 text-2xl font-bold text-pink-500">新規登録</h1>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="mb-1 block text-sm font-medium text-slate-700">ユーザー名</label>
            <input
              type="text"
              required
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-pink-400 focus:outline-none"
            />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium text-slate-700">メールアドレス</label>
            <input
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-pink-400 focus:outline-none"
            />
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium text-slate-700">パスワード</label>
            <div className="relative">
              <input
                type={showPassword ? "text" : "password"}
                required
                minLength={8}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full rounded-lg border border-slate-300 px-3 py-2 pr-10 text-sm focus:border-pink-400 focus:outline-none"
              />
              <button
                type="button"
                onClick={() => setShowPassword((v) => !v)}
                className="absolute inset-y-0 right-2 flex items-center text-slate-400 hover:text-slate-600"
                aria-label={showPassword ? "パスワードを隠す" : "パスワードを表示する"}
              >
                {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
              </button>
            </div>
            <p className="mt-1 text-xs text-slate-400">8文字以上</p>
          </div>
          <div>
            <label className="mb-1 block text-sm font-medium text-slate-700">パスワード（確認）</label>
            <div className="relative">
              <input
                type={showConfirmPassword ? "text" : "password"}
                required
                minLength={8}
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                className="w-full rounded-lg border border-slate-300 px-3 py-2 pr-10 text-sm focus:border-pink-400 focus:outline-none"
              />
              <button
                type="button"
                onClick={() => setShowConfirmPassword((v) => !v)}
                className="absolute inset-y-0 right-2 flex items-center text-slate-400 hover:text-slate-600"
                aria-label={showConfirmPassword ? "パスワードを隠す" : "パスワードを表示する"}
              >
                {showConfirmPassword ? <EyeOff size={18} /> : <Eye size={18} />}
              </button>
            </div>
            {confirmPassword && password !== confirmPassword && (
              <p className="mt-1 text-xs text-red-500">パスワードが一致しません。</p>
            )}
            {confirmPassword && password === confirmPassword && (
              <p className="mt-1 text-xs text-green-500">パスワードが一致しています。</p>
            )}
          </div>

          {error && (
            <p className="rounded-lg bg-red-50 px-3 py-2 text-sm text-red-600">{error}</p>
          )}

          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-lg bg-pink-500 py-2 text-sm font-semibold text-white hover:bg-pink-600 disabled:opacity-50"
          >
            {loading ? "処理中..." : "登録する"}
          </button>
        </form>

        <p className="mt-4 text-center text-sm text-slate-500">
          すでにアカウントをお持ちの方は{" "}
          <Link href="/auth/login" className="text-pink-500 hover:underline">
            ログイン
          </Link>
        </p>
      </div>
    </div>
  );
}
