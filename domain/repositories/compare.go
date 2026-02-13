package repositories

import (
	"context"

	"github.com/CSKU-Lab/config-server/domain/models"
	"github.com/CSKU-Lab/config-server/domain/requests"
)

type CompareRepository interface {
	Create(ctx context.Context, ID string, body *requests.CreateCompare) error
	GetAll(ctx context.Context) ([]models.Compare, error)
	GetPagination(ctx context.Context, req *requests.GetPagination) ([]models.Compare, error)
	Count(ctx context.Context) (int, error)
	GetByID(ctx context.Context, ID string) (*models.Compare, error)
	UpdateByID(ctx context.Context, ID string, body *requests.UpdateCompare) error
	DeleteByID(ctx context.Context, ID string) error
}
