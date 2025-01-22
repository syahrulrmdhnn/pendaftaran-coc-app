package exception

type BadRequestHandler struct {
	Message string
}

func NewBadRequestHandler(message string) error {
	return BadRequestHandler{Message: message}
}

func (e BadRequestHandler) Error() string {
	return e.Message
}