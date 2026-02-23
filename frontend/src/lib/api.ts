import type { Course, CoursesResponse } from "@/types/course";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8000";

export async function fetchCourses(): Promise<Course[]> {
  const res = await fetch(`${API_BASE_URL}/api/v1/courses`, {
    cache: "no-store",
  });
  if (!res.ok) {
    throw new Error("failed to fetch courses");
  }

  const data = (await res.json()) as Partial<CoursesResponse>;
  return Array.isArray(data.courses) ? data.courses : [];
}
