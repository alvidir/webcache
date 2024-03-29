package main

import (
	"encoding/base64"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"log"

	wcache "github.com/alvidir/webcache"
	"go.uber.org/zap"

	"github.com/joho/godotenv"
)

const (
	ENV_SERVICE_ADDR = "SERVICE_ADDR"
	ENV_SERVICE_NETW = "SERVICE_NETW"
	ENV_CONFIG_PATH  = "CONFIG_PATH"
	ENV_REDIS_DSN    = "REDIS_DSN"
	ENV_CACHE_TTL    = "CACHE_TTL"
	ENV_CACHE_SIZE   = "CACHE_SIZE"
	YAML_REGEX       = "^\\w*\\.(yaml|yml|YAML|YML)*$"
)

var (
	serviceAddr = "0.0.0.0:8000"
	serviceNetw = "tcp"
	configPath  = "/etc/webcache/"
	cacheTTL    = 10 * time.Minute
	cacheSize   = 1024
)

func setupConfiguration(logger *zap.Logger) *wcache.ConfigGroup {
	if path, exists := os.LookupEnv(ENV_CONFIG_PATH); exists {
		configPath = path
	}

	config, err := wcache.NewConfigGroup(configPath, logger)
	if err != nil {
		logger.Fatal("setting up configuration",
			zap.String("path", configPath),
			zap.Error(err))
	}

	return config
}

func setupCache(logger *zap.Logger) *wcache.RedisCache {
	addr, exists := os.LookupEnv(ENV_REDIS_DSN)
	if !exists {
		logger.Fatal("redis dsn must be set")
	}

	if value, exists := os.LookupEnv(ENV_CACHE_TTL); exists {
		if ttl, err := time.ParseDuration(value); err != nil {
			logger.Fatal("invalid cache ttl",
				zap.String("value", value),
				zap.Error(err))
		} else {
			cacheTTL = ttl
		}
	}

	if value, exists := os.LookupEnv(ENV_CACHE_SIZE); exists {
		if size, err := strconv.Atoi(value); err != nil {
			logger.Fatal("invalid cache size",
				zap.String("value", value),
				zap.Error(err))
		} else {
			cacheSize = size
		}
	}

	cache, err := wcache.NewRedisCache(addr, cacheSize, cacheTTL)
	if err != nil {
		log.Fatalf("cache setup has failed: %s", err)
	}

	return cache
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if err := godotenv.Load(); err != nil {
		logger.Warn("no dotenv file has been found",
			zap.Error(err))
	}

	cache := setupCache(logger)
	config := setupConfiguration(logger)
	proxy := wcache.NewReverseProxy(config, cache, logger)

	proxy.DigestRequest = func(req *http.Request) (string, error) {
		digestBytes := wcache.DigestRequest(req, []string{wcache.HTTP_LOCATION_HEADER})
		digest := base64.RawStdEncoding.EncodeToString(digestBytes)
		return digest, nil
	}

	if addr, exists := os.LookupEnv(ENV_SERVICE_ADDR); exists {
		serviceAddr = addr
	}

	if netw, exists := os.LookupEnv(ENV_SERVICE_NETW); exists {
		serviceNetw = netw
	}

	lis, err := net.Listen(serviceNetw, serviceAddr)
	if err != nil {
		logger.Panic("failed to listen: %v",
			zap.Error(err))
	}

	logger.Info("server ready to accept connections",
		zap.String("address", serviceAddr))

	http.HandleFunc("/", proxy.ServeHTTP)
	if err := http.Serve(lis, nil); err != nil {
		log.Fatalf("server abruptly terminated: %s", err.Error())
	}
}
