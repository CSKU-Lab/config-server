package requests

import "github.com/CSKU-Lab/config-server/domain/models"

type CreateRunner struct {
	Name        string
	Description string
}

type UpdateRunner struct {
	Name         *string       `bson:"name"`
	Description  *string       `bson:"description"`
	BuildScript  *string       `bson:"build_script"`
	RunScript    *string       `bson:"run_script"`
	InitialFiles []models.File `bson:"initial_files"`
}
