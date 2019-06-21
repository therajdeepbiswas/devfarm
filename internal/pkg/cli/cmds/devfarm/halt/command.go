package halt

import (
	"fmt"
	"github.com/dena/devfarm/internal/pkg/cli"
	"github.com/dena/devfarm/internal/pkg/cli/formatter"
	"github.com/dena/devfarm/internal/pkg/cli/subcmd"
	"github.com/dena/devfarm/internal/pkg/platforms"
	"github.com/dena/devfarm/internal/pkg/platforms/all"
)

var CommandDef = subcmd.SubCommandDef{
	Desc:    "halt all devices in the specified instance group",
	Command: Command,
}

func Command(args []string, procInout cli.ProcessInout) cli.ExitStatus {
	opts, optsErr := takeOptions(args)
	if optsErr != nil {
		_, _ = fmt.Fprint(procInout.Stderr, optsErr.Error())
		return cli.ExitAbnormal
	}

	bag := cli.ComposeBag(procInout, opts.verbose, opts.dryRun)

	table, haltErr := all.HaltAll(opts.instanceGroupName, bag)

	successMsg := "halting"
	if _, err := fmt.Fprint(procInout.Stdout, formatter.PrettyTSV(table.TextTable(successMsg))); err != nil {
		_, _ = fmt.Fprintln(procInout.Stderr, err.Error())
		return cli.ExitAbnormal
	}

	if haltErr != nil {
		return cli.ExitAbnormal
	}

	return cli.ExitNormal
}

type options struct {
	verbose           bool
	dryRun            bool
	instanceGroupName platforms.InstanceGroupName
}

func takeOptions(args []string) (*options, *cli.ErrorAndUsage) {
	flags, usageBuf := cli.NewFlagSet([]string{})

	verbose := cli.DefineVerboseOpts(flags)
	dryRun := cli.DefineDryRunOpts(flags)
	unsafeInstanceGroupName := cli.DefineInstanceGroupNameOpts(flags)

	if err := cli.Parse(args, flags, usageBuf); err != nil {
		return nil, err
	}

	instanceGroupName, err := platforms.NewInstanceGroupName(*unsafeInstanceGroupName)
	if err != nil {
		flags.Usage()
		return nil, cli.NewErrorAndUsage(
			fmt.Sprintf("--instance-group: %s", err.Error()),
			usageBuf.String(),
		)
	}

	return &options{
		verbose:           *verbose,
		dryRun:            *dryRun,
		instanceGroupName: instanceGroupName,
	}, nil
}