package requests

import "github.com/CSKU-Lab/config-server/domain/models"

type CreateCompare struct {
	Name        string
	Files       []models.File
	BuildScript string
	RunScript   string
	RunName     string
	Description string
}

type UpdateCompare struct {
	Name        *string       `bson:"name"`
	Files       []models.File `bson:"files"`
	BuildScript *string       `bson:"build_script"`
	RunScript   *string       `bson:"run_script"`
	RunName     *string       `bson:"run_name"`
	Description *string       `bson:"description"`
}
