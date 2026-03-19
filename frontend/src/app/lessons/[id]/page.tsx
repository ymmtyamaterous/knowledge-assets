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
import NoteDrawer from "@/components/NoteDrawer";

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
  const [noteDrawerOpen, setNoteDrawerOpen] = useState(false);

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
        <div className="flex items-start justify-between gap-4">
          <h1 className="text-2xl font-bold text-slate-800">{lesson.title}</h1>
          {user && (
            <button
              onClick={() => { setNoteDrawerOpen(true); setNoteSaved(false); }}
              className="flex shrink-0 items-center gap-2 rounded-lg border border-slate-200 px-4 py-2 text-sm font-semibold text-slate-600 hover:bg-slate-50 hover:border-pink-300 hover:text-pink-600 transition-colors"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.8}
                stroke="currentColor"
                className="h-4 w-4"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10"
                />
              </svg>
              メモ
            </button>
          )}
        </div>

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
      </div>

      {user && (
        <NoteDrawer
          isOpen={noteDrawerOpen}
          onClose={() => setNoteDrawerOpen(false)}
          noteContent={noteContent}
          onChange={(v) => { setNoteContent(v); setNoteSaved(false); }}
          onSave={handleSaveNote}
          saving={savingNote}
          saved={noteSaved}
        />
      )}
    </main>
  );
}
