import { fetchCourses } from "@/lib/api";

export default async function HomePage() {
  const courses = await fetchCourses().catch(() => []);

  return (
    <main className="mx-auto max-w-5xl p-6">
      <section className="mb-8 rounded-2xl border border-pink-200 bg-white p-6 shadow-sm">
        <h1 className="text-3xl font-bold text-pink-500">アセナレ</h1>
        <p className="mt-2 text-sm text-slate-600">桜色のやさしいUIで、お金の基礎を体系的に学習できます。</p>
      </section>

      <section>
        <h2 className="mb-4 text-xl font-semibold">コース一覧</h2>
        <div className="grid gap-4 md:grid-cols-3">
          {courses.map((course) => (
            <article key={course.id} className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
              <h3 className="font-semibold">{course.title}</h3>
              <p className="mt-2 text-sm text-slate-600">{course.description}</p>
              <p className="mt-3 text-xs text-slate-500">想定学習時間: {course.estimatedHour}h</p>
            </article>
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
