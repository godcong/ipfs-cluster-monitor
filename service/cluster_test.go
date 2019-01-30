package service

import (
	"context"
	"github.com/godcong/ipfs-cluster-monitor/config"
	"github.com/godcong/ipfs-cluster-monitor/proto"
	"testing"
)

// TestClusterMonitor_Initialized ...
func TestClusterMonitor_Initialized(t *testing.T) {
	grpc := NewMonitorGRPC(config.DefaultConfig())

	client := MonitorClient(grpc)
	reply, err := client.MonitorInit(context.Background(), &proto.MonitorInitRequest{
		//Bootstrap:            "",
		//Secret:               "",
		//Path:        "",
		//ClusterPath: "",
		//XXX_NoUnkeyedLiteral: struct{}{},
		//XXX_unrecognized:     nil,
		//XXX_sizecache:        0,
	})
	t.Log(reply, err)
}
