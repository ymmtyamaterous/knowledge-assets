import type { Course, CoursesResponse } from "@/types/course";
import type { Section, Lesson, SectionsResponse, LessonsResponse } from "@/types/lesson";
import type { User, AuthResponse } from "@/types/user";
import type {
  ProgressResponse,
  UserLessonProgress,
  CourseProgress,
  CourseProgressResponse,
  UserStreak,
  UserStats,
  UserCalendar,
} from "@/types/progress";
import type {
  GlossaryResponse,
  GlossaryTerm,
  GlossaryTag,
  GlossaryTagsResponse,
} from "@/types/glossary";
import type {
  Quiz,
  QuizDetail,
  SubmitQuizAnswer,
  SubmitQuizResponse,
  QuizResultsResponse,
  UserQuizResult,
} from "@/types/quiz";
import type { UserNote, NoteResponse, NotesResponse } from "@/types/note";
import { authHeaders } from "@/lib/token";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8000";

async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...authHeaders(),
      ...(init?.headers ?? {}),
    },
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({})) as { error?: string };
    throw new Error(body.error ?? `HTTP ${res.status}`);
  }
  return res.json() as Promise<T>;
}

// ---- Auth ----

export async function registerUser(email: string, password: string, username: string): Promise<AuthResponse> {
  return apiFetch<AuthResponse>("/api/v1/auth/register", {
    method: "POST",
    body: JSON.stringify({ email, password, username }),
  });
}

export async function loginUser(email: string, password: string): Promise<AuthResponse> {
  return apiFetch<AuthResponse>("/api/v1/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
}

// ---- Users ----

export async function fetchMe(): Promise<User> {
  return apiFetch<User>("/api/v1/users/me");
}

export async function updateMe(data: { username?: string; avatarUrl?: string }): Promise<User> {
  return apiFetch<User>("/api/v1/users/me", {
    method: "PUT",
    body: JSON.stringify(data),
  });
}

export async function changePassword(currentPassword: string, newPassword: string): Promise<void> {
  await apiFetch<{ status: string }>("/api/v1/users/me/password", {
    method: "PUT",
    body: JSON.stringify({ currentPassword, newPassword }),
  });
}

// ---- Courses ----

export async function fetchCourses(): Promise<Course[]> {
  const data = await apiFetch<Partial<CoursesResponse>>("/api/v1/courses", {
    cache: "no-store",
  });
  return Array.isArray(data.courses) ? data.courses : [];
}

export async function fetchCourse(id: string): Promise<Course> {
  return apiFetch<Course>(`/api/v1/courses/${id}`);
}

// ---- Sections ----

export async function fetchSections(courseId: string): Promise<Section[]> {
  const data = await apiFetch<Partial<SectionsResponse>>(`/api/v1/courses/${courseId}/sections`);
  return Array.isArray(data.sections) ? data.sections : [];
}

// ---- Lessons ----

export async function fetchLessons(sectionId: string): Promise<Lesson[]> {
  const data = await apiFetch<Partial<LessonsResponse>>(`/api/v1/sections/${sectionId}/lessons`);
  return Array.isArray(data.lessons) ? data.lessons : [];
}

export async function fetchLesson(id: string): Promise<Lesson> {
  return apiFetch<Lesson>(`/api/v1/lessons/${id}`);
}

export async function completeLesson(lessonId: string): Promise<UserLessonProgress> {
  return apiFetch<UserLessonProgress>(`/api/v1/lessons/${lessonId}/complete`, {
    method: "POST",
  });
}

export async function uncompleteLesson(lessonId: string): Promise<{ status: string }> {
  return apiFetch<{ status: string }>(`/api/v1/lessons/${lessonId}/complete`, {
    method: "DELETE",
  });
}

// ---- Progress ----

export async function fetchMyProgress(): Promise<UserLessonProgress[]> {
  const data = await apiFetch<Partial<ProgressResponse>>("/api/v1/users/me/progress");
  return Array.isArray(data.progress) ? data.progress : [];
}

// ---- Glossary ----

export async function fetchGlossary(tagId?: string): Promise<GlossaryTerm[]> {
  const query = tagId ? `?tagId=${encodeURIComponent(tagId)}` : "";
  const data = await apiFetch<Partial<GlossaryResponse>>(`/api/v1/glossary${query}`);
  return Array.isArray(data.terms) ? data.terms : [];
}

export async function fetchGlossaryTerm(id: string): Promise<GlossaryTerm> {
  return apiFetch<GlossaryTerm>(`/api/v1/glossary/${id}`);
}

export async function fetchGlossaryTags(): Promise<GlossaryTag[]> {
  const data = await apiFetch<Partial<GlossaryTagsResponse>>("/api/v1/glossary/tags");
  return Array.isArray(data.tags) ? data.tags : [];
}

export async function fetchDailyTerm(): Promise<GlossaryTerm | null> {
  return apiFetch<GlossaryTerm>("/api/v1/glossary/daily").catch(() => null);
}

export async function fetchCourseProgress(): Promise<CourseProgress[]> {
  const data = await apiFetch<Partial<CourseProgressResponse>>("/api/v1/users/me/course-progress");
  return Array.isArray(data.courseProgress) ? data.courseProgress : [];
}

export async function fetchMyStreak(): Promise<UserStreak> {
  const data = await apiFetch<Partial<UserStreak>>("/api/v1/users/me/streak");
  return {
    currentStreak: data.currentStreak ?? 0,
    longestStreak: data.longestStreak ?? 0,
    lastStudiedAt: data.lastStudiedAt ?? "",
  };
}

export async function fetchMyStats(): Promise<UserStats> {
  const data = await apiFetch<Partial<UserStats>>("/api/v1/users/me/stats");
  return {
    totalCompletedLessons: data.totalCompletedLessons ?? 0,
    totalStudyDays: data.totalStudyDays ?? 0,
    totalNotes: data.totalNotes ?? 0,
    averageQuizScore: data.averageQuizScore ?? 0,
  };
}

export async function fetchMyCalendar(year?: number): Promise<UserCalendar> {
  const query = year ? `?year=${year}` : "";
  const data = await apiFetch<Partial<UserCalendar>>(`/api/v1/users/me/calendar${query}`);
  return { days: Array.isArray(data.days) ? data.days : [] };
}

export async function fetchLessonQuiz(lessonId: string): Promise<Quiz> {
  return apiFetch<Quiz>(`/api/v1/lessons/${lessonId}/quiz`);
}

export async function fetchQuiz(id: string): Promise<QuizDetail> {
  return apiFetch<QuizDetail>(`/api/v1/quizzes/${id}`);
}

export async function submitQuiz(id: string, answers: SubmitQuizAnswer[]): Promise<SubmitQuizResponse> {
  return apiFetch<SubmitQuizResponse>(`/api/v1/quizzes/${id}/submit`, {
    method: "POST",
    body: JSON.stringify({ answers }),
  });
}

export async function fetchMyQuizResults(): Promise<UserQuizResult[]> {
  const data = await apiFetch<Partial<QuizResultsResponse>>("/api/v1/users/me/quiz-results");
  return Array.isArray(data.results) ? data.results : [];
}

// ---- Notes ----

export async function fetchLessonNote(lessonId: string): Promise<UserNote | null> {
  const data = await apiFetch<Partial<NoteResponse>>(`/api/v1/lessons/${lessonId}/note`);
  return data.note ?? null;
}

export async function saveLessonNote(lessonId: string, content: string): Promise<UserNote> {
  const data = await apiFetch<{ note: UserNote }>(`/api/v1/lessons/${lessonId}/note`, {
    method: "PUT",
    body: JSON.stringify({ content }),
  });
  return data.note;
}

export async function fetchMyNotes(): Promise<UserNote[]> {
  const data = await apiFetch<Partial<NotesResponse>>("/api/v1/users/me/notes");
  return Array.isArray(data.notes) ? data.notes : [];
}


