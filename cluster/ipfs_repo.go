package cluster

type FileConfig struct {
	Identity struct {
		PeerID  string `json:"PeerID"`
		PrivKey string `json:"PrivKey"`
	} `json:"Identity"`
	Datastore struct {
		StorageMax         string `json:"StorageMax"`
		StorageGCWatermark int    `json:"StorageGCWatermark"`
		GCPeriod           string `json:"GCPeriod"`
		Spec               struct {
			Mounts []struct {
				Child struct {
					Path      string `json:"path"`
					ShardFunc string `json:"shardFunc"`
					Sync      bool   `json:"sync"`
					Type      string `json:"type"`
				} `json:"child"`
				Mountpoint string `json:"mountpoint"`
				Prefix     string `json:"prefix"`
				Type       string `json:"type"`
			} `json:"mounts"`
			Type string `json:"type"`
		} `json:"Spec"`
		HashOnRead      bool `json:"HashOnRead"`
		BloomFilterSize int  `json:"BloomFilterSize"`
	} `json:"Datastore"`
	Addresses struct {
		Swarm      []string      `json:"Swarm"`
		Announce   []interface{} `json:"Announce"`
		NoAnnounce []interface{} `json:"NoAnnounce"`
		API        string        `json:"API"`
		Gateway    string        `json:"Gateway"`
	} `json:"Addresses"`
	Mounts struct {
		IPFS           string `json:"IPFS"`
		IPNS           string `json:"IPNS"`
		FuseAllowOther bool   `json:"FuseAllowOther"`
	} `json:"Mounts"`
	Discovery struct {
		MDNS struct {
			Enabled  bool `json:"Enabled"`
			Interval int  `json:"Interval"`
		} `json:"MDNS"`
	} `json:"Discovery"`
	Routing struct {
		Type string `json:"Type"`
	} `json:"Routing"`
	Ipns struct {
		RepublishPeriod  string `json:"RepublishPeriod"`
		RecordLifetime   string `json:"RecordLifetime"`
		ResolveCacheSize int    `json:"ResolveCacheSize"`
	} `json:"Ipns"`
	Bootstrap []string `json:"Bootstrap"`
	Gateway   struct {
		HTTPHeaders struct {
			AccessControlAllowHeaders []string `json:"Access-Control-Allow-Headers"`
			AccessControlAllowMethods []string `json:"Access-Control-Allow-Methods"`
			AccessControlAllowOrigin  []string `json:"Access-Control-Allow-Origin"`
		} `json:"HTTPHeaders"`
		RootRedirect string        `json:"RootRedirect"`
		Writable     bool          `json:"Writable"`
		PathPrefixes []interface{} `json:"PathPrefixes"`
		APICommands  []interface{} `json:"APICommands"`
		NoFetch      bool          `json:"NoFetch"`
	} `json:"Gateway"`
	API struct {
		HTTPHeaders struct {
		} `json:"HTTPHeaders"`
	} `json:"API"`
	Swarm struct {
		AddrFilters             interface{} `json:"AddrFilters"`
		DisableBandwidthMetrics bool        `json:"DisableBandwidthMetrics"`
		DisableNatPortMap       bool        `json:"DisableNatPortMap"`
		DisableRelay            bool        `json:"DisableRelay"`
		EnableRelayHop          bool        `json:"EnableRelayHop"`
		EnableAutoRelay         bool        `json:"EnableAutoRelay"`
		EnableAutoNATService    bool        `json:"EnableAutoNATService"`
		ConnMgr                 struct {
			Type        string `json:"Type"`
			LowWater    int    `json:"LowWater"`
			HighWater   int    `json:"HighWater"`
			GracePeriod string `json:"GracePeriod"`
		} `json:"ConnMgr"`
	} `json:"Swarm"`
	Pubsub struct {
		Router                      string `json:"Router"`
		DisableSigning              bool   `json:"DisableSigning"`
		StrictSignatureVerification bool   `json:"StrictSignatureVerification"`
	} `json:"Pubsub"`
	Reprovider struct {
		Interval string `json:"Interval"`
		Strategy string `json:"Strategy"`
	} `json:"Reprovider"`
	Experimental struct {
		FilestoreEnabled     bool `json:"FilestoreEnabled"`
		UrlstoreEnabled      bool `json:"UrlstoreEnabled"`
		ShardingEnabled      bool `json:"ShardingEnabled"`
		Libp2PStreamMounting bool `json:"Libp2pStreamMounting"`
		P2PHTTPProxy         bool `json:"P2pHttpProxy"`
		QUIC                 bool `json:"QUIC"`
	} `json:"Experimental"`
}
