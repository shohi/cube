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
	ErrInvalidRemoteAddr = errors.New("invalid remote address")
	ErrInvalidRemotePath = errors.New("invalid remote path")
	ErrInvalidLocalPath  = errors.New("invalid local path")
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
	err := cmd.Run()

	return err
}
