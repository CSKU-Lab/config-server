package models

type Compare struct {
	ID          string `bson:"_id"`
	Name        string `bson:"name"`
	Files       []File `bson:"files"`
	BuildScript string `bson:"build_script"`
	RunScript   string `bson:"run_script"`
	RunName     string `bson:"run_name"`
	Description string `bson:"description"`
}
