package api

import (
	"github.com/gin-gonic/gin"
	"github.com/godcong/ipfs-cluster-monitor/cluster"
	"github.com/json-iterator/go"
	"io/ioutil"
	"log"
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
func InitPost(s string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		remote := ctx.PostForm("Remote")
		secret := ctx.PostForm("Secret")
		clusterSecret := ctx.PostForm("CLUSTER_SECRET")
		ipfs := ctx.PostForm("IPFS_PATH")
		service := ctx.PostForm("IPFS_CLUSTER_PATH")
		//monitor := ctx.PostForm("IPFS_CLUSTER_MONITOR")
		if !cluster.IsInitialized() {
			if remote != "" {
				cluster.Config().SetClient(remote)
				cluster.Config().Secret = secret
			} else {
				cluster.Config().Secret = prefix + cluster.GenerateRandomString(64)
			}
			cluster.Config().SetEnv(cluster.EnvironSecret(clusterSecret))
			cluster.Config().SetEnv(cluster.EnvironIPFS(ipfs))
			cluster.Config().SetEnv(cluster.EnvironService(service))
			cluster.Config().Make()

			log.Println("host initialized")
			success(ctx, cluster.Config())
			return
		}

		failed(ctx, "host is initialized")
	}
}

// HeartBeatGet ...
func HeartBeatGet(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		success(ctx, cluster.GetPeers())
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

// ResetGet ...
func ResetGet(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := cluster.Reset()
		if err != nil {
			failed(ctx, err.Error())
			return
		}

		//time.Sleep(5 * time.Second)
		success(ctx, nil)
	}
}

// DeleteGet ...
func DeleteGet(s string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pn := ctx.Param("pn")
		peers := cluster.GetPeers()
		size := len(peers)
		for i := 0; i < size; i++ {
			if peers[0].Peername == pn {
				err := cluster.DeletePeers(peers[0].ID)
				if err != nil {
					failed(ctx, err.Error())
					return
				}
				success(ctx, nil)
				return
			}
		}
		failed(ctx, "peers not found")
	}

}

// KillGet ...
func KillGet(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		failed(ctx, "I will live forever")
	}
}

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

		name := ctx.PostForm("name")
		address := ctx.PostForm("address")
		if address == "" {
			address = ctx.Request.RemoteAddr
		}
		cluster.JoinFromClient(name, address)
		success(ctx, nil)
	}
}
