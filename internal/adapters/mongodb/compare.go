package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/CSKU-Lab/config-server/domain/cerrors"
	"github.com/CSKU-Lab/config-server/domain/models"
	"github.com/CSKU-Lab/config-server/domain/repositories"
	"github.com/CSKU-Lab/config-server/domain/requests"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type compareRepo struct {
	col        *mongo.Collection
	dynamicCol func(path string) *mongo.Collection
}

type compareDoc struct {
	ID          string        `bson:"_id"`
	Name        string        `bson:"name"`
	Files       []models.File `bson:"files"`
	BuildScript string        `bson:"build_script"`
	RunScript   string        `bson:"run_script"`
	RunName     string        `bson:"run_name"`
	Description string        `bson:"description"`
}

func NewCompareRepo(db *mongo.Database) repositories.CompareRepository {
	return &compareRepo{
		col: db.Collection("compares"),
		dynamicCol: func(path string) *mongo.Collection {
			return db.Collection(fmt.Sprintf("compares/%s", path))
		},
	}
}

func (c *compareRepo) GetPagination(ctx context.Context, req *requests.GetPagination) ([]models.Compare, error) {
	orderMap := map[string]int{
		"desc": -1,
		"asc":  1,
	}
	order, ok := orderMap[req.SortOrder]
	if !ok {
		order = -1
	}
	opts := options.Find().
		SetSkip(int64((req.Page - 1) * req.PageSize)).
		SetLimit(int64(req.PageSize)).
		SetSort(bson.D{{Key: "name", Value: order}})

	cursor, err := c.col.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var compares []models.Compare
	err = cursor.All(ctx, &compares)
	if err != nil {
		return nil, cerrors.New(cerrors.CANNOT_GET_DATA)
	}

	return compares, nil
}

func (c *compareRepo) Count(ctx context.Context) (int, error) {
	count, err := c.col.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (c *compareRepo) Create(ctx context.Context, ID string, body *requests.CreateCompare) error {
	compare := &compareDoc{
		ID:          ID,
		Name:        body.Name,
		Files:       body.Files,
		BuildScript: body.BuildScript,
		RunScript:   body.RunScript,
		RunName:     body.RunName,
		Description: body.Description,
	}

	_, err := c.col.InsertOne(ctx, compare)
	if err != nil {
		var mongoErr mongo.WriteException
		if errors.As(err, &mongoErr) {
			return cerrors.New(cerrors.DUPLICATE_DATA)
		}
		return cerrors.New(cerrors.UNKNOWN_ERROR)
	}
	return nil
}

func (c *compareRepo) GetAll(ctx context.Context) ([]models.Compare, error) {
	cursor, err := c.col.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("Cannot get compares : %v", err)
	}
	defer cursor.Close(ctx)

	var compares []models.Compare
	err = cursor.All(ctx, &compares)
	if err != nil {
		return nil, cerrors.New(cerrors.CANNOT_GET_DATA)
	}

	return compares, nil
}

func (c *compareRepo) GetByID(ctx context.Context, ID string) (*models.Compare, error) {
	var compare models.Compare
	err := c.col.FindOne(ctx, bson.M{"_id": ID}).Decode(&compare)
	if err != nil {
		return nil, cerrors.New(cerrors.CANNOT_GET_DATA)
	}
	return &compare, nil
}

func (c *compareRepo) UpdateByID(ctx context.Context, ID string, body *requests.UpdateCompare) error {
	updatedFields := getUpdatedFields(body)
	_, err := c.col.UpdateOne(ctx, bson.M{"_id": ID}, bson.D{{Key: "$set", Value: updatedFields}})
	if err != nil {
		return err
	}

	return nil
}

func (c *compareRepo) DeleteByID(ctx context.Context, ID string) error {
	_, err := c.col.DeleteOne(ctx, bson.M{"_id": ID})
	return err
}
