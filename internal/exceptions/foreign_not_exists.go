package exceptions

const ForeignNotFoundExceptionName = "foreign_not_found"

type ForeignNotFoundException struct {
	Name    string
	Message string
}

func NewForeignNotFoundException(msg string) *ForeignNotFoundException {
	return &ForeignNotFoundException{
		Name:    ForeignNotFoundExceptionName,
		Message: msg,
	}
}

func (impl *ForeignNotFoundException) Error() string {
	return impl.Message
}
