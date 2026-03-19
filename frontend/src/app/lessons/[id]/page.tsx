"use client";

import { useCallback, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { completeLesson, fetchLesson, fetchLessonQuiz, fetchLessonNote, saveLessonNote, fetchMyProgress, uncompleteLesson } from "@/lib/api";
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
  const [savingProgress, setSavingProgress] = useState(false);
  const [notFound, setNotFound] = useState(false);
  const [quiz, setQuiz] = useState<Quiz | null>(null);
  const [noteContent, setNoteContent] = useState("");
  const [savingNote, setSavingNote] = useState(false);
  const [noteSaved, setNoteSaved] = useState(false);

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

  useEffect(() => {
    if (!user) return;
    fetchLessonNote(id)
      .then((note) => {
        if (note) setNoteContent(note.content);
      })
      .catch(() => {});
  }, [user, id]);

  const handleSaveNote = async () => {
    if (!user) return;
    setSavingNote(true);
    setNoteSaved(false);
    try {
      await saveLessonNote(id, noteContent);
      setNoteSaved(true);
    } finally {
      setSavingNote(false);
    }
  };

  const handleComplete = useCallback(async () => {
    if (!user) {
      router.push("/auth/login");
      return;
    }
    setSavingProgress(true);
    try {
      await completeLesson(id);
      setCompleted(true);
    } finally {
      setSavingProgress(false);
    }
  }, [user, id, router]);

  const handleUncomplete = useCallback(async () => {
    if (!user) {
      router.push("/auth/login");
      return;
    }
    setSavingProgress(true);
    try {
      await uncompleteLesson(id);
      setCompleted(false);
    } finally {
      setSavingProgress(false);
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
            <div className="flex flex-wrap items-center gap-3">
              <span className="flex items-center gap-2 rounded-full bg-green-100 px-4 py-2 text-sm font-semibold text-green-700">
                ✓ 完了済み
              </span>
              <button
                onClick={handleUncomplete}
                disabled={savingProgress}
                className="rounded-lg border border-slate-300 px-4 py-2 text-sm font-semibold text-slate-600 hover:bg-slate-50 disabled:opacity-50"
              >
                {savingProgress ? "更新中..." : "完了を取り消す"}
              </button>
            </div>
          ) : (
            <button
              onClick={handleComplete}
              disabled={savingProgress}
              className="rounded-lg bg-pink-500 px-6 py-2 text-sm font-semibold text-white hover:bg-pink-600 disabled:opacity-50"
            >
              {savingProgress ? "登録中..." : "レッスンを完了する"}
            </button>
          )}

          {quiz && (
            <Link
              href={`/quiz/${quiz.id}`}
              className="rounded-lg border border-pink-200 px-4 py-2 text-sm font-semibold text-pink-600 hover:bg-pink-50"
            >
              このレッスンの確認クイズへ
            </Link>
          )}
        </div>

        {user && (
          <div className="mt-8 border-t border-slate-100 pt-6">
            <h2 className="mb-2 text-sm font-semibold text-slate-700">📝 メモ</h2>
            <textarea
              value={noteContent}
              onChange={(e) => { setNoteContent(e.target.value); setNoteSaved(false); }}
              placeholder="このレッスンのメモを入力..."
              rows={4}
              className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm text-slate-700 focus:border-pink-400 focus:outline-none focus:ring-1 focus:ring-pink-400"
            />
            <div className="mt-2 flex items-center gap-3">
              <button
                onClick={handleSaveNote}
                disabled={savingNote}
                className="rounded-lg bg-pink-500 px-4 py-2 text-sm font-semibold text-white hover:bg-pink-600 disabled:opacity-50"
              >
                {savingNote ? "保存中..." : "メモを保存"}
              </button>
              {noteSaved && <span className="text-sm text-green-600">✓ 保存しました</span>}
            </div>
          </div>
        )}
      </div>
    </main>
  );
}
