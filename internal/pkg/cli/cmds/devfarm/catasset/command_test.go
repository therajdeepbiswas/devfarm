package catasset

import (
	"github.com/dena/devfarm/internal/pkg/cli"
	"testing"
)

func TestCommand(t *testing.T) {
	procInout := cli.AnyProcInout()

	got := Command([]string{"something/asset/id"}, procInout)

	// NOTE: should fail because the asset ID is invalid.
	if got != cli.ExitAbnormal {
		t.Errorf("got %v, want %v", got, cli.ExitAbnormal)
		return
	}
}
