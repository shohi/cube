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
	inKC   *clientcmdapi.Config

	inClusterName string
	inCluster     *clientcmdapi.Cluster

	inCtxName string
	inCtx     *clientcmdapi.Context

	inUser *clientcmdapi.AuthInfo

	updatedClusterName string
	inAPIServer        string
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
	exist, isDir := FileExists(configPath)

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

	inKC, err := m.d.Download()
	if err != nil {
		return err
	}

	m.inKC = inKC
	m.extractInKC()

	if err := m.doMerge(); err != nil {
		return err
	}

	return nil
}

func (m *merger) LocalPort() int {
	return m.localPort
}

func (m *merger) RemoteAPIAddr() string {
	return genRemoteAPIAddr(m.opts.RemoteAddr, m.inAPIServer)
}

func (m *merger) Result() *clientcmdapi.Config {
	return m.mainKC
}

// extractInKC extracts Cluster/User/Context info from `inKC`.
func (m *merger) extractInKC() {
	// NOTE: Only care about `clusters/contexts/users` sections
	for ck, v := range m.inKC.Clusters {
		// Take a snapshot of incoming cluster info
		// NOTE: It's ok to set `inAPIAddr` in the loop as
		// only one element is in the map.
		m.inClusterName, m.inCluster = ck, v
		m.inCtxName, m.inCtx = getContext(m.inKC, ck)
		m.inUser = getUser(m.inKC, m.inCtx.AuthInfo)

		m.inAPIServer = v.Server
	}
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
	cluster, found := findCluster(m.mainKC, m.inCluster)
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

	// check
	m.mainKC.Clusters[m.inCtx.Cluster] = m.inCluster
	m.mainKC.AuthInfos[m.inCtx.AuthInfo] = m.inUser
	m.mainKC.Contexts[m.inCtxName] = m.inCtx

	return nil
}

// normalizeInName normalized local names for remote cluster.
// The names for remote cluster should follow below conventions:
// 1. cluster name: `kubernetes` + `-` + nameSuffix
// 2. user name: `kubernetes-admin` + `-` + nameSuffix
// 3. context name: `kubernetes-admin` + `@` + `remoteIP:remotePort` + `-` + nameSuffix
func (m *merger) normalizeInName() {
	remoteHost := base.GetHost(m.opts.RemoteAddr)
	remotePort, _ := base.GetPort(m.inAPIServer)

	m.inCtx.AuthInfo = "kubernetes-" + m.opts.NameSuffix
	m.inCtx.Cluster = "kubernetes-" + m.opts.NameSuffix
	m.inCtxName = fmt.Sprintf("%s@%s:%v-%s", "kubernetes-admin",
		remoteHost, remotePort,
		m.opts.NameSuffix)

	m.updatedClusterName = m.inCtx.Cluster

	// update cluster's Server address
	m.inCluster.Server = fmt.Sprintf("https://%s:%d", DefaultHost, m.localPort)
}

func (m *merger) checkBeforeUpdate() error {
	if _, ok := m.mainKC.Clusters[m.inCtx.Cluster]; ok {
		return errors.Wrapf(ErrClusterAlreadyExists, "name: %v", m.inCtx.Cluster)
	}

	if _, ok := m.mainKC.AuthInfos[m.inCtx.AuthInfo]; ok {
		return errors.Wrapf(ErrUserAlreadyExists, "name: %v", m.inCtx.AuthInfo)
	}

	if _, ok := m.mainKC.Contexts[m.inCtxName]; ok {
		return errors.Wrapf(ErrContextAlreadyExists, "name: %v", m.inCtxName)
	}

	return nil
}
