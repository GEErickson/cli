package context

import (
	"fmt"
	"os"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/context/store"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

func newUseCommand(dockerCli command.Cli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use CONTEXT",
		Short: "Set the current docker context",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			return RunUse(dockerCli, name)
		},
	}
	return cmd
}

// RunUse set the current Docker context
func RunUse(dockerCli command.Cli, name string) error {
	if err := store.ValidateContextName(name); err != nil && name != "default" {
		return err
	}
	if _, err := dockerCli.ContextStore().GetMetadata(name); err != nil && name != "default" {
		return err
	}
	configValue := name
	if configValue == "default" {
		configValue = ""
	}
	dockerConfig := dockerCli.ConfigFile()
	dockerConfig.CurrentContext = configValue
	if err := dockerConfig.Save(); err != nil {
		return err
	}
	fmt.Fprintln(dockerCli.Out(), name)
	fmt.Fprintf(dockerCli.Err(), "Current context is now %q\n", name)
	if os.Getenv(client.EnvOverrideHost) != "" {
		fmt.Fprintf(dockerCli.Err(), "Warning: %[1]s environment variable overrides the active context. "+
			"To use %[2]q, either set the global --context flag, or unset %[1]s environment variable.\n", client.EnvOverrideHost, name)
	}
	return nil
}
