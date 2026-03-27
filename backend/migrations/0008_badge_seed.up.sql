-- badges テーブルに condition_id カラムを追加
ALTER TABLE badges ADD COLUMN IF NOT EXISTS condition_id TEXT NOT NULL DEFAULT '';

-- バッジマスターデータ投入
INSERT INTO badges (id, name, description, image_url, condition_type, condition_id) VALUES
  -- FP3級 各セクション完了バッジ
  ('badge-fp3-s1', 'FP3級 第1章マスター', 'ライフプランと資金計画を修了しました', '', 'section_complete', 'fp3-s1'),
  ('badge-fp3-s2', 'FP3級 第2章マスター', '保険の基礎を修了しました', '', 'section_complete', 'fp3-s2'),
  ('badge-fp3-s3', 'FP3級 第3章マスター', '金融資産運用を修了しました', '', 'section_complete', 'fp3-s3'),
  ('badge-fp3-s4', 'FP3級 第4章マスター', 'タックスプランニングを修了しました', '', 'section_complete', 'fp3-s4'),
  ('badge-fp3-s5', 'FP3級 第5章マスター', '不動産と相続を修了しました', '', 'section_complete', 'fp3-s5'),
  -- 簿記3級 各セクション完了バッジ
  ('badge-boki3-s1', '簿記3級 第1章マスター', '簿記の基本を修了しました', '', 'section_complete', 'boki3-s1'),
  ('badge-boki3-s2', '簿記3級 第2章マスター', '主要な勘定科目を修了しました', '', 'section_complete', 'boki3-s2'),
  ('badge-boki3-s3', '簿記3級 第3章マスター', '決算と財務諸表を修了しました', '', 'section_complete', 'boki3-s3'),
  -- 資産運用検定3級 各セクション完了バッジ
  ('badge-asset3-s1', '資産運用 第1章マスター', '投資の基本を修了しました', '', 'section_complete', 'asset3-s1'),
  ('badge-asset3-s2', '資産運用 第2章マスター', '債券と投資信託を修了しました', '', 'section_complete', 'asset3-s2'),
  ('badge-asset3-s3', '資産運用 第3章マスター', 'ポートフォリオとリスク管理を修了しました', '', 'section_complete', 'asset3-s3'),
  -- コース完了バッジ
  ('badge-fp3',    'FP3級コース修了',    'FP3級コースの全レッスンを修了しました',       '', 'course_complete', 'fp3'),
  ('badge-boki3',  '簿記3級コース修了',  '簿記3級コースの全レッスンを修了しました',     '', 'course_complete', 'boki3'),
  ('badge-asset3', '資産運用コース修了', '資産運用検定3級コースの全レッスンを修了しました', '', 'course_complete', 'asset3')
ON CONFLICT (id) DO NOTHING;
