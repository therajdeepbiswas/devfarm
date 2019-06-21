package foreverall

import (
	"fmt"
	"github.com/dena/devfarm/internal/pkg/cli"
	"github.com/dena/devfarm/internal/pkg/cli/formatter"
	"github.com/dena/devfarm/internal/pkg/cli/planfile"
	"github.com/dena/devfarm/internal/pkg/cli/subcmd"
	"github.com/dena/devfarm/internal/pkg/platforms/all"
)

var CommandDef = subcmd.SubCommandDef{
	Desc:    "launches multiple iOS/Android apps and restarts automatically if crashed during the lifetime",
	Command: Command,
}

func Command(args []string, procInout cli.ProcessInout) cli.ExitStatus {
	opts, optsErr := takeOptions(args)
	if optsErr != nil {
		_, _ = fmt.Fprintln(procInout.Stderr, optsErr.Error())
		return cli.ExitAbnormal
	}

	bag := cli.ComposeBag(procInout, opts.verbose, opts.dryRun)

	planFile, planFileErr := planfile.Read(opts.planFile, bag.GetFileOpener())
	if planFileErr != nil {
		_, _ = fmt.Fprintln(procInout.Stderr, fmt.Sprintf("invalid plan file:\n%s", planFileErr.Error()))
		return cli.ExitAbnormal
	}

	table, foreverErr := all.ForeverAll(planFile.Plans(), bag)

	successMsg := "launching"
	if _, err := fmt.Fprint(procInout.Stdout, formatter.PrettyTSV(table.TextTable(successMsg))); err != nil {
		return cli.ExitAbnormal
	}

	if foreverErr != nil {
		return cli.ExitAbnormal
	}

	return cli.ExitNormal
}