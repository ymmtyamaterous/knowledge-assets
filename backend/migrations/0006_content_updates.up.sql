-- ============================================================
-- F-3: fp3-s1-l1（ライフプランとは）にキャッシュフロー表の具体例を追記
-- ============================================================
UPDATE lessons SET content = $$# ライフプランとは

ライフプランとは、人生のさまざまなイベントを見据えて、長期的な生活設計を行うことです。

## 主なライフイベント
- 結婚・出産
- マイホーム購入
- 子どもの教育
- 老後の生活

各イベントに必要な資金を試算し、収入・支出のバランスを取ることが重要です。

## ライフプランの6つのステップ
1. 現状把握（収入・支出・資産・負債の整理）
2. 目標設定（いつまでに何をしたいか）
3. 必要資金の算出
4. 資金調達方法の検討
5. 計画の実行
6. 定期的な見直し

## キャッシュフロー表
ライフプランを数字で見える化するために「キャッシュフロー表」を作成します。年ごとの収入・支出・貯蓄残高を一覧にすることで、いつ資金不足になるかが一目でわかります。

## キャッシュフロー表の例

夫(30歳)・妻(28歳)の家庭を例にしたキャッシュフロー表です。

| 年 | 家族の変化 | 収入合計 | 支出合計 | 年間収支 | 貯蓄残高 |
|----|----------|---------|---------|---------|---------|
| 2026年 | 現在 | 550万円 | 480万円 | +70万円 | 200万円 |
| 2027年 | 第1子誕生 | 530万円 | 520万円 | +10万円 | 210万円 |
| 2028年 | 育児費用増加 | 550万円 | 540万円 | +10万円 | 220万円 |
| 2029年 | 妻職場復帰 | 650万円 | 520万円 | +130万円 | 350万円 |
| 2030年 | マイホーム購入 | 660万円 | 900万円 | -240万円 | 110万円 |
| 2031年 | ローン返済開始 | 670万円 | 600万円 | +70万円 | 180万円 |

### 見方のポイント
- **2030年**：マイホーム購入で年間収支がマイナスになり、貯蓄残高が大きく減少
- **貯蓄残高がマイナスになる年**が「資金不足の危険時期」→ 事前に準備が必要
- インフレ率（1〜2%）を考慮すると、より現実的なシミュレーションができる
- 定期的に見直しを行い、実績と差異があれば修正する
$$, updated_at = NOW()
WHERE id = 'fp3-s1-l1';

-- ============================================================
-- F-4: クイズコンテンツの充実
-- 各セクション（章）に5問以上のセクション復習テストを追加
-- ============================================================

-- ---- FP3級 第1章 セクション復習テスト ----
INSERT INTO quizzes (id, lesson_id, section_id, is_mock_exam, time_limit_minutes) VALUES
  ('quiz-fp3-s1', NULL, 'fp3-s1', FALSE, 15)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_questions (id, quiz_id, question_text, explanation, "order") VALUES
  ('q-fp3-s1-1', 'quiz-fp3-s1', 'キャッシュフロー表に記載する主な項目はどれですか？',
   'キャッシュフロー表には年ごとの収入・支出・年間収支・貯蓄残高を記載します。', 1),
  ('q-fp3-s1-2', 'quiz-fp3-s1', 'ライフプラン作成時に「固定費」として分類されるものはどれですか？',
   '家賃・住宅ローン・保険料・通信費など毎月一定の支出が固定費です。', 2),
  ('q-fp3-s1-3', 'quiz-fp3-s1', '子どもの教育費として最も高額になるのはどの段階ですか（私立の場合）？',
   '私立小学校の6年間は約959万円で最も高額です。', 3),
  ('q-fp3-s1-4', 'quiz-fp3-s1', '住宅ローン控除（住宅借入金等特別控除）の控除率は年末ローン残高の何％ですか（2022年改正後）？',
   '2022年の税制改正により控除率は1.0%から0.7%に引き下げられました。', 4),
  ('q-fp3-s1-5', 'quiz-fp3-s1', '純資産の計算式として正しいものはどれですか？',
   '純資産 = 資産合計 − 負債合計 です。バランスシートの基本公式です。', 5)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_choices (id, question_id, choice_text, is_correct) VALUES
  ('q-fp3-s1-1-c1', 'q-fp3-s1-1', '収入・支出・年間収支・貯蓄残高', TRUE),
  ('q-fp3-s1-1-c2', 'q-fp3-s1-1', '株価・金利・為替・物価', FALSE),
  ('q-fp3-s1-1-c3', 'q-fp3-s1-1', '資産・負債・純資産のみ', FALSE),
  ('q-fp3-s1-1-c4', 'q-fp3-s1-1', '保険料と積立金だけ', FALSE),
  ('q-fp3-s1-2-c1', 'q-fp3-s1-2', '食費', FALSE),
  ('q-fp3-s1-2-c2', 'q-fp3-s1-2', '娯楽費', FALSE),
  ('q-fp3-s1-2-c3', 'q-fp3-s1-2', '通信費', TRUE),
  ('q-fp3-s1-2-c4', 'q-fp3-s1-2', '被服費', FALSE),
  ('q-fp3-s1-3-c1', 'q-fp3-s1-3', '私立幼稚園（3年間）', FALSE),
  ('q-fp3-s1-3-c2', 'q-fp3-s1-3', '私立中学校（3年間）', FALSE),
  ('q-fp3-s1-3-c3', 'q-fp3-s1-3', '私立小学校（6年間）', TRUE),
  ('q-fp3-s1-3-c4', 'q-fp3-s1-3', '私立高校（3年間）', FALSE),
  ('q-fp3-s1-4-c1', 'q-fp3-s1-4', '0.5%', FALSE),
  ('q-fp3-s1-4-c2', 'q-fp3-s1-4', '0.7%', TRUE),
  ('q-fp3-s1-4-c3', 'q-fp3-s1-4', '1.0%', FALSE),
  ('q-fp3-s1-4-c4', 'q-fp3-s1-4', '1.5%', FALSE),
  ('q-fp3-s1-5-c1', 'q-fp3-s1-5', '資産合計 ÷ 負債合計', FALSE),
  ('q-fp3-s1-5-c2', 'q-fp3-s1-5', '負債合計 − 資産合計', FALSE),
  ('q-fp3-s1-5-c3', 'q-fp3-s1-5', '資産合計 − 負債合計', TRUE),
  ('q-fp3-s1-5-c4', 'q-fp3-s1-5', '収入合計 − 支出合計', FALSE)
ON CONFLICT (id) DO NOTHING;

-- ---- FP3級 第2章 セクション復習テスト ----
INSERT INTO quizzes (id, lesson_id, section_id, is_mock_exam, time_limit_minutes) VALUES
  ('quiz-fp3-s2', NULL, 'fp3-s2', FALSE, 15)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_questions (id, quiz_id, question_text, explanation, "order") VALUES
  ('q-fp3-s2-1', 'quiz-fp3-s2', '掛け捨て型で保険料が安く、一定期間の死亡保障を提供する保険はどれですか？',
   '定期保険は掛け捨て型で保険料が安く、一定期間の死亡保障に特化しています。', 1),
  ('q-fp3-s2-2', 'quiz-fp3-s2', '地震保険はどの保険とセットで加入することが義務付けられていますか？',
   '地震保険は火災保険とセットでのみ加入できます。単独では加入できません。', 2),
  ('q-fp3-s2-3', 'quiz-fp3-s2', '雇用保険の主な給付内容として正しいものはどれですか？',
   '雇用保険は失業時の給付や育児休業給付を提供します。', 3),
  ('q-fp3-s2-4', 'quiz-fp3-s2', '必要保障額の計算式として正しいものはどれですか？',
   '必要保障額 = 遺族の生活費の総計 − 公的年金・配偶者収入などの収入総計 です。', 4),
  ('q-fp3-s2-5', 'quiz-fp3-s2', '介護保険に加入が義務付けられているのは何歳以上からですか？',
   '介護保険は40歳以上が加入対象で、保険料の支払い義務があります。', 5)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_choices (id, question_id, choice_text, is_correct) VALUES
  ('q-fp3-s2-1-c1', 'q-fp3-s2-1', '終身保険', FALSE),
  ('q-fp3-s2-1-c2', 'q-fp3-s2-1', '養老保険', FALSE),
  ('q-fp3-s2-1-c3', 'q-fp3-s2-1', '定期保険', TRUE),
  ('q-fp3-s2-1-c4', 'q-fp3-s2-1', '介護保険', FALSE),
  ('q-fp3-s2-2-c1', 'q-fp3-s2-2', '自動車保険', FALSE),
  ('q-fp3-s2-2-c2', 'q-fp3-s2-2', '火災保険', TRUE),
  ('q-fp3-s2-2-c3', 'q-fp3-s2-2', '傷害保険', FALSE),
  ('q-fp3-s2-2-c4', 'q-fp3-s2-2', '個人賠償責任保険', FALSE),
  ('q-fp3-s2-3-c1', 'q-fp3-s2-3', '病気・ケガの医療費補助', FALSE),
  ('q-fp3-s2-3-c2', 'q-fp3-s2-3', '失業時の給付・育児休業給付', TRUE),
  ('q-fp3-s2-3-c3', 'q-fp3-s2-3', '業務上の事故への給付', FALSE),
  ('q-fp3-s2-3-c4', 'q-fp3-s2-3', '介護サービスの費用補助', FALSE),
  ('q-fp3-s2-4-c1', 'q-fp3-s2-4', '遺族の生活費の総計 × 保険料率', FALSE),
  ('q-fp3-s2-4-c2', 'q-fp3-s2-4', '遺族の生活費の総計 + 公的年金等の収入', FALSE),
  ('q-fp3-s2-4-c3', 'q-fp3-s2-4', '遺族の生活費の総計 − 公的年金・配偶者収入などの収入総計', TRUE),
  ('q-fp3-s2-4-c4', 'q-fp3-s2-4', '資産合計 − 負債合計', FALSE),
  ('q-fp3-s2-5-c1', 'q-fp3-s2-5', '20歳以上', FALSE),
  ('q-fp3-s2-5-c2', 'q-fp3-s2-5', '30歳以上', FALSE),
  ('q-fp3-s2-5-c3', 'q-fp3-s2-5', '40歳以上', TRUE),
  ('q-fp3-s2-5-c4', 'q-fp3-s2-5', '65歳以上', FALSE)
ON CONFLICT (id) DO NOTHING;

-- ---- FP3級 第3章 セクション復習テスト ----
INSERT INTO quizzes (id, lesson_id, section_id, is_mock_exam, time_limit_minutes) VALUES
  ('quiz-fp3-s3', NULL, 'fp3-s3', FALSE, 15)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_questions (id, quiz_id, question_text, explanation, "order") VALUES
  ('q-fp3-s3-1', 'quiz-fp3-s3', '株式投資において「配当金」とは何ですか？',
   '配当金は企業が得た利益の一部を株主に分配するものです。', 1),
  ('q-fp3-s3-2', 'quiz-fp3-s3', '債券の利回りと価格の関係として正しいものはどれですか？',
   '債券は価格が下がると利回りが上がり、価格が上がると利回りが下がる逆相関の関係です。', 2),
  ('q-fp3-s3-3', 'quiz-fp3-s3', '投資信託の説明として正しいものはどれですか？',
   '投資信託は多くの投資家から資金を集め、専門家が運用する金融商品です。', 3),
  ('q-fp3-s3-4', 'quiz-fp3-s3', '普通預金と定期預金の違いとして正しいものはどれですか？',
   '定期預金は満期まで引き出しに制限があるが金利が高く、普通預金はいつでも引き出せるが金利が低いです。', 4),
  ('q-fp3-s3-5', 'quiz-fp3-s3', '分散投資の目的として最も適切なものはどれですか？',
   '分散投資はリスクを分散させ、一つの資産の損失が全体に与える影響を小さくすることが目的です。', 5)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_choices (id, question_id, choice_text, is_correct) VALUES
  ('q-fp3-s3-1-c1', 'q-fp3-s3-1', '株を売却した際の利益', FALSE),
  ('q-fp3-s3-1-c2', 'q-fp3-s3-1', '企業が利益の一部を株主に分配するもの', TRUE),
  ('q-fp3-s3-1-c3', 'q-fp3-s3-1', '株式の購入手数料', FALSE),
  ('q-fp3-s3-1-c4', 'q-fp3-s3-1', '株価の上昇分', FALSE),
  ('q-fp3-s3-2-c1', 'q-fp3-s3-2', '価格が上がると利回りも上がる', FALSE),
  ('q-fp3-s3-2-c2', 'q-fp3-s3-2', '価格と利回りは無関係', FALSE),
  ('q-fp3-s3-2-c3', 'q-fp3-s3-2', '価格が下がると利回りが上がる', TRUE),
  ('q-fp3-s3-2-c4', 'q-fp3-s3-2', '価格が下がると利回りも下がる', FALSE),
  ('q-fp3-s3-3-c1', 'q-fp3-s3-3', '個人が直接株式を購入する商品', FALSE),
  ('q-fp3-s3-3-c2', 'q-fp3-s3-3', '多くの投資家から資金を集め専門家が運用する商品', TRUE),
  ('q-fp3-s3-3-c3', 'q-fp3-s3-3', '元本保証のある預金商品', FALSE),
  ('q-fp3-s3-3-c4', 'q-fp3-s3-3', '特定の企業にだけ投資する商品', FALSE),
  ('q-fp3-s3-4-c1', 'q-fp3-s3-4', '普通預金の方が金利が高い', FALSE),
  ('q-fp3-s3-4-c2', 'q-fp3-s3-4', '定期預金はいつでも引き出せる', FALSE),
  ('q-fp3-s3-4-c3', 'q-fp3-s3-4', '定期預金は満期まで制限があるが金利が高い', TRUE),
  ('q-fp3-s3-4-c4', 'q-fp3-s3-4', '両者の金利は同じ', FALSE),
  ('q-fp3-s3-5-c1', 'q-fp3-s3-5', '短期で最大利益を得ること', FALSE),
  ('q-fp3-s3-5-c2', 'q-fp3-s3-5', 'リスクを分散し一つの資産の損失影響を小さくすること', TRUE),
  ('q-fp3-s3-5-c3', 'q-fp3-s3-5', '特定の有望銘柄に集中投資すること', FALSE),
  ('q-fp3-s3-5-c4', 'q-fp3-s3-5', '元本を保証すること', FALSE)
ON CONFLICT (id) DO NOTHING;

-- ---- FP3級 第4章 セクション復習テスト ----
INSERT INTO quizzes (id, lesson_id, section_id, is_mock_exam, time_limit_minutes) VALUES
  ('quiz-fp3-s4', NULL, 'fp3-s4', FALSE, 15)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_questions (id, quiz_id, question_text, explanation, "order") VALUES
  ('q-fp3-s4-1', 'quiz-fp3-s4', '所得税の計算で「課税所得」を求めるための計算式はどれですか？',
   '課税所得 = 収入 − 必要経費 − 所得控除 です。', 1),
  ('q-fp3-s4-2', 'quiz-fp3-s4', 'つみたてNISAの非課税期間は何年ですか（旧制度）？',
   '旧つみたてNISAの非課税期間は20年でした。2024年以降の新NISAは無期限です。', 2),
  ('q-fp3-s4-3', 'quiz-fp3-s4', 'iDeCoの最大の特徴として正しいものはどれですか？',
   'iDeCoは掛金全額が所得控除の対象となり、税制優遇が受けられる年金制度です。', 3),
  ('q-fp3-s4-4', 'quiz-fp3-s4', '扶養控除の対象となる扶養親族の年齢要件はどれですか？',
   '扶養控除は16歳以上の親族が対象です（年間所得48万円以下）。', 4),
  ('q-fp3-s4-5', 'quiz-fp3-s4', '所得税の税率の特徴として正しいものはどれですか？',
   '所得税は超過累進税率で、課税所得が多いほど高い税率が適用されます。', 5)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_choices (id, question_id, choice_text, is_correct) VALUES
  ('q-fp3-s4-1-c1', 'q-fp3-s4-1', '収入 − 所得控除', FALSE),
  ('q-fp3-s4-1-c2', 'q-fp3-s4-1', '収入 − 必要経費 − 所得控除', TRUE),
  ('q-fp3-s4-1-c3', 'q-fp3-s4-1', '収入 × 税率', FALSE),
  ('q-fp3-s4-1-c4', 'q-fp3-s4-1', '収入 − 税額控除', FALSE),
  ('q-fp3-s4-2-c1', 'q-fp3-s4-2', '5年', FALSE),
  ('q-fp3-s4-2-c2', 'q-fp3-s4-2', '10年', FALSE),
  ('q-fp3-s4-2-c3', 'q-fp3-s4-2', '20年', TRUE),
  ('q-fp3-s4-2-c4', 'q-fp3-s4-2', '無期限', FALSE),
  ('q-fp3-s4-3-c1', 'q-fp3-s4-3', '運用益が非課税になる', FALSE),
  ('q-fp3-s4-3-c2', 'q-fp3-s4-3', '掛金全額が所得控除の対象になる', TRUE),
  ('q-fp3-s4-3-c3', 'q-fp3-s4-3', '元本が保証される', FALSE),
  ('q-fp3-s4-3-c4', 'q-fp3-s4-3', '60歳前でも自由に引き出せる', FALSE),
  ('q-fp3-s4-4-c1', 'q-fp3-s4-4', '15歳以上', FALSE),
  ('q-fp3-s4-4-c2', 'q-fp3-s4-4', '16歳以上', TRUE),
  ('q-fp3-s4-4-c3', 'q-fp3-s4-4', '18歳以上', FALSE),
  ('q-fp3-s4-4-c4', 'q-fp3-s4-4', '20歳以上', FALSE),
  ('q-fp3-s4-5-c1', 'q-fp3-s4-5', '全員が同じ一律の税率が適用される', FALSE),
  ('q-fp3-s4-5-c2', 'q-fp3-s4-5', '所得が低いほど高い税率が適用される', FALSE),
  ('q-fp3-s4-5-c3', 'q-fp3-s4-5', '課税所得が多いほど高い税率が適用される', TRUE),
  ('q-fp3-s4-5-c4', 'q-fp3-s4-5', '法人と同じ税率が適用される', FALSE)
ON CONFLICT (id) DO NOTHING;

-- ---- FP3級 第5章 セクション復習テスト ----
INSERT INTO quizzes (id, lesson_id, section_id, is_mock_exam, time_limit_minutes) VALUES
  ('quiz-fp3-s5', NULL, 'fp3-s5', FALSE, 15)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_questions (id, quiz_id, question_text, explanation, "order") VALUES
  ('q-fp3-s5-1', 'quiz-fp3-s5', '不動産の売買において「登記」の目的として最も正しいものはどれですか？',
   '登記は不動産の権利関係（所有権など）を公示することで、第三者への対抗要件となります。', 1),
  ('q-fp3-s5-2', 'quiz-fp3-s5', '法定相続人の範囲として正しいものはどれですか？',
   '法定相続人は配偶者と血族相続人（子・親・兄弟姉妹）です。', 2),
  ('q-fp3-s5-3', 'quiz-fp3-s5', '贈与税の基礎控除額（暦年贈与）として正しいものはどれですか？',
   '暦年贈与の基礎控除額は年間110万円です。110万円以下の贈与は非課税です。', 3),
  ('q-fp3-s5-4', 'quiz-fp3-s5', '相続税の基礎控除額の計算式はどれですか？',
   '相続税の基礎控除額 = 3,000万円 + 600万円 × 法定相続人の数 です。', 4),
  ('q-fp3-s5-5', 'quiz-fp3-s5', '固定資産税の課税対象として正しいものはどれですか？',
   '固定資産税は土地・家屋・償却資産（事業用資産）が課税対象です。', 5)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_choices (id, question_id, choice_text, is_correct) VALUES
  ('q-fp3-s5-1-c1', 'q-fp3-s5-1', '建物の品質を保証するため', FALSE),
  ('q-fp3-s5-1-c2', 'q-fp3-s5-1', '不動産の権利関係を公示するため', TRUE),
  ('q-fp3-s5-1-c3', 'q-fp3-s5-1', '税金を計算するため', FALSE),
  ('q-fp3-s5-1-c4', 'q-fp3-s5-1', 'ローンの金利を決めるため', FALSE),
  ('q-fp3-s5-2-c1', 'q-fp3-s5-2', '配偶者のみ', FALSE),
  ('q-fp3-s5-2-c2', 'q-fp3-s5-2', '配偶者と血族相続人（子・親・兄弟姉妹）', TRUE),
  ('q-fp3-s5-2-c3', 'q-fp3-s5-2', '子どものみ', FALSE),
  ('q-fp3-s5-2-c4', 'q-fp3-s5-2', '6親等内の血族すべて', FALSE),
  ('q-fp3-s5-3-c1', 'q-fp3-s5-3', '年間50万円', FALSE),
  ('q-fp3-s5-3-c2', 'q-fp3-s5-3', '年間100万円', FALSE),
  ('q-fp3-s5-3-c3', 'q-fp3-s5-3', '年間110万円', TRUE),
  ('q-fp3-s5-3-c4', 'q-fp3-s5-3', '年間200万円', FALSE),
  ('q-fp3-s5-4-c1', 'q-fp3-s5-4', '1,000万円 + 300万円 × 法定相続人の数', FALSE),
  ('q-fp3-s5-4-c2', 'q-fp3-s5-4', '3,000万円 + 600万円 × 法定相続人の数', TRUE),
  ('q-fp3-s5-4-c3', 'q-fp3-s5-4', '5,000万円 + 1,000万円 × 法定相続人の数', FALSE),
  ('q-fp3-s5-4-c4', 'q-fp3-s5-4', '一律3,000万円', FALSE),
  ('q-fp3-s5-5-c1', 'q-fp3-s5-5', '株式・債券', FALSE),
  ('q-fp3-s5-5-c2', 'q-fp3-s5-5', '土地・家屋・償却資産', TRUE),
  ('q-fp3-s5-5-c3', 'q-fp3-s5-5', '預貯金', FALSE),
  ('q-fp3-s5-5-c4', 'q-fp3-s5-5', '自動車のみ', FALSE)
ON CONFLICT (id) DO NOTHING;

-- ---- 簿記3級 第1章 セクション復習テスト ----
INSERT INTO quizzes (id, lesson_id, section_id, is_mock_exam, time_limit_minutes) VALUES
  ('quiz-boki3-s1', NULL, 'boki3-s1', FALSE, 15)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_questions (id, quiz_id, question_text, explanation, "order") VALUES
  ('q-boki3-s1-1', 'quiz-boki3-s1', '仕訳において資産が増加したとき、どちらに記録しますか？',
   '資産の増加は借方（左側）に記録します。これは簿記の基本ルールです。', 1),
  ('q-boki3-s1-2', 'quiz-boki3-s1', '「売掛金」は簿記上どのような科目に分類されますか？',
   '売掛金は商品を掛けで販売したときの代金を回収する権利で、資産科目です。', 2),
  ('q-boki3-s1-3', 'quiz-boki3-s1', '仕訳の基本ルールで「負債が増加」したとき、どちらに記録しますか？',
   '負債の増加は貸方（右側）に記録します。', 3),
  ('q-boki3-s1-4', 'quiz-boki3-s1', '「総勘定元帳」の役割として正しいものはどれですか？',
   '総勘定元帳は仕訳帳の内容を勘定科目ごとに集計・転記した帳簿です。', 4),
  ('q-boki3-s1-5', 'quiz-boki3-s1', '次のうち費用科目はどれですか？',
   '給料・仕入・広告宣伝費・減価償却費などが費用科目です。資本金は純資産科目です。', 5)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_choices (id, question_id, choice_text, is_correct) VALUES
  ('q-boki3-s1-1-c1', 'q-boki3-s1-1', '借方（左側）', TRUE),
  ('q-boki3-s1-1-c2', 'q-boki3-s1-1', '貸方（右側）', FALSE),
  ('q-boki3-s1-1-c3', 'q-boki3-s1-1', '資産の増加は記録しない', FALSE),
  ('q-boki3-s1-1-c4', 'q-boki3-s1-1', 'どちらでもよい', FALSE),
  ('q-boki3-s1-2-c1', 'q-boki3-s1-2', '負債科目', FALSE),
  ('q-boki3-s1-2-c2', 'q-boki3-s1-2', '費用科目', FALSE),
  ('q-boki3-s1-2-c3', 'q-boki3-s1-2', '資産科目', TRUE),
  ('q-boki3-s1-2-c4', 'q-boki3-s1-2', '収益科目', FALSE),
  ('q-boki3-s1-3-c1', 'q-boki3-s1-3', '借方（左側）', FALSE),
  ('q-boki3-s1-3-c2', 'q-boki3-s1-3', '貸方（右側）', TRUE),
  ('q-boki3-s1-3-c3', 'q-boki3-s1-3', '負債の増加は記録しない', FALSE),
  ('q-boki3-s1-3-c4', 'q-boki3-s1-3', 'どちらでもよい', FALSE),
  ('q-boki3-s1-4-c1', 'q-boki3-s1-4', '現金の入出金を記録する帳簿', FALSE),
  ('q-boki3-s1-4-c2', 'q-boki3-s1-4', '勘定科目ごとに集計・転記した帳簿', TRUE),
  ('q-boki3-s1-4-c3', 'q-boki3-s1-4', '月次の損益を集計する帳簿', FALSE),
  ('q-boki3-s1-4-c4', 'q-boki3-s1-4', '売上だけを管理する帳簿', FALSE),
  ('q-boki3-s1-5-c1', 'q-boki3-s1-5', '資本金', FALSE),
  ('q-boki3-s1-5-c2', 'q-boki3-s1-5', '売掛金', FALSE),
  ('q-boki3-s1-5-c3', 'q-boki3-s1-5', '給料', TRUE),
  ('q-boki3-s1-5-c4', 'q-boki3-s1-5', '受取手数料', FALSE)
ON CONFLICT (id) DO NOTHING;

-- ---- 簿記3級 第2章 セクション復習テスト ----
INSERT INTO quizzes (id, lesson_id, section_id, is_mock_exam, time_limit_minutes) VALUES
  ('quiz-boki3-s2', NULL, 'boki3-s2', FALSE, 15)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_questions (id, quiz_id, question_text, explanation, "order") VALUES
  ('q-boki3-s2-1', 'quiz-boki3-s2', '減価償却の目的として正しいものはどれですか？',
   '減価償却は固定資産の取得原価を耐用年数にわたって費用として配分するための手続きです。', 1),
  ('q-boki3-s2-2', 'quiz-boki3-s2', '「買掛金」が発生する取引はどれですか？',
   '買掛金は商品を掛けで仕入れたときに発生する負債です。', 2),
  ('q-boki3-s2-3', 'quiz-boki3-s2', '定額法による減価償却費の計算式はどれですか？',
   '定額法は (取得原価 − 残存価額) ÷ 耐用年数 で毎年同額を計上します。', 3),
  ('q-boki3-s2-4', 'quiz-boki3-s2', '小口現金の説明として正しいものはどれですか？',
   '小口現金は日常的な少額支払いのために手元に置く現金で、用途ごとに経費として記録します。', 4),
  ('q-boki3-s2-5', 'quiz-boki3-s2', '売掛金が回収されたとき、仕訳はどうなりますか？',
   '売掛金（資産）が減少するため貸方に記録し、現金（資産）が増加するため借方に記録します。', 5)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_choices (id, question_id, choice_text, is_correct) VALUES
  ('q-boki3-s2-1-c1', 'q-boki3-s2-1', '固定資産の価値を増加させるため', FALSE),
  ('q-boki3-s2-1-c2', 'q-boki3-s2-1', '取得原価を耐用年数にわたって費用として配分するため', TRUE),
  ('q-boki3-s2-1-c3', 'q-boki3-s2-1', '税金を減らすためだけの処理', FALSE),
  ('q-boki3-s2-1-c4', 'q-boki3-s2-1', '現金の流出を記録するため', FALSE),
  ('q-boki3-s2-2-c1', 'q-boki3-s2-2', '商品を現金で仕入れたとき', FALSE),
  ('q-boki3-s2-2-c2', 'q-boki3-s2-2', '商品を掛けで仕入れたとき', TRUE),
  ('q-boki3-s2-2-c3', 'q-boki3-s2-2', '商品を掛けで販売したとき', FALSE),
  ('q-boki3-s2-2-c4', 'q-boki3-s2-2', '給料を支払ったとき', FALSE),
  ('q-boki3-s2-3-c1', 'q-boki3-s2-3', '取得原価 × 耐用年数', FALSE),
  ('q-boki3-s2-3-c2', 'q-boki3-s2-3', '取得原価 ÷ 耐用年数', FALSE),
  ('q-boki3-s2-3-c3', 'q-boki3-s2-3', '(取得原価 − 残存価額) ÷ 耐用年数', TRUE),
  ('q-boki3-s2-3-c4', 'q-boki3-s2-3', '取得原価 × 定率', FALSE),
  ('q-boki3-s2-4-c1', 'q-boki3-s2-4', '銀行口座に預けた大口資金', FALSE),
  ('q-boki3-s2-4-c2', 'q-boki3-s2-4', '日常的な少額支払いのための手元現金', TRUE),
  ('q-boki3-s2-4-c3', 'q-boki3-s2-4', '売掛金の担保として保持する現金', FALSE),
  ('q-boki3-s2-4-c4', 'q-boki3-s2-4', '給料支払い専用の現金', FALSE),
  ('q-boki3-s2-5-c1', 'q-boki3-s2-5', '借方：売掛金 / 貸方：現金', FALSE),
  ('q-boki3-s2-5-c2', 'q-boki3-s2-5', '借方：現金 / 貸方：売掛金', TRUE),
  ('q-boki3-s2-5-c3', 'q-boki3-s2-5', '借方：売上 / 貸方：現金', FALSE),
  ('q-boki3-s2-5-c4', 'q-boki3-s2-5', '借方：買掛金 / 貸方：現金', FALSE)
ON CONFLICT (id) DO NOTHING;

-- ---- 簿記3級 第3章 セクション復習テスト ----
INSERT INTO quizzes (id, lesson_id, section_id, is_mock_exam, time_limit_minutes) VALUES
  ('quiz-boki3-s3', NULL, 'boki3-s3', FALSE, 15)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_questions (id, quiz_id, question_text, explanation, "order") VALUES
  ('q-boki3-s3-1', 'quiz-boki3-s3', '試算表の主な役割はどれですか？',
   '試算表は仕訳・転記の正確性を確認するための一覧表で、借方合計と貸方合計が一致します。', 1),
  ('q-boki3-s3-2', 'quiz-boki3-s3', '損益計算書（P/L）が示すものとして正しいものはどれですか？',
   '損益計算書は一定期間の収益・費用・利益（損失）を示す財務諸表です。', 2),
  ('q-boki3-s3-3', 'quiz-boki3-s3', '貸借対照表（B/S）の「純資産の部」に含まれるものはどれですか？',
   '資本金・繰越利益剰余金などが純資産の部に含まれます。', 3),
  ('q-boki3-s3-4', 'quiz-boki3-s3', '売上総利益の計算式はどれですか？',
   '売上総利益 = 売上高 − 売上原価 です。', 4),
  ('q-boki3-s3-5', 'quiz-boki3-s3', '決算整理仕訳の目的はどれですか？',
   '決算整理仕訳は期末に資産・負債・収益・費用を正確に計上するために行います（減価償却・棚卸など）。', 5)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_choices (id, question_id, choice_text, is_correct) VALUES
  ('q-boki3-s3-1-c1', 'q-boki3-s3-1', '将来の収益を予測するための表', FALSE),
  ('q-boki3-s3-1-c2', 'q-boki3-s3-1', '仕訳・転記の正確性を確認するための一覧表', TRUE),
  ('q-boki3-s3-1-c3', 'q-boki3-s3-1', '顧客への請求書', FALSE),
  ('q-boki3-s3-1-c4', 'q-boki3-s3-1', '税務申告書類', FALSE),
  ('q-boki3-s3-2-c1', 'q-boki3-s3-2', '特定時点の資産・負債・純資産を示す', FALSE),
  ('q-boki3-s3-2-c2', 'q-boki3-s3-2', '一定期間の収益・費用・利益を示す', TRUE),
  ('q-boki3-s3-2-c3', 'q-boki3-s3-2', '将来のキャッシュフローを予測する', FALSE),
  ('q-boki3-s3-2-c4', 'q-boki3-s3-2', '従業員の給与明細', FALSE),
  ('q-boki3-s3-3-c1', 'q-boki3-s3-3', '借入金', FALSE),
  ('q-boki3-s3-3-c2', 'q-boki3-s3-3', '売掛金', FALSE),
  ('q-boki3-s3-3-c3', 'q-boki3-s3-3', '資本金', TRUE),
  ('q-boki3-s3-3-c4', 'q-boki3-s3-3', '買掛金', FALSE),
  ('q-boki3-s3-4-c1', 'q-boki3-s3-4', '売上高 − 販売費及び一般管理費', FALSE),
  ('q-boki3-s3-4-c2', 'q-boki3-s3-4', '売上高 − 売上原価', TRUE),
  ('q-boki3-s3-4-c3', 'q-boki3-s3-4', '収益合計 − 費用合計', FALSE),
  ('q-boki3-s3-4-c4', 'q-boki3-s3-4', '営業利益 + 営業外収益', FALSE),
  ('q-boki3-s3-5-c1', 'q-boki3-s3-5', '期首に行う仕訳', FALSE),
  ('q-boki3-s3-5-c2', 'q-boki3-s3-5', '期末に資産・負債・収益・費用を正確に計上するための仕訳', TRUE),
  ('q-boki3-s3-5-c3', 'q-boki3-s3-5', '売上の記録をやり直す手続き', FALSE),
  ('q-boki3-s3-5-c4', 'q-boki3-s3-5', '税務署への申告のための仕訳', FALSE)
ON CONFLICT (id) DO NOTHING;

-- ---- 資産運用検定3級 第1章 セクション復習テスト ----
INSERT INTO quizzes (id, lesson_id, section_id, is_mock_exam, time_limit_minutes) VALUES
  ('quiz-asset3-s1', NULL, 'asset3-s1', FALSE, 15)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_questions (id, quiz_id, question_text, explanation, "order") VALUES
  ('q-asset3-s1-1', 'quiz-asset3-s1', '株式投資で得られる利益の種類として正しい組み合わせはどれですか？',
   '株式投資ではキャピタルゲイン（売却益）とインカムゲイン（配当金）の2種類の利益があります。', 1),
  ('q-asset3-s1-2', 'quiz-asset3-s1', '成行注文の説明として正しいものはどれですか？',
   '成行注文は価格を指定せず、現在の市場価格で即時取引を行う注文方法です。', 2),
  ('q-asset3-s1-3', 'quiz-asset3-s1', '「リスク」の金融的な定義として正しいものはどれですか？',
   '金融においてリスクとは損失だけでなく、収益のブレ（不確実性）全体を指します。', 3),
  ('q-asset3-s1-4', 'quiz-asset3-s1', 'PER（株価収益率）の計算式はどれですか？',
   'PER = 株価 ÷ 1株あたり純利益（EPS）です。割安・割高の目安に使われます。', 4),
  ('q-asset3-s1-5', 'quiz-asset3-s1', '証券総合口座を開設できる機関はどれですか？',
   '証券総合口座は証券会社（銀行や郵便局では開設できません）で開設します。', 5)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_choices (id, question_id, choice_text, is_correct) VALUES
  ('q-asset3-s1-1-c1', 'q-asset3-s1-1', 'キャピタルゲインのみ', FALSE),
  ('q-asset3-s1-1-c2', 'q-asset3-s1-1', 'インカムゲインのみ', FALSE),
  ('q-asset3-s1-1-c3', 'q-asset3-s1-1', 'キャピタルゲインとインカムゲイン', TRUE),
  ('q-asset3-s1-1-c4', 'q-asset3-s1-1', '利子収入のみ', FALSE),
  ('q-asset3-s1-2-c1', 'q-asset3-s1-2', '希望の価格を指定して注文する方法', FALSE),
  ('q-asset3-s1-2-c2', 'q-asset3-s1-2', '価格を指定せず市場価格で即時取引する方法', TRUE),
  ('q-asset3-s1-2-c3', 'q-asset3-s1-2', '一定期間後に取引する方法', FALSE),
  ('q-asset3-s1-2-c4', 'q-asset3-s1-2', '複数の証券会社に同時注文する方法', FALSE),
  ('q-asset3-s1-3-c1', 'q-asset3-s1-3', '損失が出る可能性のみ', FALSE),
  ('q-asset3-s1-3-c2', 'q-asset3-s1-3', '収益のブレ（不確実性）全体', TRUE),
  ('q-asset3-s1-3-c3', 'q-asset3-s1-3', '投資した金額のすべて', FALSE),
  ('q-asset3-s1-3-c4', 'q-asset3-s1-3', '手数料の合計', FALSE),
  ('q-asset3-s1-4-c1', 'q-asset3-s1-4', '株価 ÷ 1株あたり純資産（BPS）', FALSE),
  ('q-asset3-s1-4-c2', 'q-asset3-s1-4', '株価 ÷ 1株あたり純利益（EPS）', TRUE),
  ('q-asset3-s1-4-c3', 'q-asset3-s1-4', '株価 × 発行済株式数', FALSE),
  ('q-asset3-s1-4-c4', 'q-asset3-s1-4', '利益 ÷ 売上高', FALSE),
  ('q-asset3-s1-5-c1', 'q-asset3-s1-5', '銀行', FALSE),
  ('q-asset3-s1-5-c2', 'q-asset3-s1-5', '郵便局（ゆうちょ銀行）', FALSE),
  ('q-asset3-s1-5-c3', 'q-asset3-s1-5', '証券会社', TRUE),
  ('q-asset3-s1-5-c4', 'q-asset3-s1-5', '信用金庫', FALSE)
ON CONFLICT (id) DO NOTHING;

-- ---- 資産運用検定3級 第2章 セクション復習テスト ----
INSERT INTO quizzes (id, lesson_id, section_id, is_mock_exam, time_limit_minutes) VALUES
  ('quiz-asset3-s2', NULL, 'asset3-s2', FALSE, 15)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_questions (id, quiz_id, question_text, explanation, "order") VALUES
  ('q-asset3-s2-1', 'quiz-asset3-s2', '債券の「クーポン（利息）」について正しいものはどれですか？',
   'クーポンは債券の額面に対して定められた利率で定期的に支払われる利息です。', 1),
  ('q-asset3-s2-2', 'quiz-asset3-s2', 'インデックスファンドの特徴として正しいものはどれですか？',
   'インデックスファンドは日経平均やS&P500などの市場指数に連動することを目指す投資信託です。', 2),
  ('q-asset3-s2-3', 'quiz-asset3-s2', 'ETF（上場投資信託）と通常の投資信託の違いとして正しいものはどれですか？',
   'ETFは証券取引所に上場されており、株式と同様にリアルタイムで売買が可能です。', 3),
  ('q-asset3-s2-4', 'quiz-asset3-s2', '投資信託の「信託報酬」とは何ですか？',
   '信託報酬は投資信託の運用・管理にかかる費用で、毎年一定割合が差し引かれます。', 4),
  ('q-asset3-s2-5', 'quiz-asset3-s2', '格付けが高い債券の特徴として正しいものはどれですか？',
   '格付けが高い（AAA等）債券はデフォルトリスクが低い代わりに利回りも低くなります。', 5)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_choices (id, question_id, choice_text, is_correct) VALUES
  ('q-asset3-s2-1-c1', 'q-asset3-s2-1', '債券売却時に得られる差益', FALSE),
  ('q-asset3-s2-1-c2', 'q-asset3-s2-1', '額面に対して定期的に支払われる利息', TRUE),
  ('q-asset3-s2-1-c3', 'q-asset3-s2-1', '投資信託の信託報酬', FALSE),
  ('q-asset3-s2-1-c4', 'q-asset3-s2-1', '株式の配当金', FALSE),
  ('q-asset3-s2-2-c1', 'q-asset3-s2-2', '市場平均を上回ることを目指す投資信託', FALSE),
  ('q-asset3-s2-2-c2', 'q-asset3-s2-2', '市場指数に連動することを目指す投資信託', TRUE),
  ('q-asset3-s2-2-c3', 'q-asset3-s2-2', '元本保証のある投資信託', FALSE),
  ('q-asset3-s2-2-c4', 'q-asset3-s2-2', '特定の1銘柄にだけ投資する商品', FALSE),
  ('q-asset3-s2-3-c1', 'q-asset3-s2-3', 'ETFの方が手数料が高い', FALSE),
  ('q-asset3-s2-3-c2', 'q-asset3-s2-3', 'ETFは証券取引所に上場しリアルタイムで売買できる', TRUE),
  ('q-asset3-s2-3-c3', 'q-asset3-s2-3', 'ETFは元本保証がある', FALSE),
  ('q-asset3-s2-3-c4', 'q-asset3-s2-3', 'ETFは1日1回だけ取引できる', FALSE),
  ('q-asset3-s2-4-c1', 'q-asset3-s2-4', '投資信託の購入時に一度だけ払う手数料', FALSE),
  ('q-asset3-s2-4-c2', 'q-asset3-s2-4', '運用・管理にかかる費用で毎年差し引かれるもの', TRUE),
  ('q-asset3-s2-4-c3', 'q-asset3-s2-4', '解約時に支払うペナルティ', FALSE),
  ('q-asset3-s2-4-c4', 'q-asset3-s2-4', '分配金の受取手数料', FALSE),
  ('q-asset3-s2-5-c1', 'q-asset3-s2-5', '利回りが高くリスクも高い', FALSE),
  ('q-asset3-s2-5-c2', 'q-asset3-s2-5', 'デフォルトリスクが低く利回りも低い', TRUE),
  ('q-asset3-s2-5-c3', 'q-asset3-s2-5', 'デフォルトリスクが高く利回りも低い', FALSE),
  ('q-asset3-s2-5-c4', 'q-asset3-s2-5', '格付けは利回りと無関係', FALSE)
ON CONFLICT (id) DO NOTHING;

-- ---- 資産運用検定3級 第3章 セクション復習テスト ----
INSERT INTO quizzes (id, lesson_id, section_id, is_mock_exam, time_limit_minutes) VALUES
  ('quiz-asset3-s3', NULL, 'asset3-s3', FALSE, 15)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_questions (id, quiz_id, question_text, explanation, "order") VALUES
  ('q-asset3-s3-1', 'quiz-asset3-s3', '「相関係数 -1」の2つの資産を組み合わせた場合のリスクはどうなりますか？',
   '相関係数 -1 の場合、一方が上がれば他方が下がるため、組み合わせることでリスクを最大限に低減できます。', 1),
  ('q-asset3-s3-2', 'quiz-asset3-s3', '長期投資の最大のメリットとして正しいものはどれですか？',
   '長期投資では複利効果が働き、時間の経過とともに資産が指数関数的に増加します。', 2),
  ('q-asset3-s3-3', 'quiz-asset3-s3', 'ドルコスト平均法の説明として正しいものはどれですか？',
   'ドルコスト平均法は毎月一定金額で購入し、高い時は少なく・安い時は多く買うことで平均購入単価を下げる手法です。', 3),
  ('q-asset3-s3-4', 'quiz-asset3-s3', 'ポートフォリオの「リバランス」とは何ですか？',
   'リバランスは資産配分が目標から乖離したときに、売買して元の配分比率に戻すことです。', 4),
  ('q-asset3-s3-5', 'quiz-asset3-s3', '「システマティックリスク」の説明として正しいものはどれですか？',
   'システマティックリスクは市場全体に影響する分散投資で消せないリスク（景気・金利・為替など）です。', 5)
ON CONFLICT (id) DO NOTHING;

INSERT INTO quiz_choices (id, question_id, choice_text, is_correct) VALUES
  ('q-asset3-s3-1-c1', 'q-asset3-s3-1', 'リスクが増大する', FALSE),
  ('q-asset3-s3-1-c2', 'q-asset3-s3-1', 'リスクを最大限に低減できる', TRUE),
  ('q-asset3-s3-1-c3', 'q-asset3-s3-1', 'リスクは変わらない', FALSE),
  ('q-asset3-s3-1-c4', 'q-asset3-s3-1', 'リターンも半減する', FALSE),
  ('q-asset3-s3-2-c1', 'q-asset3-s3-2', '短期間で大きな利益が得られる', FALSE),
  ('q-asset3-s3-2-c2', 'q-asset3-s3-2', '複利効果で資産が指数関数的に増加する', TRUE),
  ('q-asset3-s3-2-c3', 'q-asset3-s3-2', '元本が保証される', FALSE),
  ('q-asset3-s3-2-c4', 'q-asset3-s3-2', '手数料が無料になる', FALSE),
  ('q-asset3-s3-3-c1', 'q-asset3-s3-3', '毎月一定数量を購入する手法', FALSE),
  ('q-asset3-s3-3-c2', 'q-asset3-s3-3', '株価が安いときだけ大量購入する手法', FALSE),
  ('q-asset3-s3-3-c3', 'q-asset3-s3-3', '毎月一定金額で購入し平均購入単価を下げる手法', TRUE),
  ('q-asset3-s3-3-c4', 'q-asset3-s3-3', '株価を予測して最適なタイミングで購入する手法', FALSE),
  ('q-asset3-s3-4-c1', 'q-asset3-s3-4', '全資産を売却して再投資すること', FALSE),
  ('q-asset3-s3-4-c2', 'q-asset3-s3-4', '資産配分が目標から乖離した際に元の比率に戻すこと', TRUE),
  ('q-asset3-s3-4-c3', 'q-asset3-s3-4', '新しい投資先に乗り換えること', FALSE),
  ('q-asset3-s3-4-c4', 'q-asset3-s3-4', '損失を確定して節税する手法', FALSE),
  ('q-asset3-s3-5-c1', 'q-asset3-s3-5', '個別企業固有のリスクで分散で消せるもの', FALSE),
  ('q-asset3-s3-5-c2', 'q-asset3-s3-5', '市場全体に影響する分散投資でも消せないリスク', TRUE),
  ('q-asset3-s3-5-c3', 'q-asset3-s3-5', '短期売買にのみ関係するリスク', FALSE),
  ('q-asset3-s3-5-c4', 'q-asset3-s3-5', '信用取引のリスクのみ', FALSE)
ON CONFLICT (id) DO NOTHING;
