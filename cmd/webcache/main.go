package main

import (
	"context"
	"encoding/base64"
	"net"
	"net/http"
	"strconv"
	"time"

	"log"

	util "github.com/alvidir/go-util"
	wcache "github.com/alvidir/webcache"

	"github.com/joho/godotenv"
)

const (
	EnvServiceAddr    = "SERVICE_ADDR"
	EnvServiceNetw    = "SERVICE_NETW"
	EnvConfigPath     = "CONFIG_PATH"
	EnvWatchConfig    = "WATCH_CONFIG"
	EnvRedisAddr      = "REDIS_ADDR"
	EnvCacheSize      = "CACHE_SIZE"
	EnvCacheTimeout   = "CACHE_TIMEOUT"
	DefaultConfigPath = "/etc/webcache/"
	YamlRegex         = "^\\w*\\.(yaml|yml|YAML|YML)*$"
)

func setupBrowser(ctx context.Context) *wcache.Browser {
	configPath := DefaultConfigPath
	if path, err := util.LookupEnv(EnvConfigPath); err == nil {
		configPath = path
	}

	decoder := util.YamlEncoder.Unmarshaler()
	browser, err := wcache.NewBrowser(YamlRegex, decoder)
	if err != nil {
		log.Fatalf("browser setup has failed: %s", err)
	}

	if err := browser.ReadPath(configPath); err != nil {
		log.Fatalf("read path has failed: %s", err)
	}

	value, err := util.LookupEnv(EnvWatchConfig)
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
	addr, err := util.LookupEnv(EnvRedisAddr)
	if err != nil {
		log.Fatalf("%s: %s", EnvRedisAddr, err)
	}

	sizeStr, err := util.LookupEnv(EnvCacheSize)
	if err != nil {
		log.Fatalf("%s: %s", EnvCacheSize, err)
	}

	size, err := strconv.ParseInt(sizeStr, 0, 0)
	if err != nil {
		log.Fatalf("%s: %s", EnvCacheSize, err)
	}

	timeoutStr, err := util.LookupEnv(EnvCacheTimeout)
	if err != nil {
		log.Fatalf("%s: %s", EnvCacheTimeout, err)
	}

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		log.Fatalf("%s: %s", EnvCacheTimeout, err)
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

	proxy.DigestRequest = func(req *http.Request) (string, error) {
		digestBytes := wcache.DigestRequest(req, []string{wcache.HTTP_LOCATION_HEADER})
		digest := base64.RawStdEncoding.EncodeToString(digestBytes)
		return digest, nil
	}

	network, err := util.LookupEnv(EnvServiceNetw)
	if err != nil {
		log.Fatalf("%s: %s", EnvServiceNetw, err)
	}

	address, err := util.LookupEnv(EnvServiceAddr)
	if err != nil {
		log.Fatalf("%s: %s", EnvServiceAddr, err)
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
