package error_util

import "fmt"

type CustomError struct {
	Message    string
	HttpStatus int
	ErrorName  string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("Error : %s", e.Message)
}
