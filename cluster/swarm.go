package cluster

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

var swarm = []string{
	`/key/swarm/psk/1.0.0/`,
	`/bin/`,
	`082e236c727fb4da9ed1652ee3d974feadf4ddf6ee21e3438cf71ae5665a409b`,
}

// WriterSwarm ...
func WriterSwarm(path string) {
	_, e := os.Stat(path)
	if e != nil {
		return
	}
	openFile, e := os.OpenFile(filepath.Join(path, "swarm.key"), os.O_CREATE|os.O_SYNC|os.O_RDWR, os.ModePerm)
	if e != nil {
		log.Errorln(e)
		return
	}
	for _, v := range swarm {
		openFile.WriteString(v)
		openFile.WriteString("\n")
	}
}
