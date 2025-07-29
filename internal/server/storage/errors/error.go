package errors

import "errors"

var (
	ErrWrongUsernameOrPassword = errors.New("wrong username or password")
	ErrRecordNotFound          = errors.New("record not found")
	ErrUsernameAlreadyExists   = errors.New("username already exists")
	ErrNoMetadataSet           = errors.New("no metadata set")
	ErrNotValidateToken        = errors.New("not validate token")
	ErrNameAlreadyExists       = errors.New("name already exists")
	ErrFileNotExists           = errors.New("file not exists")
)
