# DDD / Clean Architecture 設計ガイド

このファイルには `internal/` 配下のDDDに準じたClean Architecture設計の規約とガイドラインを記載しています。

## 設計原則

### 依存方向ルール
依存は外側→内側の一方向のみ。Domain層は他のレイヤーを知らない。

```
Presentation → UseCase → Domain (Service → Repository → Model)
```

- **Domain → Interface ← Infrastructure**（依存性逆転原則）
- Repository Interfaceはdomain層に定義し、implementsはinfrastructure層に配置
- Domain層でのGin/Entインポートは**禁止**（基本的にGo標準ライブラリのみ）

### YAGNI
- interfaceは必要になった時点で実装する
- ファイル分割も必要になってから実施

### Fail Fast
- VOのコンストラクタでバリデーションを行い、無効な値はその場で`error`として返す
- 後のランタイムエラーを防ぐ

## ディレクトリ構成

```
internal/
├─ domain/
│  ├─ model/                        # 共有ドメインモデル（Date, Email等）
│  ├─ transaction/                 # Transaction Manager interface
│  └─ {domain}/                   # ドメイン単位構成（フラット配置）
│     ├─ {domain}_*.go             # Entity, Value Object, Domain Service
│     └─ {domain}_*_repository.go  # Repository Interface
├─ usecase/
│  ├─ apperror/                    # UseCaseエラー定義（ドメイン横断）
│  └─ {domain}/                   # UseCase実装
├─ infrastructure/
│  ├─ ent/
│  │  ├─ {domain}/                 # ドメイン別ディレクトリ
│  │  │  ├─ {domain}_repository.go # Ent Repository実装
│  │  │  ├─ {domain}_adapter.go    # Ent Model ↔ Domain Entity 変換
│  │  │  └─ schema/                # Ent スキーマ定義
│  │  └─ transaction/              # Transaction Manager Ent実装
│  └─ redis/
│     ├─ client.go                 # 共通Redisクライアント
│     └─ {domain}/                 # ドメイン別ディレクトリ
│        ├─ {domain}_repository.go # Cache Repository実装
│        └─ {domain}_adapter.go    # Cache Model ↔ Domain 変換（必要な場合）
├─ presentation/
│  └─ http/
│     ├─ common/                   # 共通レスポンス型
│     └─ {domain}/                 # Controller, Adapter
└─ infrastructure/
   └─ di/
      └─ injector.go               # DI Container（samber/do）
```

## 命名規約

### パッケージ
- 1ディレクトリ = 1パッケージ
- 短く簡潔な名前（Go公式スタイルに準拠）
- ドメインパッケージ内はフラット構成（`entity/`, `repository/` 等のサブディレクトリは作らない）
- `model/` は Value Object を含むドメインモデル用パッケージ

### メソッド名
- **単一取得**: `Get{Resource}` (ID指定)
- **複数取得**: `List{Resource}s` (一覧)
- **検索**: `Search{Resource}s` (条件検索)
- **作成**: `Create{Resource}`
- **更新**: `Update{Resource}`
- **削除**: `Delete{Resource}`

### UseCase分割タイミング
- **検討**: 1 UseCaseに4〜5メソッドを超え
- **必須**: 1 UseCaseに10メソッドを超える

### ファイル名
- ドメイン固有コードは固有名: `domain/note/note.go`
- UseCaseはファイル名で責務を表現:
  - まとめ型: `usecase/note/note_usecase.go`
  - 機能分割型: `usecase/note/create_note.go`, `search_note.go`

## レイヤー別の責務

### Domain層
- Entity: IDを持ち、同一性で比較。ビジネスロジックを内包
- Value Object: 不変、値等価性。コンストラクタでバリデーション
- Repository Interface: データアクセスの抽象（引数・戻り値はドメイン層の型のみ、インフラ層の型は禁止）
- Domain Service: 複数のrepositoryや集約をまたぐ操作。stateless

### UseCase層
- アプリケーションのユースケースを表現
- ドメインオブジェクトの作成とビジネスバリデーション
- 権限チェック、アクセス制限
- Domain Serviceへのビジネスロジック委譲

### Infrastructure層
- Repository Interfaceの実装
- Ent Model ↔ Domain Entity の変換（Adapterパターン）
- 外部サービスクライアント

### Presentation層
- HTTPリクエスト処理と形式バリデーション（必須フィールド、データ型）
- Domain Entity ↔ HTTP Response の変換
- ビジネスルール検証はNG（UseCaseで実施）

## Transaction Manager

詳細は [Transaction Manager ガイド](./knowledge-transaction.md) を参照。

- Domain層に `Manager` interface、Infrastructure層にEnt実装
- `context.Context` でトランザクションを伝搬、Repository は `GetTx` で透過的に参加
- UseCase内で明示的に制御するパターン（推奨）と、Middleware で自動適用するパターンの2つ

## DI統合

`internal/infrastructure/di/injector.go` に `do.Provide(...)` で依存を登録:

```go
func NewInjector() *do.Injector {
    injector := do.New()
    do.Provide(injector, NewTodoRepository)
    do.Provide(injector, NewCreateTodoUseCase)
    do.Provide(injector, NewTodoHandler)
    return injector
}
```

## 開発フロー

1. **モデリング**: Entity・VO・Repository Interfaceを定義
2. **Interface作成**: Domain層にinterfaceを配置
3. **Implements実装**: Infrastructure層にinterfaceの実装を配置

## 既存コードとの共存（移行時）

レガシーコードと新アーキテクチャを共存させる場合の指針:

- 新規機能は `/internal/` で実装
- 既存機能の修正時は「ボーイスカウトルール」で段階的に移行
- DI Container は既存の依存と並列に登録し後方互換性を維持
- API Interfaceの変更はなし（既存URL・レスポンス形式を維持）
