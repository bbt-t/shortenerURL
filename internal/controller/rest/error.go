package rest

import "errors"

var errHttpConflict = errors.New("conflict: this URL has already been shortened before")
