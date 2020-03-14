package kube

import (
	"fmt"
	"os"

	"github.com/shohi/cube/pkg/base"
	"github.com/shohi/cube/pkg/config"
	"github.com/shohi/cube/pkg/scp"
)

// Add adds new kubectl config
func Add(conf config.Config) error {
	remoteAddr := getRemoteAddr(conf.RemoteUser, conf.RemoteIP)
	p := base.GenLocalPath(remoteAddr)

	// TODO: check whether the file is empty
	err := scp.TransferFile(scp.TransferConfig{
		LocalPath:  p,
		RemoteAddr: remoteAddr,
		RemotePath: DefaultKubeConfigPath,
	})

	if err != nil {
		return err
	}

	// TODO: handle duplicated
	km := newKubeManager(kubeOptions{
		mainPath:   base.GetLocalKubePath(),
		inPath:     p,
		action:     getAction(conf.Purge, conf.PrintSSHForwarding),
		localPort:  conf.LocalPort,
		nameSuffix: conf.NameSuffix,
		force:      conf.Force,
		remoteAddr: remoteAddr,
	})

	if err := km.Do(); err != nil {
		return err
	}

	inAPIAddr := genInAPIAddr(km.opts.remoteAddr, km.inAPIServer)
	sshCmd := getPortForwardingCmd(km.opts.localPort, inAPIAddr, conf.SSHVia)
	if conf.PrintSSHForwarding {
		// NOTE: Print SSH forwarding setting
		fmt.Fprintf(os.Stdout, "# ssh forwarding command\n%s\n", sshCmd)
		return nil
	}

	// always output updated kubeconfig.
	content, err := km.Write()
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "# updated config\n%v\n", string(content))
	fmt.Fprintf(os.Stdout, "# ssh forwarding command\n%s\n", sshCmd)

	if !conf.DryRun {
		km.WriteToFile()
	}

	return nil
}
