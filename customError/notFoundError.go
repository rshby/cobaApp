package customError

type NotFoundError struct {
	s string
}

// function new not found error
func NewNotFoundError(s string) error {
	return &NotFoundError{s: s}
}

func (n *NotFoundError) Error() string {
	return n.s
}
