package storage

import "errors"

var errDBNotSelected = errors.New("database not selected")
var errDBFileDoesNotExist = errors.New("file does not exist")
var errDBUnknownID = errors.New("no such id in DB")
