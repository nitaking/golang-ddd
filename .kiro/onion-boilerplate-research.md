# Onion Architecture Boilerplate Research

## Request

`/Users/satoshi.nitawaki/code/api4/internal` を参考実装として、DDDに準じたOnion Architectureの実践的なボイラープレートを作る。

## Reference Implementation

参照先の主な構成:

```tree
internal/
├─ application/
├─ domain/
├─ infrastructure/
├─ presentation/
└─ README.md
```

参考実装で確認した代表例:

- `domain/email_otp`: Entity / VO / Domain Service / Repository Interface
- `application/email_otp`: UseCase / Application Error変換
- `infrastructure/redis/email_otp`: Repository実装 / Adapter
- `infrastructure/gorm/transaction`: context-carried transaction pattern
- `presentation/http/email_otp`: Gin Controller / HTTP Error変換
- `infrastructure/di`: 依存組み立て

## Sample Domain Decision

ユーザーの記憶では `linear api` または `github` が候補。

今回のボイラープレートでは `linear` を採用した。

理由:

- Issue / Team / Status が小さく、DDDのEntity・VO・Repository Interfaceを説明しやすい
- 外部API連携に発展しやすく、Infrastructure差し替え例として自然
- GitHubより業務ドメイン寄りで、Onion Architectureの境界を説明しやすい

## Implementation Policy

- `domain/application/infrastructure/presentation` のOnion構成にする
- 本物のLinear APIやDBには依存させず、最初は並行安全なインメモリRepositoryで動かす
- `application` と `domain` はGinやDB実装を知らない
- InfrastructureはAdapterで外部保存形式とDomain Entityを変換する
- HTTP層はリクエスト形式検証とApplication ErrorのHTTP変換だけを担当する

## Concurrency Note

GinなどのHTTPサーバーでは複数リクエストが同時にRepositoryへアクセスするため、インメモリ実装でも `sync.RWMutex` で保護する。

DBや本物のLinear APIへ置き換える場合も、同時アクセス・キャンセル・タイムアウトを考慮し、Repository Interfaceは `context.Context` を受け取る。
