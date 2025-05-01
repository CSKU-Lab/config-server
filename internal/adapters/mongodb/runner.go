package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/CSKU-Lab/config-server/domain/cerrors"
	"github.com/CSKU-Lab/config-server/domain/models/runner"
	"github.com/CSKU-Lab/config-server/domain/repositories"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type runnerRepo struct {
	col *mongo.Collection
}

func NewRunnerRepo(db *mongo.Database) repositories.RunnerRepository {
	return &runnerRepo{
		col: db.Collection("runners"),
	}
}

func (l *runnerRepo) Add(ctx context.Context, body *runner.Runner) error {
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

func (l *runnerRepo) GetAll(ctx context.Context) ([]runner.Runner, error) {
	cursor, err := l.col.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("Cannot get runners : %v", err)
	}
	defer cursor.Close(ctx)

	var runners []runner.Runner
	err = cursor.All(ctx, &runners)
	if err != nil {
		return nil, cerrors.New(cerrors.CANNOT_GET_DATA)
	}

	return runners, nil
}

func (l *runnerRepo) GetByID(ctx context.Context, ID string) (*runner.Runner, error) {
	var _runner runner.Runner
	err := l.col.FindOne(ctx, bson.M{"_id": ID}).Decode(&_runner)
	if err != nil {
		return nil, err
	}
	return &_runner, nil
}

func (l *runnerRepo) UpdateByID(ctx context.Context, ID string, body *runner.UpdateRunner) error {
	updatedFields := getUpdatedFields(body)
	_, err := l.col.UpdateOne(ctx, bson.M{"_id": ID}, bson.D{{"$set", updatedFields}})
	if err != nil {
		fmt.Printf("%T", err)
		return err
	}

	return nil
}

func (l *runnerRepo) DeleteByID(ctx context.Context, ID string) error {
	_, err := l.col.DeleteOne(ctx, bson.M{"_id": ID})
	return err
}
