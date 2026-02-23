export type UserLessonProgress = {
  id: string;
  userId: string;
  lessonId: string;
  completedAt: string;
};

export type ProgressResponse = {
  progress: UserLessonProgress[];
};
