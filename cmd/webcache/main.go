package main

import (
	"net"
	"net/http"
	"regexp"

	"log"

	util "github.com/alvidir/go-util"
	wcache "github.com/alvidir/webcache"
	"github.com/fsnotify/fsnotify"

	"github.com/joho/godotenv"
)

const (
	envAddrKey     = "SERVICE_ADDR"
	envNetwKey     = "SERVICE_NETW"
	envConfPath    = "CONFIG_PATH"
	envWatchConfig = "WATCH_CONFIG"
)

var (
	configPath = "/etc/webcache"
	config     = wcache.NewConfig()
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("no dotenv file has been found")
	}

	if path, err := util.LookupEnv(envConfPath); err == nil {
		configPath = path
	}

	if err := config.ReadFiles(configPath); err != nil {
		log.Println(err)
	}

	if watch, err := util.LookupEnv(envWatchConfig); err == nil {
		if match, err := regexp.MatchString("^(True|true)$", watch); err == nil && match {
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				log.Fatal(err)
			}

			defer watcher.Close()
			go config.AttachWatcher(watcher)

			err = watcher.Add(configPath)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	network, err := util.LookupEnv(envNetwKey)
	if err != nil {
		log.Fatalf("%s: %s", envNetwKey, err)
	}

	address, err := util.LookupEnv(envAddrKey)
	if err != nil {
		log.Fatalf("%s: %s", envAddrKey, err)
	}

	lis, err := net.Listen(network, address)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", wcache.NewHandler(config))
	log.Printf("server listening on %s", address)
	if err := http.Serve(lis, nil); err != nil {
		log.Fatal(err)
	}
}
