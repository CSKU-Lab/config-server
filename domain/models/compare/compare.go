package compare

import "strings"

type Compare struct {
	ID          string `bson:"id"`
	Name        string `bson:"name"`
	ScriptName  string `bson:"script_name"`
	Script      string `bson:"script"`
	BuildScript string `bson:"build_script"`
	RunScript   string `bson:"run_script"`
	RunName     string `bson:"run_name"`
}

type Option struct {
	Name        string
	Script      string
	ScriptName  string
	BuildScript string
	RunScript   string
	RunName     string
}

func New(option *Option) *Compare {
	return &Compare{
		ID:          strings.ToLower(option.Name),
		Name:        option.Name,
		Script:      option.Script,
		ScriptName:  option.ScriptName,
		BuildScript: option.BuildScript,
		RunScript:   option.RunScript,
		RunName:     option.RunName,
	}
}

type UpdateCompare struct {
	ID          *string `bson:"id"`
	Name        *string `bson:"name"`
	Script      *string `bson:"script"`
	ScriptName  *string `bson:"script_name"`
	BuildScript *string `bson:"build_script"`
	RunScript   *string `bson:"run_script"`
	RunName     *string `bson:"run_name"`
}

type PartialOption struct {
	Name        *string
	Script      *string
	ScriptName  *string
	BuildScript *string
	RunScript   *string
	RunName     *string
}

func NewUpdate(option *PartialOption) *UpdateCompare {
	var id *string = nil
	if option.Name != nil {
		_id := strings.ToLower(*option.Name)
		id = &_id
	}

	return &UpdateCompare{
		ID:          id,
		Name:        option.Name,
		Script:      option.Script,
		ScriptName:  option.ScriptName,
		BuildScript: option.BuildScript,
		RunScript:   option.RunScript,
		RunName:     option.RunName,
	}
}
