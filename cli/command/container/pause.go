package container

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/command/completion"
	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type pauseOptions struct {
	containers []string
}

// NewPauseCommand creates a new cobra.Command for `docker pause`
func NewPauseCommand(dockerCli command.Cli) *cobra.Command {
	var opts pauseOptions

	return &cobra.Command{
		Use:   "pause CONTAINER [CONTAINER...]",
		Short: "Pause all processes within one or more containers",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.containers = args
			return runPause(dockerCli, &opts)
		},
		ValidArgsFunction: completion.ContainerNames(dockerCli, false, func(container types.Container) bool {
			return container.State != "paused"
		}),
	}
}

func runPause(dockerCli command.Cli, opts *pauseOptions) error {
	ctx := context.Background()

	var errs []string
	errChan := parallelOperation(ctx, opts.containers, dockerCli.Client().ContainerPause)
	for _, container := range opts.containers {
		if err := <-errChan; err != nil {
			errs = append(errs, err.Error())
			continue
		}
		fmt.Fprintln(dockerCli.Out(), container)
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}
