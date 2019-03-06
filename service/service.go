package service

import "github.com/godcong/ipfs-cluster-monitor/config"

// service ...
type service struct {
	config  *config.Configure
	grpc    *GRPCServer
	client  *GRPCClient
	monitor *Monitor
}

var server *service

// Start ...
func Start() {
	cfg := config.Config()

	server = &service{
		config:  cfg,
		client:  NewMonitorGRPC(cfg),
		monitor: NewMonitor(cfg),
		grpc:    NewGRPCServer(cfg),
		//rest:    NewRestServer(cfg),
		//queue: NewQueueServer(cfg),
	}

	server.grpc.Start()
	server.monitor.Start()

}

// Stop ...
func Stop() {
	//server.rest.Stop()
	server.grpc.Stop()
	server.monitor.Stop()
}
