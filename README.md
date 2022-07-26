# webcache

[![Go version](https://img.shields.io/badge/Go-go1.18-blue.svg)](https://go.dev/) [![tests](https://github.com/alvidir/webcache/actions/workflows/test.yaml/badge.svg?branch=main)](https://github.com/alvidir/webcache/actions/workflows/test.yaml)
[![webcache](https://img.shields.io/github/v/release/alvidir/webcache.svg)](https://github.com/alvidir/webcache)

Reverse proxy as an HTTP cache and requests decorator

# About
Webcache is a reverse proxy that caches http responses in order to provide them rapidly each time a same request is performed. Through its configuration file, webcache allows blocking hosts and methods (responding with _403: Forbidden_ or _405: Method not allowed_ respectively) and lets to specify new headers for all or some specific requests.

By default, Webcache uses a Redis server as a cache, and has its own configuration structure. However, all of this can be easily customized by just implementing the corresponding interfaces **Config** and **Cache**. 

# Configuration
By default, the server will expect to find any **.yaml** file in the  `/etc/webcache/` path, or any other defined by `CONFIG_PATH` environment variable. If no config file is found, or none of them follows the structure from the example down below, no request or method will be allowed by the webcache. 

``` yaml
methods:
  - name: default   # default methods configuration for all endpoints listed in this file
    enabled: true   # since non listed methods are disabled by default, enable all them
    cached: true    # since cache is disabled by default, enable for all methods
    timeout: 1h     # for how long a cached response is valid (10 minutes by default)
    headers:        # custom headers for all the endpoints listed in this file
      X_GLOBAL_HEADER: global_header
      
  - name: DELETE
    enabled: false  # block all DELETE requests (405 - method not allowed)
  - name: POST
    enabled: false  # block all POST requests (405 - method not allowed)
  - name: PUT
    cached: false   # do not catch any PUT requests

router:
  - endpoints: # afected endpoints for the current configuration (regex)
      - https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,32}\/?$

    methods: # overwrite the default methods configuration for the endpoints above
      - name: default
        headers: # custom headers for the endpoints above
          X_A_CUSTOM_HEADER: header_value
          X_ANOTHER_HEADER: another_value
      - name: GET
        headers:
          X_JUST_IN_GET_HEADER: get_header_value
```

# Environment variables

The application requires a set of environment variables that describes how the service must perform. These environment variables are those listed down below:

``` bash
# Service endpoint and network
SERVICE_ADDR=:8000
SERVICE_NETWORK=tcp

# Configuration path 
CONFIG_PATH=.config/

# Redis datasource
REDIS_DSN=webcache-redis:6379

# In-memory cache configuration 
CACHE_SIZE=1024 # how many entries the local cache can have
CACHE_TTL=10m # for how long an entry is stored in the local cache
```
