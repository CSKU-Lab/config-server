package repositories

import (
	"context"

	"github.com/CSKU-Lab/config-server/models/language"
)

type LanguageRepository interface {
	Add(ctx context.Context, body *language.Language) error
	GetAll(ctx context.Context) ([]language.Language, error)
	GetByID(ctx context.Context, ID string) (*language.Language, error)
	UpdateByID(ctx context.Context, ID string, body *language.UpdateLanguage) error
	DeleteByID(ctx context.Context, ID string) error
}
