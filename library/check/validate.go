package check

import (
	exp "ashe/exception"
	valid "gopkg.in/go-playground/validator.v9"
)

var check = valid.New()

func Check(data interface{}, msg exp.ErrString) error {
	if err := check.Struct(data); err != nil {
		return msg.New(err)
	}
	return nil
}
