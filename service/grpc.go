package service

import (
	"context"
	"fmt"
	"github.com/godcong/ipfs-cluster-monitor/config"
	"github.com/godcong/ipfs-cluster-monitor/proto"
	"github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
	"syscall"
)

// GRPCServer ...
type GRPCServer struct {
	config                   *config.Configure
	server                   *grpc.Server
	requireTransportSecurity bool
	Type                     string
	Port                     string
	Path                     string
}

// MonitorAddress ...
func (s *GRPCServer) MonitorAddress(context.Context, *proto.MonitorRequest) (*proto.MonitorAddressReply, error) {
	return &proto.MonitorAddressReply{}, nil
}

//MonitorManager ...
func (s *GRPCServer) MonitorManager(ctx context.Context, in *proto.MonitorManagerRequest) (*proto.MonitorReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "无Token认证信息")
	}
	log.Println(md)
	log.Println(in.Type, in.Data)
	return Result("")
}

// MonitorBootstrap ...
func (s *GRPCServer) MonitorBootstrap(context.Context, *proto.MonitorRequest) (*proto.MonitorBootstrapReply, error) {
	return &proto.MonitorBootstrapReply{
		Bootstraps: nil,
	}, nil
}

// MonitorPin ...
func (s *GRPCServer) MonitorPin(context.Context, *proto.MonitorRequest) (*proto.MonitorPinReply, error) {
	return &proto.MonitorPinReply{
		Pins: nil,
	}, nil
}

// MonitorInit ...
func (s *GRPCServer) MonitorInit(ctx context.Context, req *proto.MonitorInitRequest) (*proto.MonitorReply, error) {
	log.Println("monitor init call")
	monitor := config.MustMonitor(req.StartMode, req.Secret, req.Bootstrap, req.Workspace)
	log.Printf("%+v", monitor)
	config.SetMonitor(monitor)
	log.Printf("%+v", s.config)
	if s.config.Initialize {
		return &proto.MonitorReply{
			Code:    0,
			Message: "initialized",
			Detail:  "",
		}, nil
	}
	err := server.cluster.InitMaker(monitor)
	if err != nil {
		log.Println(err)
		return nil, xerrors.Errorf("monitor init %w", err)
	}
	return Result("")
}

// MonitorProc ...
func (s *GRPCServer) MonitorProc(ctx context.Context, req *proto.MonitorProcRequest) (*proto.MonitorReply, error) {
	if req.Type == proto.MonitorType_Init {
		return &proto.MonitorReply{}, nil
	} else if req.Type == proto.MonitorType_Reset {
		err := server.cluster.Reset()
		if err != nil {
			log.Println(err)
			return nil, xerrors.Errorf("monitor proc %w", err)
		}
	}
	return Result("")
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
		server: grpc.NewServer(),
		Type:   config.DefaultString("", GRPCType),
		Port:   config.DefaultString("", ":7784"),
		Path:   config.DefaultString("", "/tmp/monitor.sock"),
	}
}

// GRPCClient ...
type GRPCClient struct {
	config                   *config.Configure
	requireTransportSecurity bool
	Type                     string
	Port                     string
	Addr                     string
}

// GetRequestMetadata ...
func (c *GRPCClient) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"token": c.config.Monitor.Token,
	}, nil
}

// RequireTransportSecurity ...
func (c *GRPCClient) RequireTransportSecurity() bool {
	return c.requireTransportSecurity
}

// Conn ...
func (c *GRPCClient) Conn() (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var opts []grpc.DialOption
	var err error
	if c.RequireTransportSecurity() {
		cred, err := credentials.NewClientTLSFromFile("./keys/server.pem", "GodCong")
		if err != nil {
			grpclog.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(cred), grpc.WithPerRPCCredentials(c))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	if c.Type == "unix" {
		conn, err = grpc.Dial("passthrough:///unix://"+c.Addr, opts...)
	} else {
		conn, err = grpc.Dial(c.Addr+c.Port, opts...)
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
		Type:   config.DefaultString(cfg.Monitor.Type, GRPCType),
		Port:   config.DefaultString(cfg.Monitor.Port, ":7784"),
		Addr:   config.DefaultString(cfg.Monitor.Addr, "/tmp/monitor.sock"),
	}
}

// NewServerMonitorGRPC ...
func NewServerMonitorGRPC(cfg *config.Monitor) *GRPCClient {
	return &GRPCClient{
		Type: config.DefaultString(cfg.Type, GRPCType),
		Port: config.DefaultString(cfg.Port, ":7774"),
		Addr: config.DefaultString(cfg.Addr, "/tmp/server-monitor.sock"),
	}
}

// NewNodeGRPC ...
func NewNodeGRPC(cfg *config.Configure) *GRPCClient {
	return &GRPCClient{
		config: cfg,
		Type:   config.DefaultString("", GRPCType),
		Port:   config.DefaultString("", ":7787"),
		Addr:   config.DefaultString("", "/tmp/node.sock"),
	}
}

// NewManagerGRPC ...
func NewManagerGRPC(cfg *config.Configure) *GRPCClient {
	return &GRPCClient{
		config: cfg,
		Type:   config.DefaultString("", GRPCType),
		Port:   config.DefaultString("", ":7781"),
		Addr:   config.DefaultString("", "/tmp/manager.sock"),
	}
}

// NewCensorGRPC ...
func NewCensorGRPC(cfg *config.Configure) *GRPCClient {
	return &GRPCClient{
		config: cfg,
		Type:   config.DefaultString("", GRPCType),
		Port:   config.DefaultString("", ":7785"),
		Addr:   config.DefaultString("", "/tmp/censor.sock"),
	}
}

// Start ...
func (s *GRPCServer) Start() {
	if !s.config.GRPC.Enable {
		return
	}
	var err error
	var lis net.Listener
	var port string

	if s.requireTransportSecurity {
		cred, err := credentials.NewServerTLSFromFile("./keys/server.pem", "./keys/server.key")
		if err != nil {
			grpclog.Fatalf("Failed to generate credentials %v", err)
		}
		s.server = grpc.NewServer(grpc.Creds(cred))
	} else {
		s.server = grpc.NewServer()
	}

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
