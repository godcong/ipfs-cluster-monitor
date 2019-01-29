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
	err := Make(config.DefaultConfig(), "config/root/toml/config.toml")
	t.Log(err)
}
