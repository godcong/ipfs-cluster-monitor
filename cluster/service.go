package cluster

import (
	"context"
	"github.com/json-iterator/go"
	"github.com/juju/errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

// ServicePeer ...
type ServicePeer struct {
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

func firstRunService() {
	cmd := exec.Command(cfg.ServiceCommandName, "init")
	cmd.Env = cfg.GetEnv()

	bytes, err := cmd.CombinedOutput()
	log.Println(string(bytes))
	if err != nil {
		errors.ErrorStack(err)
		panic(err)
	}
}

func optimizationFirstRunService(ctx context.Context) {
	err := optimizeRunCMD(cfg.ServiceCommandName, "init")
	if err != nil {
		panic(err)
	}
}

// runService ...
func runService(ctx context.Context) {
	if isClient() {
		boot := getServiceBootstrap()
		if boot != "" {
			log.Println("bootstrap")
			go optimizeRunCMD(cfg.ServiceCommandName, "daemon", "--bootstrap", boot)
			return
		}
	}
	go optimizeRunCMD(cfg.ServiceCommandName, "daemon")
}

func getPeers() ([]ServicePeer, error) {
	response, err := http.Get("http://localhost:9094/peers")
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var peers []ServicePeer

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

// GetPeers ...
func GetPeers() []ServicePeer {
	if peers, b := monitor.Load(MonitorPeers); b {
		return peers.([]ServicePeer)
	}
	return nil
}

func servicePath() string {
	if cfg.Environ.Service != "" {
		return string(cfg.Environ.Service)
	}
	return defaultPath(".ipfs-cluster")
}

func GetServiceConfig() (*ServiceConfig, error) {
	var serviceConfig ServiceConfig

	file := filepath.Join(servicePath(), "service.json")
	openFile, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
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
