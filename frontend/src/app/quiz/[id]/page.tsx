"use client";

import { use, useEffect, useMemo, useState } from "react";
import { fetchQuiz, submitQuiz, completeLesson, fetchMyProgress } from "@/lib/api";
import type { QuizDetail, SubmitQuizAnswer } from "@/types/quiz";
import Link from "next/link";
import { useAuth } from "@/features/auth/AuthContext";

type Props = {
  params: Promise<{ id: string }>;
};

export default function QuizPage({ params }: Props) {
  const { id } = use(params);
  const { user } = useAuth();

  const [detail, setDetail] = useState<QuizDetail | null>(null);
  const [currentIndex, setCurrentIndex] = useState(0);
  const effectiveLessonId = detail?.quiz.lessonId ?? "";
  const [answers, setAnswers] = useState<Record<string, string>>({});
  const [checked, setChecked] = useState<Record<string, boolean>>({});
  const [submitted, setSubmitted] = useState(false);
  const [finalScore, setFinalScore] = useState<{ score: number; total: number } | null>(null);
  const [lessonCompleted, setLessonCompleted] = useState(false);
  const [completingLesson, setCompletingLesson] = useState(false);
  const answeredCount = detail ? Object.keys(answers).filter((qId) => !!answers[qId]).length : 0;

  useEffect(() => {
    fetchQuiz(id)
      .then((res) => setDetail(res))
      .catch(() => setDetail(null));
  }, [id]);

  useEffect(() => {
    if (!user || !effectiveLessonId) return;
    fetchMyProgress()
      .then((list) => {
        setLessonCompleted(list.some((p) => p.lessonId === effectiveLessonId));
      })
      .catch(() => {});
  }, [user, effectiveLessonId]);

  const handleCompleteLesson = async () => {
    if (!effectiveLessonId) return;
    setCompletingLesson(true);
    try {
      await completeLesson(effectiveLessonId);
      setLessonCompleted(true);
    } finally {
      setCompletingLesson(false);
    }
  };

  const currentQuestion = detail?.questions[currentIndex] ?? null;

  const questionResult = useMemo(() => {
    if (!currentQuestion) return null;
    const selectedChoiceId = answers[currentQuestion.id];
    const selectedChoice = currentQuestion.choices.find((c) => c.id === selectedChoiceId);
    if (!selectedChoice) return null;
    return selectedChoice.isCorrect;
  }, [answers, currentQuestion]);

  const handleSelectChoice = (questionId: string, choiceId: string) => {
    setAnswers((prev) => ({ ...prev, [questionId]: choiceId }));
    setChecked((prev) => ({ ...prev, [questionId]: true }));
  };

  const handleFinish = async () => {
    if (!detail || !user) return;

    const payload: SubmitQuizAnswer[] = detail.questions
      .map((question) => ({ questionId: question.id, choiceId: answers[question.id] }))
      .filter((item) => item.choiceId);

    const res = await submitQuiz(detail.quiz.id, payload);
    setSubmitted(true);
    setFinalScore({ score: res.result.score, total: res.result.total });
  };

  if (!detail) {
    return (
      <main className="mx-auto w-full max-w-3xl">
        <p className="text-slate-500">クイズを読み込めませんでした。</p>
      </main>
    );
  }

  if (!user) {
    return (
      <main className="mx-auto w-full max-w-3xl">
        <h1 className="mb-4 text-2xl font-bold text-slate-800">クイズ</h1>
        <p className="text-slate-600">クイズ回答を記録するにはログインが必要です。</p>
        <Link href="/auth/login" className="mt-4 inline-block rounded-md bg-pink-500 px-4 py-2 text-sm font-semibold text-white hover:bg-pink-600">
          ログインへ
        </Link>
      </main>
    );
  }

  if (submitted && finalScore) {
    return (
      <main className="mx-auto w-full max-w-3xl">
        <div className="rounded-xl border border-pink-200 bg-white p-6 shadow-sm">
          <h1 className="text-2xl font-bold text-slate-800">結果</h1>
          <p className="mt-3 text-lg font-semibold text-pink-600">
            {finalScore.score} / {finalScore.total} 問 正解
          </p>
          <p className="mt-1 text-sm text-slate-500">
            正答率: {Math.round((finalScore.score / Math.max(finalScore.total, 1)) * 100)}%
          </p>
          {effectiveLessonId && (
            <div className="mt-6">
              {lessonCompleted ? (
                <span className="inline-flex items-center gap-2 rounded-full bg-green-100 px-4 py-2 text-sm font-semibold text-green-700">
                  ✓ レッスン完了済み
                </span>
              ) : (
                <button
                  onClick={handleCompleteLesson}
                  disabled={completingLesson}
                  className="rounded-lg bg-green-500 px-5 py-2 text-sm font-semibold text-white hover:bg-green-600 disabled:opacity-50"
                >
                  {completingLesson ? "登録中..." : "レッスンを完了にする"}
                </button>
              )}
            </div>
          )}
          <Link href="/courses" className="mt-4 inline-block rounded-md bg-pink-500 px-4 py-2 text-sm font-semibold text-white hover:bg-pink-600">
            コース一覧へ戻る
          </Link>
        </div>
      </main>
    );
  }

  if (!currentQuestion) {
    return (
      <main className="mx-auto w-full max-w-3xl">
        <p className="text-slate-500">問題がありません。</p>
      </main>
    );
  }

  return (
    <main className="mx-auto w-full max-w-3xl">
      <div className="mb-3 flex items-center justify-between text-sm text-slate-500">
        <span>問題 {currentIndex + 1} / {detail.questions.length}</span>
        <span>回答済み {answeredCount} / {detail.questions.length}</span>
      </div>
      <div className="mb-4 h-2 w-full rounded-full bg-slate-100">
        <div
          className="h-2 rounded-full bg-pink-500 transition-all"
          style={{ width: `${((currentIndex + 1) / Math.max(detail.questions.length, 1)) * 100}%` }}
        />
      </div>

      <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm sm:p-6">
        <h1 className="text-xl font-bold text-slate-800">{currentQuestion.questionText}</h1>

        <div className="mt-4 space-y-2">
          {currentQuestion.choices.map((choice) => {
            const selected = answers[currentQuestion.id] === choice.id;
            return (
              <button
                key={choice.id}
                onClick={() => handleSelectChoice(currentQuestion.id, choice.id)}
                className={`block w-full rounded-lg border px-4 py-3 text-left text-sm ${selected ? "border-pink-400 bg-pink-50 text-pink-700" : "border-slate-200 hover:bg-slate-50"}`}
              >
                {choice.choiceText}
              </button>
            );
          })}
        </div>

        {checked[currentQuestion.id] && (
          <div className={`mt-4 rounded-lg px-4 py-3 text-sm ${questionResult ? "bg-green-50 text-green-700" : "bg-rose-50 text-rose-700"}`}>
            {questionResult ? "正解です。" : "不正解です。"}
            <p className="mt-1">{currentQuestion.explanation}</p>
          </div>
        )}

        <div className="mt-6 flex items-center justify-between">
          <button
            onClick={() => setCurrentIndex((prev) => Math.max(prev - 1, 0))}
            disabled={currentIndex === 0}
            className="rounded-md border border-slate-300 px-3 py-2 text-sm text-slate-600 disabled:opacity-40"
          >
            前へ
          </button>

          {currentIndex < detail.questions.length - 1 ? (
            <button
              onClick={() => setCurrentIndex((prev) => Math.min(prev + 1, detail.questions.length - 1))}
              disabled={!answers[currentQuestion.id]}
              className="rounded-md bg-pink-500 px-4 py-2 text-sm font-semibold text-white hover:bg-pink-600 disabled:cursor-not-allowed disabled:opacity-50"
            >
              次へ
            </button>
          ) : (
            <button
              onClick={handleFinish}
              disabled={answeredCount < detail.questions.length}
              className="rounded-md bg-pink-500 px-4 py-2 text-sm font-semibold text-white hover:bg-pink-600 disabled:cursor-not-allowed disabled:opacity-50"
            >
              解答を送信
            </button>
          )}
        </div>
      </div>
    </main>
  );
}
