package main

import (
	"net"
	"net/http"

	"log"

	util "github.com/alvidir/go-util"
	wcache "github.com/alvidir/webcache"
	"github.com/fsnotify/fsnotify"

	"github.com/joho/godotenv"
)

const (
	envAddrKey  = "SERVICE_ADDR"
	envNetwKey  = "SERVICE_NETW"
	envConfPath = "CONFIG_PATH"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No dotenv file has been found")
	}

	config, err := util.LookupNempEnv(envConfPath)
	if err != nil {
		log.Fatalf("%s: %s", envConfPath, err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf(err.Error())
	}

	defer watcher.Close()
	go wcache.HandleConfigWatcher(watcher)

	err = watcher.Add(config)
	if err != nil {
		log.Fatalf(err.Error())
	}

	network, err := util.LookupNempEnv(envNetwKey)
	if err != nil {
		log.Fatalf("%s: %s", envNetwKey, err)
	}

	address, err := util.LookupNempEnv(envAddrKey)
	if err != nil {
		log.Fatalf("%s: %s", envAddrKey, err)
	}

	lis, err := net.Listen(network, address)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("server listening on %s", address)
	if err := http.Serve(lis, nil); err != nil {
		log.Fatalf(err.Error())
	}
}
