# Implementation Plan 1: MVP Vertical Slice

## 1. MVP の目標
「AI（Gemini）との対話を通じてその日の出来事を整理し、ローカル環境に Markdown 形式で保存・閲覧できる」というコア体験を、最小構成で実現する。

### 実現範囲 (Scope)
- **モード**: セルフホストモードのみ（認証はローカル固定、複雑な SSO は除外）。
- **AI 連携**: Built-In 方式（DiaLogCoreの環境変数 `GEMINI_API_KEY` を使用）。
- **日記作成**: 
    - AI とのリアルタイムチャット（ストリーミング応答）。
    - チャット終了後の日記本文（Markdown）の自動生成。
- **ストレージ**: ローカルファイルシステムへの保存（YAML Frontmatter 付き Markdown）。
- **閲覧**: 日記の一覧表示と詳細表示（Markdown レンダリング）。
- **技術スタック**:
    - **DiaLogCore**: Go (Connect-RPC)
    - **UIServer**: Next.js (TypeScript, Vanilla CSS)
    - **Protocol**: Protocol Buffers (Connect-RPC)

### 今回対象外とするもの (Deferred)
- クラウドモード（Google SSO, Firestore）。
- 週間/月間まとめ機能。
- 天気・ニュース API 連携（まずは純粋な対話のみ）。
- BYOK 機能（UI からの API キー設定と保存）。
- ドキュメントサーバの独立プロセス化（DiaLogCoreに内包させる）。

---

## 2. 具体的な実装手順

### Phase 1: プロジェクト基盤の構築
1. **ディレクトリ構造の決定**: `DiaLogCore/`, `UIServer/`, `proto/` の構成を作成。
2. **Protocol Buffers の定義**: `proto/diary/v1/diary.proto` を作成し、MVP に必要なメソッド（`Chat`, `GetDiary`, `ListDiaries`, `UpdateDiary`）を定義。
3. **ボイラープレート生成**: `buf` を用いて Go および TypeScript のコードを生成。

### Phase 2: DiaLogCoreのコア実装 (Go)
1. **Connect-RPC サーバの起動**: 基本的なルーティングの設定。
2. **AI Provider 実装**: 環境変数から API キーを読み込み、Gemini API と連携（ストリーミング対応）。
3. **LocalStorage 実装**: 
    - 指定ディレクトリへの `.md` ファイル書き込み機能。
    - YAML Frontmatter のパースと書き出し。
4. **簡易インメモリインデックス**: 起動時にファイルをスキャンし、一覧表示用のメタデータをメモリに保持。

### Phase 3: UIServerのコア実装 (Next.js)
1. **API クライアントの設定**: `connect-es` を用いたDiaLogCoreとの通信設定。
2. **チャット画面 (Diary Chat)**:
    - ストリーミング応答の表示。
    - 入力フォームと「日記を確定する」ボタン。
3. **ダッシュボード**: 「対話を始める」ボタンと最近の日記へのリンク。

### Phase 4: 閲覧・管理機能の実装
1. **日記一覧画面 (Diary List)**: ローカルに保存された日記を日付順にリスト表示。
2. **日記詳細画面 (Diary Detail)**: 生成された日記（Markdown）のレンダリング表示。
3. **簡易編集機能**: 生成された日記本文の直接編集・保存。

---

## 3. マイルストーン
- [ ] **Phase 1**: プロトコル定義とDiaLogCore/UIServerの疎通確認。
- [ ] **Phase 2**: AI とのチャット・日記生成・ローカル保存の完了。
- [ ] **Phase 3**: 一覧・詳細画面の実装。
