package scp

import (
	"errors"
	"fmt"
	"os/exec"
)

type TransferDirect int

const (
	ToRemote TransferDirect = iota
	ToLocal
)

type TransferConfig struct {
	Direct     TransferDirect
	RemoteAddr string
	RemotePath string

	LocalPath string
}

var (
	ErrInvalidRemoteAddr = errors.New("scp: invalid remote address")
	ErrInvalidRemotePath = errors.New("scp: invalid remote path")
	ErrInvalidLocalPath  = errors.New("scp: invalid local path")
)

func checkConfig(conf TransferConfig) error {
	if conf.RemoteAddr == "" {
		return ErrInvalidRemoteAddr
	}

	if conf.RemotePath == "" {
		return ErrInvalidRemotePath
	}

	if conf.LocalPath == "" {
		return ErrInvalidLocalPath
	}

	return nil
}

// TransferFile transfers file between two hosts using scp.
// FIXME: support context
func TransferFile(conf TransferConfig) error {
	if err := checkConfig(conf); err != nil {
		return err
	}

	// Fast check whether file has already been downloaded.
	if hasContent(conf.LocalPath) {
		return nil
	}

	var args []string
	remoteLoc := fmt.Sprintf("%s:%s", conf.RemoteAddr, conf.RemotePath)

	switch conf.Direct {
	case ToRemote:
		args = []string{remoteLoc, conf.LocalPath}
	default:
		args = []string{conf.LocalPath, remoteLoc}
	}

	// TODO: add timeout control
	cmd := exec.Command("scp", args...)

	// TODO: dump log
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("scp: transfer error - %v", err)
	}

	return nil
}
