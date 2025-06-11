package cmd

import (
	"bytes"
	"testing"

	"github.com/evberrypi/ssh-config/version"
	"github.com/spf13/cobra"
)

func TestVersionCmd(t *testing.T) {
	buf := &bytes.Buffer{}
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	VersionCmd.Run(cmd, []string{})

	output := buf.String()
	expected := version.String() + "\n"
	if output != expected {
		t.Errorf("expected %q, got %q", expected, output)
	}
}
