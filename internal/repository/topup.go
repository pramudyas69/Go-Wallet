package repository

import (
	"context"
	"database/sql"
	"e-wallet/domain"
	"github.com/doug-martin/goqu/v9"
)

type topupRepository struct {
	db *goqu.Database
}

func NewTopUp(con *sql.DB) domain.TopUpRepository {
	return &topupRepository{
		db: goqu.New("default", con),
	}
}

func (t topupRepository) FindByID(ctx context.Context, id string) (topup domain.TopUp, err error) {
	dataset := t.db.From("topup").Where(goqu.Ex{
		"id": id,
	})
	_, err = dataset.ScanStructContext(ctx, &topup)
	return
}

func (t topupRepository) Insert(ctx context.Context, topUp *domain.TopUp) error {
	executor := t.db.Insert("topup").Rows(goqu.Record{
		"id":       topUp.ID,
		"user_id":  topUp.UserID,
		"amount":   topUp.Amount,
		"status":   topUp.Status,
		"snap_url": topUp.SnapURL,
	}).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}

func (t topupRepository) Update(ctx context.Context, topUp *domain.TopUp) error {
	executor := t.db.Update("topup").Set(goqu.Record{
		"amount":   topUp.Amount,
		"status":   topUp.Status,
		"snap_url": topUp.SnapURL,
	}).Where(goqu.Ex{
		"id": topUp.ID,
	}).Executor()
	_, err := executor.ExecContext(ctx)
	return err
}
