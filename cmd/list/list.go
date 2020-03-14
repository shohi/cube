package list

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/shohi/cube/pkg/kube"
)

// New creates a new `list` subcommand.
// Output format
// cluster_name : `ssh-forwarding`
func New() *cobra.Command {
	c := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list all clusters",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := kube.ListAllClusters()
			if err != nil {
				log.Printf("list clusters error, err: %v\n", err)
				return err
			}

			content, _ := json.MarshalIndent(l, "", "  ")
			fmt.Println(string(content))
			return nil
		},
	}

	return c
}
