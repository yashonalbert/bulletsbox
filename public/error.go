package public

import (
	"errors"
)

// ConnError indicates that a name was malformed and the specific error
type ConnError struct {
	Name string
	Err  error
}

func (e ConnError) Error() string {
	return e.Name + ": " + e.Err.Error()
}

var (
	errUnknown    = errors.New("unknown error")
	errEmpty      = errors.New("name is empty")
	errTooLong    = errors.New("name is too long")
	errSystem     = errors.New("system error")
	errUnknownCmd = errors.New("unknown command")
	errFail       = errors.New("Fail")
)

// CheckName func
func CheckName(s string) error {
	switch {
	case len(s) == 0:
		return ConnError{"NameError", errEmpty}
	case len(s) >= 256:
		return ConnError{"NameError", errTooLong}
	}
	return nil
}

var resError = map[uint8]error{
	ResSystemErr: errSystem,
	ResUnknowCmd: errUnknownCmd,
	ResFail: errFail,
}

// ParseResError func
func ParseResError(resCode uint8) error {
	if _, exists := resError[resCode]; exists {
		return ConnError{"ResponseError", resError[resCode]}
	}
	return ConnError{"ResponseError", errUnknown}
}
