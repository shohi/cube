package kube

import (
	"fmt"
	"log"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/shohi/cube/pkg/local"
)

func newKubeOptionsForTest(remoteIP string) kubeOptions {
	configPath, err := homedir.Expand("~/.kube/config")
	if err != nil {
		panic(fmt.Sprintf("failed to expand path, err: %v", err))
	}

	inPath := configPath
	if remoteIP != "" {
		inPath = local.GetLocalPath(remoteIP)
	}

	return kubeOptions{
		mainPath: configPath,
		inPath:   inPath,
	}
}

func TestKM_Init(t *testing.T) {
	opts := newKubeOptionsForTest("")
	log.Printf("config path: %v", opts.mainPath)

	km := newKubeManager(opts)
	err := km.init()

	if err != nil {
		t.Fatalf("failed to read config, err: %v", err)
	}

	log.Printf("kube config: [%+v]", km.mainKC)
}

func TestKM_Merge(t *testing.T) {
	km := newKubeManager(kubeOptions{
		mainPath:   getLocalKubePath(),
		inPath:     getLocalPath("core@172.31.10.34"),
		localPort:  7003,
		nameSuffix: "test",
	})

	err := km.init()
	if err != nil {
		t.Fatalf("failed to init KubeManager, err: %v", err)
	}

	err = km.merge()
	if err != nil {
		t.Fatalf("KubeManager failed to merge, err: %v", err)
	}

	log.Printf("kube main config: [%+v]", km.mainKC)
	log.Printf("kube in config: [%+v]", km.inKC)
}

func TestKM_K8s(t *testing.T) {
	//
	opts := newKubeOptionsForTest("172.31.4.206")
	opts.nameSuffix = "k8s"

	km := newKubeManager(opts)
	if err := km.init(); err != nil {
		t.Fatalf("failed to init KubeManager, err: %v", err)
	}

	if err := km.extractInKC(); err != nil {
		t.Fatalf("failed to extract infor from APIConfig, err: %v", err)
	}

	log.Printf("user info: [%v]", string(km.inUser.ClientCertificateData))
	log.Printf("auth info: [%v]", string(km.inCluster.CertificateAuthorityData))
}
