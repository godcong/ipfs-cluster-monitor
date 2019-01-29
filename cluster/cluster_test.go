package cluster

import (
	"github.com/godcong/ipfs-cluster-monitor/config"
	"testing"
)

// TestFirstRun ...
func TestFirstRunIPFS(t *testing.T) {

	firstRunIPFS()

}

// TestRun ...
func TestRun(t *testing.T) {
	err := InitMaker(config.DefaultConfig(), "config.toml")
	t.Log(err)
}
