# 実装計画書：add-spec03 対応

作成日: 2026-03-18

## 対象仕様（add-spec03.md）

| # | 種別 | 内容 |
|---|------|------|
| F-1 | 追加仕様 | 個人メモ機能 |
| F-2 | 追加仕様 | 生成AIコンテンツによる注意書き |
| F-3 | 追加仕様 | キャッシュフロー表の例 |
| F-4 | 追加仕様 | クイズコンテンツの充実 |
| F-5 | 追加仕様 | トップページに今日の用語を表示 |
| F-6 | 追加仕様 | クイズ結果画面でカリキュラム完了ができるようにする |

---

## 現状分析

### バックエンド
- `domain/models.go`：`User`, `Course`, `Section`, `Lesson`, `GlossaryTerm`, `Quiz`, `QuizQuestion`, `UserQuizResult`, `UserLessonProgress` などが定義済み
- `repository/repository.go`：`UserRepository`, `CourseRepository`, `SectionRepository`, `LessonRepository`, `ProgressRepository`, `GlossaryRepository`, `QuizRepository` が定義済み
- 用語辞典の全件取得 (`GlossaryRepository.List`) は実装済み。日次ランダム取得エンドポイントは未実装
- メモ機能に関わるテーブル・ドメイン・リポジトリ・ユースケース・ハンドラーは未実装

### フロントエンド
- `app/page.tsx`（トップページ）：コース一覧のみを表示。「今日の用語」セクションは未実装
- `app/lessons/[id]/page.tsx`：レッスン本文表示・完了ボタンあり。メモ入力UIは未実装
- `app/profile/page.tsx`：プロフィール編集フォームのみ。メモ一覧表示は未実装
- `app/quiz/[id]/page.tsx`：クイズ結果画面にはレッスン完了ボタンが未実装。`detail.quiz.lessonId` はすでにデータとして取得できる
- AIコンテンツ注意書きは存在しない
- `lib/api.ts`：メモ用・今日の用語用の API 関数は未実装

---

## 実装計画

---

### F-1: 個人メモ機能

#### F-1-1: バックエンド

**マイグレーション（新規作成）**
- `backend/migrations/0005_user_notes.up.sql`
- `backend/migrations/0005_user_notes.down.sql`

```sql
-- 0005_user_notes.up.sql
CREATE TABLE IF NOT EXISTS user_notes (
  id         TEXT PRIMARY KEY,
  user_id    TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  lesson_id  TEXT NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
  content    TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(user_id, lesson_id)
);
```

**ドメインモデル（追記）**
- `backend/internal/domain/models.go` に以下を追加

```go
type UserNote struct {
    ID        string    `json:"id"`
    UserID    string    `json:"userId"`
    LessonID  string    `json:"lessonId"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
```

**リポジトリインターフェース（追記）**
- `backend/internal/repository/repository.go` に以下を追加

```go
type NoteRepository interface {
    FindByUserAndLesson(userID, lessonID string) (domain.UserNote, bool, error)
    Upsert(note domain.UserNote) (domain.UserNote, error)
    ListByUserID(userID string) ([]domain.UserNote, error)
}
```

**リポジトリ実装（新規作成）**
- `backend/internal/repository/note_postgres.go`：`FindByUserAndLesson`, `Upsert`（INSERT ON CONFLICT DO UPDATE）, `ListByUserID` を実装
- `backend/internal/repository/note_memory.go`：テスト用インメモリ実装

**ユースケース（新規作成）**
- `backend/internal/usecase/note.go`
  - `GetNote(userID, lessonID string) (domain.UserNote, bool, error)`
  - `SaveNote(userID, lessonID, content string) (domain.UserNote, error)`
  - `ListNotes(userID string) ([]domain.UserNote, error)`

**ハンドラー（新規作成）**
- `backend/internal/handler/note.go`
  - JWT 認証必須
  - `GetByLesson(w, r)`：GET `/api/v1/lessons/{lessonId}/note`
  - `Save(w, r)`：PUT `/api/v1/lessons/{lessonId}/note`（リクエストボディ `{"content": "..."}`）
  - `ListAll(w, r)`：GET `/api/v1/users/me/notes`

**サーバー設定（追記）**
- `backend/cmd/server/main.go` に `NoteRepository`, `NoteUseCase`, `NoteHandler` の初期化とルーティングを追加

**API エンドポイント**

| メソッド | パス | 認証 | 説明 |
|----------|------|------|------|
| GET  | `/api/v1/lessons/{lessonId}/note` | 必要 | 対象レッスンのメモ取得 |
| PUT  | `/api/v1/lessons/{lessonId}/note` | 必要 | 対象レッスンのメモ保存（作成・更新） |
| GET  | `/api/v1/users/me/notes` | 必要 | 自分の全メモ一覧 |

**OpenAPI 定義（追記）**
- `backend/openapi.yaml` に上記エンドポイントと `UserNote` スキーマを追加

---

#### F-1-2: フロントエンド

**型定義（新規作成）**
- `frontend/src/types/note.ts`

```ts
export type UserNote = {
  id: string;
  userId: string;
  lessonId: string;
  content: string;
  createdAt: string;
  updatedAt: string;
};

export type NotesResponse = {
  notes: UserNote[];
};
```

**API 関数（追記）**
- `frontend/src/lib/api.ts` に以下を追加

```ts
export async function fetchLessonNote(lessonId: string): Promise<UserNote | null>
export async function saveLessonNote(lessonId: string, content: string): Promise<UserNote>
export async function fetchMyNotes(): Promise<UserNote[]>
```

**レッスンページ（追記）**
- `frontend/src/app/lessons/[id]/page.tsx` に「メモ」エリアを追加
  - レッスンコンテンツの下部に `<textarea>` を配置
  - ログイン済みユーザーのみ表示
  - 初回表示時に `fetchLessonNote(lessonId)` で既存メモを取得・表示
  - 「保存」ボタン押下で `saveLessonNote(lessonId, content)` を呼び出し
  - 保存成功時にフィードバックメッセージを表示

**プロフィールページ（追記）**
- `frontend/src/app/profile/page.tsx` に「メモ一覧」セクションを追加
  - `fetchMyNotes()` でメモ一覧を取得して表示
  - メモごとにレッスンIDとメモ内容を表示
  - 対応レッスンへのリンクを表示

---

### F-2: 生成AIコンテンツによる注意書き

**対象ファイル**
- `frontend/src/app/layout.tsx`（または汎用バナーコンポーネント）

**実装方針**
- レイアウトのフッター部分（または Header 下の目立たない位置）に AI 生成コンテンツであることを示すバナーを追加
- バックエンド変更なし

**表示内容（例）**
```
⚠️ このサイトのコンテンツは生成AIによって作成されています。
情報の正確性については十分ご確認ください。
```

**実装詳細**
- `frontend/src/components/AIGeneratedNotice.tsx` を新規作成
- `layout.tsx` の `<div className="mx-auto max-w-5xl px-4 py-6">` の前後に配置
- スタイル：黄色系（`bg-amber-50 border border-amber-200 text-amber-800 text-xs`）のシンプルなバナー

---

### F-3: キャッシュフロー表の例

**対象ファイル**
- 新規マイグレーション `backend/migrations/0006_content_updates.up.sql`（または `0005` に統合）

**実装方針**
- `fp3-s1-l1`（ライフプランとは）がキャッシュフロー表に言及しているが、具体的な数値例がない
- **新規レッスンは作成せず**、既存の `fp3-s1-l1` の `content` を UPDATE して具体例を追記する
- マイグレーション SQL で `UPDATE lessons SET content = $$ ... $$ WHERE id = 'fp3-s1-l1'` を実行

**追記するコンテンツ（既存の「## キャッシュフロー表」セクションに以下を追記）**

```markdown
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
```

---

### F-4: クイズコンテンツの充実

**対象ファイル**
- 新規マイグレーション（`0006_content_updates.up.sql` または `0005` に統合）

**実装方針**
- 現状はシードデータにクイズが存在するか確認する必要あり（`0004_seed_data.up.sql` に一部含まれる場合は把握した上で追加）
- 各章に最低1つのクイズを確保する
- レッスン単体に紐づく確認クイズを優先的に追加

**追加するクイズのスコープ**

| コース | セクション | 追加するクイズ種類 |
|--------|----------|-----------------|
| FP3級 | 第1章・第2章・第3章・第4章・第5章 | 各章に5問以上のセクション復習テスト |
| 簿記3級 | 第1章・第2章・第3章 | 各章に5問以上のセクション復習テスト |
| 資産運用検定3級 | 第1章・第2章・第3章 | 各章に5問以上のセクション復習テスト |

**実装方法**
- マイグレーション SQL でクイズデータを INSERT
- `quizzes`→`quiz_questions`→`quiz_choices` の順で INSERT

---

### F-5: トップページに今日の用語を表示

#### F-5-1: バックエンド

**対象ファイル**
- `backend/internal/usecase/glossary.go`（メソッド追加）
- `backend/internal/handler/glossary.go`（ハンドラー追加）
- `backend/internal/repository/glossary_postgres.go`（全件カウントメソッド追加 or 既存の List 活用）
- `backend/cmd/server/main.go`（ルーティング追加）
- `backend/openapi.yaml`（エンドポイント追加）

**実装方針**
- `GlossaryUseCase` に `GetDailyTerm() (domain.GlossaryTerm, error)` メソッドを追加
- 実装：全用語を取得し、`日付のシリアル番号 mod 用語数` でインデックスを決定（日ごとに同じ用語が返り、ランダム性も保たれる）
- `GlossaryHandler` に `GetDaily` メソッドを追加
- ルーティング：`GET /api/v1/glossary/daily`

**API エンドポイント**

| メソッド | パス | 認証 | 説明 |
|----------|------|------|------|
| GET | `/api/v1/glossary/daily` | 不要 | 今日の用語（日付ベースで決定論的に選択） |

**OpenAPI 定義（追記）**
- `backend/openapi.yaml` に追加

---

#### F-5-2: フロントエンド

**型定義**
- 既存の `frontend/src/types/glossary.ts` の `GlossaryTerm` 型を再利用

**API 関数（追記）**
- `frontend/src/lib/api.ts` に以下を追加

```ts
export async function fetchDailyTerm(): Promise<GlossaryTerm | null>
```

**トップページ（追記）**
- `frontend/src/app/page.tsx` に「今日の用語」セクションを追加
  - `fetchDailyTerm()` でデータを取得
  - 取得できた場合：用語名・読み方・定義を表示（カード形式）
  - 取得できなかった場合：セクション自体を非表示
  - 用語辞典ページへのリンクを表示
  - スタイル：ピンク系バッジで「今日の用語」ラベルを表示

---

### F-6: クイズ結果画面でカリキュラム完了

**対象ファイル**
- `frontend/src/app/quiz/[id]/page.tsx`

**実装方針**
- クイズ結果画面（`submitted && finalScore` の条件ブランチ）に「レッスンを完了する」ボタンを追加
- `detail.quiz.lessonId` が存在する場合のみボタンを表示
- `completeLesson(lessonId)` を呼び出し、完了後はボタンを完了済み表示に切り替え
- バックエンド変更なし

**実装詳細**

```tsx
// 結果画面に追加するステート
const [lessonCompleted, setLessonCompleted] = useState(false);
const [completingLesson, setCompletingLesson] = useState(false);

// 完了ハンドラー
const handleCompleteLesson = async () => {
  if (!detail?.quiz.lessonId) return;
  setCompletingLesson(true);
  try {
    await completeLesson(detail.quiz.lessonId);
    setLessonCompleted(true);
  } finally {
    setCompletingLesson(false);
  }
};
```

- 完了ボタンは「コース一覧へ戻る」リンクの上に配置
- 完了済み状態では「✓ レッスン完了済み」テキストを表示

---

## 実装順序と優先度

| 優先度 | # | 機能 | 理由 |
|--------|---|------|------|
| 高 | F-6 | クイズ結果からカリキュラム完了 | フロントエンドのみ・影響範囲が小さく即効性あり |
| 高 | F-2 | AI注意書き | フロントエンドのみ・信頼性確保に必要 |
| 高 | F-5 | 今日の用語（トップページ） | バックエンド+フロントエンド・実装コスト小 |
| 中 | F-3 | キャッシュフロー表の例 | マイグレーション追加のみ |
| 中 | F-4 | クイズコンテンツの充実 | マイグレーション追加のみ（量が多い） |
| 低 | F-1 | 個人メモ機能 | フルスタック実装で最も工数が大きい |

---

## ファイル変更一覧

### 新規作成ファイル

| ファイル | 内容 |
|---------|------|
| `backend/migrations/0005_user_notes.up.sql` | user_notes テーブル作成 |
| `backend/migrations/0005_user_notes.down.sql` | user_notes テーブル削除 |
| `backend/migrations/0006_content_updates.up.sql` | fp3-s1-l1 へのキャッシュフロー表具体例追記・クイズ追加 |
| `backend/migrations/0006_content_updates.down.sql` | 上記ロールバック |
| `backend/internal/repository/note_postgres.go` | NoteRepository Postgres 実装 |
| `backend/internal/repository/note_memory.go` | NoteRepository インメモリ実装 |
| `backend/internal/usecase/note.go` | NoteUseCase |
| `backend/internal/handler/note.go` | NoteHandler |
| `frontend/src/types/note.ts` | UserNote 型定義 |
| `frontend/src/components/AIGeneratedNotice.tsx` | AI注意書きコンポーネント |

### 既存ファイルへの追記・修正

| ファイル | 変更内容 |
|---------|---------|
| `backend/internal/domain/models.go` | UserNote モデル追加 |
| `backend/internal/repository/repository.go` | NoteRepository インターフェース追加 |
| `backend/internal/usecase/glossary.go` | GetDailyTerm メソッド追加 |
| `backend/internal/handler/glossary.go` | GetDaily ハンドラー追加 |
| `backend/cmd/server/main.go` | Note/DailyTerm の初期化・ルーティング追加 |
| `backend/openapi.yaml` | メモ・今日の用語エンドポイント・スキーマ追加 |
| `frontend/src/lib/api.ts` | fetchLessonNote, saveLessonNote, fetchMyNotes, fetchDailyTerm 追加 |
| `frontend/src/app/page.tsx` | 今日の用語セクション追加 |
| `frontend/src/app/layout.tsx` | AIGeneratedNotice コンポーネント追加 |
| `frontend/src/app/lessons/[id]/page.tsx` | メモ入力エリア追加 |
| `frontend/src/app/profile/page.tsx` | メモ一覧セクション追加 |
| `frontend/src/app/quiz/[id]/page.tsx` | 結果画面にレッスン完了ボタン追加 |
