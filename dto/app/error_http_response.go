package app

type ErrorHttpResponse struct {
	HttpStatus int
	ErrorName  string
	Message    string
	Data       interface{}
}
