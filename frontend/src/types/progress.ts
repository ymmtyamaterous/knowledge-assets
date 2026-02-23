export type UserLessonProgress = {
  id: string;
  userId: string;
  lessonId: string;
  completedAt: string;
};

export type ProgressResponse = {
  progress: UserLessonProgress[];
};

export type SectionProgress = {
  sectionId: string;
  sectionTitle: string;
  totalLessons: number;
  completedLessons: number;
  progressRate: number;
};

export type CourseProgress = {
  courseId: string;
  courseTitle: string;
  totalLessons: number;
  completedLessons: number;
  progressRate: number;
  sections: SectionProgress[];
};

export type CourseProgressResponse = {
  courseProgress: CourseProgress[];
};
