-- バッジマスターデータ削除
DELETE FROM badges WHERE id IN (
  'badge-fp3-s1', 'badge-fp3-s2', 'badge-fp3-s3', 'badge-fp3-s4', 'badge-fp3-s5',
  'badge-boki3-s1', 'badge-boki3-s2', 'badge-boki3-s3',
  'badge-asset3-s1', 'badge-asset3-s2', 'badge-asset3-s3',
  'badge-fp3', 'badge-boki3', 'badge-asset3'
);

-- condition_id カラム削除
ALTER TABLE badges DROP COLUMN IF EXISTS condition_id;
