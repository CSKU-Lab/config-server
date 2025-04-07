package cerrors

import "errors"

type CError string

const (
	DUPLICATE_DATA  CError = "duplicate data"
	CANNOT_GET_DATA CError = "cannot get the data"
	UNKNOWN_ERROR   CError = "unknown error"
)

func New(err CError) error {
	return errors.New(string(err))
}
