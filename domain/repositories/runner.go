package repositories

import (
	"context"

	"github.com/CSKU-Lab/config-server/domain/models"
	"github.com/CSKU-Lab/config-server/domain/requests"
)

type RunnerRepository interface {
	Create(ctx context.Context, ID string, body *requests.CreateRunner) error
	GetAll(ctx context.Context) ([]models.Runner, error)
	GetByID(ctx context.Context, ID string) (*models.Runner, error)
	UpdateByID(ctx context.Context, ID string, body *requests.UpdateRunner) error
	DeleteByID(ctx context.Context, ID string) error
}
