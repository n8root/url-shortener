package entities

type EntityError struct {
	Status  int
	Code    string
	Message string
	Errors  any
}

func (e EntityError) Error() string {
	return e.Message
}

func (e *EntityError) GetStatus() int {
	return e.Status
}

func (e *EntityError) GetCode() string {
	return e.Code
}

func (e *EntityError) GetMessage() string {
	return e.Message
}

func (e *EntityError) GetErrors() any {
	return e.Errors
}
