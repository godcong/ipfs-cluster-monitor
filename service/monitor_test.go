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

	reply, err := client.MonitorProc(context.Background(), &proto.MonitorProcRequest{
		Type: proto.MonitorType_Reset,
	})

	reply, err = client.MonitorInit(context.Background(), &proto.MonitorInitRequest{})
	//BootStrap:            "",
	//Secret:               "",
	//IpfsPath:             "",
	//IpfsClusterPath:      "",
	//XXX_NoUnkeyedLiteral: struct{}{},
	//XXX_unrecognized:     nil,
	//XXX_sizecache:        0,

	t.Log(reply, err)
}
