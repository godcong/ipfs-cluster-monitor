package cluster

import "github.com/godcong/ipfs-cluster-monitor/config"

// IsInitialized ...
func IsInitialized() bool {
	if cluster.isInitialized == false {
		cluster.isInitialized = config.Config().MonitorProperty.CheckExist()
	}
	return cluster.isInitialized
}
