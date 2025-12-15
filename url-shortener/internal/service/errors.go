package service

import (
	"errors"
)

var (
	ErrURLCollision = errors.New("URL collision")
	ErrURLNotFound  = errors.New("URL not found")
	ErrStorage      = errors.New("storage failure")
)
