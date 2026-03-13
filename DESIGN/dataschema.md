# Data Schema

本プロジェクトにおけるデータモデル（エンティティ、RPCメッセージ等）のフィールド定義については、 **本ドキュメントを唯一の正解（SSOT）とする。** 

gRPC/Connect-RPC のスキーマ定義 (`diary.proto`) は、本ドキュメントの定義に厳密に従うこと。

## 基本型

- `UUID`: 一意識別子
- `ISO8601`: 日時文字列
- `Cursor`: ページング用の不透明文字列

## 列挙・制約

- `Speaker`: `"User" | "Agent"`
- `Feedback`: `-1 | 0 | 1`
- `SummaryGranularity`: `"weekly" | "monthly"`

## エンティティ

### ChatEntry

- `chatId`: UUID
- `speaker`: Speaker
- `content`: string
- `timestamp`: ISO8601
- `feedback`: Feedback

### DiaryEntry

- `diaryId`: UUID
- `topic`: string
- `createdAt`: ISO8601
- `updatedAt`: ISO8601
- `chatHistory`: ChatEntry[]
- `content`: string
- `weather`: WeatherInfo (optional)
- `news`: NewsArticle[] (optional)

実装の都合上以下のフィールドを含む可能性がある:
`latestChatId`: UUID - ストリーミング制御（論理的中断チェック）に使用。

### WeatherInfo

- `summary`: string (e.g., "晴れ", "雨")
- `temperature`: number (Celsius)
- `iconUrl`: string (optional)

### NewsArticle

- `title`: string
- `url`: string
- `source`: string
- `publishedAt`: ISO8601

### DiaryListItem

- `diaryId`: UUID
- `topic`: string
- `createdAt`: ISO8601
- `updatedAt`: ISO8601
- `snippet`: string

### Summary

- `summaryId`: UUID
- `granularity`: SummaryGranularity
- `summarizedDiaries`: UUID[]
- `summaryContent`: string
- `topic`: string
- `periodStart`: ISO8601
- `periodEnd`: ISO8601
- `createdAt`: ISO8601
- `updatedAt`: ISO8601

### UserProfile

- `userId`: UUID
- `location`: string (天気情報の取得用。デフォルト "Tokyo")
- `persona`: PersonaData
- `conversationContext`: ConversationContextData
- `updatedAt`: ISO8601

### PersonaData

- `nickname`: string (optional)
- `traits`: string[] (optional)
- `description`: string (optional)

### ConversationContextData

- `extractedTraits`: string[]
- `contextSummary`: string
- `lastExtractedAt`: ISO8601

### ErrorResponse

API Layer からクライアントに返却される標準的なエラー形式。
内部的には `APIError` オブジェクトとして扱われ、下位層のエラーを合成して保持する。

- `code`: `Code` (Internal | Unauthenticated | NotFound | ResourceExhausted | InvalidArgument)
- `message`: string (クライアントに表示可能な安全なメッセージ)
- `requestId`: string (ログ追跡用ID)
- `details`: object (optional, **非公開・デバッグ用**) - 下記の `DomainError` 構造を含む

### (内部型) DomainError

ドメイン層（Core Logic）で発生する、ビジネスルールに関連するエラー。

- `code`: `DomainErrorCode` (DIARY_NOT_FOUND | BUDGET_EXCEEDED | INVALID_OPERATION | INTEGRITY_ERROR)
- `message`: string (ドメイン文脈のメッセージ)
- `cause`: `ProviderError` (optional, 下位エラーをラップ)

### (内部型) ProviderError

外部リソースアクセス（Adapter）層で発生する、技術的なエラー。

- `code`: `ProviderErrorCode` (NETWORK_ERROR | AUTH_FAILED | RATE_LIMIT | DISK_FULL | UPSTREAM_5XX)
- `message`: string (外部ライブラリ等の生メッセージ)
- `raw`: object (optional, 生の例外オブジェクト)

### ExportResponse

- `filename`: string
- `content`: string (Markdown format with YAML Frontmatter)

### MigrationPackageResponse

- `downloadUrl`: string (ZIP archive URL)
- `summary`: MigrationPackageSummary

### MigrationPackageSummary

- `totalDiaries`: number
- `totalSummaries`: number
- `exportedAt`: ISO8601
- `filters`: PeriodFilter (optional)

## 補助データ構造

### PeriodFilter

- `after`: ISO8601 (optional)
- `before`: ISO8601 (optional)

制約:
- `after <= before` を満たすこと

### PagingQuery

- `limit`: number (optional, default 20, max 100)
- `cursor`: Cursor (optional)

## 関連

- `DiaryEntry.chatHistory` は `ChatEntry` の配列。
- `Summary.summarizedDiaries` は `DiaryEntry.diaryId` の配列。
- `DiaryListItem` は `DiaryEntry` の一覧向け投影（サブセット＋ `snippet`）として扱う。
- `DiaryListItem.snippet` は `DiaryEntry.content` の先頭抜粋（長さは実装依存）。

## Markdown ストレージ・エクスポート形式

セルフホストモードにおける保存および全モードでのエクスポートにおいて、**YAML Frontmatter 形式を厳格に遵守する**。これにより、外部ツール（Obsidian等）での閲覧性と、システムへの再インポートの互換性を担保する。

### 1. ディレクトリ構造 (Self-Host)

```text
storage_root/
├── diaries/             # 日々の対話・日記（DiaryEntry）
│   ├── 2026-03-12_550e8400.md  # YYYY-MM-DD_{UUIDの一部}.md
│   └── ...
├── summaries/           # 要約データ（Summary）
│   ├── weekly/          # 週間まとめ
│   │   ├── 2026-W11_770e8400.md
│   │   └── ...
│   └── monthly/         # 月間まとめ
│       ├── 2026-03_880f9511.md
│       └── ...
└── user/                # ユーザープロファイル・設定
    ├── profile.json     # ペルソナ、会話コンテキスト（UserProfile）
    └── settings.json    # システム動作設定（Mode等）
```

### 2. 共通の YAML Frontmatter 要件
以下のフィールドを必須とする：
- `id`: UUID
- `topic`: string
- `created_at`: ISO8601
- `updated_at`: ISO8601
- `chat_history`: (対話ログのリスト)

### 3. マイグレーション・パッケージ (Migration/Bulk Export)

Cloud モードから Self-host モードへの移行、または完全なバックアップを目的としたパッケージ形式。
以下のディレクトリ構造を持つ ZIP ファイルとして提供され、展開することでそのままセルフホストモードの `storage_root` として機能する。

```text
storage_root/ (ZIP Root)
├── export_metadata.json  # エクスポート日時、対象件数、ユーザーID等のメタデータ
├── diaries/              # 全ての日記エントリ (.md)
├── summaries/            # 全ての要約データ (.md)
│   ├── weekly/
│   └── monthly/
└── user/                 # ユーザーデータ
    ├── profile.json      # ペルソナ、会話コンテキスト
    └── settings.json     # システム設定
```

#### export_metadata.json の構造
- `version`: エクスポートフォーマットのバージョン
- `exportedAt`: ISO8601
- `totalDiaries`: 件数
- `totalSummaries`: 件数
- `includesUserProfile`: boolean (プロファイルが含まれているか)
- `sourceMode`: "cloud" | "self-host"


### 4. Diary 形式
- **YAML Frontmatter**: `id`, `topic`, `created_at`, `updated_at`, `weather`, `news`, `chat_history` を保持。
- **Markdown Sections**: `# トピック`, `## Weather`, `## News`, `## Diary` (本文) を保持。

### 5. Summary 形式
- **YAML Frontmatter**: `id`, `granularity`, `period_start`, `period_end`, `topic`, `summarized_diaries` (UUIDリスト), `created_at`, `updated_at` を保持。
- **Markdown Sections**: `# トピック`, `## Period`, `## Summarized Diaries` (元日記へのリンク), `## Summary` (本文) を保持。

### 4. ファイル構造例 (`summaries/weekly/2026-W11_770e8400.md`)

```markdown
---
id: "770e8400-e29b-41d4-a716-446655449999"
granularity: "weekly"
period_start: "2026-03-09T00:00:00Z"
period_end: "2026-03-15T23:59:59Z"
topic: "春の訪れと新しい趣味の始まり"
summarized_diaries:
  - "550e8400-e29b-41d4-a716-446655440000"
  - "660f9511-f30c-52d5-b827-557766551111"
created_at: "2026-03-16T09:00:00Z"
updated_at: "2026-03-16T09:00:00Z"
---

# 春の訪れと新しい趣味の始まり

## Period
2026-03-09 〜 2026-03-15 (Week 11)

## Summarized Diaries
- [穏やかな午後のコーヒー](../../diaries/2026-03-12_550e8400.md)
- [週末のキャンプ計画](../../diaries/2026-03-14_660f9511.md)

## Summary
今週は気温が上がり、春の訪れを強く感じる一週間でした。
前半は新しいコーヒー豆の開拓に熱心で、後半は週末に向けたキャンプの準備をAIと相談しながら進めました。
全体として、日常の小さな楽しみに焦点を当てた充実した日々でした。
```

## クラウド実装形式 (Firestore)

Cloudモードにおける Firestore コレクション構造とデータ型の定義。

### 1. コレクション構造

Firestore は階層化されたドキュメント指向 DB として、各ユーザーのサンドボックス構造を採用する。

```text
/users/{userId} (Document: UserProfile)
    ├── /diaries/{diaryId} (Document: DiaryEntry)
    │       └── /chats/{chatId} (Document: ChatEntry)
    └── /summaries/{summaryId} (Document: Summary)
```

### 2. データ型のマッピング

| 論理型 (`dataschema.md`) | Firestore 型 | 特記点 |
| :--- | :--- | :--- |
| `UUID` | `String` | ドキュメントの ID 自体として使用する |
| `ISO8601` | `Timestamp` | クエリ（時系列ソート）の効率を最大化するため |
| `Speaker` | `String` | `"User"` または `"Agent"` |
| `Feedback` | `Integer` | `-1`, `0`, `1` |
| `WeatherInfo` | `Map` | `summary` (string), `temperature` (number) 等を保持 |
| `NewsArticle[]` | `Array` | Map オブジェクトの配列 |
| `chatHistory` | **サブコレクション** | 1ドキュメントあたりのサイズ制限（1MB）回避のため |

### 3. セキュリティルールの基本方針

Firebase Authentication と連携し、自身のデータのみにアクセスできるルール（RBAC）を適用する。

```javascript
service cloud.firestore {
  match /databases/{database}/documents {
    // ユーザー個別のデータパスへのアクセス制限
    match /users/{userId}/{document=**} {
      allow read, write: if request.auth != null && request.auth.uid == userId;
    }
  }
}
```

### 4. インデックス設計の要件

効率的な取得のため、以下の複合インデックスが必要となる。
- `diaries`: `createdAt` DESC (一覧表示用)
- `summaries`: `granularity` ASC + `periodStart` DESC (期間別まとめ表示用)
- `chats`: `timestamp` ASC (対話履歴の再現用)
