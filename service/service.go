package service

import "github.com/godcong/ipfs-cluster-monitor/config"

// service ...
type service struct {
	config  *config.Configure
	grpc    *GRPCServer
	monitor *GRPCClient
	cluster *Monitor
}

var server *service

// Start ...
func Start() {
	cfg := config.Config()

	server = &service{
		config:  cfg,
		monitor: NewMonitorGRPC(cfg),
		cluster: NewMonitor(cfg),
		grpc:    NewGRPCServer(cfg),
		//rest:    NewRestServer(cfg),
		//queue: NewQueueServer(cfg),
	}

	server.grpc.Start()
	server.cluster.Start()

}

// Stop ...
func Stop() {
	//server.rest.Stop()
	server.grpc.Stop()
	server.cluster.Stop()
}
