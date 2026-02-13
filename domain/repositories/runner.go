package repositories

import (
	"context"

	"github.com/CSKU-Lab/config-server/domain/models"
	"github.com/CSKU-Lab/config-server/domain/requests"
)

type RunnerRepository interface {
	Create(ctx context.Context, ID string, body *requests.CreateRunner) error
	GetAll(ctx context.Context) ([]models.Runner, error)
	GetPagination(ctx context.Context, req *requests.GetPagination) ([]models.Runner, error)
	Count(ctx context.Context) (int, error)
	GetByID(ctx context.Context, ID string) (*models.Runner, error)
	UpdateByID(ctx context.Context, ID string, body *requests.UpdateRunner) error
	DeleteByID(ctx context.Context, ID string) error
}
