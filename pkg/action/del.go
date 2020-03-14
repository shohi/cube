package action

import (
	"fmt"
	"os"

	"github.com/shohi/cube/pkg/base"
	"github.com/shohi/cube/pkg/kube"
)

type DelConfig struct {
	Name   string
	All    bool
	DryRun bool
}

// Del remove specified kubectl config
func Del(conf DelConfig) error {
	opts := kube.PurgeOptions{
		Name: conf.Name,
		All:  conf.All,
	}

	p := kube.NewPurger(opts)
	if err := p.Purge(); err != nil {
		return err
	}

	// always output updated kubeconfig.
	content, err := kube.Write(p.Result())
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "# updated config\n%v\n", string(content))
	if !conf.DryRun {
		kube.WriteToFile(p.Result(), base.GetLocalKubePath())
	}

	fmt.Fprintf(os.Stdout, "# cluster deleted\n%v\n", p.Deleted())

	return nil
}
