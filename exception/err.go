package exception

import "fmt"

type ErrString string

func (e ErrString) New(s interface{}) ErrString {
	return ErrString(fmt.Sprintf("%v, error: %v", e, s))
}

func (e ErrString) Error() string {
	return string(e)
}
