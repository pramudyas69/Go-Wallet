package repository

import (
	"context"
	"database/sql"
	"e-wallet/domain"
	"github.com/doug-martin/goqu/v9"
)

type factorRepository struct {
	db *goqu.Database
}

func NewFactor(con *sql.DB) domain.FactorRepository {
	return &factorRepository{
		db: goqu.New("default", con),
	}
}

func (f factorRepository) FindByUserID(ctx context.Context, userID int64) (factor domain.Factor, err error) {
	dataset := f.db.From("factors").Where(goqu.Ex{
		"user_id": userID,
	})
	_, err = dataset.ScanStructContext(ctx, &factor)
	return
}

func (f factorRepository) Insert(ctx context.Context, factor *domain.Factor) error {
	executor := f.db.Insert("factors").Rows(goqu.Record{
		"user_id": factor.UserID,
		"pin":     factor.Pin,
	}).Returning("id").Executor()
	_, err := executor.ScanStructContext(ctx, factor)
	return err
}
