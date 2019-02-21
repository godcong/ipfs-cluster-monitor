package main

import (
	"context"
	"flag"
	"github.com/godcong/ipfs-cluster-monitor/config"
	"github.com/godcong/ipfs-cluster-monitor/proto"
	"github.com/godcong/ipfs-cluster-monitor/service"
	log "github.com/sirupsen/logrus"
	"os"
)

var bootstrap = flag.String("bootstrap", "", "monitor manager bootstrap address")
var pin = flag.String("pin", "", "monitor manager pin address")
var del = flag.Bool("delete", false, "is delete")
var configPath = flag.String("config", "config.toml", "config path")

func main() {
	flag.Parse()
	err := config.Initialize(os.Args[0], *configPath)
	if err != nil {
		panic(err)
	}
	grpc := service.NewMonitorGRPC(config.Config())
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
