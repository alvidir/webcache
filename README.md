# webcache
Web requests cache as a middleware

# About
Webcache is a middleware as a service that caches http responses in order to provide them fastly each time a same request is performed. As it is, it also works as a requests decorator able to modify http headers before the requests is redirected to the original target. 

# Configuration
By default, the server will watch for any file event into the  `/etc/webcache` directory. There, the application expects to find a set of configuration files **.yaml** with the following structure:

``` yaml
endpoints: # afected endpoints for the current configuration (regex)
  - https?://(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,32}/?$

headers: # global headers to add to any request
  XXX_A_CUSTOM_HEADER: header_value
  XXX_ANOTHER_HEADER: another_value

methods:
  - name: GET
    headers: # custom headers for GET requests
      XXX_JUST_IN_GET_HEADER: ghi 
  - name: POST
    enabled: false # block all POST requests (405 - method not allowed)
  - name: PUT
    enabled: false # block all PUT requests (405 - method not allowed)
  - name: default
    enabled: true # non listed method are disabled by default

cache:
  enabled: true
  timeout: 3600 # for how long a cached response is valid (s)
  capacity: 32 # how many keys (rows) can store the cache
  on_methods:
    - GET # catch only responses of GET requests

timeout: 3000 # how much long can it takes a request to get back a response (ms)
```

> Any addition, deletion or update inside the config directory will trigger an update over the server configuration. This behaviour can be turned off setting the environment variable `WATCH_CONFIG` to _False_