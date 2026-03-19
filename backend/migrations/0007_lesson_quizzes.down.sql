-- 0007_lesson_quizzes.down.sql
-- 追加したレッスン単位クイズのみ削除

DELETE FROM quiz_choices
WHERE id IN ('q-fp3-1-c4', 'q-fp3-2-c4');

DELETE FROM quiz_questions
WHERE id IN ('q-fp3-s1-l1-3');

DELETE FROM quizzes
WHERE id IN (
  'quiz-fp3-s1-l2',
  'quiz-fp3-s1-l3',
  'quiz-fp3-s1-l4',
  'quiz-fp3-s2-l1',
  'quiz-fp3-s2-l2',
  'quiz-fp3-s2-l3',
  'quiz-fp3-s3-l1',
  'quiz-fp3-s3-l2',
  'quiz-fp3-s3-l3',
  'quiz-fp3-s4-l1',
  'quiz-fp3-s4-l2',
  'quiz-fp3-s4-l3',
  'quiz-fp3-s5-l1',
  'quiz-fp3-s5-l2',
  'quiz-boki3-s1-l1',
  'quiz-boki3-s1-l2',
  'quiz-boki3-s1-l3',
  'quiz-boki3-s2-l1',
  'quiz-boki3-s2-l2',
  'quiz-boki3-s2-l3',
  'quiz-boki3-s3-l1',
  'quiz-boki3-s3-l2',
  'quiz-boki3-s3-l3',
  'quiz-asset3-s1-l1',
  'quiz-asset3-s1-l2',
  'quiz-asset3-s1-l3',
  'quiz-asset3-s2-l1',
  'quiz-asset3-s2-l2',
  'quiz-asset3-s2-l3',
  'quiz-asset3-s3-l1',
  'quiz-asset3-s3-l2'
);
