package service

import "github.com/godcong/ipfs-cluster-monitor/config"

// service ...
type service struct {
	config *config.Configure
	grpc   *GRPCServer
	//rest    *RestServer
	cluster *ClusterMonitor
}

var server *service

// Start ...
func Start() {
	cfg := config.Config()

	server = &service{
		config:  cfg,
		cluster: NewClusterMonitor(cfg),
		grpc:    NewGRPCServer(cfg),
		//rest:    NewRestServer(cfg),
		//queue: NewQueueServer(cfg),
	}

	//server.rest.Start()
	server.grpc.Start()
	server.cluster.Start()
}

// Stop ...
func Stop() {
	//server.rest.Stop()
	server.grpc.Stop()
	server.cluster.Stop()
}
