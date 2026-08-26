package model

import "errors"

var (
	ErrInvalidRelation = errors.New("invalid entity relation")
	ErrInvalidState    = errors.New("invalid entity state")
)
