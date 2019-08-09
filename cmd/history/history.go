package history

import (
	"log"

	"github.com/shohi/cube/pkg/history"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	c := &cobra.Command{
		Use:   "history",
		Short: "show cube commands history",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := history.Read(); err != nil {
				log.Printf("get cube history error, err: %v\n", err)
			}

			return nil
		},
	}

	return c
}
