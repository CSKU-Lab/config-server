package models

type Runner struct {
	ID           string `bson:"_id"`
	Name         string `bson:"name"`
	Description  string `bson:"description"`
	BuildScript  string `bson:"build_script"`
	RunScript    string `bson:"run_script"`
	InitialFiles []File `bson:"initial_files"`
}
