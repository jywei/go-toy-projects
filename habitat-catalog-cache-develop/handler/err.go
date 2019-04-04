package handler

import "net/http"

const (
	// SuccessCode means request success = 1000
	SuccessCode = iota + 1000
	// ServerInternalErrCode means 500 internal server error = 1001
	ServerInternalErrCode
	// InvalidAttributeErrCode means 400 bad request = 1002
	InvalidAttributeErrCode
	// RecordNotFoundErrCode means 404 not found = 1003
	RecordNotFoundErrCode
	// MissingParametersErrCode means 400 bad request = 1004
	MissingParametersErrCode
	// UnauthorizedErrCode means 401 unauthorized = 1005
	UnauthorizedErrCode
)

const (
	// SuccessMsg is the SuccessCode message
	SuccessMsg = "Success"
	// ServerInternalErrMsg is the ServerInternalErrCode message
	ServerInternalErrMsg = "Internal Server Error"
	// InvalidAttributeErrMsg is the InvalidAttributeErrCode message
	InvalidAttributeErrMsg = "You passed an invalid value for the attributes."
	// RecordNotFoundErrMsg is the RecordNotFoundErrCode message
	RecordNotFoundErrMsg = "Record Not Found"
	// MissingParametersErrMsg is the MissingParametersErrCode message
	MissingParametersErrMsg = "Missing Parameters"
	// UnauthorizedErrMsg is the UnauthorizedErrCode message
	UnauthorizedErrMsg = "Unauthorized"
)

// Error represents an error with an associated ExternalAPI status code.
type Error struct {
	internalErr error
	Status      int    `json:"-"`
	OutputErr   string `json:"error"`
}

// NewErr returns a Error instance.
func NewErr(code int, err error) *Error {
	e := &Error{
		internalErr: err,
	}
	e.setup(code)
	return e
}

// Error allows handler Error struct to satisfy the build-in error interface.
func (e *Error) Error() string {
	return e.internalErr.Error()
}

func (e *Error) setup(code int) {
	switch code {
	case SuccessCode:
		e.Status = http.StatusOK
		e.OutputErr = SuccessMsg
	case InvalidAttributeErrCode:
		e.Status = http.StatusBadRequest
		e.OutputErr = InvalidAttributeErrMsg
	case RecordNotFoundErrCode:
		e.Status = http.StatusNotFound
		e.OutputErr = RecordNotFoundErrMsg
	case MissingParametersErrCode:
		e.Status = http.StatusBadRequest
		e.OutputErr = MissingParametersErrMsg
	case UnauthorizedErrCode:
		e.Status = http.StatusUnauthorized
		e.OutputErr = UnauthorizedErrMsg
	default:
		e.Status = http.StatusInternalServerError
		e.OutputErr = ServerInternalErrMsg
	}
}
