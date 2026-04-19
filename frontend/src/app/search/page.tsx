"use client";

import { useEffect, useState, Suspense } from "react";
import { useSearchParams } from "next/navigation";
import Link from "next/link";
import { searchContent } from "@/lib/api";
import type { SearchResult } from "@/types/search";

function SearchContent() {
  const searchParams = useSearchParams();
  const q = searchParams.get("q") ?? "";
  const [result, setResult] = useState<SearchResult | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    if (q.length < 2) {
      setResult(null);
      setError("検索キーワードは2文字以上入力してください。");
      return;
    }
    setLoading(true);
    setError("");
    searchContent(q)
      .then(setResult)
      .catch(() => setError("検索中にエラーが発生しました。"))
      .finally(() => setLoading(false));
  }, [q]);

  return (
    <main className="max-w-2xl">
      <h1 className="mb-2 text-2xl font-bold text-slate-800">検索結果</h1>
      <p className="mb-6 text-sm text-slate-500">
        {q ? (
          <>「<span className="font-semibold text-pink-500">{q}</span>」の検索結果</>
        ) : (
          "キーワードを入力してください。"
        )}
      </p>

      {loading && (
        <div className="space-y-3">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="h-14 animate-pulse rounded-xl bg-slate-100" />
          ))}
        </div>
      )}

      {error && (
        <div className="rounded-xl border border-slate-200 bg-white p-6 text-center text-sm text-slate-500">
          {error}
        </div>
      )}

      {!loading && result && (
        <div className="space-y-8">
          {/* レッスン結果 */}
          <section>
            <h2 className="mb-3 text-base font-semibold text-slate-700">
              📚 レッスン（{result.lessons.length}件）
            </h2>
            {result.lessons.length === 0 ? (
              <p className="rounded-xl border border-slate-100 bg-slate-50 p-4 text-sm text-slate-400">
                レッスンは見つかりませんでした。
              </p>
            ) : (
              <div className="space-y-2">
                {result.lessons.map((lesson) => (
                  <Link
                    key={lesson.id}
                    href={`/lessons/${lesson.id}`}
                    className="flex items-center justify-between rounded-xl border border-slate-200 bg-white px-4 py-3 shadow-sm hover:border-pink-200 hover:bg-pink-50 transition"
                  >
                    <span className="text-sm font-medium text-slate-700">{lesson.title}</span>
                    <span className="text-xs text-pink-400">→</span>
                  </Link>
                ))}
              </div>
            )}
          </section>

          {/* 用語結果 */}
          <section>
            <h2 className="mb-3 text-base font-semibold text-slate-700">
              📖 用語（{result.terms.length}件）
            </h2>
            {result.terms.length === 0 ? (
              <p className="rounded-xl border border-slate-100 bg-slate-50 p-4 text-sm text-slate-400">
                用語は見つかりませんでした。
              </p>
            ) : (
              <div className="space-y-2">
                {result.terms.map((term) => (
                  <Link
                    key={term.id}
                    href={`/glossary/${term.id}`}
                    className="flex items-center justify-between rounded-xl border border-slate-200 bg-white px-4 py-3 shadow-sm hover:border-pink-200 hover:bg-pink-50 transition"
                  >
                    <div>
                      <span className="text-sm font-medium text-slate-700">{term.term}</span>
                      {term.reading && (
                        <span className="ml-2 text-xs text-slate-400">（{term.reading}）</span>
                      )}
                    </div>
                    <span className="text-xs text-pink-400">→</span>
                  </Link>
                ))}
              </div>
            )}
          </section>

          {result.lessons.length === 0 && result.terms.length === 0 && (
            <div className="rounded-xl border border-slate-200 bg-white p-8 text-center">
              <p className="text-slate-500">
                「{q}」に一致する結果が見つかりませんでした。
              </p>
              <p className="mt-2 text-sm text-slate-400">
                別のキーワードで検索してみてください。
              </p>
            </div>
          )}
        </div>
      )}
    </main>
  );
}

export default function SearchPage() {
  return (
    <Suspense fallback={<div className="max-w-2xl"><div className="h-8 animate-pulse rounded bg-slate-100" /></div>}>
      <SearchContent />
    </Suspense>
  );
}
