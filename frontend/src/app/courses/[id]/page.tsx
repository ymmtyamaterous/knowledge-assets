"use client";

import { use, useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import {
  completeLesson,
  fetchCourse,
  fetchSections,
  fetchLessons,
  fetchMyProgress,
  uncompleteLesson,
} from "@/lib/api";
import type { Lesson, Section } from "@/types/lesson";
import type { Course } from "@/types/course";
import { useAuth } from "@/features/auth/AuthContext";
import Link from "next/link";

type Props = {
  params: Promise<{ id: string }>;
};

function SectionBlock({
  section,
  lessons,
  completedLessonIds,
  canToggle,
  savingLessonIds,
  onToggle,
}: {
  section: Section;
  lessons: Lesson[];
  completedLessonIds: Set<string>;
  canToggle: boolean;
  savingLessonIds: Set<string>;
  onToggle: (lessonId: string, completed: boolean) => void;
}) {
  return (
    <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
      <h3 className="font-semibold text-slate-800">{section.title}</h3>
      {section.description && (
        <p className="mt-1 text-sm text-slate-500">{section.description}</p>
      )}
      {lessons.length > 0 && (
        <ul className="mt-3 space-y-1">
          {lessons.map((lesson) => (
            <li key={lesson.id} className="flex items-center gap-2 rounded-lg">
              <Link
                href={`/lessons/${lesson.id}`}
                className="flex min-w-0 flex-1 items-start gap-2 rounded-lg px-3 py-2 text-sm text-slate-700 hover:bg-pink-50 hover:text-pink-600"
              >
                <span className="mt-0.5 text-pink-300">▷</span>
                <span className="leading-6">{lesson.title}</span>
                {completedLessonIds.has(lesson.id) && (
                  <span className="ml-auto rounded-full bg-pink-100 px-2 py-0.5 text-xs font-bold text-pink-600">✓ 完了</span>
                )}
              </Link>

              {canToggle && (
                <button
                  onClick={() => onToggle(lesson.id, completedLessonIds.has(lesson.id))}
                  disabled={savingLessonIds.has(lesson.id)}
                  className="shrink-0 rounded-md border border-slate-300 px-2 py-1 text-xs font-semibold text-slate-600 hover:bg-slate-50 disabled:opacity-50"
                >
                  {savingLessonIds.has(lesson.id)
                    ? "更新中..."
                    : completedLessonIds.has(lesson.id)
                      ? "未完了へ"
                      : "完了へ"}
                </button>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

export default function CourseDetailPage({ params }: Props) {
  const { id } = use(params);
  const { user } = useAuth();
  const router = useRouter();

  const [course, setCourse] = useState<Course | null>(null);
  const [sections, setSections] = useState<Section[]>([]);
  const [lessonsBySection, setLessonsBySection] = useState<Record<string, Lesson[]>>({});
  const [completedLessonIds, setCompletedLessonIds] = useState<Set<string>>(new Set());
  const [savingLessonIds, setSavingLessonIds] = useState<Set<string>>(new Set());
  const [notFound, setNotFound] = useState(false);

  useEffect(() => {
    fetchCourse(id)
      .then(setCourse)
      .catch(() => setNotFound(true));
  }, [id]);

  useEffect(() => {
    fetchSections(id)
      .then((res) => setSections(res ?? []))
      .catch(() => setSections([]));
  }, [id]);

  useEffect(() => {
    if (sections.length === 0) {
      setLessonsBySection({});
      return;
    }

    Promise.all(
      sections.map(async (section) => ({
        sectionId: section.id,
        lessons: await fetchLessons(section.id).catch(() => []),
      })),
    ).then((pairs) => {
      const next: Record<string, Lesson[]> = {};
      pairs.forEach((pair) => {
        next[pair.sectionId] = pair.lessons ?? [];
      });
      setLessonsBySection(next);
    });
  }, [sections]);

  useEffect(() => {
    if (!user) {
      setCompletedLessonIds(new Set());
      return;
    }
    fetchMyProgress()
      .then((progress) => {
        const ids = new Set((progress ?? []).map((p) => p.lessonId));
        setCompletedLessonIds(ids);
      })
      .catch(() => setCompletedLessonIds(new Set()));
  }, [user]);

  const difficultyLabel = useMemo<Record<string, string>>(() => ({
    beginner: "初級",
    intermediate: "中級",
    advanced: "上級",
  }), []);

  const handleToggleLessonComplete = async (lessonId: string, completed: boolean) => {
    if (!user) {
      router.push("/auth/login");
      return;
    }

    setSavingLessonIds((prev) => {
      const next = new Set(prev);
      next.add(lessonId);
      return next;
    });

    try {
      if (completed) {
        await uncompleteLesson(lessonId);
        setCompletedLessonIds((prev) => {
          const next = new Set(prev);
          next.delete(lessonId);
          return next;
        });
      } else {
        await completeLesson(lessonId);
        setCompletedLessonIds((prev) => {
          const next = new Set(prev);
          next.add(lessonId);
          return next;
        });
      }
    } finally {
      setSavingLessonIds((prev) => {
        const next = new Set(prev);
        next.delete(lessonId);
        return next;
      });
    }
  };

  if (notFound) {
    return (
      <main>
        <p className="text-slate-500">コースが見つかりませんでした。</p>
        <Link href="/" className="mt-4 inline-block text-sm text-pink-400 hover:underline">
          ← トップに戻る
        </Link>
      </main>
    );
  }

  if (!course) {
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
        <Link href="/" className="text-sm text-pink-400 hover:underline">
          ← トップに戻る
        </Link>
      </div>

      <div className="mb-6 rounded-2xl border border-pink-100 bg-white p-6 shadow-sm">
        <h1 className="text-2xl font-bold text-slate-800">{course.title}</h1>
        <p className="mt-2 text-slate-600">{course.description}</p>
        <div className="mt-3 flex flex-wrap gap-3 text-sm text-slate-500">
          <span className="rounded-full bg-pink-50 px-3 py-1">
            難易度: {difficultyLabel[course.difficulty] ?? course.difficulty}
          </span>
          <span className="rounded-full bg-pink-50 px-3 py-1">
            想定学習時間: {course.estimatedHour}h
          </span>
        </div>
      </div>

      <section>
        <h2 className="mb-4 text-lg font-semibold">カリキュラム</h2>
        {sections.length > 0 ? (
          <div className="space-y-4">
            {sections.map((section) => (
              <SectionBlock
                key={section.id}
                section={section}
                lessons={lessonsBySection[section.id] ?? []}
                completedLessonIds={completedLessonIds}
                canToggle={!!user}
                savingLessonIds={savingLessonIds}
                onToggle={handleToggleLessonComplete}
              />
            ))}
          </div>
        ) : (
          <p className="rounded-lg border border-dashed border-slate-300 p-4 text-sm text-slate-500">
            セクションがまだありません。
          </p>
        )}
      </section>
    </main>
  );
}
