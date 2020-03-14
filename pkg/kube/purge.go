package kube

import (
	"github.com/pkg/errors"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/shohi/cube/pkg/base"
)

var (
	ErrClusterNotFound       = errors.New("cube: cluster not found for purging")
	ErrMultipleClustersFound = errors.New("cube: multiple clusters found for purging")
)

// Purger deletes Kubernetes configs under given conditions
type Purger interface {
	Purge() error
	Result() *clientcmdapi.Config
	Deleted() []string
}

// PurgeOptions represent options for purge.
type PurgeOptions struct {
	Name string
	All  bool
}

type purger struct {
	opts   PurgeOptions
	mainKC *clientcmdapi.Config

	selectedCtxs map[string]*clientcmdapi.Context
}

func NewPurger(opts PurgeOptions) Purger {
	return &purger{
		opts: opts,
	}
}

// Purge delete Kubernetes config whose context name matches the given pattern.
func (p *purger) Purge() error {
	mainKC, err := Load(base.GetLocalKubePath())
	if err != nil {
		return err
	}
	p.mainKC = mainKC

	p.selectedCtxs = FindContextsByName(p.mainKC, p.opts.Name, nil)
	if len(p.selectedCtxs) == 0 {
		return ErrClusterNotFound
	}

	if len(p.selectedCtxs) > 1 && !p.opts.All {
		return errors.Wrapf(ErrMultipleClustersFound, "list: %v", p.contextList())
	}

	for k, v := range p.selectedCtxs {
		delete(p.mainKC.Clusters, v.Cluster)
		delete(p.mainKC.AuthInfos, v.AuthInfo)
		delete(p.mainKC.Contexts, k)
	}

	return nil
}

func (p *purger) contextList() []string {
	var ret = make([]string, 0, len(p.selectedCtxs))

	for k := range p.selectedCtxs {
		ret = append(ret, k)
	}

	return ret
}

func (p *purger) Result() *clientcmdapi.Config {
	return p.mainKC
}

func (p *purger) Deleted() []string {
	return p.contextList()
}

/*
// inferLocalPort infers local port from default kubeconfig if local-port not provided.
func (m *merger) inferLocalPort() error {
	if m.opts.LocalPort > 0 && m.opts.LocalPort < 65535 {
		return nil
	}

	cluster, found := findCluster(m.mainKC, m.inCluster)
	if !found {
		return errors.Wrapf(ErrClusterNotFound, "cluster: [%v]", k.inCluster)
	}

	c := k.mainKC.Clusters[cluster]
	k.opts.localPort, _ = base.GetPort(c.Server)

	return nil
}
*/
