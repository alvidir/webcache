package main

import (
	"net"
	"net/http"

	util "github.com/alvidir/go-util"
	"github.com/alvidir/webcache/logger"

	"github.com/joho/godotenv"
)

const (
	envPortKey = "SERVICE_PORT"
	envNetwKey = "SERVICE_NETW"
)

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Log.Fatalf(err.Error())
	}
}

func main() {
	network, err := util.LookupNempEnv(envNetwKey)
	if err != nil {

	}

	address, err := util.LookupNempEnv(envPortKey)
	if err != nil {

	}

	lis, err := net.Listen(network, address)
	if err != nil {

	}

	http.ListenAndServe(":8080", nil)
}
