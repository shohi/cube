package action

import (
	"fmt"
	"os"

	"github.com/shohi/cube/pkg/base"
	"github.com/shohi/cube/pkg/kube"
)

type AddConfig struct {
	RemoteUser string
	RemoteIP   string

	LocalPort  int
	SSHVia     string
	NameSuffix string

	DryRun bool
	Force  bool

	PrintSSHForwarding bool
}

// TODO: test
// Add adds new kubectl config.
func Add(conf AddConfig) error {
	remoteAddr := base.SshHost(conf.RemoteUser, conf.RemoteIP)

	opts := kube.MergeOptions{
		RemoteAddr: remoteAddr,
		NameSuffix: conf.NameSuffix,
		LocalPort:  conf.LocalPort,
		Force:      conf.Force,
	}
	m := kube.NewMerger(opts)
	if err := m.Merge(); err != nil {
		return err
	}

	localPort := m.LocalPort()
	apiAddr := m.RemoteAPIAddr()
	sshCmd := kube.GetPortForwardingCmd(localPort, apiAddr, conf.SSHVia)
	if conf.PrintSSHForwarding {
		// NOTE: Print SSH forwarding setting
		fmt.Fprintf(os.Stdout, "# ssh forwarding command\n%s\n", sshCmd)
		return nil
	}

	// always output updated kubeconfig.
	content, err := kube.Write(m.Result())
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "# updated config\n%v\n", string(content))
	fmt.Fprintf(os.Stdout, "# ssh forwarding command\n%s\n", sshCmd)

	if !conf.DryRun {
		kube.WriteToFile(m.Result(), base.GetLocalKubePath())
	}

	return nil
}
