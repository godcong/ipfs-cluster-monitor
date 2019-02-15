package main

import (
	"context"
	"flag"
	"github.com/godcong/ipfs-cluster-monitor/config"
	"github.com/godcong/ipfs-cluster-monitor/proto"
	"github.com/godcong/ipfs-cluster-monitor/service"
	log "github.com/sirupsen/logrus"
)

var bootstrap = flag.String("bootstrap", "", "monitor manager bootstrap address")
var pin = flag.String("pin", "", "monitor manager pin address")
var del = flag.Bool("delete", false, "is delete")

func main() {
	flag.Parse()
	cfg := config.DefaultConfig()
	//cfg.Monitor.Addr = "192.168.1.183"
	//cfg.Monitor.Type = "tcp"
	grpc := service.NewMonitorGRPC(cfg)
	var e error
	var reply *proto.MonitorReply
	client := service.MonitorClient(grpc)
	if *del {
		if *pin != "" {
			reply, e = client.MonitorManager(context.Background(), &proto.MonitorManagerRequest{
				Type: proto.ManagerType_PinAdd,
				Data: []string{*pin},
			})
			log.Println(reply, e)
			return
		}
		if *bootstrap != "" {
			reply, e = client.MonitorManager(context.Background(), &proto.MonitorManagerRequest{
				Type: proto.ManagerType_BootstrapAdd,
				Data: []string{*bootstrap},
			})
			log.Println(reply, e)
			return
		}
	} else {
		if *pin != "" {
			reply, e = client.MonitorManager(context.Background(), &proto.MonitorManagerRequest{
				Type: proto.ManagerType_PinRemove,
				Data: []string{*pin},
			})
			log.Println(reply, e)
			return
		}
		if *bootstrap != "" {
			reply, e = client.MonitorManager(context.Background(), &proto.MonitorManagerRequest{
				Type: proto.ManagerType_BootstrapRemove,
				Data: []string{*bootstrap},
			})
			log.Println(reply, e)
			return
		}
	}

}
