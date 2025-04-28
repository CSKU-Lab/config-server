package compare

import (
	"strings"

	"github.com/CSKU-Lab/config-server/domain/models/file"
)

type Compare struct {
	ID          string      `bson:"_id"`
	Name        string      `bson:"name"`
	Files       []file.File `bson:"files"`
	Script      string      `bson:"script"`
	BuildScript string      `bson:"build_script"`
	RunScript   string      `bson:"run_script"`
	RunName     string      `bson:"run_name"`
	Description string      `bson:"description"`
}

type Option struct {
	Name        string
	Script      string
	Files       []file.File
	BuildScript string
	RunScript   string
	RunName     string
	Description string
}

func New(option *Option) *Compare {
	var id string
	if option.Name != "" {
		splitted := strings.Split(option.Name, " ")
		for i := range len(splitted) {
			splitted[i] = strings.ToLower(splitted[i])
		}
		id = strings.Join(splitted, "_")
	}

	return &Compare{
		ID:          id,
		Name:        option.Name,
		Script:      option.Script,
		Files:       option.Files,
		BuildScript: option.BuildScript,
		RunScript:   option.RunScript,
		RunName:     option.RunName,
		Description: option.Description,
	}
}

type UpdateCompare struct {
	ID          *string     `bson:"_id"`
	Name        *string     `bson:"name"`
	Script      *string     `bson:"script"`
	Files       []file.File `bson:"files"`
	BuildScript *string     `bson:"build_script"`
	RunScript   *string     `bson:"run_script"`
	RunName     *string     `bson:"run_name"`
	Description *string     `bson:"description"`
}

type PartialOption struct {
	Name        *string
	Script      *string
	Files       []file.File
	BuildScript *string
	RunScript   *string
	RunName     *string
	Description *string
}

func NewUpdate(option *PartialOption) *UpdateCompare {
	var id *string = nil
	if option.Name != nil {
		_id := generateID(*option.Name)
		id = &_id
	}

	return &UpdateCompare{
		ID:          id,
		Name:        option.Name,
		Script:      option.Script,
		Files:       option.Files,
		BuildScript: option.BuildScript,
		RunScript:   option.RunScript,
		RunName:     option.RunName,
		Description: option.Description,
	}
}

func generateID(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "_"))
}
