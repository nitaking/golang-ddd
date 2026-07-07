# Transaction Manager ガイド

`internal/` の新アーキテクチャにおけるトランザクション管理の設計と使い分けを記載しています。

## 構成

Domain層にinterfaceを定義し、Infrastructure層でEnt実装を提供する。
`context.Context` を通じてトランザクションを伝搬し、Repository が透過的に利用する（context-carried transaction pattern）。

引数で都度渡す場合、Repository Interfaceに実装が漏れるので、contextパターンを採用。

```go
// Domain層: internal/domain/transaction/manager.go
type Manager interface {
    RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
```

```go
// Infrastructure層: internal/infrastructure/ent/transaction/context.go
// context から tx を取得。なければ fallback client を返す
func GetTx(ctx context.Context, fallback *ent.Client) *ent.Tx
```

## Repository 側の実装（共通）

どちらのパターンでも Repository の書き方は同じ: 

```go
func (r *userRepo) Create(ctx context.Context, user *entity.User) error {
    client := transaction.GetTx(ctx, r.client)
    _, err := client.User.Create().
        SetName(user.Name).
        Save(ctx)
    return err
}
```

## パターン1: UseCase 内で明示的に制御する

```go
func (uc *UserUseCase) Register(ctx context.Context, input RegisterInput) error {
    return uc.txManager.RunInTransaction(ctx, func(txCtx context.Context) error {
        user, err := uc.userRepo.Create(txCtx, input.ToEntity())
        if err != nil { return err }
        return uc.profileRepo.Create(txCtx, user.ID, input.Profile)
    })
}
```

**向いているケース:**
- 1リクエスト内で複数のトランザクション境界が必要な場合
- 一部の処理だけトランザクション内で実行したい場合
- 非HTTPコンテキスト（CLI, batch, worker）からも呼ばれる場合
- トランザクションの範囲をビジネスロジックに合わせて細かく制御したい場合

## パターン2: Middleware で自動適用する

```go
func (m TxMiddleware) Handle() gin.HandlerFunc {
    return func(c *gin.Context) {
        m.txManager.RunInTransaction(c.Request.Context(), func(txCtx context.Context) error {
            c.Request = c.Request.WithContext(txCtx)
            c.Next()
            if c.Writer.Status() >= 400 {
                return errors.New("request failed")
            }
            return nil
        })
    }
}
```

Controller は context を渡すだけ:

```go
func (ctrl *UserController) Create(c *gin.Context) {
    ctx := c.Request.Context()  // middleware により tx が入っている
    result, err := ctrl.useCase.Create(ctx, input)
    ...
}
```

**向いているケース:**
- 1エンドポイント = 1トランザクションで十分な場合
- エンドポイント全体を原子的に扱いたい場合（CUD操作など）
- 既存の middleware パターンに慣れたチームでの移行初期

## 使い分けの判断基準

| 判断基準 | パターン1 (UseCase内)                 | パターン2 (Middleware) |
|---------|----------------------------------|----------------------|
| トランザクション境界 | ビジネスロジックに合わせて自由                  | エンドポイント単位で固定 |
| Commit/Rollback 判定 | error の有無                        | HTTP ステータスコード |
| 非HTTP利用 | 可能                               | 不可（Gin依存） |
| 実装コスト | UseCase ごとに明示的に書く                | route に middleware 追加のみ |
| 推奨 | **デフォルトはこちら** | シンプルなCUD操作で十分な場合 |

> [!NOTE] \
> Q. なぜUseCase制御が推奨なのか？\
> トランザクション（原子性）の要件はビジネスロジックであり、その責務はUseCase層であり、HTTP(middleware)の責務ではないため。

## なぜ context.Context ベースなのか

フレームワーク固有の Context（例: `gin.Context.Set()`）でトランザクションを伝搬するアプローチもあるが、以下の問題がある:

```go
// アンチパターン: フレームワーク Context に依存 → Domain/UseCase 層から使えない
trx := c.MustGet("db_transaction").(*ent.Client)
service.DoSomething(trx, ...)
```

`context.Context` ベースにすることで、Domain/UseCase 層がフレームワークに依存しなくなる。
