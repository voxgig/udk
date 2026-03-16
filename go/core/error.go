package core

type UniversalError struct {
	IsUniversalError bool
	Sdk              string
	Code             string
	Msg              string
	Ctx              *Context
	Result           any
	Spec             any
}

func NewUniversalError(code string, msg string, ctx *Context) *UniversalError {
	return &UniversalError{
		IsUniversalError: true,
		Sdk:              "Universal",
		Code:             code,
		Msg:              msg,
		Ctx:              ctx,
	}
}

func (e *UniversalError) Error() string {
	return e.Msg
}
