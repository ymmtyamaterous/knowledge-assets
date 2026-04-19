"use client";

import { use, useEffect, useMemo, useState, useRef } from "react";
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

  // 模擬試験専用ステート
  const [mockStarted, setMockStarted] = useState(false);
  const [timeLeft, setTimeLeft] = useState(0);
  const [elapsedSeconds, setElapsedSeconds] = useState(0);
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

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

  // 模擬試験タイマー
  useEffect(() => {
    if (!mockStarted || !detail?.quiz.isMockExam) return;
    const totalSeconds = (detail.quiz.timeLimitMinutes ?? 60) * 60;
    setTimeLeft(totalSeconds);
    setElapsedSeconds(0);
    const startedAt = Date.now();
    timerRef.current = setInterval(() => {
      const elapsed = Math.floor((Date.now() - startedAt) / 1000);
      const remaining = Math.max(totalSeconds - elapsed, 0);
      setTimeLeft(remaining);
      setElapsedSeconds(elapsed);
      if (remaining === 0) {
        if (timerRef.current) clearInterval(timerRef.current);
        handleFinish();
      }
    }, 1000);
    return () => {
      if (timerRef.current) clearInterval(timerRef.current);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [mockStarted]);

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
    if (timerRef.current) clearInterval(timerRef.current);

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
    const isMock = detail?.quiz.isMockExam ?? false;
    const elapsedMin = String(Math.floor(elapsedSeconds / 60)).padStart(2, "0");
    const elapsedSec = String(elapsedSeconds % 60).padStart(2, "0");
    return (
      <main className="mx-auto w-full max-w-3xl">
        <div className="rounded-xl border border-pink-200 bg-white p-6 shadow-sm">
          <h1 className="text-2xl font-bold text-slate-800">
            {isMock ? "模擬試験 結果" : "結果"}
          </h1>
          <p className="mt-3 text-lg font-semibold text-pink-600">
            {finalScore.score} / {finalScore.total} 問 正解
          </p>
          <p className="mt-1 text-sm text-slate-500">
            正答率: {Math.round((finalScore.score / Math.max(finalScore.total, 1)) * 100)}%
          </p>
          {isMock && (
            <p className="mt-1 text-sm text-slate-500">
              解答時間: {elapsedMin}:{elapsedSec}
            </p>
          )}
          {/* 模擬試験: 問題ごとの正誤一覧 */}
          {isMock && detail && (
            <div className="mt-6 overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-slate-200 text-left text-xs text-slate-500">
                    <th className="pb-2 pr-3">No.</th>
                    <th className="pb-2 pr-3">問題</th>
                    <th className="pb-2 pr-3">あなたの回答</th>
                    <th className="pb-2 pr-3">正誤</th>
                    <th className="pb-2">解説</th>
                  </tr>
                </thead>
                <tbody>
                  {detail.questions.map((q, i) => {
                    const selectedId = answers[q.id];
                    const selectedChoice = q.choices.find((c) => c.id === selectedId);
                    const isCorrect = selectedChoice?.isCorrect ?? false;
                    const correctChoice = q.choices.find((c) => c.isCorrect);
                    return (
                      <tr key={q.id} className="border-b border-slate-100">
                        <td className="py-2 pr-3 text-slate-400">{i + 1}</td>
                        <td className="py-2 pr-3 text-slate-700">{q.questionText}</td>
                        <td className="py-2 pr-3 text-slate-600">
                          {selectedChoice?.choiceText ?? <span className="text-slate-400">未回答</span>}
                        </td>
                        <td className="py-2 pr-3">
                          {isCorrect ? (
                            <span className="font-bold text-green-600">○</span>
                          ) : (
                            <div>
                              <span className="font-bold text-red-500">✕</span>
                              <p className="text-xs text-slate-500">正答: {correctChoice?.choiceText}</p>
                            </div>
                          )}
                        </td>
                        <td className="py-2 text-xs text-slate-500">{q.explanation}</td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          )}
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

  // 模擬試験: 開始確認画面
  if (detail.quiz.isMockExam && !mockStarted) {
    return (
      <main className="mx-auto w-full max-w-3xl">
        <div className="rounded-xl border border-amber-200 bg-amber-50 p-6 shadow-sm">
          <h1 className="text-2xl font-bold text-slate-800">模擬試験</h1>
          <div className="mt-4 space-y-2 text-sm text-slate-700">
            <p>📋 問題数: <span className="font-semibold">{detail.questions.length} 問</span></p>
            <p>⏱ 制限時間: <span className="font-semibold">{detail.quiz.timeLimitMinutes} 分</span></p>
          </div>
          <div className="mt-4 rounded-lg border border-amber-300 bg-white p-4 text-sm text-amber-800">
            <p className="font-semibold">注意事項</p>
            <ul className="mt-2 list-inside list-disc space-y-1">
              <li>試験開始後は全問が一覧表示されます。</li>
              <li>制限時間が終了すると自動で送信されます。</li>
              <li>途中でページを離れると回答が失われます。</li>
            </ul>
          </div>
          <button
            onClick={() => setMockStarted(true)}
            className="mt-6 w-full rounded-lg bg-amber-500 py-3 text-sm font-semibold text-white hover:bg-amber-600 transition"
          >
            試験を開始する
          </button>
        </div>
      </main>
    );
  }

  // 模擬試験: 全問表示モード
  if (detail.quiz.isMockExam && mockStarted) {
    const mm = String(Math.floor(timeLeft / 60)).padStart(2, "0");
    const ss = String(timeLeft % 60).padStart(2, "0");
    const timeWarning = timeLeft <= 300; // 残り5分以下
    return (
      <main className="mx-auto w-full max-w-3xl pb-16">
        {/* 固定タイマー */}
        <div className={`fixed right-4 top-16 z-50 rounded-xl border px-4 py-2 text-sm font-bold shadow-md transition-colors ${timeWarning ? "border-red-300 bg-red-50 text-red-700" : "border-amber-300 bg-amber-50 text-amber-700"}`}>
          ⏱ {mm}:{ss}
        </div>
        <h1 className="mb-6 text-2xl font-bold text-slate-800">模擬試験</h1>
        <div className="space-y-8">
          {detail.questions.map((q, i) => {
            const selectedChoiceId = answers[q.id];
            return (
              <div key={q.id} className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
                <p className="text-sm text-slate-400 mb-1">問題 {i + 1}</p>
                <p className="font-semibold text-slate-800">{q.questionText}</p>
                <div className="mt-3 space-y-2">
                  {q.choices.map((c) => {
                    const selected = selectedChoiceId === c.id;
                    return (
                      <button
                        key={c.id}
                        onClick={() => setAnswers((prev) => ({ ...prev, [q.id]: c.id }))}
                        className={`block w-full rounded-lg border px-4 py-2.5 text-left text-sm transition ${selected ? "border-amber-400 bg-amber-50 text-amber-800 font-medium" : "border-slate-200 hover:bg-slate-50"}`}
                      >
                        {c.choiceText}
                      </button>
                    );
                  })}
                </div>
              </div>
            );
          })}
        </div>
        <div className="mt-8">
          <p className="mb-3 text-center text-sm text-slate-500">
            回答済み {answeredCount} / {detail.questions.length} 問
          </p>
          <button
            onClick={handleFinish}
            disabled={answeredCount < detail.questions.length}
            className="w-full rounded-lg bg-amber-500 py-3 text-sm font-semibold text-white hover:bg-amber-600 disabled:cursor-not-allowed disabled:opacity-50 transition"
          >
            解答を送信する
          </button>
        </div>
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
