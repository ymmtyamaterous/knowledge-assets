"use client";

import { use, useEffect, useState } from "react";
import Link from "next/link";
import { fetchGlossaryTerm } from "@/lib/api";
import type { GlossaryTerm } from "@/types/glossary";

type Props = {
  params: Promise<{ id: string }>;
};

export default function GlossaryTermPage({ params }: Props) {
  const { id } = use(params);
  const [term, setTerm] = useState<GlossaryTerm | null>(null);
  const [notFound, setNotFound] = useState(false);

  useEffect(() => {
    fetchGlossaryTerm(id)
      .then(setTerm)
      .catch(() => setNotFound(true));
  }, [id]);

  if (notFound) {
    return (
      <main>
        <p className="text-slate-500">用語が見つかりませんでした。</p>
        <Link href="/glossary" className="mt-4 inline-block text-sm text-pink-400 hover:underline">
          ← 用語辞典へ戻る
        </Link>
      </main>
    );
  }

  if (!term) {
    return (
      <main>
        <div className="h-8 w-48 animate-pulse rounded bg-slate-200" />
        <div className="mt-4 h-32 animate-pulse rounded bg-slate-100" />
      </main>
    );
  }

  return (
    <main className="mx-auto w-full max-w-3xl">
      <div className="mb-4">
        <Link href="/glossary" className="text-sm text-pink-400 hover:underline">
          ← 用語辞典へ戻る
        </Link>
      </div>

      <div className="rounded-2xl border border-pink-100 bg-white p-6 shadow-sm">
        <div className="flex flex-wrap items-baseline gap-3">
          <h1 className="text-2xl font-bold text-slate-800">{term.term}</h1>
          {term.reading && (
            <span className="rounded-md bg-slate-50 px-2 py-0.5 text-sm font-medium text-slate-500">
              {term.reading}
            </span>
          )}
        </div>

        {term.tags && term.tags.length > 0 && (
          <div className="mt-3 flex flex-wrap gap-2">
            {term.tags.map((tag) => (
              <Link
                key={tag.id}
                href={`/glossary?tagId=${tag.id}`}
                className="rounded-full bg-pink-50 px-2.5 py-1 text-xs font-semibold text-pink-600 hover:bg-pink-100"
              >
                #{tag.name}
              </Link>
            ))}
          </div>
        )}

        <div className="mt-5 border-t border-slate-100 pt-5">
          <p className="whitespace-pre-wrap text-base leading-relaxed text-slate-700">
            {term.definition}
          </p>
        </div>
      </div>
    </main>
  );
}
