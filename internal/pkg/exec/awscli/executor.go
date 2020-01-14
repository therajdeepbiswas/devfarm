package awscli

import (
	"github.com/dena/devfarm/internal/pkg/exec"
)

type Executor func(args ...string) (exec.Result, error)

func NewExecutor(execute exec.Executor) Executor {
	return func(args ...string) (exec.Result, error) {
		request := exec.NewRequest("aws", args)
		return execute(request)
	}
}
