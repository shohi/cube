package list

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"

	"github.com/shohi/cube/pkg/kube"
)

// Options for list subcommand
type Options struct {
	Name string // name filter
}

// New creates a new `list` subcommand.
// Output format
// cluster_name : `ssh-forwarding`
func New() *cobra.Command {
	var opts = &Options{}

	c := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list clusters",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := kube.ListAllClusters()
			if err != nil {
				log.Printf("list clusters error, err: %v\n", err)
				return err
			}

			if len(args) > 0 {
				opts.Name = args[0]
			}

			selected := filter(l, ByName(opts.Name))

			content, _ := json.MarshalIndent(selected, "", "  ")
			fmt.Println(string(content))
			return nil
		},
	}

	// setupFlags(c, opts)

	return c
}

func filter(s kube.ClusterInfos, fn FilterFunc) kube.ClusterInfos {
	result := make(kube.ClusterInfos, 0, len(s))

	for _, c := range s {
		if fn(c) {
			result = append(result, c)
		}
	}

	return result
}

type FilterFunc func(kube.ClusterInfo) bool

func EnableAll(_ kube.ClusterInfo) bool {
	return true
}

func ByName(name string) FilterFunc {
	n := strings.TrimSpace(name)

	if n == "" {
		return EnableAll
	}

	return func(c kube.ClusterInfo) bool {
		return strings.Contains(c.Name, n)
	}
}

// setupFlags sets flags for comand line
func setupFlags(cmd *cobra.Command, opts *Options) {
	flagSet := cmd.Flags()

	flagSet.StringVar(&opts.Name, "name", "", "cluster name pattern")
}
