package services

import (
	"context"

	"github.com/CSKU-Lab/config-server/domain/models/runner"
	"github.com/CSKU-Lab/config-server/domain/repositories"
)

type runnerService struct {
	repo repositories.RunnerRepository
}

type RunnerService interface {
	Add(ctx context.Context, body *runner.Runner) error
	GetAll(ctx context.Context) ([]runner.Runner, error)
	GetByID(ctx context.Context, ID string) (*runner.Runner, error)
	UpdateByID(ctx context.Context, ID string, body *runner.PartialOptions) (*runner.Runner, error)
	DeleteByID(ctx context.Context, ID string) error
}

func NewLanguageService(repo repositories.RunnerRepository) *runnerService {
	return &runnerService{
		repo: repo,
	}
}

func (l *runnerService) Add(ctx context.Context, body *runner.Runner) error {
	return l.repo.Add(ctx, body)
}

func (l *runnerService) GetAll(ctx context.Context) ([]runner.Runner, error) {
	return l.repo.GetAll(ctx)
}

func (l *runnerService) GetByID(ctx context.Context, ID string) (*runner.Runner, error) {
	return l.repo.GetByID(ctx, ID)
}

func (l *runnerService) UpdateByID(ctx context.Context, ID string, body *runner.PartialOptions) (*runner.Runner, error) {
	lang, err := l.repo.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	if body.Name == nil {
		body.Name = &lang.Name
	}

	modRunner := runner.NewUpdate(&runner.PartialOptions{
		Name:        body.Name,
		BuildScript: body.BuildScript,
		RunScript:   body.RunScript,
	})

	err = l.repo.UpdateByID(ctx, ID, modRunner)
	if err != nil {
		return nil, err
	}

	buildScript := lang.BuildScript
	if body.BuildScript != nil {
		buildScript = *body.BuildScript
	}

	runScript := lang.BuildScript
	if body.RunScript != nil {
		runScript = *body.RunScript
	}

	return &runner.Runner{
		ID:          *modRunner.ID,
		Name:        *modRunner.Name,
		BuildScript: buildScript,
		RunScript:   runScript,
	}, nil
}

func (l *runnerService) DeleteByID(ctx context.Context, ID string) error {
	return l.repo.DeleteByID(ctx, ID)
}
