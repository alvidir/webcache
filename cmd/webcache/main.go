package main

import (
	"context"
	"encoding/base64"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"log"

	util "github.com/alvidir/go-util"
	wcache "github.com/alvidir/webcache"

	"github.com/joho/godotenv"
)

const (
	ENV_SERVICE_ADDR    = "SERVICE_ADDR"
	ENV_SERVICE_NETWORK = "SERVICE_NETWORK"
	ENV_CONFIG_PATH     = "CONFIG_PATH"
	ENV_WATCH_CONFIG    = "WATCH_CONFIG"
	ENV_REDIS_ADDR      = "REDIS_ADDR"
	ENV_CACHE_SIZE      = "CACHE_SIZE"
	ENV_CACHE_TIMEOUT   = "CACHE_TIMEOUT"
	DEFAULT_CONFIG_PATH = "/etc/webcache/"
	YAML_REGEX          = "^\\w*\\.(yaml|yml|YAML|YML)*$"
)

func setupBrowser(ctx context.Context) *wcache.Browser {
	configPath := DEFAULT_CONFIG_PATH
	if path, err := util.LookupEnv(ENV_CONFIG_PATH); err == nil {
		configPath = path
	}

	decoder := util.YamlEncoder.Unmarshaler()
	browser, err := wcache.NewBrowser(YAML_REGEX, decoder)
	if err != nil {
		log.Fatalf("browser setup has failed: %s", err)
	}

	if err := browser.ReadPath(configPath); err != nil {
		log.Fatalf("read path has failed: %s", err)
	}

	value, err := util.LookupEnv(ENV_WATCH_CONFIG)
	if err != nil {
		return browser
	}

	watch, err := strconv.ParseBool(value)
	if err != nil {
		return browser
	}

	if watch {
		if err := browser.WatchPath(ctx, configPath); err != nil {
			log.Fatalf("watch path has failed: %s", err)
		}
	}

	return browser
}

func setupCache() *wcache.RedisCache {
	addr, err := util.LookupEnv(ENV_REDIS_ADDR)
	if err != nil {
		log.Fatalf("%s: %s", ENV_REDIS_ADDR, err)
	}

	sizeStr, err := util.LookupEnv(ENV_CACHE_SIZE)
	if err != nil {
		log.Fatalf("%s: %s", ENV_CACHE_SIZE, err)
	}

	size, err := strconv.ParseInt(sizeStr, 0, 0)
	if err != nil {
		log.Fatalf("%s: %s", ENV_CACHE_SIZE, err)
	}

	timeoutStr, err := util.LookupEnv(ENV_CACHE_TIMEOUT)
	if err != nil {
		log.Fatalf("%s: %s", ENV_CACHE_TIMEOUT, err)
	}

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		log.Fatalf("%s: %s", ENV_CACHE_TIMEOUT, err)
	}

	cache, err := wcache.NewRedisCache(addr, int(size), timeout)
	if err != nil {
		log.Fatalf("cache setup has failed: %s", err)
	}

	return cache
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("no dotenv file has been found")
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	config := setupBrowser(ctx)
	cache := setupCache()
	proxy := wcache.NewReverseProxy(config, cache)
	proxy.TargetURI = func(req *http.Request) (string, error) {
		if target := req.URL.Query().Get("target"); len(target) != 0 {
			return target, nil
		}

		return "", wcache.ErrNoContent
	}

	proxy.DigestRequest = func(req *http.Request) (string, error) {
		digestBytes := wcache.DigestRequest(req, []string{"target"}, nil)
		digest := base64.RawStdEncoding.EncodeToString(digestBytes)
		return digest, nil
	}

	proxy.DecorateRequest = func(req *http.Request) {
		req.Host = strings.Split(req.Host, ":")[0]
	}

	network, err := util.LookupEnv(ENV_SERVICE_NETWORK)
	if err != nil {
		log.Fatalf("%s: %s", ENV_SERVICE_NETWORK, err)
	}

	address, err := util.LookupEnv(ENV_SERVICE_ADDR)
	if err != nil {
		log.Fatalf("%s: %s", ENV_SERVICE_ADDR, err)
	}

	lis, err := net.Listen(network, address)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", proxy.ServeHTTP)
	log.Printf("server listening on %s", address)
	if err := http.Serve(lis, nil); err != nil {
		log.Fatalf("server abruptly terminated: %s", err.Error())
	}
}
