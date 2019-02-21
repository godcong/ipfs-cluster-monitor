package service

import (
	"context"
	"github.com/godcong/ipfs-cluster-monitor/config"
	"github.com/godcong/ipfs-cluster-monitor/proto"
	"testing"
)

// TestMonitorClient_Init ...
func TestMonitorClient_Init(t *testing.T) {
	cfg := config.DefaultConfig("")
	cfg.Monitor.Addr = "localhost"
	cfg.Monitor.Type = "tcp"
	grpc := NewMonitorGRPC(cfg)

	client := MonitorClient(grpc)
	reply, err := client.MonitorInit(context.Background(), &proto.MonitorInitRequest{
		StartMode: proto.StartMode_Simple,
		Host:      "http://localhost:8081",
		Bootstrap: "/ip4/47.101.169.94/tcp/9096/ipfs/QmdpBCokb3XBZL5o79X8MaxatPQWxPhBZmmV7pGP13gRmL",
		Secret:    "",
		Workspace: "d:\\workspace\\ipfs2",
		//Workspace: "/storage/1A247F77247F54AB/ws/",
	})
	t.Log(reply, err)
}

// TestClusterMonitor_Initialized ...
func TestMonitorClient_Proc(t *testing.T) {
	cfg := config.DefaultConfig("")
	//cfg.Monitor.Addr = "192.168.1.183"
	//cfg.Monitor.Type = "tcp"
	grpc := NewMonitorGRPC(cfg)
	client := MonitorClient(grpc)
	reply, err := client.MonitorProc(context.Background(), &proto.MonitorProcRequest{
		Type: proto.MonitorType_Reset,
		//	Workspace: "d:\\workspace\\ipfs2",
		//	BootStrap: "/ip4/47.101.169.94/tcp/9096/ipfs/Qmc8XTmaXivEuFQLL4m2GSw8BurZKvf34rEXQu8PLfCWii",
	})

	t.Log(reply, err)
}

// TestMonitorClient_Manager ...
func TestMonitorClient_Manager(t *testing.T) {
	cfg := config.DefaultConfig("")
	//cfg.Monitor.Addr = "192.168.1.183"
	cfg.Monitor.Type = "tcp"
	grpc := NewMonitorGRPC(cfg)
	client := MonitorClient(grpc)
	reply, err := client.MonitorManager(context.Background(), &proto.MonitorManagerRequest{
		Type: proto.ManagerType_BootstrapAdd,
		Data: []string{"data1", "data2"},
	})

	t.Log(reply, err)
}
