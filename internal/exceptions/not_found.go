package exceptions

const NotFoundExceptionName = "not_found"

type NotFoundException struct {
	Name    string
	Message string
}

func NewNotFoundException(msg string) *NotFoundException {
	return &NotFoundException{
		Name:    NotFoundExceptionName,
		Message: msg,
	}
}

func (impl *NotFoundException) Error() string {
	return impl.Message
}
