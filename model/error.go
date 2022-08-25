package model

import "fmt"

const (
	MarshalErrorType = "Marshal"
	ConvertErrorType = "Convert"
)

type CustomError struct {
	ErrorType string
	TypeName  string
	Message   string
}

func NewCustomError(errorType, typeName, message string) *CustomError {
	return &CustomError{
		ErrorType: errorType,
		TypeName:  typeName,
		Message:   message,
	}
}

func (c *CustomError) Error() string {
	return fmt.Sprintf("failed to %s %s, error : %s", c.ErrorType, c.TypeName, c.Message)
}
