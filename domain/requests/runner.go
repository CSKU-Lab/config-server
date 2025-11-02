package requests

type CreateRunner struct {
	Name        string
	Tags        []string
	BuildScript string
	RunScript   string
}

type UpdateRunner struct {
	Name        *string   `bson:"name"`
	Tags        *[]string `bson:"tags"`
	BuildScript *string   `bson:"build_script,omitempty"`
	RunScript   *string   `bson:"run_script"`
}
