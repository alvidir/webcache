package webcache

import (
	"errors"
)

var (
	ErrUnknown   = errors.New("E001")
	ErrNotFound  = errors.New("E002")
	ErrNoContent = errors.New("E003")
)
