package add

import (
	"log"

	"github.com/shohi/cube/pkg/config"
	hist "github.com/shohi/cube/pkg/history"
	"github.com/shohi/cube/pkg/kube"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	var conf = config.Config{
		Purge: false,
	}

	c := &cobra.Command{
		Use:   "add",
		Short: "add remote cluster to kube config",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := hist.Write(); err != nil {
				log.Printf("failed to write history, err: %v\n", err)
			}
			return kube.Dispatch(conf)
		},
	}

	setupFlags(c, &conf)

	return c

}

// setupFlags sets flags for comand line
func setupFlags(cmd *cobra.Command, conf *config.Config) {
	flagSet := cmd.Flags()

	flagSet.StringVar(&conf.RemoteUser, "remote-user", "core", "remote user")
	flagSet.StringVar(&conf.RemoteIP, "remote-ip", "", "remote master private ip")

	flagSet.IntVar(&conf.LocalPort, "local-port", 0, "local forwarding port")
	flagSet.StringVar(&conf.SSHVia, "ssh-via", "", "ssh jump server, e.g. user@jump. If not set, SSH_VIA env will be used ")
	flagSet.StringVar(&conf.NameSuffix, "name-suffix", "", "cluster name suffix")

	flagSet.BoolVar(&conf.DryRun, "dry-run", false, "dry-run mode. validate config and then exit")
	flagSet.BoolVar(&conf.Force, "force", false, "merge configuration forcedly. Only take effect when cluster name is unique")
	flagSet.BoolVar(&conf.PrintSSHForwarding, "print-ssh-forwarding", false, "print ssh forwarding command and exit")

	cmd.MarkFlagRequired("remote-ip")
}
