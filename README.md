
```shell
mise install
```

```shell
# 1. モジュール初期化済みとし、依存を取得
go mod tidy

# 2. サーバ起動（カレントディレクトリはプロジェクトルート）
go run ./cmd/api

# 3. 別ターミナルでエンドポイント呼び出し例
# ─────────────────────────────────────
# 3-1. ノート作成
curl -s -X POST http://localhost:8080/notes \
  -H "Content-Type: application/json" \
  -d '{"title":"My First Note","content":"Hello, Onion!"}' | jq

# 3-2. 検索 (キーワード “First” を含むノート一覧)
curl -s "http://localhost:8080/notes?q=First" | jq

# 3-3. ノート編集（例: {ID} は作成時に返った値を流用）
curl -s -X PUT http://localhost:8080/notes/{ID} \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Title","content":"Updated!"}' | jq

# 3-4. ノート削除
curl -s -X DELETE http://localhost:8080/notes/{ID} | jq

# 3-5. ノートリンク (from: {ID1}, to: {ID2})
curl -s -X POST http://localhost:8080/notes/{ID1}/links \
  -H "Content-Type: application/json" \
  -d '{"to":"{ID2}","label":"See also"}' | jq
```

```tree
cmd/           ← エントリポイント（Composition Root）
internal/
├─ domain/           ← ドメイン層（エンティティ・VO・ドメインサービス）
├─ usecase/          ← アプリケーション層（ユースケース実装）
├─ infrastructure/   ← 外部依存（DB, メモリ, MQ Adapter）
└─ presentation/     ← プレゼンテーション層（HTTP/gRPC ハンドラ）
```