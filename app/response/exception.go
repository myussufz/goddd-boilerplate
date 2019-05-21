package response

import (
	"bytes"
	"fmt"

	"goddd-boilerplate/app/response/errcode"
)

// Exception :
type Exception struct {
	Code   string
	Detail string
	Error  error
}

// MarshalJSON :
func (e Exception) MarshalJSON() ([]byte, error) {
	bb := new(bytes.Buffer)
	bb.WriteString(`{`)
	bb.WriteString(`"error":{`)
	bb.WriteString(`"code":`)
	bb.WriteString(fmt.Sprintf("%q", fmt.Sprintf("%s", e.Code)))
	message, isExist := errcode.Message.Load(e.Code)
	if isExist {
		bb.WriteString(",")
		bb.WriteString(`"message":`)
		bb.WriteString(fmt.Sprintf("%q", message))
	}
	if e.Detail != "" {
		bb.WriteString(",")
		bb.WriteString(`"detail":`)
		bb.WriteString(fmt.Sprintf("%q", e.Detail))
	}
	if e.Error != nil {
		bb.WriteString(",")
		bb.WriteString(`"stackTrace":`)
		bb.WriteString(fmt.Sprintf("%q", fmt.Sprintf("%+v", e.Error)))
	}
	bb.WriteString(`}`)
	bb.WriteString(`}`)

	return bb.Bytes(), nil
}
