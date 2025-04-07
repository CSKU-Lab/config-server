package services

import (
	"context"

	"github.com/CSKU-Lab/config-server/domain/repositories"
	"github.com/CSKU-Lab/config-server/models/language"
)

type languageService struct {
	repo repositories.LanguageRepository
}

type LanguageService interface {
	Add(ctx context.Context, body *language.Language) error
	GetAll(ctx context.Context) ([]language.Language, error)
	GetByID(ctx context.Context, ID string) (*language.Language, error)
	UpdateByID(ctx context.Context, ID string, body *language.PartialOptions) (*language.Language, error)
	DeleteByID(ctx context.Context, ID string) error
}

func NewLanguageService(repo repositories.LanguageRepository) *languageService {
	return &languageService{
		repo: repo,
	}
}

func (l *languageService) Add(ctx context.Context, body *language.Language) error {
	return l.repo.Add(ctx, body)
}

func (l *languageService) GetAll(ctx context.Context) ([]language.Language, error) {
	return l.repo.GetAll(ctx)
}

func (l *languageService) GetByID(ctx context.Context, ID string) (*language.Language, error) {
	return l.repo.GetByID(ctx, ID)
}

func (l *languageService) UpdateByID(ctx context.Context, ID string, body *language.PartialOptions) (*language.Language, error) {
	lang, err := l.repo.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	if body.Name == nil {
		body.Name = &lang.Name
	}

	if body.Version == nil {
		body.Version = &lang.Version
	}

	modLang := language.NewUpdate(&language.PartialOptions{
		Name:        body.Name,
		Version:     body.Version,
		BuildScript: body.BuildScript,
		RunScript:   body.RunScript,
	})

	err = l.repo.UpdateByID(ctx, ID, modLang)
	if err != nil {
		return nil, err
	}

	return &language.Language{
		ID:          lang.ID,
		Name:        *body.Name,
		Version:     *body.Version,
		BuildScript: lang.BuildScript,
		RunScript:   lang.RunScript,
	}, nil
}

func (l *languageService) DeleteByID(ctx context.Context, ID string) error {
	return l.repo.DeleteByID(ctx, ID)
}
