package url

import "errors"

var (
	ErrURLNotFound    = errors.New("url not found")
	ErrLongURLExists  = errors.New("long url already exists")
	ErrShortURLExists = errors.New("short url already exists")
)
