export type SearchLesson = {
  id: string;
  title: string;
  sectionId: string;
};

export type SearchTerm = {
  id: string;
  term: string;
  reading: string;
};

export type SearchResult = {
  lessons: SearchLesson[];
  terms: SearchTerm[];
};
