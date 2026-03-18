export type UserNote = {
  id: string;
  userId: string;
  lessonId: string;
  content: string;
  createdAt: string;
  updatedAt: string;
};

export type NoteResponse = {
  note: UserNote | null;
};

export type NotesResponse = {
  notes: UserNote[];
};
