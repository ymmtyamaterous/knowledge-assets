export type Section = {
  id: string;
  courseId: string;
  title: string;
  description: string;
  order: number;
};

export type Lesson = {
  id: string;
  sectionId: string;
  title: string;
  content: string;
  order: number;
};

export type SectionsResponse = {
  sections: Section[];
};

export type LessonsResponse = {
  lessons: Lesson[];
};
