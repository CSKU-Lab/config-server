package services

import (
	"context"

	"github.com/CSKU-Lab/config-server/domain/models/language"
	"github.com/CSKU-Lab/config-server/domain/repositories"
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

	buildScript := lang.BuildScript
	if body.BuildScript != nil {
		buildScript = *body.BuildScript
	}

	runScript := lang.BuildScript
	if body.RunScript != nil {
		runScript = *body.RunScript
	}

	return &language.Language{
		ID:          *modLang.ID,
		Name:        *modLang.Name,
		Version:     *modLang.Version,
		BuildScript: buildScript,
		RunScript:   runScript,
	}, nil
}

func (l *languageService) DeleteByID(ctx context.Context, ID string) error {
	return l.repo.DeleteByID(ctx, ID)
}
