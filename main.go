package main

import (
	proto "github.com/micro/go-config/source/grpc/proto"
	"github.com/micro/go-log"
	"sync"
)

var (
	mux        sync.RWMutex
	configMaps = make(map[string]*proto.ChangeSet)
	apps       = []string{"micro"}
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			log.Logf("[main] Recovered in f %v", r)
		}
	}()
	err := loadAndWatchConfigFile()
	if err != nil {
		log.Fatal(err)
	}

}
