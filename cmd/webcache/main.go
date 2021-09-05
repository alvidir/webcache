package main

import (
	"net"
	"net/http"

	util "github.com/alvidir/go-util"
	wcache "github.com/alvidir/webcache"
	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

const (
	envAddrKey = "SERVICE_ADDR"
	envNetwKey = "SERVICE_NETW"
)

func main() {
	if err := godotenv.Load(); err != nil {
		wcache.Log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("No dotenv file has been found")
	}

	network, err := util.LookupNempEnv(envNetwKey)
	if err != nil {
		wcache.Log.WithFields(log.Fields{
			"error": err.Error(),
			"var":   envNetwKey,
		}).Fatal("Environment variable must be set")
	}

	address, err := util.LookupNempEnv(envAddrKey)
	if err != nil {
		wcache.Log.WithFields(log.Fields{
			"error": err.Error(),
			"var":   envAddrKey,
		}).Fatal("Environment variable must be set")
	}

	lis, err := net.Listen(network, address)
	if err != nil {
		wcache.Log.WithFields(log.Fields{
			"error":   err.Error(),
			"network": network,
			"address": address,
		}).Fatal("Failed on listening")
	}

	wcache.Log.WithFields(log.Fields{
		"network": network,
		"address": address,
	}).Info("Server listening and ready")

	if err := http.Serve(lis, nil); err != nil {
		wcache.Log.WithFields(log.Fields{
			"error":   err.Error(),
			"network": network,
			"address": address,
		}).Fatal("Failed serving")
	}
}
