"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/features/auth/AuthContext";
import { updateMe, fetchMyNotes, fetchMyQuizResults, changePassword, fetchMyStreak, fetchMyStats, fetchMyBadges } from "@/lib/api";
import type { UserNote } from "@/types/note";
import type { UserQuizResult } from "@/types/quiz";
import type { UserStreak, UserStats } from "@/types/progress";
import type { UserBadge } from "@/types/badge";
import Link from "next/link";

export default function ProfilePage() {
  const { user, setUser, loading: authLoading } = useAuth();
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState("");
  const [notes, setNotes] = useState<UserNote[]>([]);
  const [quizResults, setQuizResults] = useState<UserQuizResult[]>([]);
  const [streak, setStreak] = useState<UserStreak>({ currentStreak: 0, longestStreak: 0, lastStudiedAt: "" });
  const [stats, setStats] = useState<UserStats | null>(null);
  const [badges, setBadges] = useState<UserBadge[]>([]);
  const [pwCurrent, setPwCurrent] = useState("");
  const [pwNew, setPwNew] = useState("");
  const [pwConfirm, setPwConfirm] = useState("");
  const [pwLoading, setPwLoading] = useState(false);
  const [pwSuccess, setPwSuccess] = useState(false);
  const [pwError, setPwError] = useState("");

  useEffect(() => {
    if (authLoading) return;
    if (user === null) {
      router.replace("/auth/login");
      return;
    }
    setUsername(user.username ?? "");
  }, [user, authLoading, router]);

  useEffect(() => {
    if (!user) return;
    fetchMyNotes().then(setNotes).catch(() => {});
    fetchMyQuizResults().then(setQuizResults).catch(() => {});
    fetchMyStreak().then(setStreak).catch(() => {});
    fetchMyStats().then(setStats).catch(() => {});
    fetchMyBadges().then((res) => setBadges(res.badges)).catch(() => {});
  }, [user]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setSuccess(false);
    if (!username.trim()) {
      setError("ユーザー名を入力してください。");
      return;
    }
    setLoading(true);
    try {
      const updated = await updateMe({ username });
      setUser(updated);
      setSuccess(true);
    } catch {
      setError("更新に失敗しました。もう一度お試しください。");
    } finally {
      setLoading(false);
    }
  };

  const handleChangePassword = async (e: React.FormEvent) => {
    e.preventDefault();
    setPwError("");
    setPwSuccess(false);
    if (!pwCurrent || !pwNew || !pwConfirm) {
      setPwError("すべての項目を入力してください。");
      return;
    }
    if (pwNew !== pwConfirm) {
      setPwError("新しいパスワードが一致しません。");
      return;
    }
    if (pwNew.length < 8) {
      setPwError("新しいパスワードは8文字以上にしてください。");
      return;
    }
    setPwLoading(true);
    try {
      await changePassword(pwCurrent, pwNew);
      setPwSuccess(true);
      setPwCurrent("");
      setPwNew("");
      setPwConfirm("");
    } catch {
      setPwError("現在のパスワードが正しくありません。");
    } finally {
      setPwLoading(false);
    }
  };

  const handleExportNotes = () => {
    const md = notes
      .map((note) => {
        const title = note.lessonTitle || `レッスン ${note.lessonId}`;
        return `# ${title}\n\n${note.content}\n\n---\n`;
      })
      .join("\n");
    const blob = new Blob([md], { type: "text/markdown" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "my-notes.md";
    a.click();
    URL.revokeObjectURL(url);
  };

  const getScoreBadgeClass = (score: number, total: number): string => {
    if (total === 0) return "bg-slate-100 text-slate-600";
    const rate = (score / total) * 100;
    if (rate >= 80) return "bg-pink-100 text-pink-700";
    if (rate >= 60) return "bg-yellow-100 text-yellow-700";
    return "bg-slate-100 text-slate-600";
  };

  if (authLoading) {
    return (
      <main>
        <div className="h-8 w-40 animate-pulse rounded bg-slate-200" />
      </main>
    );
  }

  return (
    <main className="max-w-md">
      <h1 className="mb-6 text-2xl font-bold text-slate-800">プロフィール</h1>

      {/* ストリークバナー */}
      {streak.currentStreak > 0 && (
        <div className="mb-6 rounded-xl border border-pink-200 bg-gradient-to-r from-pink-50 to-rose-50 p-5 shadow-sm">
          <div className="flex items-center gap-2">
            <span className="text-2xl">🔥</span>
            <div>
              <p className="text-lg font-bold text-pink-600">{streak.currentStreak}日連続学習中！</p>
              <p className="text-xs text-slate-500">最長記録: {streak.longestStreak}日</p>
            </div>
          </div>
        </div>
      )}
      {streak.currentStreak === 0 && streak.longestStreak > 0 && (
        <div className="mb-6 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          <div className="flex items-center gap-2">
            <span className="text-2xl">📚</span>
            <div>
              <p className="text-sm font-medium text-slate-700">今日も学習してストリークをつなごう！</p>
              <p className="text-xs text-slate-500">最長記録: {streak.longestStreak}日</p>
            </div>
          </div>
        </div>
      )}

      <div className="mb-6 rounded-xl border border-slate-200 bg-white p-6 shadow-sm">
        <p className="text-sm text-slate-500">メールアドレス</p>
        <p className="mt-1 font-medium text-slate-800">{user?.email}</p>
      </div>

      {/* 学習統計 */}
      <div className="mb-6 rounded-xl border border-slate-200 bg-white p-6 shadow-sm">
        <h2 className="mb-4 text-lg font-semibold text-slate-800">📊 学習統計</h2>
        {stats ? (
          <div className="grid grid-cols-2 gap-3">
            <div className="rounded-lg bg-pink-50 p-4 text-center">
              <p className="text-2xl font-bold text-pink-600">{stats.totalCompletedLessons}</p>
              <p className="mt-1 text-xs text-slate-500">📚 完了レッスン</p>
            </div>
            <div className="rounded-lg bg-blue-50 p-4 text-center">
              <p className="text-2xl font-bold text-blue-600">{stats.totalStudyDays}</p>
              <p className="mt-1 text-xs text-slate-500">📅 学習日数</p>
            </div>
            <div className="rounded-lg bg-green-50 p-4 text-center">
              <p className="text-2xl font-bold text-green-600">{stats.totalNotes}</p>
              <p className="mt-1 text-xs text-slate-500">📝 メモ件数</p>
            </div>
            <div className="rounded-lg bg-amber-50 p-4 text-center">
              <p className="text-2xl font-bold text-amber-600">{stats.averageQuizScore.toFixed(1)}%</p>
              <p className="mt-1 text-xs text-slate-500">🎯 クイズ正答率</p>
            </div>
          </div>
        ) : (
          <div className="grid grid-cols-2 gap-3">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="h-20 animate-pulse rounded-lg bg-slate-100" />
            ))}
          </div>
        )}
      </div>

      <form onSubmit={handleSubmit} className="rounded-xl border border-slate-200 bg-white p-6 shadow-sm space-y-4">
        <h2 className="text-lg font-semibold text-slate-800">情報の編集</h2>
        <div>
          <label htmlFor="username" className="block text-sm font-medium text-slate-700">
            ユーザー名
          </label>
          <input
            id="username"
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-pink-400 focus:outline-none focus:ring-1 focus:ring-pink-400"
          />
        </div>
        {error && <p className="text-sm text-red-500">{error}</p>}
        {success && <p className="text-sm text-green-500">更新しました！</p>}
        <button
          type="submit"
          disabled={loading}
          className="w-full rounded-lg bg-pink-400 py-2 text-sm font-semibold text-white hover:bg-pink-500 disabled:opacity-50 transition"
        >
          {loading ? "更新中..." : "変更を保存"}
        </button>
      </form>

      {/* パスワード変更 */}
      <form onSubmit={handleChangePassword} className="mt-6 rounded-xl border border-slate-200 bg-white p-6 shadow-sm space-y-4">
        <h2 className="text-lg font-semibold text-slate-800">パスワード変更</h2>
        <div>
          <label htmlFor="pw-current" className="block text-sm font-medium text-slate-700">
            現在のパスワード
          </label>
          <input
            id="pw-current"
            type="password"
            value={pwCurrent}
            onChange={(e) => setPwCurrent(e.target.value)}
            className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-pink-400 focus:outline-none focus:ring-1 focus:ring-pink-400"
          />
        </div>
        <div>
          <label htmlFor="pw-new" className="block text-sm font-medium text-slate-700">
            新しいパスワード（8文字以上）
          </label>
          <input
            id="pw-new"
            type="password"
            value={pwNew}
            onChange={(e) => setPwNew(e.target.value)}
            className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-pink-400 focus:outline-none focus:ring-1 focus:ring-pink-400"
          />
        </div>
        <div>
          <label htmlFor="pw-confirm" className="block text-sm font-medium text-slate-700">
            新しいパスワード（確認）
          </label>
          <input
            id="pw-confirm"
            type="password"
            value={pwConfirm}
            onChange={(e) => setPwConfirm(e.target.value)}
            className="mt-1 w-full rounded-lg border border-slate-300 px-3 py-2 text-sm focus:border-pink-400 focus:outline-none focus:ring-1 focus:ring-pink-400"
          />
        </div>
        {pwError && <p className="text-sm text-red-500">{pwError}</p>}
        {pwSuccess && <p className="text-sm text-green-500">パスワードを変更しました！</p>}
        <button
          type="submit"
          disabled={pwLoading}
          className="w-full rounded-lg bg-slate-600 py-2 text-sm font-semibold text-white hover:bg-slate-700 disabled:opacity-50 transition"
        >
          {pwLoading ? "変更中..." : "パスワードを変更する"}
        </button>
      </form>

      {/* クイズ解答履歴 */}
      {quizResults.length > 0 && (
        <div className="mt-8">
          <h2 className="mb-4 text-lg font-semibold text-slate-800">🎯 クイズ解答履歴</h2>
          <div className="space-y-2">
            {quizResults.map((result, index) => {
              const rate = result.total > 0 ? Math.round((result.score / result.total) * 100) : 0;
              return (
                <div key={result.id} className="flex items-center justify-between rounded-xl border border-slate-200 bg-white px-4 py-3 shadow-sm">
                  <div>
                    <p className="text-sm font-medium text-slate-700">
                      {result.isMockExam
                        ? `模擬試験：${result.lessonTitle || "不明"}`
                        : result.lessonTitle
                          ? `${result.lessonTitle}`
                          : `クイズ結果 #${quizResults.length - index}`}
                    </p>
                    <p className="mt-0.5 text-xs text-slate-400">
                      {new Date(result.takenAt).toLocaleString("ja-JP")}
                    </p>
                  </div>
                  <span className={`rounded-full px-3 py-1 text-xs font-bold ${getScoreBadgeClass(result.score, result.total)}`}>
                    {result.score} / {result.total}問 ({rate}%)
                  </span>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* バッジ一覧 */}
      {badges.length > 0 && (
        <div className="mt-8">
          <h2 className="mb-4 text-lg font-semibold text-slate-800">🏅 獲得バッジ</h2>
          <div className="grid grid-cols-3 gap-3">
            {badges.map((ub) => (
              <div
                key={ub.id}
                className="rounded-xl border border-slate-200 bg-white p-3 text-center shadow-sm"
              >
                <p className="text-3xl">🏅</p>
                <p className="mt-1 text-xs font-semibold text-slate-700 leading-tight">{ub.badge.name}</p>
                <p className="mt-0.5 text-xs text-slate-400">
                  {new Date(ub.earnedAt).toLocaleDateString("ja-JP")}
                </p>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* メモ一覧 */}
      {notes.length > 0 && (
        <div className="mt-8">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-lg font-semibold text-slate-800">📝 メモ一覧</h2>
            <button
              onClick={handleExportNotes}
              className="flex items-center gap-1.5 rounded-lg border border-slate-300 px-3 py-1.5 text-xs font-semibold text-slate-600 hover:bg-slate-50 transition"
            >
              📥 Markdownでエクスポート
            </button>
          </div>
          <div className="space-y-3">
            {notes.map((note) => (
              <div key={note.id} className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
                <Link
                  href={`/lessons/${note.lessonId}`}
                  className="text-sm font-medium text-pink-500 hover:underline"
                >
                  {note.lessonTitle || "レッスンへ移動"}
                </Link>
                <p className="mt-2 whitespace-pre-wrap text-sm text-slate-600">{note.content}</p>
                <p className="mt-1 text-xs text-slate-400">
                  {new Date(note.updatedAt).toLocaleString("ja-JP")}
                </p>
              </div>
            ))}
          </div>
        </div>
      )}
    </main>
  );
}
