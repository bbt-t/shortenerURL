package rest

import "errors"

var ErrHTTPConflict = errors.New("conflict: this URL has already been shortened before")
