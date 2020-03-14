package forward

import (
	"github.com/shohi/cube/pkg/action"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var conf action.ForwardConfig

	c := &cobra.Command{
		Use:   "forward",
		Short: "run local ssh port forwarding for remote cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return action.Forward(conf)

		},
	}
	setupFlags(c, &conf)

	return c
}

// setupFlags sets flags for comand line
func setupFlags(cmd *cobra.Command, conf *action.ForwardConfig) {
	flagSet := cmd.Flags()

	flagSet.StringVar(&conf.Name, "name", "", "cluster name")
	flagSet.StringVar(&conf.Operation, "op", "print", "operation, avaliable options: print/run/stop")
	flagSet.StringVar(&conf.SSHVia, "ssh-via", "", "ssh jump server, e.g. user@jump. If not set, SSH_VIA env will be used ")

	cmd.MarkFlagRequired("name")
}
