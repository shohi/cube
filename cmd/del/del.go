package del

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/shohi/cube/pkg/action"
	hist "github.com/shohi/cube/pkg/history"
)

func New() *cobra.Command {
	var conf = action.DelConfig{}

	c := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del"},
		Short:   "delete kubectl config for specified cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := hist.Write(); err != nil {
				log.Printf("failed to write history, err: %v\n", err)
			}
			return action.Del(conf)
		},
	}

	setupFlags(c, &conf)

	return c

}

// setupFlags sets flags for comand line
func setupFlags(cmd *cobra.Command, conf *action.DelConfig) {
	flagSet := cmd.Flags()

	flagSet.StringVar(&conf.Name, "name", "", "cluster name to delete")
	flagSet.BoolVar(&conf.DryRun, "dry-run", false, "dry-run mode. print modified config and exit")
	flagSet.BoolVar(&conf.All, "all", false, "delete all matched cluster.")

	cmd.MarkFlagRequired("name")
}
