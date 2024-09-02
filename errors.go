package squeue

import "fmt"

const (
	ErrUnmarshal int = -1
	ErrMarshal   int = -2
	ErrDriver    int = -3
)

type Err struct {
	message string
	code    int
	data    []byte
	inner   error
}

func (e *Err) Error() string {
	return fmt.Sprintf("%s: %+v", e.message, e.inner)
}

func (e *Err) Code() int {
	return e.code
}

func (e *Err) Data() []byte {
	return e.data
}

func (e *Err) Unwrap() error {
	return e.inner
}

func wrapErr(err error, code int, data []byte) *Err {
	var errMessages = map[int]string{
		ErrUnmarshal: "unmarshalling error",
		ErrMarshal:   "marshalling error",
		ErrDriver:    "driver error",
	}

	if message, ok := errMessages[code]; ok {
		return &Err{
			message: message,
			code:    code,
			data:    data,
			inner:   err,
		}
	}

	return &Err{
		message: "unknown error",
		code:    code,
		data:    data,
		inner:   err,
	}
}
