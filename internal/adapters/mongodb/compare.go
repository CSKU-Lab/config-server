package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/CSKU-Lab/config-server/domain/cerrors"
	"github.com/CSKU-Lab/config-server/domain/models/compare"
	"github.com/CSKU-Lab/config-server/domain/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type compareRepo struct {
	col *mongo.Collection
}

func NewCompareRepo(db *mongo.Database) repositories.CompareRepository {
	return &compareRepo{
		col: db.Collection("compares"),
	}
}

func (c *compareRepo) Add(ctx context.Context, body *compare.Compare) error {
	_, err := c.col.InsertOne(ctx, body)
	if err != nil {
		var mongoErr mongo.WriteException
		if errors.As(err, &mongoErr) {
			return cerrors.New(cerrors.DUPLICATE_DATA)
		}
		return cerrors.New(cerrors.UNKNOWN_ERROR)
	}
	return nil
}

func (c *compareRepo) GetAll(ctx context.Context) ([]compare.Compare, error) {
	cursor, err := c.col.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("Cannot get compares : %v", err)
	}

	var compares []compare.Compare
	err = cursor.All(ctx, &compares)
	if err != nil {
		return nil, cerrors.New(cerrors.CANNOT_GET_DATA)
	}

	return compares, nil
}

func (c *compareRepo) GetByID(ctx context.Context, ID string) (*compare.Compare, error) {
	var compare compare.Compare
	err := c.col.FindOne(ctx, bson.M{"id": ID}).Decode(&compare)
	if err != nil {
		return nil, cerrors.New(cerrors.CANNOT_GET_DATA)
	}
	return &compare, nil
}

func (c *compareRepo) UpdateByID(ctx context.Context, ID string, body *compare.UpdateCompare) error {
	updatedFields := getUpdatedFields(body)
	_, err := c.col.UpdateOne(ctx, bson.M{"id": ID}, bson.D{{Key: "$set", Value: updatedFields}})
	if err != nil {
		fmt.Printf("%T", err)
		return err
	}

	return nil
}

func (c *compareRepo) DeleteByID(ctx context.Context, ID string) error {
	_, err := c.col.DeleteOne(ctx, bson.M{"id": ID})
	return err
}
