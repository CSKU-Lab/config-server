package compare

import "strings"

type Compare struct {
	ID          string `bson:"id"`
	Name        string `bson:"name"`
	Script      string `bson:"script"`
	BuildScript string `bson:"build_script"`
	RunScript   string `bson:"run_script"`
}

type Option struct {
	Name        string
	Script      string
	BuildScript string
	RunScript   string
}

func New(option *Option) *Compare {
	return &Compare{
		ID:          strings.ToLower(option.Name),
		Name:        option.Name,
		Script:      option.Script,
		BuildScript: option.BuildScript,
		RunScript:   option.RunScript,
	}
}

type UpdateCompare struct {
	ID          *string `bson:"id"`
	Name        *string `bson:"name"`
	Script      *string `bson:"script"`
	BuildScript *string `bson:"build_script"`
	RunScript   *string `bson:"run_script"`
}

type PartialOption struct {
	Name        *string
	Script      *string
	BuildScript *string
	RunScript   *string
}

func NewUpdate(option *PartialOption) *UpdateCompare {
	id := strings.ToLower(*option.Name)
	return &UpdateCompare{
		ID:          &id,
		Name:        option.Name,
		Script:      option.Script,
		BuildScript: option.BuildScript,
		RunScript:   option.RunScript,
	}
}
