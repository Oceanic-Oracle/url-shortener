package dto

import "errors"

var (
	ErrParseURL               = errors.New("invalid URL format")
	ErrURLMissingSchemaOrHost = errors.New("URL schema and host is required")
	ErrURLUnsupportedScheme   = errors.New("only http and https schemes are allowed")
	ErrURLPointsToService     = errors.New("URL must not point to this service")
)
