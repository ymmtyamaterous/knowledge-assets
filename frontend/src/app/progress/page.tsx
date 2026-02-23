"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { fetchCourseProgress } from "@/lib/api";
import type { CourseProgress } from "@/types/progress";
import { useAuth } from "@/features/auth/AuthContext";

export default function ProgressPage() {
  const { user, loading } = useAuth();
  const [courseProgress, setCourseProgress] = useState<CourseProgress[]>([]);
  const [openCourseId, setOpenCourseId] = useState<string>("");

  useEffect(() => {
    if (!user) {
      setCourseProgress([]);
      return;
    }
    fetchCourseProgress()
      .then((res) => {
        const list = res ?? [];
        setCourseProgress(list);
        if (list.length > 0) {
          setOpenCourseId(list[0].courseId);
        }
      })
      .catch(() => setCourseProgress([]));
  }, [user]);

  if (!loading && !user) {
    return (
      <main>
        <h1 className="mb-4 text-2xl font-bold text-slate-800">学習進捗</h1>
        <p className="text-slate-600">進捗を表示するにはログインしてください。</p>
        <Link href="/auth/login" className="mt-4 inline-block rounded-md bg-pink-500 px-4 py-2 text-sm font-semibold text-white hover:bg-pink-600">
          ログインへ
        </Link>
      </main>
    );
  }

  return (
    <main>
      <h1 className="mb-6 text-2xl font-bold text-slate-800">学習進捗</h1>

      <div className="space-y-4">
        {courseProgress.map((cp) => (
          <div key={cp.courseId} className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
            <button
              onClick={() => setOpenCourseId((prev) => (prev === cp.courseId ? "" : cp.courseId))}
              className="flex w-full items-center justify-between text-left"
            >
              <div>
                <h2 className="text-lg font-semibold text-slate-800">{cp.courseTitle}</h2>
                <p className="text-xs text-slate-500">
                  {cp.completedLessons} / {cp.totalLessons} レッスン完了
                </p>
              </div>
              <span className="text-sm font-bold text-pink-600">{cp.progressRate}%</span>
            </button>

            <div className="mt-3 h-2 rounded-full bg-slate-100">
              <div
                className="h-2 rounded-full bg-pink-500"
                style={{ width: `${cp.progressRate}%` }}
              />
            </div>

            {openCourseId === cp.courseId && (
              <div className="mt-4 space-y-2 border-t border-slate-100 pt-3">
                {cp.sections.map((section) => (
                  <div key={section.sectionId} className="rounded-lg bg-slate-50 px-3 py-2">
                    <div className="flex items-center justify-between">
                      <span className="text-sm font-medium text-slate-700">{section.sectionTitle}</span>
                      <span className="text-xs font-semibold text-pink-600">{section.progressRate}%</span>
                    </div>
                    <p className="mt-1 text-xs text-slate-500">
                      {section.completedLessons} / {section.totalLessons} 完了
                    </p>
                  </div>
                ))}
              </div>
            )}
          </div>
        ))}
      </div>
    </main>
  );
}
