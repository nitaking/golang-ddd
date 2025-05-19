package transaction

type Transaction interface {
	Do(fn func() error) error
}
