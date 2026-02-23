"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { completeLesson, fetchLesson, fetchMyProgress } from "@/lib/api";
import type { Lesson } from "@/types/lesson";
import { useAuth } from "@/features/auth/AuthContext";
import { use } from "react";

type Props = {
  params: Promise<{ id: string }>;
};

function MarkdownContent({ content }: { content: string }) {
  // 簡易Markdown→HTML変換（実用ではreact-markdownを使う想定）
  const html = content
    .replace(/^### (.+)$/gm, "<h3 class='text-lg font-semibold mt-4 mb-1'>$1</h3>")
    .replace(/^## (.+)$/gm, "<h2 class='text-xl font-bold mt-6 mb-2 text-pink-500'>$1</h2>")
    .replace(/^# (.+)$/gm, "<h1 class='text-2xl font-bold mt-6 mb-2'>$1</h1>")
    .replace(/\*\*(.+?)\*\*/g, "<strong>$1</strong>")
    .replace(/^- (.+)$/gm, "<li class='ml-4 list-disc'>$1</li>")
    .replace(/```[\s\S]*?```/g, (m) => {
      const code = m.replace(/^```\w*\n?/, "").replace(/```$/, "");
      return `<pre class='bg-slate-100 rounded p-3 text-sm overflow-x-auto my-3'><code>${code}</code></pre>`;
    })
    .replace(/\n\n/g, "</p><p class='mt-3'>")
    .replace(/^\|(.+)\|$/gm, (row) => {
      const cells = row
        .split("|")
        .filter((c) => c.trim() !== "")
        .map((c) => `<td class='border border-slate-300 px-3 py-1 text-sm'>${c.trim()}</td>`)
        .join("");
      return `<tr>${cells}</tr>`;
    });

  return (
    <div
      className="prose prose-slate max-w-none"
      dangerouslySetInnerHTML={{ __html: `<p class='mt-3'>${html}</p>` }}
    />
  );
}

export default function LessonPage({ params }: Props) {
  const { id } = use(params);
  const { user } = useAuth();
  const router = useRouter();
  const [lesson, setLesson] = useState<Lesson | null>(null);
  const [completed, setCompleted] = useState(false);
  const [completing, setCompleting] = useState(false);
  const [notFound, setNotFound] = useState(false);

  useEffect(() => {
    fetchLesson(id)
      .then(setLesson)
      .catch(() => setNotFound(true));
  }, [id]);

  useEffect(() => {
    if (!user) return;
    fetchMyProgress()
      .then((list) => {
        setCompleted(list.some((p) => p.lessonId === id));
      })
      .catch(() => {});
  }, [user, id]);

  const handleComplete = useCallback(async () => {
    if (!user) {
      router.push("/auth/login");
      return;
    }
    setCompleting(true);
    try {
      await completeLesson(id);
      setCompleted(true);
    } finally {
      setCompleting(false);
    }
  }, [user, id, router]);

  if (notFound) {
    return (
      <main>
        <p className="text-slate-500">レッスンが見つかりませんでした。</p>
        <Link href="/" className="mt-4 inline-block text-sm text-pink-400 hover:underline">
          ← トップに戻る
        </Link>
      </main>
    );
  }

  if (!lesson) {
    return (
      <main>
        <div className="h-8 w-48 animate-pulse rounded bg-slate-200" />
        <div className="mt-4 h-64 animate-pulse rounded bg-slate-100" />
      </main>
    );
  }

  return (
    <main>
      <div className="mb-2">
        <button
          onClick={() => router.back()}
          className="text-sm text-pink-400 hover:underline"
        >
          ← 戻る
        </button>
      </div>

      <div className="rounded-2xl border border-pink-100 bg-white p-6 shadow-sm">
        <h1 className="text-2xl font-bold text-slate-800">{lesson.title}</h1>

        <div className="mt-4 border-t border-slate-100 pt-4">
          <MarkdownContent content={lesson.content} />
        </div>

        <div className="mt-8 flex items-center gap-4">
          {completed ? (
            <span className="flex items-center gap-2 rounded-full bg-green-100 px-4 py-2 text-sm font-semibold text-green-700">
              ✓ 完了済み
            </span>
          ) : (
            <button
              onClick={handleComplete}
              disabled={completing}
              className="rounded-lg bg-pink-500 px-6 py-2 text-sm font-semibold text-white hover:bg-pink-600 disabled:opacity-50"
            >
              {completing ? "登録中..." : "レッスンを完了する"}
            </button>
          )}
        </div>
      </div>
    </main>
  );
}
