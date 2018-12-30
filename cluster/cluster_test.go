package cluster

import (
	"golang.org/x/net/context"
	"testing"
)

// TestFirstRun ...
func TestFirstRunIPFS(t *testing.T) {

	firstRunIPFS()

}

// TestRun ...
func TestRun(t *testing.T) {
	Start(context.Background())
}
