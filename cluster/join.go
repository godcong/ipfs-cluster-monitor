package cluster

import (
	"context"
	"log"
)

//runJoin join remote server
func runJoin(ctx context.Context) {
	//	ServerCheckInterval
}

func joinToServer() {

}

// JoinFromClient ...
func JoinFromClient(key, val string) {
	if _, loaded := clients.LoadOrStore(key, val); loaded {
		log.Println(key, "is already joined")
	} else {
		log.Println(key, "is joined")
	}
}
