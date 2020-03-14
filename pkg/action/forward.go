package action

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/shirou/gopsutil/process"
	"github.com/shohi/cube/pkg/base"
	"github.com/shohi/cube/pkg/kube"
)

var (
	errProcessNotFound   = errors.New("process not found")
	errSSHViaEnvNotFound = errors.New("cube: SSH_VIA env not set")
)

type Operation int

const (
	OpPrint Operation = iota
	OpRun
	OpStop
)

func parseOp(op string) Operation {
	switch strings.ToLower(op) {
	case "run":
		return OpRun
	case "stop":
		return OpStop
	default:
		return OpPrint
	}
}

type ForwardConfig struct {
	Name      string
	Operation string
	SSHVia    string
}

var (
	errEmptyName            = errors.New("forward: empty cluster name")
	errClusterNotFound      = errors.New("forward: cluster not found")
	errMultipleClusterFound = errors.New("forward: multiple clusters found")
)

var (
	filterOutRe = regexp.MustCompile(`.*(kind|minikube).*`)
)

func filter(name string) bool {
	return !filterOutRe.MatchString(name)
}

func Forward(conf ForwardConfig) error {
	op := parseOp(conf.Operation)
	if conf.Name == "" {
		return errEmptyName
	}

	if err := setSSHVia(conf.SSHVia); err != nil {
		return err
	}

	kc, err := kube.Load(base.GetLocalKubePath())
	if err != nil {
		return err
	}

	ctxs := kube.FindContextsByName(kc, conf.Name, filter)
	if len(ctxs) == 0 {
		return errClusterNotFound
	}

	if op == OpRun && len(ctxs) > 1 {
		return errMultipleClusterFound
	}

	for k := range ctxs {
		info, err := kube.ParseContext(kc, k)
		if err != nil {
			return err
		}

		err = doSSHForwarding(op, k, info.SSHForward)
		if err != nil {
			return err
		}
	}

	return nil
}

func doSSHForwarding(op Operation, ctxName, forwardCmd string) error {
	fmt.Fprintf(os.Stdout, "# context - %v\n", ctxName)

	switch op {
	case OpPrint:
		fmt.Fprintln(os.Stdout, forwardCmd)
		return nil
	case OpRun:
		fmt.Fprintf(os.Stdout, "# %v\n", forwardCmd)
		err := startPortFowarding(forwardCmd)
		if err != nil {
			return err
		}
		fmt.Println("start ssh local port forwarding successfully.")
		return nil
	case OpStop:
		fmt.Fprintf(os.Stdout, "# %v\n", forwardCmd)
		err := stopPortForwarding(forwardCmd)
		if err == nil {
			fmt.Println("stop ssh local port forwarding successfully.")
		}

		if err == errProcessNotFound {
			fmt.Println("process not found")
			return nil
		}

		return err
	}

	return nil
}

func startPortFowarding(cmdStr string) error {
	cmdArgs := strings.Fields(cmdStr)

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	return cmd.Run()
}

func stopPortForwarding(cmdStr string) error {
	cmdArgs := strings.Fields(cmdStr)
	formalCmd := strings.Join(cmdArgs, " ")

	var pfProcess *process.Process
	lst, err := process.Processes()
	if err != nil {
		return err
	}

	for _, p := range lst {
		cli, err := p.Cmdline()
		if err != nil {
			return err
		}
		if strings.Contains(cli, formalCmd) {
			pfProcess = p
		}
	}

	if pfProcess == nil {
		return errProcessNotFound
	}

	return pfProcess.Kill()
}

func setSSHVia(via string) error {
	if via != "" {
		return os.Setenv("SSH_VIA", via)
	}
	sv := os.Getenv("SSH_VIA")
	if len(sv) == 0 {
		return errSSHViaEnvNotFound
	}

	return nil
}
