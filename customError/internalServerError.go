package customError

type InternalServerError struct {
	s string
}

// function create new internal server error
func NewInternalServerError(s string) error {
	return &InternalServerError{s}
}

func (i *InternalServerError) Error() string {
	return i.s
}
