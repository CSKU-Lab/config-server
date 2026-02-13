package services

import (
	"context"

	"github.com/CSKU-Lab/config-server/domain/models"
	"github.com/CSKU-Lab/config-server/domain/repositories"
	"github.com/CSKU-Lab/config-server/domain/requests"
	"github.com/google/uuid"
)

type compareService struct {
	repo repositories.CompareRepository
}

type CompareService interface {
	Create(ctx context.Context, body *requests.CreateCompare) (string, error)
	GetAll(ctx context.Context) ([]models.Compare, error)
	GetPagination(ctx context.Context, req *requests.GetPagination) ([]models.Compare, int, error)
	GetByID(ctx context.Context, ID string) (*models.Compare, error)
	UpdateByID(ctx context.Context, ID string, body *requests.UpdateCompare) error
	DeleteByID(ctx context.Context, ID string) error
}

func NewCompareService(repo repositories.CompareRepository) CompareService {
	return &compareService{
		repo: repo,
	}
}

func (c *compareService) GetPagination(ctx context.Context, req *requests.GetPagination) ([]models.Compare, int, error) {
	pagination, err := c.repo.GetPagination(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	count, err := c.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return pagination, count, nil
}

func (c *compareService) Create(ctx context.Context, body *requests.CreateCompare) (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	err = c.repo.Create(ctx, id.String(), body)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (c *compareService) GetAll(ctx context.Context) ([]models.Compare, error) {
	return c.repo.GetAll(ctx)
}

func (c *compareService) GetByID(ctx context.Context, ID string) (*models.Compare, error) {
	return c.repo.GetByID(ctx, ID)
}

func (c *compareService) UpdateByID(ctx context.Context, ID string, body *requests.UpdateCompare) error {
	return c.repo.UpdateByID(ctx, ID, body)
}

func (c *compareService) DeleteByID(ctx context.Context, ID string) error {
	return c.repo.DeleteByID(ctx, ID)
}
