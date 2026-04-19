# 実装計画書：add-spec04 対応

作成日: 2026-03-27

## 対象仕様（add-spec04.md）

### フェーズ1 実装対象（優先度高・フロントエンドのみ）

| # | 種別 | 内容 | 工数感 |
|---|------|------|--------|
| A-1 | 追加仕様 | 次のレッスンへのナビゲーション | 小 |
| A-2 | 追加仕様 | クイズ解答履歴ページ | 小 |
| A-3 | 追加仕様 | 用語詳細ページ | 小 |
| B-3 | 追加仕様 | メモのエクスポート機能 | 小 |

### フェーズ2 実装対象（優先度中・フルスタック）

| # | 種別 | 内容 | 工数感 |
|---|------|------|--------|
| A-4 | 追加仕様 | パスワード変更 | 中 |
| B-1 | 追加仕様 | 学習連続日数（ストリーク）表示 | 中 |

---

## 現状分析

### フロントエンド
- `app/lessons/[id]/page.tsx`：レッスン本文・完了ボタン・クイズリンク・メモドロワーは実装済み。「次のレッスン」ナビゲーションは未実装。
- `app/quiz/[id]/page.tsx`：クイズ解答・採点・結果表示は実装済み。過去の結果一覧ページは未実装。
- `app/glossary/page.tsx`：用語一覧・タグフィルタ・50音・検索は実装済み。各用語のリンクが `<div>` のため詳細ページに遷移できない。
- `app/profile/page.tsx`：メモ一覧は実装済み。エクスポートボタン・パスワード変更フォームは未実装。
- `lib/api.ts`：`fetchMyQuizResults()`, `fetchGlossaryTerm(id)` は実装済み。`changePassword`, `fetchMyStreak` は未実装。

### バックエンド
- `GET /api/v1/users/me/quiz-results`：実装済み。
- `GET /api/v1/glossary/{id}`：実装済み。
- パスワード変更エンドポイント（`PUT /api/v1/users/me/password`）：未実装。
- ストリークエンドポイント（`GET /api/v1/users/me/streak`）：未実装。
- `UserRepository.Update(user)` は実装済みで `password_hash` も更新対象に含まれる。
- `ProgressRepository.ListByUserID(userID)` は実装済み。ストリーク計算に使用可能。

---

## フェーズ1 実装計画

---

### A-1: 次のレッスンへのナビゲーション

**対象ファイル（変更）**
- `frontend/src/app/lessons/[id]/page.tsx`

**実装方針**
1. レッスン取得後、`lesson.sectionId` を使って `fetchLessons(lesson.sectionId)` を呼び出し、同セクション内の全レッスンを取得する
2. `lesson.order` を基準に「次のレッスン」（order が現在より1大きいもの）を特定する
3. レッスン本文カードの下部（完了ボタン・クイズリンクの後）に「次のレッスンへ →」リンクを追加する
4. 次のレッスンがない（セクション末尾）場合は「コース一覧へ戻る」リンクを表示する

**UI仕様**
```
[ ✓ 完了済み ] [ 完了を取り消す ]  [ このレッスンの確認クイズへ ]

─────────────────────────────────────────

                           [ 次のレッスンへ → ]
  （末尾の場合）           [ コース一覧へ戻る ]
```

**追加するステート**
```tsx
const [sectionLessons, setSectionLessons] = useState<Lesson[]>([]);
```

**副作用（useEffect）**
```tsx
useEffect(() => {
  if (!lesson) return;
  fetchLessons(lesson.sectionId)
    .then(setSectionLessons)
    .catch(() => setSectionLessons([]));
}, [lesson]);
```

**次のレッスン特定ロジック**
```tsx
const nextLesson = sectionLessons.find((l) => l.order === (lesson?.order ?? 0) + 1) ?? null;
```

---

### A-2: クイズ解答履歴ページ

**対象ファイル（変更）**
- `frontend/src/app/profile/page.tsx`

**実装方針**
- 独立ページではなく、プロフィールページ（`/profile`）のセクションとして追加する
- `fetchMyQuizResults()` で取得した結果を一覧表示する
- `UserQuizResult` にはレッスル名が含まれていないため、スコア・日時を中心に表示する
  - 例：「クイズ結果 #3」「4 / 5 問正解（80%）」「2026-03-26」のように表示
- スコアに応じてバッジ色を変える
  - 80%以上：ピンク（`bg-pink-100 text-pink-700`）
  - 60%以上：黄（`bg-yellow-100 text-yellow-700`）
  - 60%未満：グレー（`bg-slate-100 text-slate-600`）

**追加するステート**
```tsx
const [quizResults, setQuizResults] = useState<UserQuizResult[]>([]);
```

**副作用（useEffect）**
```tsx
useEffect(() => {
  if (!user) return;
  fetchMyQuizResults()
    .then(setQuizResults)
    .catch(() => {});
}, [user]);
```

---

### A-3: 用語詳細ページ

**対象ファイル（新規作成）**
- `frontend/src/app/glossary/[id]/page.tsx`

**対象ファイル（変更）**
- `frontend/src/app/glossary/page.tsx`

**実装方針（詳細ページ）**
1. `fetchGlossaryTerm(id)` は `api.ts` 実装済みのため、そのまま使用する
2. 表示内容：用語名・読み仮名・定義・タグバッジ一覧
3. 「← 用語辞典へ戻る」リンクを上部に配置する

**実装方針（一覧ページ変更）**
- 現在の各用語カードは `<div>` 要素。`<Link href={/glossary/${term.id}}>` で囲むことで詳細ページに遷移可能にする
- タグのクリックはリンク遷移ではなくフィルタ適用のままにするため、`<Link>` の中の `<button>` は `e.preventDefault()` で親リンクの遷移を止める

---

### B-3: メモのエクスポート機能

**対象ファイル（変更）**
- `frontend/src/app/profile/page.tsx`

**実装方針**
- 「メモ一覧」セクションの上部に「📥 Markdownでエクスポート」ボタンを追加する
- クライアントサイドでのみ処理（バックエンド変更なし）
- ボタンクリック時に以下の処理を実行する

```ts
const handleExportNotes = () => {
  const md = notes
    .map((note) => {
      const title = note.lessonTitle || `レッスン ${note.lessonId}`;
      return `# ${title}\n\n${note.content}\n\n---\n`;
    })
    .join("\n");

  const blob = new Blob([md], { type: "text/markdown" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = "my-notes.md";
  a.click();
  URL.revokeObjectURL(url);
};
```

---

## フェーズ2 実装計画

---

### A-4: パスワード変更

#### A-4-1: バックエンド

**対象ファイル（変更）**
- `backend/internal/usecase/auth.go`（メソッド追加）
- `backend/internal/handler/user.go`（ハンドラーメソッド追加 + `AuthUseCase` 依存追加）
- `backend/cmd/server/main.go`（ルーティング追加）
- `backend/openapi.yaml`（エンドポイント追加）

**AuthUseCase への追加メソッド**
```go
func (uc *AuthUseCase) ChangePassword(userID, currentPassword, newPassword string) error {
    user, ok, err := uc.users.FindByID(userID)
    if err != nil || !ok {
        return ErrInvalidCredentials
    }
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
        return ErrInvalidCredentials
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.PasswordHash = string(hash)
    _, err = uc.users.Update(user)
    return err
}
```

**UserHandler への変更**
- `authUC *usecase.AuthUseCase` フィールドを追加し、`NewUserHandler` の引数に `authUC` を加える
- `ChangePassword` メソッドを追加する

```go
type changePasswordRequest struct {
    CurrentPassword string `json:"currentPassword"`
    NewPassword     string `json:"newPassword"`
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
    userID, _ := r.Context().Value(userIDContextKey).(string)
    // ... バリデーション・usecase呼び出し
}
```

**ルーティング追加（`cmd/server/main.go`）**
```go
private.Put("/users/me/password", userHandler.ChangePassword)
```

**OpenAPI 定義（追記）**
```yaml
/api/v1/users/me/password:
  put:
    summary: パスワード変更
    security:
      - bearerAuth: []
    requestBody:
      required: true
      content:
        application/json:
          schema:
            type: object
            required: [currentPassword, newPassword]
            properties:
              currentPassword:
                type: string
              newPassword:
                type: string
                minLength: 8
    responses:
      '200':
        description: OK
      '400':
        description: パスワード不一致または形式エラー
      '401':
        description: Unauthorized
```

#### A-4-2: フロントエンド

**対象ファイル（変更）**
- `frontend/src/lib/api.ts`（関数追加）
- `frontend/src/app/profile/page.tsx`（フォーム追加）

**api.ts 追加関数**
```ts
export async function changePassword(currentPassword: string, newPassword: string): Promise<void> {
  await apiFetch<void>("/api/v1/users/me/password", {
    method: "PUT",
    body: JSON.stringify({ currentPassword, newPassword }),
  });
}
```

**プロフィールページ追加フォーム**
- 「パスワード変更」セクションを「情報の編集」フォームの下に追加する
- フォーム項目：現在のパスワード / 新しいパスワード / 確認用パスワード
- バリデーション：`新しいパスワード === 確認用パスワード`（クライアント側）・8文字以上
- 成功時：「パスワードを変更しました！」フィードバック表示
- エラー時：「現在のパスワードが正しくありません。」など適切なエラー表示

---

### B-1: 学習連続日数（ストリーク）表示

#### B-1-1: バックエンド

**対象ファイル（変更）**
- `backend/internal/domain/models.go`（`UserStreak` 型追加）
- `backend/internal/usecase/progress.go`（`GetStreak` メソッド追加）
- `backend/internal/handler/progress.go`（`GetMyStreak` ハンドラー追加）
- `backend/cmd/server/main.go`（ルーティング追加）
- `backend/openapi.yaml`（エンドポイント追加）

**domain.UserStreak 型**
```go
type UserStreak struct {
    CurrentStreak int    `json:"currentStreak"`
    LongestStreak int    `json:"longestStreak"`
    LastStudiedAt string `json:"lastStudiedAt"` // "2026-03-26" or ""
}
```

**ProgressUseCase の GetStreak メソッド**
```go
func (uc *ProgressUseCase) GetStreak(userID string) (domain.UserStreak, error) {
    list, err := uc.progress.ListByUserID(userID)
    if err != nil {
        return domain.UserStreak{}, err
    }
    // completedAt を日付単位でユニーク化・ソート
    // 今日から遡って連続日数（currentStreak）をカウント
    // 全期間を走査して最長連続（longestStreak）をカウント
}
```

**ルーティング追加（`cmd/server/main.go`）**
```go
private.Get("/users/me/streak", progressHandler.GetMyStreak)
```

**OpenAPI 定義（追記）**
```yaml
/api/v1/users/me/streak:
  get:
    summary: 学習連続日数（ストリーク）取得
    security:
      - bearerAuth: []
    responses:
      '200':
        description: OK
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UserStreak'
      '401':
        description: Unauthorized
```

#### B-1-2: フロントエンド

**対象ファイル（変更）**
- `frontend/src/types/progress.ts`（`UserStreak` 型追加）
- `frontend/src/lib/api.ts`（`fetchMyStreak` 関数追加）
- `frontend/src/app/page.tsx`（ストリーク表示追加）
- `frontend/src/app/profile/page.tsx`（統計セクション追加）

**UserStreak 型（types/progress.ts）**
```ts
export type UserStreak = {
  currentStreak: number;
  longestStreak: number;
  lastStudiedAt: string;
};
```

**トップページ（page.tsx）表示仕様**
- ログインユーザーかつ `currentStreak >= 1` の場合に表示
- ウェルカムセクションの下部に以下のバナーを追加する

```
🔥 3日連続学習中！ | 最長記録 7日
```
- スタイル：`bg-orange-50 border border-orange-200 text-orange-700`

**プロフィールページ表示仕様**
- 既存フォームの上部に「📊 学習統計」カードを追加する
- 表示項目：現在の連続日数・最長連続日数・最終学習日

---

## テスト方針

| 対象 | テスト内容 |
|------|-----------|
| A-4（バックエンド） | `AuthUseCase.ChangePassword` のユニットテスト（正常・現在PW不一致・短すぎるPWのケース） |
| B-1（バックエンド） | `ProgressUseCase.GetStreak` のユニットテスト（0日・連続1日・連続N日・途切れるケース） |

---

## ファイル変更一覧

### フェーズ1（新規作成）

| ファイル | 内容 |
|---------|------|
| `frontend/src/app/glossary/[id]/page.tsx` | 用語詳細ページ（新規） |

### フェーズ1（既存ファイル変更）

| ファイル | 変更内容 |
|---------|---------|
| `frontend/src/app/lessons/[id]/page.tsx` | 次のレッスンナビゲーション追加 |
| `frontend/src/app/profile/page.tsx` | クイズ解答履歴セクション・メモエクスポートボタン追加 |
| `frontend/src/app/glossary/page.tsx` | 各用語カードをリンクに変更 |

### フェーズ2（既存ファイル変更）

| ファイル | 変更内容 |
|---------|---------|
| `backend/internal/domain/models.go` | `UserStreak` 型追加 |
| `backend/internal/usecase/auth.go` | `ChangePassword` メソッド追加 |
| `backend/internal/usecase/progress.go` | `GetStreak` メソッド追加 |
| `backend/internal/handler/user.go` | `ChangePassword` ハンドラー追加・`authUC` 依存追加 |
| `backend/internal/handler/progress.go` | `GetMyStreak` ハンドラー追加 |
| `backend/cmd/server/main.go` | ルーティング追加（password, streak） |
| `backend/openapi.yaml` | 2エンドポイント追加（password, streak） |
| `frontend/src/types/progress.ts` | `UserStreak` 型追加 |
| `frontend/src/lib/api.ts` | `changePassword`, `fetchMyStreak` 追加 |
| `frontend/src/app/page.tsx` | ストリークバナー追加 |
| `frontend/src/app/profile/page.tsx` | パスワード変更フォーム・ストリーク統計追加 |

---

## 実装順序（フェーズ1・2）

| 順序 | 機能 | 理由 |
|------|------|------|
| 1 | A-3 用語詳細ページ | 新規ページ作成・既存への影響最小 |
| 2 | A-1 次のレッスンナビゲーション | `lessons/[id]/page.tsx` 1ファイルのみ変更 |
| 3 | A-2 クイズ解答履歴 + B-3 メモエクスポート | `profile/page.tsx` への同時追加でまとめて実施 |
| 4 | A-4 パスワード変更 | バックエンドから実装し、フロントエンドで確認 |
| 5 | B-1 ストリーク | バックエンドの計算ロジック実装後にフロントエンド実装 |

---

## フェーズ3 実装計画

---

### A-6: 模擬試験 専用UI（Q-03）

**対象ファイル（変更）**
- `frontend/src/app/quiz/[id]/page.tsx`

**実装方針**
- バックエンド変更なし。`Quiz` モデルの `isMockExam: boolean` と `timeLimitMinutes: number` を既に取得しているため、フロントエンドのみで制御する
- 既存の `/quiz/[id]` ページ内で `detail.quiz.isMockExam` を判定し、模擬試験モードと通常クイズモードを分岐させる
- 模擬試験モードでは「開始確認画面 → 全問表示 → カウントダウンタイマー → 一括送信 → 詳細結果」の流れにする

**追加するステートとロジック**

```tsx
// 模擬試験専用ステート
const [mockStarted, setMockStarted] = useState(false);     // 開始確認後
const [timeLeft, setTimeLeft] = useState<number>(0);        // 残り秒数
const [timerId, setTimerId] = useState<ReturnType<typeof setInterval> | null>(null);

// 全問同時表示モードでの回答管理
// answers ステートは既存のまま流用（questionId → choiceId）
```

**開始確認画面（`isMockExam && !mockStarted && !submitted`）**
- 制限時間・問題数を表示する
- 「試験を開始する」ボタンで `mockStarted = true`・タイマー開始
- スタイル：`bg-amber-50 border border-amber-200` の注意書きカードを表示

**タイマー処理**
```tsx
useEffect(() => {
  if (!mockStarted || !detail?.quiz.isMockExam) return;
  const seconds = (detail.quiz.timeLimitMinutes ?? 60) * 60;
  setTimeLeft(seconds);
  const id = setInterval(() => {
    setTimeLeft((prev) => {
      if (prev <= 1) {
        clearInterval(id);
        handleFinish(); // 時間切れで自動送信
        return 0;
      }
      return prev - 1;
    });
  }, 1000);
  setTimerId(id);
  return () => clearInterval(id);
}, [mockStarted]);
```

**タイマー表示（画面右上に固定）**
```tsx
{detail.quiz.isMockExam && mockStarted && !submitted && (
  <div className="fixed right-4 top-16 z-50 rounded-xl border border-amber-300 bg-amber-50 px-4 py-2 text-sm font-bold text-amber-700 shadow-md">
    ⏱ {String(Math.floor(timeLeft / 60)).padStart(2, "0")}:{String(timeLeft % 60).padStart(2, "0")}
  </div>
)}
```

**全問表示モード（`isMockExam && mockStarted && !submitted`）**
- 全問を縦にスクロールして表示する（ページング不要）
- 各問の回答選択後にチェックを入れる（正誤フィードバックは表示しない）
- 全問を表示し終えたら下部に「解答を送信する」ボタンを1つ配置する
- 送信時にタイマーを停止する

**結果画面の拡張（`isMockExam && submitted`）**
- 通常クイズの結果画面に加え、以下を追加表示する
  - かかった時間（`timeLimitMinutes * 60 - timeLeft` を分:秒表示）
  - 問題ごとの正誤一覧テーブル（問題文・選択肢・正答・解説）

---

### B-4: 学習統計サマリ

#### B-4-1: バックエンド

**対象ファイル（変更）**
- `backend/internal/domain/models.go`（`UserStats` 型追加）
- `backend/internal/usecase/progress.go`（`GetStats` メソッド追加）
- `backend/internal/handler/progress.go`（`GetMyStats` ハンドラー追加）
- `backend/cmd/server/main.go`（ルーティング追加）
- `backend/openapi.yaml`（エンドポイント追加）

**追加する依存リポジトリ**
- `ProgressUseCase` に `quizRepo repository.QuizRepository` と `noteRepo repository.NoteRepository` を追加する

`NewProgressUseCase` の引数を追加：
```go
func NewProgressUseCase(
    progress repository.ProgressRepository,
    lessons  repository.LessonRepository,
    courses  repository.CourseRepository,
    sections repository.SectionRepository,
    quizzes  repository.QuizRepository,
    notes    repository.NoteRepository,
) *ProgressUseCase
```

**domain.UserStats 型**
```go
type UserStats struct {
    TotalCompletedLessons int     `json:"totalCompletedLessons"`
    TotalStudyDays        int     `json:"totalStudyDays"`
    TotalNotes            int     `json:"totalNotes"`
    AverageQuizScore      float64 `json:"averageQuizScore"` // 0〜100
}
```

**ProgressUseCase.GetStats メソッド**
```go
func (uc *ProgressUseCase) GetStats(userID string) (domain.UserStats, error) {
    // 1. user_lesson_progress を全件取得し件数・ユニーク日数を算出
    // 2. user_quiz_results を全件取得し正答率の平均を算出
    // 3. user_notes を全件取得し件数を算出
}
```

**ProgressHandler.GetMyStats**
```go
func (h *ProgressHandler) GetMyStats(w http.ResponseWriter, r *http.Request) {
    userID, _ := r.Context().Value(userIDContextKey).(string)
    // ... usecase 呼び出し → WriteJSON
}
```

**ルーティング（`cmd/server/main.go`）**
```go
private.Get("/users/me/stats", progressHandler.GetMyStats)
```

**OpenAPI 定義（追記）**
```yaml
/api/v1/users/me/stats:
  get:
    summary: 学習統計サマリ取得
    security:
      - bearerAuth: []
    responses:
      '200':
        description: OK
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UserStats'
      '401':
        description: Unauthorized
```

#### B-4-2: フロントエンド

**対象ファイル（変更）**
- `frontend/src/types/progress.ts`（`UserStats` 型追加）
- `frontend/src/lib/api.ts`（`fetchMyStats` 関数追加）
- `frontend/src/app/profile/page.tsx`（統計カード追加）

**UserStats 型（types/progress.ts）**
```ts
export type UserStats = {
  totalCompletedLessons: number;
  totalStudyDays: number;
  totalNotes: number;
  averageQuizScore: number;
};
```

**プロフィールページ表示仕様**
- フォームの上部に「📊 学習統計」として 2×2 グリッドのカードを配置する
- ローディング中はスケルトンスクリーン（`animate-pulse`）

```
┌─────────────────┬─────────────────┐
│ 📚 完了レッスン  │ 📅 学習日数      │
│      42 本       │      18 日       │
├─────────────────┼─────────────────┤
│ 📝 メモ件数     │ 🎯 クイズ正答率   │
│       7 件       │     76.3 %       │
└─────────────────┴─────────────────┘
```

---

### B-2: 学習カレンダー（ヒートマップ）

#### B-2-1: バックエンド

**対象ファイル（変更）**
- `backend/internal/domain/models.go`（`CalendarDay`, `UserCalendar` 型追加）
- `backend/internal/usecase/progress.go`（`GetCalendar` メソッド追加）
- `backend/internal/handler/progress.go`（`GetMyCalendar` ハンドラー追加）
- `backend/cmd/server/main.go`（ルーティング追加）
- `backend/openapi.yaml`（エンドポイント追加）

**domain の追加型**
```go
type CalendarDay struct {
    Date  string `json:"date"`  // "2026-03-26"
    Count int    `json:"count"` // その日の完了レッスン数
}

type UserCalendar struct {
    Days []CalendarDay `json:"days"`
}
```

**ProgressUseCase.GetCalendar メソッド**
```go
func (uc *ProgressUseCase) GetCalendar(userID string, year int) (domain.UserCalendar, error) {
    list, err := uc.progress.ListByUserID(userID)
    if err != nil {
        return domain.UserCalendar{}, err
    }
    // completedAt を "YYYY-MM-DD" に丸め、year でフィルタして日付ごとにカウント
    counts := make(map[string]int)
    for _, p := range list {
        if p.CompletedAt.Year() != year {
            continue
        }
        key := p.CompletedAt.UTC().Format("2006-01-02")
        counts[key]++
    }
    days := make([]domain.CalendarDay, 0, len(counts))
    for date, count := range counts {
        days = append(days, domain.CalendarDay{Date: date, Count: count})
    }
    // sort by date asc
    return domain.UserCalendar{Days: days}, nil
}
```

**ルーティング（`cmd/server/main.go`）**
```go
private.Get("/users/me/calendar", progressHandler.GetMyCalendar)
```

**OpenAPI 定義（追記）**
```yaml
/api/v1/users/me/calendar:
  get:
    summary: 学習カレンダー（日別完了数）取得
    security:
      - bearerAuth: []
    parameters:
      - in: query
        name: year
        schema:
          type: integer
        required: false
        description: 対象年（省略時は今年）
    responses:
      '200':
        description: OK
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UserCalendar'
      '401':
        description: Unauthorized
```

#### B-2-2: フロントエンド

**対象ファイル（新規作成）**
- `frontend/src/components/LearningCalendar.tsx`（カレンダーコンポーネント）

**対象ファイル（変更）**
- `frontend/src/types/progress.ts`（`CalendarDay`, `UserCalendar` 型追加）
- `frontend/src/lib/api.ts`（`fetchMyCalendar` 関数追加）
- `frontend/src/app/progress/page.tsx`（カレンダーセクション追加）

**CalendarDay / UserCalendar 型（types/progress.ts）**
```ts
export type CalendarDay = {
  date: string;
  count: number;
};

export type UserCalendar = {
  days: CalendarDay[];
};
```

**api.ts 追加関数**
```ts
export async function fetchMyCalendar(year?: number): Promise<UserCalendar> {
  const query = year ? `?year=${year}` : "";
  const data = await apiFetch<UserCalendar>(`/api/v1/users/me/calendar${query}`);
  return data;
}
```

**LearningCalendar コンポーネントの仕様**
- 1年分（52週 × 7日）を GitHub Contribution グラフ形式で表示する
- ライブラリは使用せず、SVG または `div` グリッドで自前実装する
- 完了数に応じた色分け（ピンク系）

| count | 色クラス |
|-------|---------|
| 0 | `bg-slate-100` |
| 1 | `bg-pink-100` |
| 2〜3 | `bg-pink-300` |
| 4以上 | `bg-pink-500` |

- マウスホバー時に日付と完了数をツールチップで表示する
- 月ラベル（Jan, Feb ... または 1月, 2月 ...）を上部に表示する
- 曜日ラベル（月・水・金）を左側に表示する

**進捗ページへの組み込み**
- 進捗ページ（`/progress`）の既存コースプログレスカードの下に「学習カレンダー」セクションとして追加する
- 当年のカレンダーを表示し、左上に「{year}年の学習記録」テキストを表示する

---

## フェーズ4 実装計画

---

### A-5: 達成バッジ機能（P-04）

#### A-5-1: バックエンド

**前提**
- `badges` / `user_badges` テーブルは `0002_quizzes_glossary_badges.up.sql` で定義済み
- バッジ付与のトリガーはレッスン完了（`POST /api/v1/lessons/:id/complete`）時

**対象ファイル（変更・新規）**
- `backend/internal/domain/models.go`（`Badge`, `UserBadge` 型追加）
- `backend/internal/repository/repository.go`（`BadgeRepository` インターフェース追加）
- `backend/internal/repository/badge_postgres.go`（新規作成）
- `backend/internal/usecase/progress.go`（`badgeRepo` 依存追加・バッジ付与ロジック追加）
- `backend/internal/handler/progress.go`（`GetMyBadges` ハンドラー追加）
- `backend/cmd/server/main.go`（`BadgeRepository` 初期化・ルーティング追加）
- `backend/openapi.yaml`（エンドポイント追加）
- `backend/migrations/0008_badge_seed.up.sql`（バッジマスターデータ投入）
- `backend/migrations/0008_badge_seed.down.sql`（ロールバック）

**domain の追加型**
```go
type Badge struct {
    ID            string    `json:"id"`
    Name          string    `json:"name"`
    Description   string    `json:"description"`
    ImageURL      string    `json:"imageUrl"`
    ConditionType string    `json:"conditionType"` // "section_complete" | "course_complete"
    ConditionID   string    `json:"conditionId"`   // section_id または course_id
}

type UserBadge struct {
    ID       string    `json:"id"`
    UserID   string    `json:"userId"`
    Badge    Badge     `json:"badge"`
    EarnedAt time.Time `json:"earnedAt"`
}
```

**BadgeRepository インターフェース**
```go
type BadgeRepository interface {
    FindByCondition(conditionType, conditionID string) (Badge, bool, error)
    CreateUserBadge(userID, badgeID string) (UserBadge, error)
    ListByUserID(userID string) ([]UserBadge, error)
    ExistsUserBadge(userID, badgeID string) (bool, error)
}
```

**バッジ付与ロジック（ProgressUseCase.CompleteLesson 内に追記）**

レッスン完了後、以下の順で判定してバッジを付与する：
1. **セクション完了バッジ**：そのセクションの全レッスンが完了済みになったか確認 → `condition_type = 'section_complete'` かつ `condition_id = section.id` のバッジを検索 → まだ付与されていなければ付与
2. **コース完了バッジ**：コース全体の全レッスンが完了済みになったか確認 → `condition_type = 'course_complete'` かつ `condition_id = course.id` のバッジを検索 → 同様に付与

`CompleteLesson` の戻り値を変更し、新たに取得したバッジ一覧も返すことで、フロントエンドでトースト表示に使用できるようにする：

```go
type CompleteLessonResult struct {
    Progress    domain.UserLessonProgress `json:"progress"`
    NewBadges   []domain.UserBadge        `json:"newBadges"`
}
```

**マイグレーション（`0008_badge_seed.up.sql`）でのバッジマスターデータ例**
```sql
INSERT INTO badges (id, name, description, image_url, condition_type, condition_id) VALUES
  ('badge-fp3-s1',    'FP3級 第1章マスター', 'ライフプランと資金計画を修了', '', 'section_complete', 'fp3-s1'),
  ('badge-fp3-s2',    'FP3級 第2章マスター', '保険の基礎を修了', '', 'section_complete', 'fp3-s2'),
  -- ... 各セクション分
  ('badge-fp3',       'FP3級コース修了',     'FP3級コース全レッスン修了', '', 'course_complete', 'fp3'),
  ('badge-boki3',     '簿記3級コース修了',   '簿記3級コース全レッスン修了', '', 'course_complete', 'boki3'),
  ('badge-asset3',    '資産運用コース修了',  '資産運用検定3級コース全レッスン修了', '', 'course_complete', 'asset3')
ON CONFLICT (id) DO NOTHING;
```

> **注意**: `badges` テーブルに `condition_id` カラムが存在しない場合、マイグレーション `0009` でカラム追加が必要。既存スキーマ（`0002`）は `condition_type TEXT NOT NULL` のみ定義のため ALTER TABLE で追加する。

**GetMyBadges ハンドラー**
- `GET /api/v1/users/me/badges`
- レスポンス: `{ "badges": [ UserBadge, ... ] }`

#### A-5-2: フロントエンド

**対象ファイル（新規作成）**
- `frontend/src/types/badge.ts`
- `frontend/src/components/BadgeToast.tsx`（バッジ取得通知）

**対象ファイル（変更）**
- `frontend/src/lib/api.ts`（`fetchMyBadges` 関数追加・`completeLesson` の戻り値型更新）
- `frontend/src/app/profile/page.tsx`（バッジ一覧セクション追加）
- `frontend/src/app/lessons/[id]/page.tsx`（完了時のトースト表示追加）

**types/badge.ts**
```ts
export type Badge = {
  id: string;
  name: string;
  description: string;
  imageUrl: string;
  conditionType: string;
};

export type UserBadge = {
  id: string;
  userId: string;
  badge: Badge;
  earnedAt: string;
};

export type UserBadgesResponse = {
  badges: UserBadge[];
};

export type CompleteLessonResult = {
  progress: UserLessonProgress; // types/progress.ts から
  newBadges: UserBadge[];
};
```

**プロフィールページ バッジ表示**
- 「🏅 獲得バッジ」セクションをメモ一覧の上に追加する
- バッジカードをグリッド表示（`grid-cols-3`）
- 取得済み：カラー表示 + 取得日
- 未取得：グレーアウト（`opacity-40 grayscale`）で全バッジを表示して達成感を演出する（全バッジ一覧取得用 API が別途必要）

**BadgeToast コンポーネント**
```tsx
// レッスン完了時に新着バッジがあった場合、右下にアニメーション付きで表示
// 3秒後に自動消去
<div className="fixed bottom-6 right-6 z-50 animate-slide-in rounded-xl
  bg-pink-500 p-4 text-white shadow-lg">
  🎉 バッジを取得しました！<br />
  <strong>{badge.name}</strong>
</div>
```

---

### A-7: コンテンツ検索（C-05）

#### A-7-1: バックエンド

**対象ファイル（新規作成）**
- `backend/internal/repository/search_postgres.go`
- `backend/internal/usecase/search.go`
- `backend/internal/handler/search.go`

**対象ファイル（変更）**
- `backend/internal/repository/repository.go`（`SearchRepository` インターフェース追加）
- `backend/cmd/server/main.go`（Search の初期化・ルーティング追加）
- `backend/openapi.yaml`（エンドポイント追加）

**SearchRepository インターフェース**
```go
type SearchRepository interface {
    SearchLessons(query string) ([]SearchLesson, error)
    SearchTerms(query string) ([]SearchTerm, error)
}

type SearchLesson struct {
    ID        string `json:"id"`
    Title     string `json:"title"`
    SectionID string `json:"sectionId"`
}

type SearchTerm struct {
    ID      string `json:"id"`
    Term    string `json:"term"`
    Reading string `json:"reading"`
}
```

**search_postgres.go の実装**
```sql
-- lessons 検索
SELECT id, title, section_id FROM lessons
WHERE title ILIKE '%' || $1 || '%' OR content ILIKE '%' || $1 || '%'
LIMIT 20;

-- glossary 検索
SELECT id, term, reading FROM glossary_terms
WHERE term ILIKE '%' || $1 || '%' OR reading ILIKE '%' || $1 || '%'
  OR definition ILIKE '%' || $1 || '%'
LIMIT 20;
```

**SearchUseCase**
```go
type SearchResult struct {
    Lessons []SearchLesson `json:"lessons"`
    Terms   []SearchTerm   `json:"terms"`
}

func (uc *SearchUseCase) Search(query string) (SearchResult, error)
```

**SearchHandler**
- `GET /api/v1/search?q={keyword}`（認証不要）
- `q` が空または2文字未満の場合は `400 Bad Request`

**OpenAPI 定義（追記）**
```yaml
/api/v1/search:
  get:
    summary: コンテンツ検索（レッスン・用語）
    parameters:
      - in: query
        name: q
        required: true
        schema:
          type: string
          minLength: 2
        description: 検索キーワード
    responses:
      '200':
        description: OK
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/SearchResult'
      '400':
        description: キーワードが短すぎる
```

#### A-7-2: フロントエンド

**対象ファイル（新規作成）**
- `frontend/src/app/search/page.tsx`（検索結果ページ）
- `frontend/src/types/search.ts`

**対象ファイル（変更）**
- `frontend/src/lib/api.ts`（`searchContent` 関数追加）
- `frontend/src/components/Header.tsx`（検索バー追加）

**types/search.ts**
```ts
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
```

**Header の検索バー仕様**
- デスクトップ: ナビリンクの右に虫眼鏡アイコンボタンを追加
- クリックで検索入力欄がインライン展開する
- `Enter` キーまたはフォーム送信で `/search?q={keyword}` に遷移する
- モバイルメニュー内にも同様の検索入力欄を追加する

**検索結果ページ（`/search?q=...`）の仕様**
- URL の `q` クエリパラメータからキーワードを取得し、`fetchContent(q)` を呼び出す
- レッスン結果・用語結果を「レッスン（X件）」「用語（X件）」の2セクションで表示する
- 0件の場合は「"{q}" に一致する結果が見つかりませんでした。」を表示する
- 各レッスン結果は `/lessons/{id}` へのリンク、各用語結果は `/glossary/{id}` へのリンク

---

## テスト方針（フェーズ3・4追加分）

| 対象 | テスト内容 |
|------|-----------|
| B-4（バックエンド） | `ProgressUseCase.GetStats` のユニットテスト（各集計値の正確性） |
| B-2（バックエンド） | `ProgressUseCase.GetCalendar` のユニットテスト（日付集計・年フィルタの正確性） |
| A-5（バックエンド） | `ProgressUseCase.CompleteLesson` のバッジ付与ロジックのユニットテスト |
| A-7（バックエンド） | `SearchUseCase.Search` のユニットテスト（キーワードマッチ・空文字のエラー処理） |

---

## ファイル変更一覧（フェーズ3・4追加分）

### フェーズ3（新規作成）

| ファイル | 内容 |
|---------|------|
| `frontend/src/components/LearningCalendar.tsx` | 学習カレンダー（ヒートマップ）コンポーネント |

### フェーズ3（既存ファイル変更）

| ファイル | 変更内容 |
|---------|---------|
| `backend/internal/domain/models.go` | `UserStats`, `CalendarDay`, `UserCalendar` 型追加 |
| `backend/internal/usecase/progress.go` | `GetStats`, `GetCalendar` メソッド追加・`quizRepo`/`noteRepo` 依存追加 |
| `backend/internal/handler/progress.go` | `GetMyStats`, `GetMyCalendar` ハンドラー追加 |
| `backend/cmd/server/main.go` | ルーティング追加（stats, calendar） |
| `backend/openapi.yaml` | 2エンドポイント追加（stats, calendar） |
| `frontend/src/types/progress.ts` | `UserStats`, `CalendarDay`, `UserCalendar` 型追加 |
| `frontend/src/lib/api.ts` | `fetchMyStats`, `fetchMyCalendar` 追加 |
| `frontend/src/app/profile/page.tsx` | 統計カード追加 |
| `frontend/src/app/progress/page.tsx` | 学習カレンダーセクション追加 |
| `frontend/src/app/quiz/[id]/page.tsx` | 模擬試験モード分岐・タイマー・全問表示追加 |

### フェーズ4（新規作成）

| ファイル | 内容 |
|---------|------|
| `backend/internal/repository/badge_postgres.go` | BadgeRepository Postgres 実装 |
| `backend/internal/usecase/search.go` | SearchUseCase |
| `backend/internal/repository/search_postgres.go` | SearchRepository Postgres 実装 |
| `backend/internal/handler/search.go` | SearchHandler |
| `backend/migrations/0008_badge_seed.up.sql` | バッジマスターデータ + condition_id カラム追加 |
| `backend/migrations/0008_badge_seed.down.sql` | ロールバック |
| `frontend/src/types/badge.ts` | Badge / UserBadge 型定義 |
| `frontend/src/types/search.ts` | SearchResult / SearchLesson / SearchTerm 型定義 |
| `frontend/src/components/BadgeToast.tsx` | バッジ取得通知コンポーネント |
| `frontend/src/app/search/page.tsx` | 検索結果ページ |

### フェーズ4（既存ファイル変更）

| ファイル | 変更内容 |
|---------|---------|
| `backend/internal/domain/models.go` | `Badge`, `UserBadge`, `SearchLesson`, `SearchTerm` 型追加 |
| `backend/internal/repository/repository.go` | `BadgeRepository`, `SearchRepository` インターフェース追加 |
| `backend/internal/usecase/progress.go` | `badgeRepo` 依存追加・`CompleteLesson` のバッジ付与ロジック追加 |
| `backend/internal/handler/progress.go` | `GetMyBadges` ハンドラー追加 |
| `backend/cmd/server/main.go` | ルーティング追加（badges, search）・SearchHandler 初期化追加 |
| `backend/openapi.yaml` | 3エンドポイント追加（badges, search, UserCalendar スキーマ） |
| `frontend/src/lib/api.ts` | `fetchMyBadges`, `searchContent`, `completeLesson` 戻り値型更新 |
| `frontend/src/app/profile/page.tsx` | バッジ一覧セクション追加 |
| `frontend/src/app/lessons/[id]/page.tsx` | バッジトースト表示追加 |
| `frontend/src/components/Header.tsx` | 検索バー追加 |

---

## 全フェーズ実装順序まとめ

| フェーズ | 順序 | 機能 | 工数感 |
|---------|------|------|--------|
| 1 | 1 | A-3 用語詳細ページ | 小 |
| 1 | 2 | A-1 次のレッスンナビゲーション | 小 |
| 1 | 3 | A-2 クイズ解答履歴 + B-3 メモエクスポート | 小 |
| 2 | 4 | A-4 パスワード変更 | 中 |
| 2 | 5 | B-1 ストリーク | 中 |
| 3 | 6 | A-6 模擬試験 専用UI | 中 |
| 3 | 7 | B-4 学習統計サマリ | 中 |
| 3 | 8 | B-2 学習カレンダー | 中 |
| 4 | 9 | A-5 達成バッジ | 大 |
| 4 | 10 | A-7 コンテンツ検索 | 大 |
