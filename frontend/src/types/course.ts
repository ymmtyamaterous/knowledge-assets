export type Course = {
  id: string;
  title: string;
  description: string;
  difficulty: string;
  estimatedHour: number;
  thumbnailUrl: string;
  order: number;
};

export type CoursesResponse = {
  courses: Course[];
};
