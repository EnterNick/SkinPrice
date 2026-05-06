package uow

import (
	"SkinPrice/skinprice/internal/adapters/database/ent"
	"context"
	"fmt"
)

type ctxClientKey struct{}
type ctxTxKey struct{}

var clientKey ctxClientKey
var txKey ctxTxKey

type Manager struct {
	root *ent.Client
}

func New(root *ent.Client) *Manager {
	return &Manager{root: root}
}

func (m *Manager) Client(ctx context.Context) *ent.Client {
	if c, ok := ctx.Value(clientKey).(*ent.Client); ok && c != nil {
		return c
	}
	return m.root
}

func (m *Manager) Tx(ctx context.Context) (*ent.Tx, bool) {
	tx, ok := ctx.Value(txKey).(*ent.Tx)
	return tx, ok && tx != nil
}

func (m *Manager) InTx(ctx context.Context) bool {
	_, ok := ctx.Value(txKey).(*ent.Tx)
	return ok
}

func (m *Manager) Do(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	if m.InTx(ctx) {
		return fn(ctx)
	}

	tx, err := m.root.Tx(ctx)
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}

	ctx = context.WithValue(ctx, clientKey, tx.Client())
	ctx = context.WithValue(ctx, txKey, tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(ctx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("%w: rollback transaction: %w", err, rerr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
