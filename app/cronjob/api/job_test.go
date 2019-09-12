package api

import (
	"ashe/app/cronjob/models"
	"ashe/protocol"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"testing"
)

func TestJob_Add(t *testing.T) {
	jb := new(Job)
	out := models.AsheJob{}
	err := jb.unmarshalData(&protocol.Request{RequestJson: `{
	"pageSize": "heheh"
}`}, &out)
	validate := validator.New()
	err = validate.Struct(out)
	fmt.Println(err)
	fmt.Println(out)

}
