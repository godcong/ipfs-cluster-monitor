package service

import (
	"context"
	"fmt"
	"github.com/godcong/ipfs-cluster-monitor/config"
	"github.com/godcong/ipfs-cluster-monitor/proto"
	"github.com/json-iterator/go"
	"golang.org/x/exp/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"syscall"
)

// GRPCServer ...
type GRPCServer struct {
	config *config.Configure
	server *grpc.Server
	Type   string
	Port   string
	Path   string
}

// MonitorInit ...
func (s *GRPCServer) MonitorInit(ctx context.Context, req *proto.MonitorInitRequest) (*proto.MonitorReply, error) {
	log.Println("monitor init call")
	monitor := config.MustMonitor(req.Secret, req.Bootstrap, req.Path, req.ClusterPath)
	err := server.cluster.InitMaker(monitor)
	if err != nil {
		log.Println(err)
		return &proto.MonitorReply{}, err
	}
	return Result("")
}

// MonitorProc ...
func (s *GRPCServer) MonitorProc(ctx context.Context, req *proto.MonitorProcRequest) (*proto.MonitorReply, error) {
	if req.Type == proto.MonitorType_Init {
		return &proto.MonitorReply{}, nil
	}
	return nil, xerrors.New("monitor proc error")
}

// Result ...
func Result(v interface{}) (*proto.MonitorReply, error) {
	detail, err := jsoniter.MarshalToString(v)
	if err != nil {
		return &proto.MonitorReply{
			Code:    -1,
			Message: err.Error(),
			Detail:  detail,
		}, err
	}
	return &proto.MonitorReply{
		Code:    0,
		Message: "success",
		Detail:  detail,
	}, nil
}

// NewGRPCServer ...
func NewGRPCServer(cfg *config.Configure) *GRPCServer {
	return &GRPCServer{
		config: cfg,
		Type:   config.DefaultString("", Type),
		Port:   config.DefaultString("", ":7784"),
		Path:   config.DefaultString("", "/tmp/monitor.sock"),
	}
}

// GRPCClient ...
type GRPCClient struct {
	config *config.Configure
	Type   string
	Port   string
	Addr   string
}

// Conn ...
func (c *GRPCClient) Conn() (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error

	if c.Type == "unix" {
		conn, err = grpc.Dial("passthrough:///unix://"+c.Addr, grpc.WithInsecure())
	} else {
		conn, err = grpc.Dial(c.Addr+c.Port, grpc.WithInsecure())
	}

	return conn, err
}

// MonitorClient ...
func MonitorClient(g *GRPCClient) proto.ClusterMonitorClient {
	clientConn, err := g.Conn()
	if err != nil {
		log.Println(err)
		return nil
	}
	return proto.NewClusterMonitorClient(clientConn)
}

// NewMonitorGRPC ...
func NewMonitorGRPC(cfg *config.Configure) *GRPCClient {
	return &GRPCClient{
		config: cfg,
		Type:   config.DefaultString(cfg.Monitor.Type, Type),
		Port:   config.DefaultString(cfg.Monitor.Port, ":7784"),
		Addr:   config.DefaultString(cfg.Monitor.Addr, "/tmp/monitor.sock"),
	}
}

// NewNodeGRPC ...
func NewNodeGRPC(cfg *config.Configure) *GRPCClient {
	return &GRPCClient{
		config: cfg,
		Type:   config.DefaultString("", Type),
		Port:   config.DefaultString("", ":7787"),
		Addr:   config.DefaultString("", "/tmp/node.sock"),
	}
}

// NewManagerGRPC ...
func NewManagerGRPC(cfg *config.Configure) *GRPCClient {
	return &GRPCClient{
		config: cfg,
		Type:   config.DefaultString("", Type),
		Port:   config.DefaultString("", ":7781"),
		Addr:   config.DefaultString("", "/tmp/manager.sock"),
	}
}

// NewCensorGRPC ...
func NewCensorGRPC(cfg *config.Configure) *GRPCClient {
	return &GRPCClient{
		config: cfg,
		Type:   config.DefaultString("", Type),
		Port:   config.DefaultString("", ":7785"),
		Addr:   config.DefaultString("", "/tmp/censor.sock"),
	}
}

// Start ...
func (s *GRPCServer) Start() {

	s.server = grpc.NewServer()
	var lis net.Listener
	var port string
	var err error
	go func() {
		if s.Type == "unix" {
			_ = syscall.Unlink(s.Path)
			lis, err = net.Listen(s.Type, s.Path)
			port = s.Path
		} else {
			lis, err = net.Listen("tcp", s.Port)
			port = s.Port
		}

		if err != nil {
			panic(fmt.Sprintf("failed to listen: %v", err))
		}

		proto.RegisterClusterMonitorServer(s.server, s)
		// Register reflection service on gRPC server.
		reflection.Register(s.server)
		log.Printf("Listening and serving TCP on %s\n", port)
		if err := s.server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

}

// Stop ...
func (s *GRPCServer) Stop() {
	s.server.Stop()
}
