export type Quiz = {
  id: string;
  lessonId: string;
  sectionId: string;
  isMockExam: boolean;
  timeLimitMinutes: number;
  createdAt: string;
};

export type QuizChoice = {
  id: string;
  questionId: string;
  choiceText: string;
  isCorrect: boolean;
};

export type QuizQuestion = {
  id: string;
  quizId: string;
  questionText: string;
  explanation: string;
  order: number;
  choices: QuizChoice[];
};

export type QuizDetail = {
  quiz: Quiz;
  questions: QuizQuestion[];
};

export type SubmitQuizAnswer = {
  questionId: string;
  choiceId: string;
};

export type UserQuizResult = {
  id: string;
  userId: string;
  quizId: string;
  score: number;
  total: number;
  takenAt: string;
};

export type SubmitQuizResponse = {
  result: UserQuizResult;
  correct: Record<string, boolean>;
};

export type QuizResultsResponse = {
  results: UserQuizResult[];
};
