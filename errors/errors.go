package errors

import "errors"

const (
	callStackDepth = 10
)

//DetailError error with details
type DetailError interface {
	error
	ErrCoder
	CallStacker
	GetErrRoot() error
}

//NewErr create a new error with error message
func NewErr(errmsg string) error {
	return errors.New(errmsg)
}

//NewDetailErr create new DetailError
func NewDetailErr(err error, errcode ErrCode, errmsg string) DetailError {
	if err == nil {
		return nil
	}

	e, ok := err.(Errors)
	if !ok {
		e.root = err
		e.errmsg = err.Error()
		e.callstack = getCallStack(0, callStackDepth)
		e.code = errcode

	}
	if errmsg != "" {
		e.errmsg = errmsg + ": " + e.errmsg
	}

	return e
}

//RootErr error root
func RootErr(err error) error {
	if err, ok := err.(DetailError); ok {
		return err.GetErrRoot()
	}
	return err
}

// Errors errors
type Errors struct {
	errmsg    string
	callstack *CallStack
	root      error
	code      ErrCode
}

//Error get error message
func (e Errors) Error() string {
	return e.errmsg
}

//GetErrCode get error code
func (e Errors) GetErrCode() ErrCode {
	return e.code
}

//GetErrRoot implement DetailError interface
func (e Errors) GetErrRoot() error {
	return e.root
}

//GetCallStack get call stack
func (e Errors) GetCallStack() *CallStack {
	return e.callstack
}
