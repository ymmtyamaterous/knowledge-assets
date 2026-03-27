"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { fetchGlossary, fetchGlossaryTags } from "@/lib/api";
import type { GlossaryTag, GlossaryTerm } from "@/types/glossary";

const INDEXES = [
  "あ", "い", "う", "え", "お",
  "か", "き", "く", "け", "こ",
  "さ", "し", "す", "せ", "そ",
  "た", "ち", "つ", "て", "と",
  "な", "に", "ぬ", "ね", "の",
  "は", "ひ", "ふ", "へ", "ほ",
  "ま", "み", "む", "め", "も",
  "や", "ゆ", "よ",
  "ら", "り", "る", "れ", "ろ",
  "わ", "を", "ん",
  "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
];

export default function GlossaryPage() {
  const [terms, setTerms] = useState<GlossaryTerm[]>([]);
  const [tags, setTags] = useState<GlossaryTag[]>([]);
  const [loadingTerms, setLoadingTerms] = useState(true);
  const [selectedTagId, setSelectedTagId] = useState<string>("");
  const [selectedIndex, setSelectedIndex] = useState<string>("");
  const [keyword, setKeyword] = useState("");

  useEffect(() => {
    setLoadingTerms(true);
    fetchGlossary(selectedTagId || undefined)
      .then((res) => setTerms(res ?? []))
      .catch(() => setTerms([]))
      .finally(() => setLoadingTerms(false));
  }, [selectedTagId]);

  useEffect(() => {
    fetchGlossaryTags()
      .then((res) => setTags(res ?? []))
      .catch(() => setTags([]));
  }, []);

  const filteredTerms = useMemo(() => {
    const k = keyword.trim().toLowerCase();
    return terms.filter((term) => {
      if (selectedIndex) {
        const first = (term.reading || term.term).trim().charAt(0).toLowerCase();
        if (first !== selectedIndex) {
          return false;
        }
      }
      if (!k) return true;
      return [term.term, term.reading, term.definition]
        .some((value) => value?.toLowerCase().includes(k));
    });
  }, [terms, selectedIndex, keyword]);

  return (
    <main>
      <h1 className="mb-6 text-2xl font-bold text-slate-800">用語辞典</h1>

      <div className="mb-4 rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
        <input
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          placeholder="用語・読み・定義で検索"
          className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm outline-none focus:border-pink-300"
        />

        <div className="mt-3 flex flex-wrap gap-2">
          <button
            onClick={() => setSelectedTagId("")}
            className={`rounded-full px-3 py-1 text-xs ${selectedTagId === "" ? "bg-pink-500 text-white" : "bg-slate-100 text-slate-600"}`}
          >
            すべて
          </button>
          {tags.map((tag) => (
            <button
              key={tag.id}
              onClick={() => setSelectedTagId(tag.id)}
              className={`rounded-full px-3 py-1 text-xs ${selectedTagId === tag.id ? "bg-pink-500 text-white" : "bg-slate-100 text-slate-600 hover:bg-pink-50"}`}
            >
              {tag.name}
            </button>
          ))}
        </div>
      </div>

      <div className="mb-5 grid grid-cols-6 gap-2 sm:grid-cols-8 md:grid-cols-12">
        <button
          onClick={() => setSelectedIndex("")}
          className={`rounded px-2 py-1 text-xs ${selectedIndex === "" ? "bg-pink-500 text-white" : "bg-slate-100 text-slate-600"}`}
        >
          全
        </button>
        {INDEXES.map((idx) => (
          <button
            key={idx}
            onClick={() => setSelectedIndex(idx)}
            className={`rounded px-2 py-1 text-xs ${selectedIndex === idx ? "bg-pink-500 text-white" : "bg-slate-100 text-slate-600 hover:bg-pink-50"}`}
          >
            {idx.toUpperCase()}
          </button>
        ))}
      </div>

      {loadingTerms ? (
        <div className="space-y-3">
          <div className="h-24 animate-pulse rounded-xl bg-slate-100" />
          <div className="h-24 animate-pulse rounded-xl bg-slate-100" />
          <div className="h-24 animate-pulse rounded-xl bg-slate-100" />
        </div>
      ) : filteredTerms.length === 0 ? (
        <p className="text-sm text-slate-500">条件に一致する用語がありません。</p>
      ) : (
        <div className="space-y-3">
          {filteredTerms.map((term) => {
            const termTags = term.tags ?? [];
            return (
              <div
                key={term.id}
                className="rounded-xl border border-slate-200 bg-white shadow-sm transition hover:border-pink-300 hover:shadow-md"
              >
                <Link href={`/glossary/${term.id}`} className="block p-4">
                  <div className="flex flex-wrap items-baseline gap-2">
                    <h2 className="text-lg font-semibold text-slate-800">{term.term}</h2>
                    {term.reading && (
                      <span className="rounded-md bg-slate-50 px-2 py-0.5 text-sm font-medium text-slate-500">{term.reading}</span>
                    )}
                  </div>
                  <p className="mt-2 text-sm text-slate-600">{term.definition}</p>
                </Link>
                {termTags.length > 0 && (
                  <div className="flex flex-wrap gap-2 border-t border-slate-100 px-4 pb-3 pt-2">
                    {termTags.map((tag) => (
                      <button
                        key={`${term.id}-${tag.id}`}
                        onClick={() => setSelectedTagId(tag.id)}
                        className="rounded-full bg-pink-50 px-2.5 py-1 text-xs font-semibold text-pink-600 hover:bg-pink-100"
                      >
                        #{tag.name}
                      </button>
                    ))}
                  </div>
                )}
              </div>
            );
          })}
        </div>
      )}
    </main>
  );
}
