package storage

import "errors"

var (
	ErrDBUnknownID  = errors.New("no such id in DB")
	ErrDBEmpty      = errors.New("db is empty")
	ErrHTTPConflict = errors.New("conflict: this URL has already been shortened before")
	ErrDeleted      = errors.New("deleted")
)
