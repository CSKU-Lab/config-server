package language

import (
	"regexp"
	"strings"
)

type Language struct {
	ID          string   `bson:"_id"`
	Name        string   `bson:"name"`
	Tags        []string `bson:"tags"`
	BuildScript string   `bson:"build_script,omitempty"`
	RunScript   string   `bson:"run_script"`
	FileNames   []string `bson:"file_names"`
}

type Options struct {
	Name        string
	Tags        []string
	BuildScript string
	RunScript   string
	FileNames   []string
}

func New(opts *Options) *Language {
	id := genID(opts.Name)
	return &Language{
		ID:          id,
		Name:        opts.Name,
		Tags:        opts.Tags,
		BuildScript: opts.BuildScript,
		RunScript:   opts.RunScript,
		FileNames:   opts.FileNames,
	}
}

type UpdateLanguage struct {
	ID          *string   `bson:"_id"`
	Name        *string   `bson:"name"`
	Tags        *[]string `bson:"tags"`
	BuildScript *string   `bson:"build_script,omitempty"`
	RunScript   *string   `bson:"run_script"`
	FileNames   *[]string `bson:"file_names"`
}

type PartialOptions struct {
	Name        *string
	Tags        *[]string
	BuildScript *string
	RunScript   *string
	FileNames   *[]string
}

func NewUpdate(opts *PartialOptions) *UpdateLanguage {
	var id *string = nil
	if opts.Name != nil {
		_id := genID(*opts.Name)
		id = &_id
	}

	return &UpdateLanguage{
		ID:          id,
		Name:        opts.Name,
		Tags:        opts.Tags,
		BuildScript: opts.BuildScript,
		RunScript:   opts.RunScript,
		FileNames:   opts.FileNames,
	}
}

func genID(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	cleanedSpecialChars := re.ReplaceAllString(name, "_")
	trimmed := strings.Trim(cleanedSpecialChars, "_")

	return strings.ToLower(trimmed)
}
