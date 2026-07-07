# `internal`

DDDに準じたClean Architectureに基づくアーキテクチャです。

> [!IMPORTANT]
> 依存方向は以下になります.<br>
> Presentation(Controller) -> UseCase -> Service -> Repository -> Model

## Overview
```tree
internal/
├─ domain/           ← Domain Layer（Entity, ValueObject, Repository Interface, Domain Service）
├─ usecase/          ← Application Layer（UseCase, Interactor）
├─ infrastructure/   ← Infrastructure Layer（Repository Implementation, External Service Client）
└─ presentation/     ← Presentation Layer（Controller, Presenter, ViewModel）
```

| レイヤー         | 責任                               | 検証内容                        |
|--------------|----------------------------------|-----------------------------|
| Domain       | ビジネスロジック                         | 有効な値の範囲、ドメインルール             |
| Usecase      | アプリケーションのユースケースを表現し、ドメインモデルを操作する | 権限チェック、アクセス制限、Serviceを使った検証 |
| Presentation | コントローラーやプレゼンターを通じてHTTPリクエストを処理する | 形式的な検証（必須フィールド、データ型など）      |

## Controller

Controllerは、HTTPリクエストを受け取り、UseCaseを呼び出してレスポンスを返す役割を持ちます。
入力データの形式的な検証は、Controllerでおこないます。

```go
type UserController struct {
    userUsecase *usecase.UserUsecase
}
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // Controller層での形式バリデーション
    if validationErrors := h.validateCreateUserRequest(req); len(validationErrors) > 0 {
        responses.BadRequest(g, constants.GetErrorMessage("YOUR ERROR CODE"))
        return
    }
}
```

- OK
    - 必須フィールドの形式チェック
    - 文字列制限
    - データ型検証
- NG
    - ビジネスルールに基づいた検証(Usecaseで実施)
    - メールアドレスの重複チェック(Usecaseで実施)

## UseCase

- UseCaseは、アプリケーションのユースケースを表現する役割です。ドメインモデルを操作し、ビジネスロジックを実行します。
- 基本的にはDomain層であるRepositoryやServiceを利用してドメインモデルの操作を行います。
  - Service層は必須ではありませんが、複数のRepositoryを参照する場合や、Modelの操作が複雑な場合に利用します。
  - 「置き場所に困る場合はServiceに移譲する」ような感覚でも問題ありません。

```go
func (uc *CreateTodoUseCase) Execute(ctx context.Context, input CreateTodoInput) error {
    // Usecase層：ドメインオブジェクトの作成とビジネスバリデーション
    todo, err := domain.NewTodo(input.Title, input.Description)
    if err != nil {
        return fmt.Errorf("invalid todo: %w", err)
    }

    // データ永続化
    if err := uc.todoRepo.Save(ctx, todo); err != nil {
        return fmt.Errorf("failed to save todo: %w", err)
    }

    return nil
}
```

> [!TIP]
> - 要約の取得のような、ファットで複数のRepositoryを参照するものはServiceがよいでしょう

## Infrastructure

RepositoryやServiceの実装を含むレイヤーです。
Interfaceは`domain`に定義し、実装はここに配置します。

### ディレクトリ構成

```tree
internal/infrastructure/
├── ent/
│   ├── {domain}/              ← ドメイン別ディレクトリ
│   │   ├── {domain}_repository.go
│   │   ├── {domain}_adapter.go    ← Ent Model ↔ Domain Entity 変換
│   │   └── schema/                ← Ent スキーマ定義
│   └── transaction/           ← Transaction Manager実装
├── redis/
│   ├── client.go              ← 共通 Redis クライアント（Interfaceベース）
│   └── {domain}/              ← ドメイン別ディレクトリ
│       ├── {domain}_repository.go
│       └── {domain}_adapter.go    ← 必要な場合
└── di/
    └── injector.go            ← DI Container（samber/do）
```

### Adapter パターン

Infrastructure層では、外部データモデル（Ent Model, Cache Model）とDomain Entityの変換にAdapterパターンを使用します。

```go
// Ent Model → Domain Entity
func (a *Adapter) ToDomain(model *ent.Todo) (*entity.Todo, error) { ... }

// Domain Entity → Ent Create Input
func (a *Adapter) ToCreateInput(entity *entity.Todo) *ent.TodoCreate { ... }
```

## 新規機能の追加手順

1. `internal/domain/{domain}/` に Entity・Model・Repository Interface を定義
2. `internal/usecase/{domain}/` に UseCase を実装
3. `internal/infrastructure/` に Repository・Adapter を配置（`ent/{domain}/`, `redis/{domain}/`）
4. `internal/presentation/http/{domain}/` に Controller・Adapter を実装
5. `internal/infrastructure/di/injector.go` に `do.Provide(...)` を追加

## 関連ドキュメント

- [DDD / Clean Architecture 設計ガイド](./knowledge-ddd-architecture.md) — 依存方向ルール、命名規約、レイヤー別責務
- [Transaction Manager ガイド](./knowledge-transaction.md) — トランザクション管理の設計と使い分け
- [Design Document](./design-architecture-migration.md) — 設計判断の経緯と根拠
- [linear ドメイン モデリング設計ドキュメント](./knowledge-linear-modeling.md) — ユビキタス言語、集約・境界、業務フロー、ER図

## Design Principles

- **YAGNI**: インターフェースや抽象化は必要になった時点で実装する
- **Fail Fast**: VOのコンストラクタでバリデーションを行い、無効な値をその場で`error`として返す
- **UseCase分割**: 1 UseCaseのメソッドが4〜5個を超えたら分割を検討、10個を超えたら必須
