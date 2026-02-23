"use client";

import { useState } from "react";
import Link from "next/link";
import { useAuth } from "@/features/auth/AuthContext";

export default function Header() {
  const { user, loading, logout } = useAuth();
  const [menuOpen, setMenuOpen] = useState(false);

  return (
    <header className="border-b border-pink-100 bg-white shadow-sm">
      <div className="mx-auto flex max-w-5xl items-center justify-between px-4 py-3 sm:px-6">
        <Link href="/" className="text-xl font-bold text-pink-500">
          🌸 アセナレ
        </Link>

        <button
          onClick={() => setMenuOpen((prev) => !prev)}
          className="rounded-md border border-slate-300 px-3 py-1 text-sm text-slate-700 md:hidden"
          aria-expanded={menuOpen}
          aria-label="メニューを開閉"
        >
          {menuOpen ? "閉じる" : "メニュー"}
        </button>

        <nav className="hidden items-center gap-4 text-sm md:flex">
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
                  <Link href="/profile" className="max-w-32 truncate text-slate-700 hover:text-pink-500">
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

      {menuOpen && (
        <div className="border-t border-pink-100 px-4 py-3 md:hidden">
          <nav className="flex flex-col gap-2 text-sm">
            <Link href="/courses" onClick={() => setMenuOpen(false)} className="rounded-md px-2 py-1 text-slate-700 hover:bg-pink-50 hover:text-pink-600">
              コース
            </Link>
            <Link href="/glossary" onClick={() => setMenuOpen(false)} className="rounded-md px-2 py-1 text-slate-700 hover:bg-pink-50 hover:text-pink-600">
              用語辞典
            </Link>
            <Link href="/progress" onClick={() => setMenuOpen(false)} className="rounded-md px-2 py-1 text-slate-700 hover:bg-pink-50 hover:text-pink-600">
              進捗
            </Link>

            {!loading && (
              <div className="mt-2 border-t border-slate-100 pt-2">
                {user ? (
                  <div className="flex flex-col gap-2">
                    <Link href="/profile" onClick={() => setMenuOpen(false)} className="rounded-md px-2 py-1 text-slate-700 hover:bg-pink-50 hover:text-pink-600">
                      {user.username}
                    </Link>
                    <button
                      onClick={() => {
                        logout();
                        setMenuOpen(false);
                      }}
                      className="rounded-md border border-slate-300 px-3 py-1 text-left text-slate-600 hover:bg-slate-50"
                    >
                      ログアウト
                    </button>
                  </div>
                ) : (
                  <div className="flex gap-2">
                    <Link
                      href="/auth/login"
                      onClick={() => setMenuOpen(false)}
                      className="rounded-md border border-pink-300 px-3 py-1 text-pink-600 hover:bg-pink-50"
                    >
                      ログイン
                    </Link>
                    <Link
                      href="/auth/register"
                      onClick={() => setMenuOpen(false)}
                      className="rounded-md bg-pink-500 px-3 py-1 text-white hover:bg-pink-600"
                    >
                      登録
                    </Link>
                  </div>
                )}
              </div>
            )}
          </nav>
        </div>
      )}
    </header>
  );
}
