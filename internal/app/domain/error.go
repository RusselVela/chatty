package domain

import "fmt"

// ErrorCode represents error information occurring while processing a request
type ErrorCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Args    []any  `json:"args"`
}

// Clone will return a copy of the ErrorCode with the same Code and Message.  The caller args will be used instead of Args if supplied.
func (ec *ErrorCode) Clone(args ...any) *ErrorCode {
	ecNew := &ErrorCode{
		Code:    ec.Code,
		Message: ec.Message,
	}

	if len(args) > 0 {
		ecNew.Args = make([]any, len(args))
		copy(ecNew.Args, args)
	} else {
		ecNew.Args = make([]any, len(ec.Args))
		copy(ecNew.Args, ec.Args)
	}
	return ecNew
}

// String will return a representation of an error with code and formatted message
func (ec *ErrorCode) String() string {
	return fmt.Sprintf("error %d; ", ec.Code) + fmt.Sprintf(ec.Message, ec.Args...)
}
