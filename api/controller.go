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
			cluster.Config().SetEnv("IPFS_PATH", ipfs)
			cluster.Config().SetEnv("IPFS_CLUSTER_PATH", service)
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
		var s cluster.Peer
		resp, err := http.Get("http://127.0.0.1:9094/id")
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
		id := ctx.Param("id")
		peers := cluster.GetPeers()
		size := len(peers)
		for i := 0; i < size; i++ {
			if peers[0].ID == id {
				err := cluster.DeletePeers(id)
				if err != nil {
					failed(ctx, err.Error())
					return
				}
				success(ctx, nil)
			}
		}
		failed(ctx, "peers not found")
	}

}
