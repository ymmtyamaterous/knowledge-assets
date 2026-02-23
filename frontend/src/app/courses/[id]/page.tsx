import { fetchCourse, fetchSections, fetchLessons } from "@/lib/api";
import type { Section } from "@/types/lesson";
import Link from "next/link";
import { notFound } from "next/navigation";

type Props = {
  params: Promise<{ id: string }>;
};

async function SectionBlock({ section }: { section: Section }) {
  const lessons = await fetchLessons(section.id).catch(() => []);

  return (
    <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
      <h3 className="font-semibold text-slate-800">{section.title}</h3>
      {section.description && (
        <p className="mt-1 text-sm text-slate-500">{section.description}</p>
      )}
      {lessons.length > 0 && (
        <ul className="mt-3 space-y-1">
          {lessons.map((lesson) => (
            <li key={lesson.id}>
              <Link
                href={`/lessons/${lesson.id}`}
                className="flex items-center gap-2 rounded-lg px-3 py-2 text-sm text-slate-700 hover:bg-pink-50 hover:text-pink-600"
              >
                <span className="text-pink-300">▷</span>
                {lesson.title}
              </Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

export default async function CourseDetailPage({ params }: Props) {
  const { id } = await params;

  const [course, sections] = await Promise.all([
    fetchCourse(id).catch(() => null),
    fetchSections(id).catch(() => []),
  ]);

  if (!course) return notFound();

  const difficultyLabel: Record<string, string> = {
    beginner: "初級",
    intermediate: "中級",
    advanced: "上級",
  };

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
        <div className="mt-3 flex gap-3 text-sm text-slate-500">
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
              <SectionBlock key={section.id} section={section} />
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
