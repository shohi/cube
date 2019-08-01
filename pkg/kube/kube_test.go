package kube

import (
	"log"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func TestKM_Init(t *testing.T) {
	configPath, err := homedir.Expand("~/.kube/config")
	if err != nil {
		t.Fatalf("failed to expand path, err: %v", err)
	}

	log.Printf("config path: %v", configPath)

	km := newKubeManager(kubeOptions{
		mainPath: configPath,
		inPath:   configPath,
	})
	err = km.init()

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
