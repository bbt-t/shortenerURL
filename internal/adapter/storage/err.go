package storage

import "errors"

var (
	errDBNotSelected      = errors.New("database not selected")
	errDBFileDoesNotExist = errors.New("file does not exist")
	errDBUnknownID        = errors.New("no such id in DB")
	errDBEmpty            = errors.New("db is empty")
	errHTTPConflict       = errors.New("conflict: this URL has already been shortened before")
)
