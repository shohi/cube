package kube

import (
	"fmt"
	"os"

	"github.com/shohi/cube/pkg/config"
	"github.com/shohi/cube/pkg/scp"
)

// Dispatch dispatches tasks according to config.
// TODO: use subcommands instead
func Dispatch(conf config.Config) error {
	remoteAddr := getRemoteAddr(conf.RemoteUser, conf.RemoteIP)
	p := getLocalPath(remoteAddr)

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
		mainPath:   getLocalKubePath(),
		inPath:     p,
		action:     getAction(conf.Purge, conf.PrintSSHForwarding),
		localPort:  conf.LocalPort,
		nameSuffix: conf.NameSuffix,
	})

	if err := km.Do(); err != nil {
		return err
	}

	sshCmd := getPortForwardingCmd(km.opts.localPort, km.inAPIAddr, conf.SSHVia)
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
