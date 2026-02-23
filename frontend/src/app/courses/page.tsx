import Link from "next/link";
import { fetchCourses } from "@/lib/api";

export default async function CoursesPage() {
  const courses = await fetchCourses().catch(() => []);

  return (
    <main>
      <h1 className="mb-6 text-2xl font-bold text-slate-800">コース一覧</h1>

      {courses.length === 0 ? (
        <p className="text-sm text-slate-500">コースデータを取得できませんでした。</p>
      ) : (
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {courses.map((course) => (
            <Link
              key={course.id}
              href={`/courses/${course.id}`}
              className="group rounded-2xl border border-slate-200 bg-white p-6 shadow-sm transition hover:shadow-md"
            >
              <h2 className="mb-2 text-lg font-bold text-slate-800 group-hover:text-pink-500">
                {course.title}
              </h2>
              <p className="mb-4 text-sm text-slate-500 line-clamp-3">{course.description}</p>
              <div className="flex items-center justify-between text-xs text-slate-400">
                <span>難易度: {course.difficulty}</span>
                <span>約 {course.estimatedHour} 時間</span>
              </div>
            </Link>
          ))}
        </div>
      )}
    </main>
  );
}
