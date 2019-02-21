package main

import (
	"context"
	"flag"
	"github.com/godcong/ipfs-cluster-monitor/config"
	"github.com/godcong/ipfs-cluster-monitor/proto"
	"github.com/godcong/ipfs-cluster-monitor/service"
	"log"
)

var ws = flag.String("ws", "/tmp/workspace", "set the workspace path")

func main() {
	flag.Parse()
	cfg := config.DefaultConfig("")
	cfg.Monitor.Addr = "localhost"
	cfg.Monitor.Type = "tcp"
	grpc := service.NewMonitorGRPC(cfg)

	client := service.MonitorClient(grpc)
	reply, err := client.MonitorInit(context.Background(), &proto.MonitorInitRequest{
		StartMode: proto.StartMode_Cluster,
		Workspace: *ws,
	})
	log.Println(reply, err)
}
