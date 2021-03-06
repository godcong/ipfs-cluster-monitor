package service

import (
	"github.com/gin-gonic/gin"
	"github.com/godcong/ipfs-cluster-monitor/cluster"
	"github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

const prefix = "ICM"

func result(ctx *gin.Context, code int, message string, detail interface{}) {
	h := gin.H{
		"code":    code,
		"message": message,
		"detail":  detail,
	}
	ctx.JSON(http.StatusOK, h)
}

func success(ctx *gin.Context, detail interface{}) {
	result(ctx, 0, "success", detail)
}

func failed(ctx *gin.Context, message string) {
	result(ctx, -1, message, nil)
}

// InitPost ...
func InitPost(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//remote := ctx.PostForm("Remote")
		//secret := ctx.PostForm("Secret")
		//clusterSecret := ctx.PostForm("CLUSTER_SECRET")
		//ipfs := ctx.PostForm("IPFS_PATH")
		//service := ctx.PostForm("IPFS_CLUSTER_PATH")
		////client := ctx.PostForm("IPFS_CLUSTER_MONITOR")
		//if !monitor {
		//	if remote != "" {
		//		monitor.Config().SetClient(remote)
		//		monitor.Config().MonitorSecret = secret
		//	} else {
		//		monitor.Config().MonitorSecret = prefix + monitor.GenerateRandomString(64)
		//	}
		//	monitor.Config().SetEnv(monitor.EnvironSecret(clusterSecret))
		//	monitor.Config().SetEnv(monitor.EnvironIPFS(ipfs))
		//	monitor.Config().SetEnv(monitor.EnvironService(service))
		//	monitor.Config().Make()
		//	monitor.Default().SetStatus("init", monitor.StatusCreated)
		//
		//	log.Println("host initialized")
		//	success(ctx, monitor.Config())
		//	return
		//}

		failed(ctx, "host is initialized")

	}
}

// HeartBeatGet ...
func HeartBeatGet(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		detail := gin.H{
			"ipfs_status":    "failed",
			"service_status": "failed",
		}
		ipfs, err := cluster.GetIpfsInfo()
		if err == nil {
			detail["ipfs_status"] = "success"
			detail["ipfs"] = ipfs
		}
		service, err := cluster.GetServiceInfo()
		if err == nil {
			detail["service_status"] = "success"
			detail["service"] = service
		}
		success(ctx, detail)
	}
}

// LogGet ...
func LogGet(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

// BootstrapGet ...
func BootstrapGet(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var s cluster.ServiceInfo
		resp, err := http.Get("http://127.0.0.1:9094/id")
		log.Println(ctx.Request.Header.Get("secret"))
		if err != nil {
			failed(ctx, err.Error())
			return
		}
		bytes, err := ioutil.ReadAll(resp.Body)
		err = jsoniter.Unmarshal(bytes, &s)
		if err != nil {
			failed(ctx, err.Error())
			return
		}
		address := ""

		for _, value := range s.Addresses {
			if strings.Index(value, "/p2p-circuit") >= 0 {
				continue
			} else if strings.Index(value, "/ip4/127.0.0.1") >= 0 {
				continue
			}
			address = value
			break
		}
		success(ctx, gin.H{"bootstrap": address})
	}
}

// WaitingGet ...
func WaitingGet(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

// ResetGet ...
func ResetGet(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		failed(ctx, "can't reset now")
	}
}

// DeleteGet ...
func DeleteGet(s string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//pn := ctx.Param("pn")
		//peers := monitor.GetPeers()
		//size := len(peers)
		//for i := 0; i < size; i++ {
		//	if peers[0].Peername == pn {
		//		err := monitor.DeletePeers(peers[0].ID)
		//		if err != nil {
		//			failed(ctx, err.Error())
		//			return
		//		}
		//		success(ctx, nil)
		//		return
		//	}
		//}
		failed(ctx, "peers not found")
	}

}

// KillGet ...
func KillGet(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		failed(ctx, "I will live forever")
	}
}

// SecretGet ...
func SecretGet(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		config, err := cluster.GetServiceConfig()
		if err != nil {
			failed(ctx, err.Error())
			return
		}
		success(ctx, gin.H{"secret": config.Cluster.Secret})
	}
}

// JoinPost ...
func JoinPost(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//name := ctx.PostForm("name")
		address := ctx.PostForm("address")
		if address == "" {
			address = ctx.Request.RemoteAddr
		}
		//monitor.AddMySon(name, address)
		success(ctx, gin.H{})
	}
}
