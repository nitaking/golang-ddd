package bolt

import (
	"go-clean-architecture-boilerplate/internal/usecase/transaction"
	"go.etcd.io/bbolt"
)

type boltTransaction struct {
	db *bbolt.DB
}

func NewTransaction(db *bbolt.DB) transaction.Transaction {
	return &boltTransaction{db: db}
}

/* <<<<<<<<<<<<<<  ✨ Windsurf Command ⭐ >>>>>>>>>>>>>>>> */
// Do executes a function within a Bolt transaction.
//
// The function is called with the transaction held, and must return
// either nil or an error. If the function returns an error, the
// transaction is rolled back and the error is returned; otherwise,
// the transaction is committed and nil is returned.
/* <<<<<<<<<<  107f0745-a046-4cc7-a1cf-0dc4f9d19a80  >>>>>>>>>>> */
func (t *boltTransaction) Do(fn func() error) error {
	return t.db.Update(func(tx *bbolt.Tx) error {
		return fn()
	})
}
