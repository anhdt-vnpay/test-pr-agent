package common

const (
	ERROR_INTERNAL int32 = 7
)

type BaseError interface {
	error
	GetCode() int32
}

type baseError struct {
	code int32
	err  error
}

func NewError(code int32, err error) *baseError {
	return &baseError{
		code: code,
		err:  err,
	}
}

func NewUnknownError(err error) *baseError {
	return &baseError{
		code: ERROR_INTERNAL,
		err:  err,
	}
}

func (e *baseError) GetCode() int32 {
	return e.code
}

func (e *baseError) Error() string {
	return e.err.Error()
}
