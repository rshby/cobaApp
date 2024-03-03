package customError

type BadRequestError struct {
	s string
}

// function create new bad request error
func NewBadRequestError(s string) error {
	return &BadRequestError{s}
}

func (b *BadRequestError) Error() string {
	return b.s
}
