package services

import (
	"context"

	"github.com/CSKU-Lab/config-server/domain/models/compare"
	"github.com/CSKU-Lab/config-server/domain/repositories"
)

type compareService struct {
	repo repositories.CompareRepository
}

type CompareService interface {
	Add(ctx context.Context, body *compare.Option) (*compare.Compare, error)
	GetAll(ctx context.Context) ([]compare.Compare, error)
	GetByID(ctx context.Context, ID string) (*compare.Compare, error)
	UpdateByID(ctx context.Context, ID string, body *compare.PartialOption) (*compare.Compare, error)
	DeleteByID(ctx context.Context, ID string) error
}

func NewCompareService(repo repositories.CompareRepository) CompareService {
	return &compareService{
		repo: repo,
	}
}

func (c *compareService) Add(ctx context.Context, body *compare.Option) (*compare.Compare, error) {
	compare := compare.New(&compare.Option{
		Name:        body.Name,
		Script:      body.Script,
		ScriptName:  body.ScriptName,
		BuildScript: body.BuildScript,
		RunScript:   body.RunScript,
		RunName:     body.RunName,
		Description: body.Description,
	})

	return compare, c.repo.Add(ctx, compare)
}

func (c *compareService) GetAll(ctx context.Context) ([]compare.Compare, error) {
	return c.repo.GetAll(ctx)
}

func (c *compareService) GetByID(ctx context.Context, ID string) (*compare.Compare, error) {
	return c.repo.GetByID(ctx, ID)
}

func (c *compareService) UpdateByID(ctx context.Context, ID string, body *compare.PartialOption) (*compare.Compare, error) {
	updatedFields := compare.NewUpdate(&compare.PartialOption{
		Name:        body.Name,
		Script:      body.Script,
		ScriptName:  body.ScriptName,
		BuildScript: body.BuildScript,
		RunScript:   body.RunScript,
		RunName:     body.RunName,
		Description: body.Description,
	})

	err := c.repo.UpdateByID(ctx, ID, updatedFields)
	if err != nil {
		return nil, err
	}

	id := ID
	if updatedFields.ID != nil {
		id = *updatedFields.ID
	}

	updated, err := c.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (c *compareService) DeleteByID(ctx context.Context, ID string) error {
	return c.repo.DeleteByID(ctx, ID)
}
