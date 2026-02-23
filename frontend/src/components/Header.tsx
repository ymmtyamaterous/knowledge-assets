"use client";

import Link from "next/link";
import { useAuth } from "@/features/auth/AuthContext";

export default function Header() {
  const { user, loading, logout } = useAuth();

  return (
    <header className="border-b border-pink-100 bg-white shadow-sm">
      <div className="mx-auto flex max-w-5xl items-center justify-between px-6 py-3">
        <Link href="/" className="text-xl font-bold text-pink-500">
          🌸 アセナレ
        </Link>

        <nav className="flex items-center gap-4 text-sm">
          <Link href="/courses" className="text-slate-600 hover:text-pink-500">
            コース
          </Link>
          <Link href="/glossary" className="text-slate-600 hover:text-pink-500">
            用語辞典
          </Link>
          <Link href="/progress" className="text-slate-600 hover:text-pink-500">
            進捗
          </Link>

          {!loading && (
            <>
              {user ? (
                <div className="flex items-center gap-3">
                  <Link href="/profile" className="text-slate-700 hover:text-pink-500">
                    {user.username}
                  </Link>
                  <button
                    onClick={logout}
                    className="rounded-md border border-slate-300 px-3 py-1 text-slate-600 hover:bg-slate-50"
                  >
                    ログアウト
                  </button>
                </div>
              ) : (
                <div className="flex items-center gap-2">
                  <Link
                    href="/auth/login"
                    className="rounded-md border border-pink-300 px-3 py-1 text-pink-600 hover:bg-pink-50"
                  >
                    ログイン
                  </Link>
                  <Link
                    href="/auth/register"
                    className="rounded-md bg-pink-500 px-3 py-1 text-white hover:bg-pink-600"
                  >
                    登録
                  </Link>
                </div>
              )}
            </>
          )}
        </nav>
      </div>
    </header>
  );
}
