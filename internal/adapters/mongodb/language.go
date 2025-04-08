package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/CSKU-Lab/config-server/domain/cerrors"
	"github.com/CSKU-Lab/config-server/domain/repositories"
	"github.com/CSKU-Lab/config-server/domain/models/language"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type languageRepo struct {
	col *mongo.Collection
}

func NewLanguageRepo(db *mongo.Database) repositories.LanguageRepository {
	return &languageRepo{
		col: db.Collection("languages"),
	}
}

func (l *languageRepo) Add(ctx context.Context, body *language.Language) error {
	_, err := l.col.InsertOne(ctx, body)
	if err != nil {
		var mongoErr mongo.WriteException
		if errors.As(err, &mongoErr) {
			return cerrors.New(cerrors.DUPLICATE_DATA)
		}
		return cerrors.New(cerrors.UNKNOWN_ERROR)
	}
	return nil
}

func (l *languageRepo) GetAll(ctx context.Context) ([]language.Language, error) {
	cursor, err := l.col.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("Cannot get languages : %v", err)
	}

	var langauges []language.Language
	err = cursor.All(ctx, &langauges)
	if err != nil {
		return nil, cerrors.New(cerrors.CANNOT_GET_DATA)
	}

	return langauges, nil
}

func (l *languageRepo) GetByID(ctx context.Context, ID string) (*language.Language, error) {
	var lang language.Language
	err := l.col.FindOne(ctx, bson.M{"id": ID}).Decode(&lang)
	if err != nil {
		return nil, err
	}
	return &lang, nil
}

func (l *languageRepo) UpdateByID(ctx context.Context, ID string, body *language.UpdateLanguage) error {
	updatedFields := getUpdatedFields(body)
	fmt.Println(updatedFields)
	_, err := l.col.UpdateOne(ctx, bson.M{"id": ID}, bson.D{{"$set", updatedFields}})
	if err != nil {
		fmt.Printf("%T", err)
		return err
	}

	return nil
}

func (l *languageRepo) DeleteByID(ctx context.Context, ID string) error {
	_, err := l.col.DeleteOne(ctx, bson.M{"id": ID})
	return err
}
