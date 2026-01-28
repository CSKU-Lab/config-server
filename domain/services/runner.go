package services

import (
	"context"

	"github.com/CSKU-Lab/config-server/domain/models"
	"github.com/CSKU-Lab/config-server/domain/repositories"
	"github.com/CSKU-Lab/config-server/domain/requests"
	"github.com/google/uuid"
)

type runnerService struct {
	repo repositories.RunnerRepository
}

type RunnerService interface {
	Create(ctx context.Context, body *requests.CreateRunner) (string, error)
	GetAll(ctx context.Context) ([]models.Runner, error)
	GetByID(ctx context.Context, ID string) (*models.Runner, error)
	UpdateByID(ctx context.Context, ID string, body *requests.UpdateRunner) error
	DeleteByID(ctx context.Context, ID string) error
}

func NewRunnerService(repo repositories.RunnerRepository) *runnerService {
	return &runnerService{
		repo: repo,
	}
}

func (l *runnerService) Create(ctx context.Context, body *requests.CreateRunner) (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	err = l.repo.Create(ctx, id.String(), body)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (l *runnerService) GetAll(ctx context.Context) ([]models.Runner, error) {
	return l.repo.GetAll(ctx)
}

func (l *runnerService) GetByID(ctx context.Context, ID string) (*models.Runner, error) {
	return l.repo.GetByID(ctx, ID)
}

func (l *runnerService) UpdateByID(ctx context.Context, ID string, body *requests.UpdateRunner) error {
	return l.repo.UpdateByID(ctx, ID, body)
}

func (l *runnerService) DeleteByID(ctx context.Context, ID string) error {
	return l.repo.DeleteByID(ctx, ID)
}
