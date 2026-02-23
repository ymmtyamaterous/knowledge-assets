# 実装計画書：add-spec01 対応

作成日: 2026-02-23

## 対象仕様（add-spec01.md）

| # | 種別 | 内容 |
|---|------|------|
| B-1 | バグ | 表や数字の箇条書き表記が崩れている |
| F-1 | 追加仕様 | 一覧でレッスンが完了しているか確認できるようにしたい |
| F-2 | 追加仕様 | クイズ機能 |
| F-3 | 追加仕様 | 用語辞典画面のUI改善 |
| F-4 | 追加仕様 | 用語にタグを設定できるようにしたい |
| F-5 | 追加仕様 | 学習進捗確認ができるようにしたい |

---

## 現状分析

### フロントエンド
- `frontend/src/app/lessons/[id]/page.tsx` に独自の簡易 Markdown → HTML 変換関数 `MarkdownContent` が実装されている
  - テーブル: 各行を `<tr>` に変換しているが `<table>` タグで囲まれていない、区切り行（`|---|`）も変換対象になっている
  - 番号付きリスト（`1. 2. 3.`）: 変換ロジックが未実装
- コース詳細ページ（`courses/[id]/page.tsx`）: レッスン一覧に完了状態の表示なし（Server Component のため進捗取得不可）
- 用語辞典ページ（`glossary/page.tsx`）: シンプルなカード一覧のみ、フィルタ・検索・50音インデックスなし
- 学習進捗の専用ページは未実装

### バックエンド
- `/api/v1/users/me/progress` エンドポイント（レッスン単位の進捗）は実装済み
- クイズ関連のDBスキーマ（`0002_quizzes_glossary_badges.up.sql`）は定義済み
  - `quizzes`, `quiz_questions`, `quiz_choices`, `user_quiz_results` テーブルが存在
  - しかしドメインモデル・リポジトリ・ユースケース・ハンドラーは未実装
- 用語タグ機能の DB スキーマ・API は未実装

---

## 実装計画

### B-1: Markdown表示バグ修正（フロントエンド）

**対象ファイル**
- `frontend/src/app/lessons/[id]/page.tsx`

**問題点と対応**

| 問題 | 原因 | 対応 |
|------|------|------|
| テーブルが崩れる | 各行を `<tr>` に変換しているが、`<table>` タグで囲んでいない。区切り行（`\|---\|`）がデータ行として変換される | テーブルブロック全体を検出して `<table>` タグで囲む処理に修正。区切り行はスキップ |
| 番号付きリストが崩れる | `^\d+\. ` のパターンが未実装 | `^\d+\. (.+)$` を `<li class='ml-4 list-decimal'>$1</li>` に変換し `<ol>` で囲む処理を追加 |

---

### F-1: レッスン完了状態の一覧表示（フロントエンド）

**対象ファイル**
- `frontend/src/app/courses/[id]/page.tsx`（Server Component → Client Component 化が必要）

**実装方針**
1. `courses/[id]/page.tsx` を Client Component に変換し、ページ読み込み時に `fetchMyProgress()` を呼び出す
2. 取得した進捗一覧（`lessonId` の配列）を `SectionBlock` に props として渡す
3. 各レッスンリンクの横に完了アイコン（✓）を表示する
4. 未ログイン時は完了状態を表示しない（進捗取得をスキップ）

**表示仕様**
- 完了済み: レッスン名の右側にチェックマーク（✓）をピンク色で表示
- 未完了: 通常の表示

---

### F-2: クイズ機能（バックエンド + フロントエンド）

#### F-2-1: バックエンド実装

**対象ファイル（新規作成）**
- `backend/internal/domain/models.go`（Quiz 関連モデルを追加）
- `backend/internal/repository/quiz_memory.go`（クイズリポジトリ）
- `backend/internal/usecase/quiz.go`（クイズユースケース）
- `backend/internal/handler/quiz.go`（クイズハンドラー）

**追加する API エンドポイント**

| メソッド | パス | 認証 | 説明 |
|----------|------|------|------|
| GET | `/api/v1/lessons/{lessonId}/quiz` | 不要 | レッスンに紐付くクイズ取得 |
| GET | `/api/v1/quizzes/{id}` | 不要 | クイズ詳細（問題・選択肢）取得 |
| POST | `/api/v1/quizzes/{id}/submit` | 必要 | クイズ回答の送信・採点 |
| GET | `/api/v1/users/me/quiz-results` | 必要 | 自分のクイズ結果履歴 |

**追加するドメインモデル**
```go
type Quiz struct {
    ID               string
    LessonID         string  // nullable
    SectionID        string  // nullable
    IsMockExam       bool
    TimeLimitMinutes int     // nullable
    CreatedAt        time.Time
}

type QuizQuestion struct {
    ID           string
    QuizID       string
    QuestionText string
    Explanation  string
    Order        int
    Choices      []QuizChoice
}

type QuizChoice struct {
    ID         string
    QuestionID string
    ChoiceText  string
    IsCorrect   bool
}

type UserQuizResult struct {
    ID      string
    UserID  string
    QuizID  string
    Score   int
    Total   int
    TakenAt time.Time
}
```

**openapi.yaml への追加**
- 上記 API エンドポイントと関連スキーマ（`Quiz`, `QuizQuestion`, `QuizChoice`, `SubmitQuizRequest`, `QuizResult`）を追加

#### F-2-2: フロントエンド実装

**対象ファイル（新規作成）**
- `frontend/src/types/quiz.ts`（クイズ型定義）
- `frontend/src/app/quiz/[id]/page.tsx`（クイズ回答ページ）

**対象ファイル（更新）**
- `frontend/src/lib/api.ts`（クイズ関連 API 関数を追加）
- `frontend/src/app/lessons/[id]/page.tsx`（レッスン末尾にクイズへのリンクを追加）

**クイズページ仕様**
- 1問ずつ表示する形式（ページング）
- 回答選択後に正誤フィードバックと解説を表示
- 全問終了後にスコアサマリを表示
- 解答送信後は `POST /api/v1/quizzes/{id}/submit` を呼び出してサーバーに記録

---

### F-3: 用語辞典 UI 改善（フロントエンド）

**対象ファイル**
- `frontend/src/app/glossary/page.tsx`

**改善内容**

| 項目 | 現状 | 改善後 |
|------|------|--------|
| インデックス | なし | 50音・A〜Z のインデックスボタンを表示 |
| 検索 | なし | キーワード検索入力フォームを追加 |
| レイアウト | シンプルなカード縦並び | インデックス + カードの2カラム or インデックスに応じたセクション分け |
| 読み仮名 | 小さいテキストで表示 | より視認性の高いスタイルに変更 |

**実装方針**
- `"use client"` に変更し、検索・フィルタ状態を `useState` で管理
- 50音インデックスは `reading` フィールドの先頭文字でフィルタ
- 検索は `term`・`reading`・`definition` を対象としたクライアントサイドフィルタ

---

### F-4: 用語タグ機能（バックエンド + フロントエンド）

#### F-4-1: バックエンド実装

**マイグレーション（新規）**
- `backend/migrations/0003_glossary_tags.up.sql`
- `backend/migrations/0003_glossary_tags.down.sql`

```sql
-- 0003_glossary_tags.up.sql
CREATE TABLE IF NOT EXISTS glossary_tags (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS glossary_term_tags (
  term_id TEXT NOT NULL REFERENCES glossary_terms(id) ON DELETE CASCADE,
  tag_id  TEXT NOT NULL REFERENCES glossary_tags(id) ON DELETE CASCADE,
  PRIMARY KEY (term_id, tag_id)
);
```

**追加する API エンドポイント**

| メソッド | パス | 認証 | 説明 |
|----------|------|------|------|
| GET | `/api/v1/glossary/tags` | 不要 | タグ一覧取得 |
| GET | `/api/v1/glossary?tagId={id}` | 不要 | タグでフィルタした用語一覧 |

**追加するドメインモデル**
```go
type GlossaryTag struct {
    ID        string
    Name      string
    CreatedAt time.Time
}
```

**更新するドメインモデル**
```go
// GlossaryTerm に Tags フィールドを追加
type GlossaryTerm struct {
    ...
    Tags []GlossaryTag `json:"tags"`
}
```

**openapi.yaml への追加**
- `GlossaryTag` スキーマ
- タグ一覧エンドポイント
- `GlossaryTerm` の `tags` フィールド追加
- glossary 一覧の `tagId` クエリパラメータ追加

#### F-4-2: フロントエンド実装

**対象ファイル（更新）**
- `frontend/src/types/glossary.ts`（`GlossaryTag` 型、`GlossaryTerm.tags` フィールドを追加）
- `frontend/src/app/glossary/page.tsx`（タグフィルタ UI を追加）
- `frontend/src/lib/api.ts`（タグ API 関数を追加）

**タグ表示仕様**
- 用語カードにタグをバッジ形式で表示
- タグクリックでフィルタが適用される

---

### F-5: 学習進捗確認ページ（フロントエンド + バックエンド）

#### F-5-1: バックエンド追加 API

**追加する API エンドポイント**

| メソッド | パス | 認証 | 説明 |
|----------|------|------|------|
| GET | `/api/v1/users/me/course-progress` | 必要 | コース・セクション単位の進捗サマリ取得 |

**レスポンス仕様**
```json
{
  "courseProgress": [
    {
      "courseId": "...",
      "courseTitle": "...",
      "totalLessons": 30,
      "completedLessons": 12,
      "progressRate": 40,
      "sections": [
        {
          "sectionId": "...",
          "sectionTitle": "...",
          "totalLessons": 10,
          "completedLessons": 4,
          "progressRate": 40
        }
      ]
    }
  ]
}
```

#### F-5-2: フロントエンド実装

**対象ファイル（新規作成）**
- `frontend/src/app/progress/page.tsx`（学習進捗ページ）
- `frontend/src/types/progress.ts`（`CourseProgress` 型を追加）

**対象ファイル（更新）**
- `frontend/src/lib/api.ts`（`fetchCourseProgress` 関数を追加）
- `frontend/src/components/Header.tsx`（「進捗」ナビゲーションリンクを追加）

**進捗ページ仕様**
- コースごとにプログレスバーを表示（進捗率 %）
- コース内の各セクションの進捗もアコーディオン形式で確認可能
- 未ログイン時はログインを促すメッセージを表示

---

## テスト方針

| 対象 | テスト内容 |
|------|-----------|
| B-1 | Markdown変換関数の単体テスト（表・番号リスト・通常リストの変換を検証） |
| F-2 | クイズユースケースの単体テスト（採点ロジック、結果保存） |
| F-4 | タグ関連のリポジトリ・ユースケースの単体テスト |
| F-5 | 進捗計算ロジックの単体テスト（進捗率の計算） |

---

## 実装優先順位・順序

以下の順序で実装を進めることを推奨します。

1. **B-1**: Markdown表示バグ修正（影響範囲が小さく即効性が高い）
2. **F-1**: レッスン完了状態の一覧表示（APIは既存のものを使用するため比較的容易）
3. **F-3**: 用語辞典 UI 改善（フロントエンドのみの変更）
4. **F-4**: 用語タグ機能（DBマイグレーション + バックエンド + フロントエンド）
5. **F-5**: 学習進捗確認ページ（新規 API + フロントエンド実装）
6. **F-2**: クイズ機能（最も複雑なため最後に実装）

---

## 影響範囲まとめ

| # | バックエンド | フロントエンド | DB変更 | openapi.yaml |
|---|-------------|--------------|--------|--------------|
| B-1 | なし | `lessons/[id]/page.tsx` | なし | なし |
| F-1 | なし | `courses/[id]/page.tsx` | なし | なし |
| F-2 | handler/usecase/repository追加 | quiz/[id]/page.tsx追加、api.ts更新 | なし（既存） | クイズ関連追加 |
| F-3 | なし | `glossary/page.tsx` | なし | なし |
| F-4 | handler/usecase/repository追加 | glossary/page.tsx、api.ts更新 | 新規マイグレーション | タグ関連追加 |
| F-5 | handler/usecase追加 | progress/page.tsx追加、api.ts更新 | なし | 進捗サマリAPI追加 |
