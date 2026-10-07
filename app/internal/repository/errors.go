package repository

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrConflict            = errors.New("conflict")
	ErrAlreadyExists       = errors.New("already exists")
	ErrUnknownMaterialType = errors.New("unknown material type")
)
