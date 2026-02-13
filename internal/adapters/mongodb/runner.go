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

type runnerRepo struct {
	col *mongo.Collection
}

type runnerDoc struct {
	ID          string   `bson:"_id"`
	Name        string   `bson:"name"`
	Tags        []string `bson:"tags"`
	BuildScript string   `bson:"build_script"`
	RunScript   string   `bson:"run_script"`
}

func NewRunnerRepo(db *mongo.Database) repositories.RunnerRepository {
	return &runnerRepo{
		col: db.Collection("runners"),
	}
}

func (l *runnerRepo) Create(ctx context.Context, ID string, body *requests.CreateRunner) error {
	runner := &runnerDoc{
		ID:          ID,
		Name:        body.Name,
		Tags:        body.Tags,
		BuildScript: body.BuildScript,
		RunScript:   body.RunScript,
	}
	_, err := l.col.InsertOne(ctx, runner)
	if err != nil {
		var mongoErr mongo.WriteException
		if errors.As(err, &mongoErr) {
			return cerrors.New(cerrors.DUPLICATE_DATA)
		}
		return cerrors.New(cerrors.UNKNOWN_ERROR)
	}
	return nil
}

func (l *runnerRepo) GetAll(ctx context.Context) ([]models.Runner, error) {
	cursor, err := l.col.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("Cannot get runners : %v", err)
	}
	defer cursor.Close(ctx)

	var runners []models.Runner
	err = cursor.All(ctx, &runners)
	if err != nil {
		return nil, cerrors.New(cerrors.CANNOT_GET_DATA)
	}

	return runners, nil
}

func (l *runnerRepo) GetPagination(ctx context.Context, req *requests.GetPagination) ([]models.Runner, error) {
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

	cursor, err := l.col.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var runners []models.Runner
	err = cursor.All(ctx, &runners)
	if err != nil {
		return nil, cerrors.New(cerrors.CANNOT_GET_DATA)
	}

	return runners, nil
}

func (l *runnerRepo) Count(ctx context.Context) (int, error) {
	count, err := l.col.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (l *runnerRepo) GetByID(ctx context.Context, ID string) (*models.Runner, error) {
	var _runner models.Runner
	err := l.col.FindOne(ctx, bson.M{"_id": ID}).Decode(&_runner)
	if err != nil {
		return nil, err
	}
	return &_runner, nil
}

func (l *runnerRepo) UpdateByID(ctx context.Context, ID string, body *requests.UpdateRunner) error {
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
