package models

type EntityError struct {
	Status  int
	Message string
	Errors  any
}

func (e EntityError) Error() string {
	return e.Message
}

func (e EntityError) GetStatus() int {
	return e.Status
}

func (e EntityError) GetMessage() string {
	return e.Message
}

func (e EntityError) GetErrors() any {
	return e.Errors
}
