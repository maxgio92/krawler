package matrix

import "fmt"

var (
	errSupportedColumnTypes        = []string{"[]string"}
	errUnsopportedPointTypeMessage = fmt.Sprintf("type of the matrix column is not supported: supported types are %s", errSupportedColumnTypes)
)

type ErrUnsopportedPointType struct {
	message string
}

func NewErrUnsopportedPointType() *ErrUnsopportedPointType {
	return &ErrUnsopportedPointType{
		message: errUnsopportedPointTypeMessage,
	}
}

func (e *ErrUnsopportedPointType) Error() string {
	return e.message
}
