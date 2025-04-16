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
	Add(ctx context.Context, body *compare.Compare) error
	GetAll(ctx context.Context) ([]compare.Compare, error)
	GetByID(ctx context.Context, ID string) (*compare.Compare, error)
	UpdateByID(ctx context.Context, ID string, body *compare.PartialOption) (*compare.UpdateCompare, error)
	DeleteByID(ctx context.Context, ID string) error
}

func NewCompareService(repo repositories.CompareRepository) CompareService {
	return &compareService{
		repo: repo,
	}
}

func (c *compareService) Add(ctx context.Context, body *compare.Compare) error {
	return c.repo.Add(ctx, body)
}

func (c *compareService) GetAll(ctx context.Context) ([]compare.Compare, error) {
	return c.repo.GetAll(ctx)
}

func (c *compareService) GetByID(ctx context.Context, ID string) (*compare.Compare, error) {
	return c.repo.GetByID(ctx, ID)
}

func (c *compareService) UpdateByID(ctx context.Context, ID string, body *compare.PartialOption) (*compare.UpdateCompare, error) {
	oldCompare, err := c.repo.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	updated := compare.NewUpdate(body)
	err = c.repo.UpdateByID(ctx, ID, updated)
	if err != nil {
		return nil, err
	}

	if updated.Name == nil {
		updated.ID = &oldCompare.ID
		updated.Name = &oldCompare.Name
	}

	if updated.Script == nil {
		updated.Script = &oldCompare.Script
	}

	if updated.ScriptName == nil {
		updated.ScriptName = &oldCompare.ScriptName
	}

	if updated.BuildScript == nil {
		updated.BuildScript = &oldCompare.BuildScript
	}

	if updated.RunScript == nil {
		updated.RunScript = &oldCompare.RunScript
	}

	return updated, nil
}

func (c *compareService) DeleteByID(ctx context.Context, ID string) error {
	return c.repo.DeleteByID(ctx, ID)
}
