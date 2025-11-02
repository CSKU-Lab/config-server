package models

type Runner struct {
	ID          string   `bson:"_id"`
	Name        string   `bson:"name"`
	Tags        []string `bson:"tags"`
	BuildScript string   `bson:"build_script,omitempty"`
	RunScript   string   `bson:"run_script"`
}
