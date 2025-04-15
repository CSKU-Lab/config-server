package repositories

import (
	"context"

	"github.com/CSKU-Lab/config-server/domain/models/compare"
)

type CompareRepository interface {
	Add(ctx context.Context, body *compare.Compare) error
	GetAll(ctx context.Context) ([]compare.Compare, error)
	GetByID(ctx context.Context, ID string) (*compare.Compare, error)
	UpdateByID(ctx context.Context, ID string, body *compare.UpdateCompare) error
	DeleteByID(ctx context.Context, ID string) error
}
