package kube

import (
	"fmt"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/pkg/errors"
	"github.com/shohi/cube/pkg/base"
)

var (
	ErrEmptyNameSuffix      = errors.New("cube: empty name-suffix for merge")
	ErrConfigAlreadyMerged  = errors.New("cube: kubeconfig already merged")
	ErrClusterAlreadyExists = errors.New("cube: cluster already exists")
	ErrContextAlreadyExists = errors.New("cube: context already exists")
	ErrUserAlreadyExists    = errors.New("cube: user already exists")
	ErrInvalidLocalPort     = errors.New("cube: invalid local port for merge")
)

// MergeOptions represents options for merge
type MergeOptions struct {
	RemoteAddr string
	NameSuffix string
	LocalPort  int
	Force      bool
}

// Merger merge remote cluster config into local `~/.kube/config`
type Merger interface {
	Merge() error
	Result() *clientcmdapi.Config
	LocalPort() int
	RemoteAPIAddr() string
}

type merger struct {
	opts MergeOptions

	d         *Downloader
	localPort int

	mainKC *clientcmdapi.Config

	inKC          *clientcmdapi.Config
	inClusterName string
	inCK          ClusterKeyInfo

	updatedClusterName string
}

func NewMerger(opts MergeOptions) Merger {
	d := NewDownloader(opts.RemoteAddr)

	m := &merger{
		opts:      opts,
		localPort: opts.LocalPort,
		d:         d,
	}

	return m
}

func (m *merger) loadMainKC() error {
	configPath := base.GetLocalKubePath()
	exist, isDir := base.FileExists(configPath)

	if !exist {
		m.mainKC = clientcmdapi.NewConfig()
		return nil
	} else if isDir {
		return fmt.Errorf("config path [%v] is a dir, not file", configPath)
	}

	var err error
	m.mainKC, err = Load(base.GetLocalKubePath())
	if err != nil {
		return err
	}

	return nil
}

func (m *merger) Merge() error {
	if m.opts.NameSuffix == "" {
		return ErrEmptyNameSuffix
	}

	if err := m.loadMainKC(); err != nil {
		return err
	}

	res, err := m.d.Download()
	if err != nil {
		return err
	}

	m.inKC = res.Kc
	m.inClusterName = res.ClusterName
	m.inCK = getClusterKeyInfo(m.inKC, m.inClusterName)

	if err := m.doMerge(); err != nil {
		return err
	}

	return nil
}

func (m *merger) LocalPort() int {
	return m.localPort
}

func (m *merger) RemoteAPIAddr() string {
	return genRemoteAPIAddr(m.opts.RemoteAddr, m.inCK.Cluster.Server)
}

func (m *merger) Result() *clientcmdapi.Config {
	return m.mainKC
}

func (m *merger) checkLocalPort() error {
	if m.localPort == 0 {
		nextPort := getNextLocalPort(m.mainKC)
		m.localPort = nextPort
	}

	if m.localPort <= 0 || m.localPort > 65535 {
		return errors.Wrapf(ErrInvalidLocalPort, "port: %v", m.localPort)
	}

	return nil
}

func (m *merger) checkExists() error {
	cluster, found := findCluster(m.mainKC, m.inCK.Cluster)
	if found && !m.opts.Force {
		return errors.Wrapf(ErrConfigAlreadyMerged, "cluster: [%v]", cluster)
	}

	return nil
}

func (m *merger) doMerge() error {
	if err := m.checkExists(); err != nil {
		return err
	}

	if err := m.checkLocalPort(); err != nil {
		return err
	}

	m.normalizeInName()

	if err := m.checkBeforeUpdate(); err != nil {
		return err
	}

	m.mainKC.Clusters[m.inCK.Ctx.Cluster] = m.inCK.Cluster
	if m.inCK.User != nil {
		m.mainKC.AuthInfos[m.inCK.Ctx.AuthInfo] = m.inCK.User
	}
	m.mainKC.Contexts[m.inCK.CtxName] = m.inCK.Ctx

	return nil
}

// normalizeInName normalized local names for remote cluster.
// The names for remote cluster should follow below conventions:
// 1. cluster name: `kubernetes` + `-` + nameSuffix
// 2. user name: `kubernetes-admin` + `-` + nameSuffix
// 3. context name: `kubernetes-admin` + `@` + `remoteIP:remotePort` + `-` + nameSuffix
func (m *merger) normalizeInName() {
	remoteHost := base.GetHost(m.opts.RemoteAddr)
	remotePort, _ := base.GetPort(m.inCK.Cluster.Server)

	m.inCK.Ctx.AuthInfo = "kubernetes-" + m.opts.NameSuffix
	m.inCK.Ctx.Cluster = "kubernetes-" + m.opts.NameSuffix
	m.inCK.CtxName = fmt.Sprintf("%s@%s:%v-%s", "kubernetes-admin",
		remoteHost, remotePort,
		m.opts.NameSuffix)

	m.updatedClusterName = m.inCK.Ctx.Cluster

	// update server address aware of http/https
	var schema = "https"
	if m.inCK.IsHTTP {
		schema = "http"
	}

	m.inCK.Cluster.Server = fmt.Sprintf("%s://%s:%d",
		schema, DefaultHost, m.localPort)
}

func (m *merger) checkBeforeUpdate() error {
	if _, ok := m.mainKC.Clusters[m.inCK.Ctx.Cluster]; ok {
		return errors.Wrapf(ErrClusterAlreadyExists, "name: %v", m.inCK.Ctx.Cluster)
	}

	if _, ok := m.mainKC.AuthInfos[m.inCK.Ctx.AuthInfo]; ok {
		return errors.Wrapf(ErrUserAlreadyExists, "name: %v", m.inCK.Ctx.AuthInfo)
	}

	if _, ok := m.mainKC.Contexts[m.inCK.CtxName]; ok {
		return errors.Wrapf(ErrContextAlreadyExists, "name: %v", m.inCK.CtxName)
	}

	return nil
}
