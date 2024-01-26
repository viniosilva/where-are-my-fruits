package exceptions

const ForbiddenExceptionName = "forbidden"

type ForbiddenException struct {
	Name    string
	Message string
}

func NewForbiddenException(msg string) *ForbiddenException {
	return &ForbiddenException{
		Name:    ForbiddenExceptionName,
		Message: msg,
	}
}

func (impl *ForbiddenException) Error() string {
	return impl.Message
}
