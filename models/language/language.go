package languages

import "strings"

type Language struct {
	ID          string `bson:"id"`
	Name        string `bson:"name"`
	Version     string `bson:"version"`
	BuildScript string `bson:"build_script,omitempty"`
	RunScript   string `bson:"run_script"`
}

type Options struct {
	Name        string
	Version     string
	BuildScript string
	RunScript   string
}

func New(opts *Options) *Language {
	lowerName := strings.ToLower(opts.Name)
	id := strings.Join([]string{lowerName, opts.Version}, "_")
	return &Language{
		ID:          id,
		Name:        opts.Name,
		Version:     opts.Version,
		BuildScript: opts.BuildScript,
		RunScript:   opts.RunScript,
	}
}
