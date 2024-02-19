package response

type LogicReturn[T any] struct {
	Response      T
	HttpErrorCode int
	ErrorMsg      error
	Count         int64
}
