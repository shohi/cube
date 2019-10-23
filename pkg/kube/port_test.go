package kube

import (
	"log"
	"testing"
)

func TestPort_getOccupiedLocalPort(t *testing.T) {
	srv := "https://kubernetes:8001"
	log.Println(getOccupiedLocalPort(srv))
}
