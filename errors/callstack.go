package errors

import (
	"bytes"
	"fmt"
	"runtime"
)

//CallStacker call stack
type CallStacker interface {
	GetCallStack() *CallStack
}

//CallStack call stack
type CallStack struct {
	Stacks []uintptr
}

//GetCallStacks get call stack
func GetCallStacks(err error) *CallStack {
	if err, ok := err.(CallStacker); ok {
		return err.GetCallStack()
	}
	return nil
}

// CallStacksString get call stack's string
func CallStacksString(call *CallStack) string {
	buf := bytes.Buffer{}
	if call == nil {
		return fmt.Sprintf("No call stack available")
	}

	for _, stack := range call.Stacks {
		f := runtime.FuncForPC(stack)
		file, line := f.FileLine(stack)
		buf.WriteString(fmt.Sprintf("%s:%d - %s\n", file, line, f.Name()))
	}

	return fmt.Sprintf("%s", buf.Bytes())
}

func getCallStack(skip int, depth int) *CallStack {
	stacks := make([]uintptr, depth)
	stacklen := runtime.Callers(skip, stacks)

	return &CallStack{
		Stacks: stacks[:stacklen],
	}
}
