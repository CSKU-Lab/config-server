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
	Description string `bson:"description"`
}

type Option struct {
	Name        string
	Script      string
	ScriptName  string
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
		ScriptName:  option.ScriptName,
		BuildScript: option.BuildScript,
		RunScript:   option.RunScript,
		RunName:     option.RunName,
		Description: option.Description,
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
	Description *string `bson:"description"`
}

type PartialOption struct {
	Name        *string
	Script      *string
	ScriptName  *string
	BuildScript *string
	RunScript   *string
	RunName     *string
	Description *string
}

func NewUpdate(option *PartialOption) *UpdateCompare {
	var id *string = nil
	if option.Name != nil {
		splitted := strings.Split(*option.Name, " ")
		for i := range len(splitted) {
			splitted[i] = strings.ToLower(splitted[i])
		}
		_id := strings.Join(splitted, "_")
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
		Description: option.Description,
	}
}
