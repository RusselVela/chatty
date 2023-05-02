package service

import (
	"fmt"
)

// ErrorCode represents error information occurring while processing a request
type ErrorCode struct {
	Status  int
	Code    int    `json:"code"`
	Message string `json:"message"`
	Args    []any  `json:"args"`
}

// Error interface for errors
func (ec *ErrorCode) Error() string {
	return ec.String()
}

// Clone will return a copy of the ErrorCode with the same Code and Message.  The caller args will be used instead of Args if supplied.
func (ec *ErrorCode) Clone(args ...any) *ErrorCode {
	ecNew := &ErrorCode{
		Status:  ec.Status,
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
