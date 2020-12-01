package list

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/shohi/cube/pkg/kube"
)

// Options for list subcommand
type Options struct {
	Filter string // name filter
}

// New creates a new `list` subcommand.
// Output format
// cluster_name : `ssh-forwarding`
func New() *cobra.Command {
	var opts = &Options{}

	c := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list clusters <name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := kube.ListAllClusters()
			if err != nil {
				log.Printf("list clusters error, err: %v\n", err)
				return err
			}
			selected := filter(l, ByName(opts.Filter))

			content, _ := json.MarshalIndent(selected, "", "  ")
			fmt.Println(string(content))
			return nil
		},
	}

	setupFlags(c, opts)

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

func ByName(pattern string) FilterFunc {
	p := strings.TrimSpace(pattern)

	if p == "" {
		return EnableAll
	}

	re := regexp.MustCompile(".*" + p + ".*")

	return func(c kube.ClusterInfo) bool {
		return re.MatchString(c.Name)
	}
}

// setupFlags sets flags for comand line
func setupFlags(cmd *cobra.Command, opts *Options) {
	flagSet := cmd.Flags()

	flagSet.StringVarP(&opts.Filter, "filter", "f", "", "cluster name pattern, default list all")
}
