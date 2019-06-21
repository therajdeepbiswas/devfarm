package devicefarm

import (
	"github.com/dena/devfarm/internal/pkg/executor"
	"github.com/dena/devfarm/internal/pkg/executor/awscli"
)

type Executor func(devicefarmArgs ...string) (executor.Result, error)

func NewExecutor(awsCmd awscli.Executor) Executor {
	// NOTE: AWS Device Farm is only supported on us-west-2 (2019-06-28).
	// https://docs.aws.amazon.com/ja_jp/general/latest/gr/rande.html#devicefarm_region
	region := awscli.RegionIsUSWest2

	return func(devicefarmArgs ...string) (executor.Result, error) {
		args := make([]string, len(devicefarmArgs)+1)
		args[0] = "devicefarm"

		for i, arg := range devicefarmArgs {
			args[i+1] = arg
		}

		implicitArgs := []string{"--region", region.Name(), "--output", "json"}
		args = append(args, implicitArgs...)

		return awsCmd(args...)
	}
}
