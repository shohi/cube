package show

import (
	"fmt"
	"io/ioutil"

	"github.com/shohi/cube/pkg/base"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	c := &cobra.Command{
		Use:   "show",
		Short: "show local kubectl config",
		Run:   showLocalKubeConfig,
	}

	return c
}

func showLocalKubeConfig(_c *cobra.Command, _args []string) {
	configPath := base.GetLocalKubePath()
	exist, isDir := base.FileExists(configPath)

	if !exist {
		fmt.Printf("config path not exits - %v\n", configPath)
		return
	} else if isDir {
		fmt.Printf("config path [%v] is a dir, not file\n", configPath)
		return
	}

	var err error
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(content))
	}
}
