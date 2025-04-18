package language

import "strings"

type Language struct {
	ID          string   `bson:"id"`
	Name        string   `bson:"name"`
	Version     string   `bson:"version"`
	BuildScript string   `bson:"build_script,omitempty"`
	RunScript   string   `bson:"run_script"`
	FileNames   []string `bson:"file_names"`
}

type Options struct {
	Name        string
	Version     string
	BuildScript string
	RunScript   string
	FileNames   []string
}

func New(opts *Options) *Language {
	id := genID(opts.Name, opts.Version)
	return &Language{
		ID:          id,
		Name:        opts.Name,
		Version:     opts.Version,
		BuildScript: opts.BuildScript,
		RunScript:   opts.RunScript,
		FileNames:   opts.FileNames,
	}
}

type UpdateLanguage struct {
	ID          *string   `bson:"id"`
	Name        *string   `bson:"name"`
	Version     *string   `bson:"version"`
	BuildScript *string   `bson:"build_script,omitempty"`
	RunScript   *string   `bson:"run_script"`
	FileNames   *[]string `bson:"file_names"`
}

type PartialOptions struct {
	Name        *string
	Version     *string
	BuildScript *string
	RunScript   *string
	FileNames   *[]string
}

func NewUpdate(opts *PartialOptions) *UpdateLanguage {
	var id *string = nil
	if opts.Name != nil && opts.Version != nil {
		_id := genID(*opts.Name, *opts.Version)
		id = &_id
	}

	return &UpdateLanguage{
		ID:          id,
		Name:        opts.Name,
		Version:     opts.Version,
		BuildScript: opts.BuildScript,
		RunScript:   opts.RunScript,
		FileNames:   opts.FileNames,
	}
}

func genID(name, version string) string {
	id := ""
	lowerName := strings.ToLower(name)
	id = strings.Join([]string{lowerName, version}, "_")

	return id
}
