package cluster

import (
	"context"
	"github.com/godcong/ipfs-cluster-monitor/config"
	"github.com/json-iterator/go"
	"github.com/juju/errors"
	"golang.org/x/exp/xerrors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// ServiceInfo ...
type ServiceInfo struct {
	ID                    string        `json:"id"`
	Addresses             []string      `json:"addresses"`
	ClusterPeers          []string      `json:"cluster_peers"`
	ClusterPeersAddresses []interface{} `json:"cluster_peers_addresses"`
	Version               string        `json:"version"`
	Commit                string        `json:"commit"`
	RPCProtocolVersion    string        `json:"rpc_protocol_version"`
	Error                 string        `json:"error"`
	Ipfs                  struct {
		ID        string   `json:"id"`
		Addresses []string `json:"addresses"`
		Error     string   `json:"error"`
	} `json:"ipfs"`
	Peername string `json:"peername"`
}

// ServiceConfig ...
type ServiceConfig struct {
	Cluster struct {
		ID                   string `json:"id"`
		Peername             string `json:"peername"`
		PrivateKey           string `json:"private_key"`
		Secret               string `json:"secret"`
		LeaveOnShutdown      bool   `json:"leave_on_shutdown"`
		ListenMultiaddress   string `json:"listen_multiaddress"`
		StateSyncInterval    string `json:"state_sync_interval"`
		IpfsSyncInterval     string `json:"ipfs_sync_interval"`
		ReplicationFactorMin int    `json:"replication_factor_min"`
		ReplicationFactorMax int    `json:"replication_factor_max"`
		MonitorPingInterval  string `json:"monitor_ping_interval"`
		PeerWatchInterval    string `json:"peer_watch_interval"`
		DisableRepinning     bool   `json:"disable_repinning"`
	} `json:"cluster"`
	Consensus struct {
		Raft struct {
			InitPeerset          []interface{} `json:"init_peerset"`
			WaitForLeaderTimeout string        `json:"wait_for_leader_timeout"`
			NetworkTimeout       string        `json:"network_timeout"`
			CommitRetries        int           `json:"commit_retries"`
			CommitRetryDelay     string        `json:"commit_retry_delay"`
			BackupsRotate        int           `json:"backups_rotate"`
			HeartbeatTimeout     string        `json:"heartbeat_timeout"`
			ElectionTimeout      string        `json:"election_timeout"`
			CommitTimeout        string        `json:"commit_timeout"`
			MaxAppendEntries     int           `json:"max_append_entries"`
			TrailingLogs         int           `json:"trailing_logs"`
			SnapshotInterval     string        `json:"snapshot_interval"`
			SnapshotThreshold    int           `json:"snapshot_threshold"`
			LeaderLeaseTimeout   string        `json:"leader_lease_timeout"`
		} `json:"raft"`
	} `json:"consensus"`
	API struct {
		Ipfsproxy struct {
			NodeMultiaddress   string `json:"node_multiaddress"`
			ListenMultiaddress string `json:"listen_multiaddress"`
			ReadTimeout        string `json:"read_timeout"`
			ReadHeaderTimeout  string `json:"read_header_timeout"`
			WriteTimeout       string `json:"write_timeout"`
			IdleTimeout        string `json:"idle_timeout"`
		} `json:"ipfsproxy"`
		Restapi struct {
			HTTPListenMultiaddress string      `json:"http_listen_multiaddress"`
			ReadTimeout            string      `json:"read_timeout"`
			ReadHeaderTimeout      string      `json:"read_header_timeout"`
			WriteTimeout           string      `json:"write_timeout"`
			IdleTimeout            string      `json:"idle_timeout"`
			BasicAuthCredentials   interface{} `json:"basic_auth_credentials"`
			Headers                struct {
				AccessControlAllowHeaders []string `json:"Access-Control-Allow-Headers"`
				AccessControlAllowMethods []string `json:"Access-Control-Allow-Methods"`
				AccessControlAllowOrigin  []string `json:"Access-Control-Allow-Origin"`
			} `json:"headers"`
		} `json:"restapi"`
	} `json:"api"`
	IpfsConnector struct {
		Ipfshttp struct {
			NodeMultiaddress   string `json:"node_multiaddress"`
			ConnectSwarmsDelay string `json:"connect_swarms_delay"`
			PinMethod          string `json:"pin_method"`
			IpfsRequestTimeout string `json:"ipfs_request_timeout"`
			PinTimeout         string `json:"pin_timeout"`
			UnpinTimeout       string `json:"unpin_timeout"`
		} `json:"ipfshttp"`
	} `json:"ipfs_connector"`
	PinTracker struct {
		Maptracker struct {
			MaxPinQueueSize int `json:"max_pin_queue_size"`
			ConcurrentPins  int `json:"concurrent_pins"`
		} `json:"maptracker"`
		Stateless struct {
			MaxPinQueueSize int `json:"max_pin_queue_size"`
			ConcurrentPins  int `json:"concurrent_pins"`
		} `json:"stateless"`
	} `json:"pin_tracker"`
	Monitor struct {
		Monbasic struct {
			CheckInterval string `json:"check_interval"`
		} `json:"monbasic"`
		Pubsubmon struct {
			CheckInterval string `json:"check_interval"`
		} `json:"pubsubmon"`
	} `json:"monitor"`
	Informer struct {
		Disk struct {
			MetricTTL  string `json:"metric_ttl"`
			MetricType string `json:"metric_type"`
		} `json:"disk"`
		Numpin struct {
			MetricTTL string `json:"metric_ttl"`
		} `json:"numpin"`
	} `json:"informer"`
}

// RunServiceInit ...
func RunServiceInit(ctx context.Context, cfg *config.Configure) error {
	cmd := exec.CommandContext(ctx, cfg.MonitorProperty.ClusterCommandName, "init")
	cmd.Env = cfg.Monitor.Env()

	bytes, err := cmd.CombinedOutput()
	log.Println(string(bytes))
	if err != nil {
		return xerrors.Errorf("first run ipfs:%w", err)
	}
	return nil
}

func optimizationFirstRunService(ctx context.Context, cfg *config.Configure) {
	err := optimizeRunCMD(ctx, cfg.MonitorProperty.ClusterCommandName, cfg.Monitor.Env(), "init")
	if err != nil {
		panic(err)
	}
}

// RunService ...
func RunService(ctx context.Context, cfg *config.Configure) {
	log.Println("bootstrap", cfg.Monitor.Bootstrap)
	go optimizeRunCMD(ctx, cfg.MonitorProperty.ClusterCommandName, cfg.Monitor.Env(), "daemon", "--bootstrap", cfg.Monitor.Bootstrap)

}

func getPeers() ([]ServiceInfo, error) {
	response, err := http.Get("http://localhost:9094/peers")
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var peers []ServiceInfo

	err = jsoniter.Unmarshal(bytes, &peers)
	if err != nil {
		return nil, err
	}

	//monitor.Store("peers", peers)

	return peers, nil
}

// DeletePeers ...
func DeletePeers(peerID string) error {
	request, err := http.NewRequest(http.MethodDelete, "http://localhost:9094/peers/"+peerID, nil)
	if err != nil {
		errors.ErrorStack(err)
		return err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		errors.ErrorStack(err)
		return err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errors.ErrorStack(err)
		return err
	}
	log.Println(string(bytes))
	return nil
}

// GetServiceInfo ...
func GetServiceInfo() (*ServiceInfo, error) {
	return getServiceInfo()
}

func getServiceInfo() (*ServiceInfo, error) {
	var service ServiceInfo
	response, err := http.Get("http://localhost:9094/id")
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = jsoniter.Unmarshal(bytes, &service)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

// WaitingService ...
func WaitingService(ctx context.Context) {
	var err error
	for {
		select {
		case <-ctx.Done():
			log.Println("monitor done ")
			return
		default:
			log.Println("waiting service")
			time.Sleep(1 * time.Second)
			_, err = getServiceInfo()
			if err == nil {
				return
			}
		}
	}
}

// GetServiceConfig ...
func GetServiceConfig() (*ServiceConfig, error) {
	var serviceConfig ServiceConfig

	file := filepath.Join(config.Config().Root, "config.toml")
	openFile, err := os.OpenFile(file, os.O_RDONLY|os.O_SYNC, os.ModePerm)
	if err != nil {
		return nil, err
	}
	dec := jsoniter.NewDecoder(openFile)
	err = dec.Decode(&serviceConfig)
	if err != nil {
		return nil, err
	}
	return &serviceConfig, nil
}
