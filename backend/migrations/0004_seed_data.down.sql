-- シードデータの削除（依存関係の順序に従って削除）
DELETE FROM quiz_choices WHERE id IN (
  'q-fp3-1-c1','q-fp3-1-c2','q-fp3-1-c3',
  'q-fp3-2-c1','q-fp3-2-c2','q-fp3-2-c3'
);
DELETE FROM quiz_questions WHERE id IN ('q-fp3-1','q-fp3-2');
DELETE FROM quizzes WHERE id = 'quiz-fp3-s1-l1';

DELETE FROM glossary_term_tags WHERE term_id IN (
  'g01','g04','g05','g06','g07','g08','g09','g10',
  'g14','g15','g16','g20','g21','g22','g23','g24','g25','g26',
  'g27','g28','g29','g30','g31','g32','g33'
);
DELETE FROM glossary_terms WHERE id IN (
  'g01','g02','g03','g04','g05','g06','g07','g08','g09','g10',
  'g11','g12','g13','g14','g15','g16','g17','g18','g19','g20',
  'g21','g22','g23','g24','g25','g26','g27','g28','g29','g30',
  'g31','g32','g33'
);
DELETE FROM glossary_tags WHERE id IN (
  'tag-fp','tag-boki','tag-asset','tag-tax','tag-insurance'
);

DELETE FROM lessons WHERE section_id IN (
  'fp3-s1','fp3-s2','fp3-s3','fp3-s4','fp3-s5',
  'boki3-s1','boki3-s2','boki3-s3',
  'asset3-s1','asset3-s2','asset3-s3'
);
DELETE FROM sections WHERE course_id IN ('fp3','boki3','asset3');
DELETE FROM courses WHERE id IN ('fp3','boki3','asset3');
