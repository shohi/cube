package forward

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var errNoClusterName = errors.New("please specify cluster name or suffix")

const usage = `
Usage:
  cube forward [flags] [cluster-name]

Flags:
  -h, --help   help for forward
`

func New() *cobra.Command {
	usageFunc := func(c *cobra.Command) error {
		_, err := fmt.Fprintf(c.OutOrStderr(), strings.TrimSpace(usage)+"\n")
		return err
	}
	c := &cobra.Command{
		Use:   "forward",
		Short: "run local ssh port forwarding for remote cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errNoClusterName
			}

			// TODO
			fmt.Printf("start ssh port forwarding for cluster - %s ...\n", args[0])
			return nil
		},
		SilenceErrors: true,
	}

	c.SetUsageFunc(usageFunc)

	return c
}
