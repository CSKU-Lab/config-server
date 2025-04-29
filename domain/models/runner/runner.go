package runner

import (
	"regexp"
	"strings"
)

type Runner struct {
	ID          string   `bson:"_id"`
	Name        string   `bson:"name"`
	Tags        []string `bson:"tags"`
	BuildScript string   `bson:"build_script,omitempty"`
	RunScript   string   `bson:"run_script"`
}

type Options struct {
	Name        string
	Tags        []string
	BuildScript string
	RunScript   string
}

func New(opts *Options) *Runner {
	id := genID(opts.Name)
	return &Runner{
		ID:          id,
		Name:        opts.Name,
		Tags:        opts.Tags,
		BuildScript: opts.BuildScript,
		RunScript:   opts.RunScript,
	}
}

type UpdateRunner struct {
	ID          *string   `bson:"_id"`
	Name        *string   `bson:"name"`
	Tags        *[]string `bson:"tags"`
	BuildScript *string   `bson:"build_script,omitempty"`
	RunScript   *string   `bson:"run_script"`
}

type PartialOptions struct {
	Name        *string
	Tags        *[]string
	BuildScript *string
	RunScript   *string
}

func NewUpdate(opts *PartialOptions) *UpdateRunner {
	var id *string = nil
	if opts.Name != nil {
		_id := genID(*opts.Name)
		id = &_id
	}

	return &UpdateRunner{
		ID:          id,
		Name:        opts.Name,
		Tags:        opts.Tags,
		BuildScript: opts.BuildScript,
		RunScript:   opts.RunScript,
	}
}

func genID(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	cleanedSpecialChars := re.ReplaceAllString(name, "_")
	trimmed := strings.Trim(cleanedSpecialChars, "_")

	return strings.ToLower(trimmed)
}
