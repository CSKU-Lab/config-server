package repositories

import (
	"context"

	"github.com/CSKU-Lab/config-server/domain/models/runner"
)

type RunnerRepository interface {
	Add(ctx context.Context, body *runner.Runner) error
	GetAll(ctx context.Context) ([]runner.Runner, error)
	GetByID(ctx context.Context, ID string) (*runner.Runner, error)
	UpdateByID(ctx context.Context, ID string, body *runner.UpdateRunner) error
	DeleteByID(ctx context.Context, ID string) error
}
