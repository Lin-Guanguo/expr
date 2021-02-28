package evaluator

import (
	"fmt"
)

type EvalError struct {
	Message string
}

func (t EvalError) Error() string {
	return fmt.Sprintf("EvalError: %s", t.Message)
}
