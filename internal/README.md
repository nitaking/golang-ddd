# `internal`

DDDに準じたOnion Architectureの実践ボイラープレートです。

サンプルドメインは `linear` です。Linear API風の Issue 作成・一覧・完了を題材に、外部DBやAPIキーなしで動く最小構成にしています。

## 依存方向

```text
presentation/http -> application -> domain <- infrastructure
```

- `domain`: Entity / Value Object / Repository Interface / Domain Error
- `application`: UseCase / Input / Output / Application Error変換
- `infrastructure`: Repository実装 / 外部モデルとのAdapter / DI
- `presentation`: HTTP Controller / Request形式検証 / HTTP Error変換

## ディレクトリ

```tree
internal/
├─ domain/
│  └─ linear/                  # Issue Entity, VO, Repository Interface
├─ application/
│  ├─ application_error/       # アプリケーション共通エラー
│  └─ linear/                  # Issue UseCase
├─ infrastructure/
│  ├─ di/                      # 依存組み立て
│  └─ memory/linear/           # 並行安全なRepository実装とAdapter
└─ presentation/http/
   └─ linear/                  # Gin Controller
```

## 実装の読み方

1. `domain/linear/issue.go` で不変条件を表現する
2. `domain/linear/issue_repository.go` に永続化の抽象を置く
3. `application/linear/issue_usecase.go` でユースケースを組み立てる
4. `infrastructure/memory/linear/issue_repository.go` でRepository Interfaceを実装する
5. `presentation/http/linear/controller.go` でHTTP入出力に変換する

## 並行処理の扱い

`memory/linear/IssueRepository` は `sync.RWMutex` で保護しています。HTTPサーバーでは複数リクエストが同時に同じRepositoryを触るため、インメモリ実装でもデータ競合を避けます。

DBや本物のLinear APIに置き換える場合も、`domain.IssueRepository` を満たす実装を `infrastructure` に追加すれば `application` と `domain` は変更不要です。
