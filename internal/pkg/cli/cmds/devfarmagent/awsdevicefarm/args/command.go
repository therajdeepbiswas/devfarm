package args

import (
	"github.com/dena/devfarm/internal/pkg/cli/cmds/devfarmagent/awsdevicefarm/args/decode"
	"github.com/dena/devfarm/internal/pkg/cli/cmds/devfarmagent/awsdevicefarm/args/encode"
	"github.com/dena/devfarm/internal/pkg/cli/subcmd"
)

var CommandDef = subcmd.SubCommandDef{
	Desc:    "encode/decode apps args to debug",
	Command: subcmd.NewSubCommand(table, subcmd.HelpFallbackCommand(table)),
}

var table = map[string]subcmd.SubCommandDef{
	"encode": encode.CommandDef,
	"decode": decode.CommandDef,
}
