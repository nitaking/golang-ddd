# Go DDD Onion Architecture Boilerplate

DDDに準じたOnion Architectureの実践的なGoボイラープレートです。

参考実装: `/Users/satoshi.nitawaki/code/api4/internal`

## Sample Domain

サンプルドメインは `linear` です。

Linear API風の Issue 作成・一覧・完了を題材に、以下の境界を確認できます。

- `domain`: Entity / Value Object / Repository Interface
- `application`: UseCase / Input / Output / Application Error
- `infrastructure`: Repository実装 / Adapter / DI
- `presentation`: Gin Controller / HTTP Request/Response

## Run

```bash
go run ./cmd/api
```

## Try

```bash
curl -X POST http://localhost:8080/linear/issues \
  -H 'Content-Type: application/json' \
  -d '{"teamId":"team-1","title":"Fix login","description":"sample issue"}'
```

```bash
curl http://localhost:8080/linear/teams/team-1/issues
```

## Structure

```tree
internal/
├─ domain/linear
├─ application/linear
├─ infrastructure/memory/linear
├─ infrastructure/di
└─ presentation/http/linear
```

詳細は `internal/README.md` と `.kiro/onion-boilerplate-research.md` を参照してください。
