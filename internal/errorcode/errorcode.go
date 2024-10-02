package errorcode

type ErrorCode string

const (
	ErrNotImplemented  ErrorCode = "NotImplemented"
	ErrInternal        ErrorCode = "Internal"
	ErrInvalidArgument ErrorCode = "InvalidArgument"
)
