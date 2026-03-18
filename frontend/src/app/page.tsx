import { fetchCourses, fetchDailyTerm } from "@/lib/api";
import Link from "next/link";

export default async function HomePage() {
  const [courses, dailyTerm] = await Promise.all([
    fetchCourses().catch(() => []),
    fetchDailyTerm().catch(() => null),
  ]);

  return (
    <main>
      <section className="mb-8 rounded-2xl border border-pink-200 bg-white p-6 shadow-sm">
        <h1 className="text-3xl font-bold text-pink-500">アセナレへようこそ</h1>
        <p className="mt-2 text-slate-600">
          桜色のやさしいUIで、FP・簿記・資産運用の基礎を体系的に学習できます。
        </p>
      </section>

      {dailyTerm && (
        <section className="mb-8">
          <div className="rounded-2xl border border-pink-100 bg-gradient-to-r from-pink-50 to-rose-50 p-5 shadow-sm">
            <div className="mb-3 flex items-center gap-2">
              <span className="rounded-full bg-pink-500 px-3 py-0.5 text-xs font-bold text-white">今日の用語</span>
            </div>
            <p className="text-lg font-bold text-slate-800">
              {dailyTerm.term}
              {dailyTerm.reading && (
                <span className="ml-2 text-sm font-normal text-slate-500">（{dailyTerm.reading}）</span>
              )}
            </p>
            <p className="mt-2 text-sm text-slate-600">{dailyTerm.definition}</p>
            <Link href="/glossary" className="mt-3 inline-block text-xs font-medium text-pink-500 hover:underline">
              用語辞典をもっと見る →
            </Link>
          </div>
        </section>
      )}

      <section>
        <h2 className="mb-4 text-xl font-semibold">コース一覧</h2>
        <div className="grid gap-4 md:grid-cols-3">
          {courses.map((course) => (
            <Link
              key={course.id}
              href={`/courses/${course.id}`}
              className="block rounded-xl border border-slate-200 bg-white p-4 shadow-sm transition hover:border-pink-300 hover:shadow-md"
            >
              <h3 className="font-semibold text-slate-800">{course.title}</h3>
              <p className="mt-2 text-sm text-slate-600">{course.description}</p>
              <p className="mt-3 text-xs text-slate-500">想定学習時間: {course.estimatedHour}h</p>
            </Link>
          ))}
          {courses.length === 0 && (
            <p className="rounded-lg border border-dashed border-slate-300 p-4 text-sm text-slate-500">
              コース情報を取得できませんでした。
            </p>
          )}
        </div>
      </section>
    </main>
  );
}
