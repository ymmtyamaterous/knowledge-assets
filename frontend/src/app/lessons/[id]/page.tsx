"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { completeLesson, fetchLesson, fetchLessonQuiz, fetchMyProgress } from "@/lib/api";
import type { Lesson } from "@/types/lesson";
import { useAuth } from "@/features/auth/AuthContext";
import { use } from "react";
import type { Quiz } from "@/types/quiz";
import { convertMarkdownToHtml } from "@/lib/markdown";

type Props = {
  params: Promise<{ id: string }>;
};

function MarkdownContent({ content }: { content: string }) {
  const html = convertMarkdownToHtml(content);

  return (
    <div
      className="prose prose-slate max-w-none"
      dangerouslySetInnerHTML={{ __html: html }}
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
  const [quiz, setQuiz] = useState<Quiz | null>(null);

  useEffect(() => {
    fetchLesson(id)
      .then(setLesson)
      .catch(() => setNotFound(true));
  }, [id]);

  useEffect(() => {
    fetchLessonQuiz(id)
      .then(setQuiz)
      .catch(() => setQuiz(null));
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

          {quiz && (
            <Link
              href={`/quiz/${quiz.id}`}
              className="rounded-lg border border-pink-200 px-4 py-2 text-sm font-semibold text-pink-600 hover:bg-pink-50"
            >
              このレッスンのクイズへ
            </Link>
          )}
        </div>
      </div>
    </main>
  );
}
