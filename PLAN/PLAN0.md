# Implementation Plan 0: Environment Setup & Hello World

## 1. 目的
開発環境を構築し、バックエンド (Go) とフロントエンド (Next.js) が Protocol Buffers (Connect-RPC) を介して型安全に通信できることを確認する。

## 2. 事前準備 (Prerequisites)
以下のツールがインストールされていることを確認してください。
- **Go**: 1.22 以上
- **Node.js**: 20 以上 (pnpm 推奨)
- **Buf CLI**: Protobuf の管理ツール ([インストール方法](https://buf.build/docs/installation))
- **Protobuf Plugins**:
  - `protoc-gen-go`
  - `protoc-gen-connect-go`
  - `@connectrpc/protoc-gen-connect-es`
  - `@bufbuild/protoc-gen-es`

---

## 3. ステップバイステップ手順

### Step 1: プロジェクト構造の作成
プロジェクトのルートディレクトリで以下の構造を作成します。
```text
.
├── proto/             # Protocol Buffers 定義
├── backend/           # Go バックエンド
└── frontend/          # Next.js フロントエンド
```

### Step 2: Protocol Buffers の定義
`proto/diary/v1/diary.proto` を作成し、疎通確認用の `Ping` メソッドを定義します。

```proto
syntax = "proto3";
package diary.v1;
option go_package = "github.com/user/hackason/gen/diary/v1;diaryv1";

message PingRequest { string name = 1; }
message PingResponse { string message = 1; }

service DiaryService {
  rpc Ping(PingRequest) returns (PingResponse);
}
```

`buf.gen.yaml` を作成し、Go と TypeScript のコード生成設定を行います。

### Step 3: バックエンドの初期化 (Go)
1. `backend/` ディレクトリで `go mod init` を実行。
2. `buf generate` で Go のコードを生成。
3. `main.go` を作成し、`DiaryService` の `Ping` メソッドを実装。
4. CORS 設定（開発用）を行い、ポート 8080 でサーバを起動。

### Step 4: フロントエンドの初期化 (Next.js)
1. `frontend/` ディレクトリで `npx create-next-app@latest` を実行。
2. `buf generate` で TypeScript のコードを生成。
3. `connect-es` を用いたクライアントを作成。
4. 画面上のボタンを押すとバックエンドの `Ping` を呼び出し、結果を表示する簡単な UI を作成。

---

## 4. 完了条件 (Definition of Done)
- [ ] `buf generate` がエラーなく実行され、Go/TS のコードが生成される。
- [ ] バックエンドサーバが起動し、`localhost:8080` でリクエストを待機している。
- [ ] フロントエンドからボタンを押すと、バックエンドから "Hello, [Name]!" というレスポンスが返り、画面に表示される。
