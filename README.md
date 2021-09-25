# webcache
Reverse proxy as an HTTP cache and requests decorator

# About
Webcache is a reverse proxy that caches http responses in order to provide them rapidly each time a same request is performed. Through its configuration file, webcache allows blocking hosts and methods (responding with _403: Forbidden_ or _405: Method not allowed_ respectively) and lets to specify new headers for all or some specific requests.

By default, Webcache uses a Redis server as a cache, and has its own configuration structure. However, all of this can be easily customized by just implementing the corresponding interfaces **Config** and **Cache**. 

# Configuration
By default, the server will expect to find any **.yaml** file in the  `/etc/webcache/` path, or any other defined by `CONFIG_PATH` environment variable. If no config file is found, or none of them follows the structure from the example down below, no request or method will be allowed by the webcache. 

``` yaml
cache:
  enabled: true
  timeout: 3600 # for how long a cached response is valid (s)
  methods:
    - GET # catch only responses of GET requests

request:
  methods: # methods configuration for any endpoint listed in this file
    - name: default
      enabled: true # since non listed method are disabled by default, enable all them
    - name: DELETE
      cached: false # do not catch any DELETE response

  headers: # global headers for any endpoint listed in this file
    - XXX_GLOBAL_HEADER: global_header

router:
  - endpoints: # afected endpoints for the current configuration (regex)
      - https?://(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,32}/?$

    headers: # custom headers for the endpoints above
      XXX_A_CUSTOM_HEADER: header_value
      XXX_ANOTHER_HEADER: another_value

    methods: # methods configuration for the endpoints above
      - name: GET
        headers: # custom headers for GET requests
          XXX_JUST_IN_GET_HEADER: get_header_value
      - name: POST
        enabled: false # block all POST requests (405 - method not allowed)
      - name: PUT
        enabled: false # block all PUT requests (405 - method not allowed)
    
    cached: true # enable caching responses from the endpoints above
```

> The webcache's configuration is static by default, meaning that once a config file is applied, any update on it will take no effect over webcache's configuration. To enable event's watching for any config file or directory, the environment variable `WATCH_CONFIG` must be set as _True_. 
