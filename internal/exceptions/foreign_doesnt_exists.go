package exceptions

const ForeignDoesntExistsExceptionName = "foreign_doesnt_exists"

type ForeignDoesntExistsException struct {
	Name    string
	Message string
}

func NewForeignDoesntExistsException(msg string) *ForeignDoesntExistsException {
	return &ForeignDoesntExistsException{
		Name:    ForeignDoesntExistsExceptionName,
		Message: msg,
	}
}

func (impl *ForeignDoesntExistsException) Error() string {
	return impl.Message
}
