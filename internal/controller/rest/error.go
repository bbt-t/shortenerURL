package rest

import "errors"

var errHTTPConflict = errors.New("conflict: this URL has already been shortened before")
