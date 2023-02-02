package storage

import "errors"

var (
	errDBUnknownID  = errors.New("no such id in DB")
	errDBEmpty      = errors.New("db is empty")
	errHTTPConflict = errors.New("conflict: this URL has already been shortened before")
	errDeleted      = errors.New("deleted")
)
