# アセナレ — システム仕様書

## 1. プロジェクト概要

| 項目 | 内容 |
|------|------|
| サービス名 | アセナレ |
| コンセプト | お金の入門的知識を体系的に学べる学習プラットフォーム |
| ターゲット | 個人・サラリーマン・起業家 |
| 対応資格 | FP3級 / 簿記3級 / 資産運用検定3級 |
| リリース時期 | 春（桜・新生活テーマのデザイン） |

---

## 2. 技術スタック

| レイヤー | 技術 |
|----------|------|
| フロントエンド | TypeScript / Next.js / TailwindCSS |
| バックエンド | Go / Air（ホットリロード） |
| データベース | PostgreSQL |
| 認証 | JWT |

---

## 3. 機能要件

### 3-1. ユーザー管理

| # | 機能 | 詳細 |
|---|------|------|
| U-01 | ユーザー登録 | メールアドレス・パスワードで新規登録 |
| U-02 | ログイン | メールアドレス・パスワードによる認証、JWT発行 |
| U-03 | ログアウト | クライアント側のトークン破棄 |
| U-04 | プロフィール表示 | ユーザー名・アバター画像・学習統計の表示 |
| U-05 | プロフィール編集 | ユーザー名・アバター画像の変更 |
| U-06 | パスワード変更 | 現在のパスワードを確認した上で変更 |

### 3-2. コース・コンテンツ管理

| # | 機能 | 詳細 |
|---|------|------|
| C-01 | コース一覧表示 | FP3級 / 簿記3級 / 資産運用検定3級の3コース |
| C-02 | コース詳細表示 | コースの説明・カリキュラム構成・難易度・学習時間の目安 |
| C-03 | セクション（章）管理 | コース内の章ごとにコンテンツをグループ化 |
| C-04 | レッスン表示 | テキスト・図解を用いた学習ページ（Markdown対応） |
| C-05 | コンテンツ検索 | キーワードによるレッスン・用語検索 |

### 3-3. 学習進捗管理

| # | 機能 | 詳細 |
|---|------|------|
| P-01 | レッスン完了マーク | レッスンを読み終えたら完了登録できる |
| P-02 | 進捗率表示 | コース・セクションごとの進捗率（%）をプログレスバーで表示 |
| P-03 | 学習履歴 | 最近学習したレッスンの履歴一覧 |
| P-04 | 達成バッジ | セクション・コース完了時にバッジを付与 |

### 3-4. クイズ・確認テスト

| # | 機能 | 詳細 |
|---|------|------|
| Q-01 | レッスン確認クイズ | レッスン末尾に設置する4択クイズ（3〜5問） |
| Q-02 | セクション復習テスト | セクション完了後に実施できる小テスト（10問前後） |
| Q-03 | 模擬試験 | 資格本試験を想定した制限時間付き模擬試験 |
| Q-04 | 回答解説 | 正答・不正答を問わず解説文を表示 |
| Q-05 | 過去の解答履歴 | 模擬試験の結果スコア・日時の一覧 |

### 3-5. 用語辞典

| # | 機能 | 詳細 |
|---|------|------|
| G-01 | 用語一覧 | 50音・アルファベット順でのインデックス表示 |
| G-02 | 用語詳細 | 用語の定義・関連レッスンへのリンク |

### 3-6. 管理機能（Admin）

| # | 機能 | 詳細 |
|---|------|------|
| A-01 | コース CRUD | コースの作成・編集・削除 |
| A-02 | セクション CRUD | セクションの作成・編集・削除・並び替え |
| A-03 | レッスン CRUD | レッスンの作成（Markdownエディタ）・編集・削除・並び替え |
| A-04 | クイズ CRUD | 問題・選択肢・正答・解説の管理 |
| A-05 | 用語 CRUD | 用語辞典のエントリ管理 |
| A-06 | ユーザー管理 | 登録ユーザーの一覧・停止・削除 |
| A-07 | 画像アップロード | レッスン内で使用する画像のアップロード（UPLOAD_DIR管理） |

---

## 4. 非機能要件

### 4-1. セキュリティ
- 認証はJWTを使用し、有効期限を設ける
- パスワードはbcryptでハッシュ化して保存
- CORS設定は環境変数 `ALLOWED_ORIGINS` で管理
- 管理機能へのアクセスはAdminロールのみ許可

### 4-2. 環境変数
以下の値は必ず環境変数から取得する：

| 変数名 | 説明 |
|--------|------|
| `HOST` | サーバーホスト名 |
| `API_PORT` | APIサーバーのポート番号 |
| `ALLOWED_ORIGINS` | CORS許可オリジン |
| `DATABASE_URL` | PostgreSQL接続URL |
| `JWT_SECRET` | JWT署名シークレット |
| `UPLOAD_DIR` | 画像アップロード先ディレクトリ |

### 4-3. レスポンシブデザイン
- モバイル・タブレット・デスクトップに対応
- TailwindCSSのブレークポイントを利用

### 4-4. パフォーマンス
- 一覧系APIはページネーション対応
- 画像はアップロード時にリサイズ・最適化を検討

---

## 5. データモデル（概要）

```
users
  id, email, password_hash, username, avatar_url, role, created_at, updated_at

courses
  id, title, description, difficulty, estimated_hours, thumbnail_url, order, created_at, updated_at

sections
  id, course_id, title, description, order, created_at, updated_at

lessons
  id, section_id, title, content(markdown), order, created_at, updated_at

quizzes
  id, lesson_id, section_id (nullable), is_mock_exam, time_limit_minutes (nullable), created_at

quiz_questions
  id, quiz_id, question_text, explanation, order

quiz_choices
  id, question_id, choice_text, is_correct

user_lesson_progress
  id, user_id, lesson_id, completed_at

user_quiz_results
  id, user_id, quiz_id, score, total, taken_at

badges
  id, name, description, image_url, condition_type

user_badges
  id, user_id, badge_id, earned_at

glossary_terms
  id, term, reading, definition, created_at, updated_at

lesson_glossary_terms
  lesson_id, term_id
```

---

## 6. API 設計方針

- RESTful API として設計し、`backend/openapi.yaml` で定義・管理する
- エンドポイントプレフィックス: `/api/v1`
- 認証が必要なエンドポイントは `Authorization: Bearer <token>` ヘッダーを使用
- レスポンス形式は JSON

### 主要エンドポイント（抜粋）

| メソッド | パス | 説明 |
|----------|------|------|
| POST | `/api/v1/auth/register` | ユーザー登録 |
| POST | `/api/v1/auth/login` | ログイン |
| GET | `/api/v1/courses` | コース一覧 |
| GET | `/api/v1/courses/:id` | コース詳細 |
| GET | `/api/v1/courses/:id/sections` | セクション一覧 |
| GET | `/api/v1/lessons/:id` | レッスン詳細 |
| POST | `/api/v1/lessons/:id/complete` | レッスン完了登録 |
| GET | `/api/v1/quizzes/:id` | クイズ取得 |
| POST | `/api/v1/quizzes/:id/submit` | クイズ回答送信 |
| GET | `/api/v1/glossary` | 用語一覧 |
| GET | `/api/v1/users/me` | 自分のプロフィール取得 |
| PUT | `/api/v1/users/me` | プロフィール更新 |
| GET | `/api/v1/users/me/progress` | 学習進捗取得 |
| GET | `/api/v1/admin/users` | 管理者：ユーザー一覧 |

---

## 7. UI/UX 方針

- **テーマ**: 桜・新生活をイメージしたフレッシュなデザイン
  - メインカラー: ピンク系（桜色）
  - アクセントカラー: 白・薄緑（新芽）
- **モーダル**: 背景は `bg-black/50` で半透過
- **レイアウト**: サイドバー（コースナビゲーション）+ メインコンテンツ エリア
- **ローディング**: スケルトンスクリーンでUXを向上

---

## 8. マイグレーション方針

- Go CLI コマンドとして実装する（例: `go run ./cmd/migrate up`）
- マイグレーションファイルは `backend/migrations/` 以下に連番で管理する

---

## 9. テスト方針

- バックエンド: Go の標準テストパッケージを使用、各ハンドラー・ユースケースにユニットテストを実施
- フロントエンド: Jest / React Testing Library を使用、主要コンポーネント・APIフックにユニットテストを実施
- 機能追加の際は必ずユニットテストを追加する

---

## 10. ディレクトリ構成（予定）

```
/workspace
├── backend/
│   ├── cmd/
│   │   ├── server/        # サーバーエントリポイント
│   │   └── migrate/       # マイグレーションCLI
│   ├── internal/
│   │   ├── handler/       # HTTPハンドラー
│   │   ├── usecase/       # ビジネスロジック
│   │   ├── repository/    # DBアクセス層
│   │   ├── domain/        # ドメインモデル・インターフェース
│   │   └── middleware/    # 認証・CORS等
│   ├── migrations/        # SQLマイグレーションファイル
│   ├── openapi.yaml       # API仕様書
│   └── .air.toml          # Airの設定
├── frontend/
│   ├── src/
│   │   ├── app/           # Next.js App Router
│   │   ├── components/    # 共通UIコンポーネント
│   │   ├── features/      # 機能ごとのコンポーネント
│   │   ├── hooks/         # カスタムフック
│   │   ├── lib/           # APIクライアント・ユーティリティ
│   │   └── types/         # 型定義
│   └── public/
├── infra/
│   └── compose.yml
└── docs/
    ├── agent/
    │   └── spec.md
    └── user/
        └── draft.md
```
