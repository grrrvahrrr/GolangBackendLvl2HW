package enteties

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

type User struct {
	Login    string
	Password string
	Age      string
	Sex      string
}
