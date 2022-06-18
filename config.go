package webcache

import (
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"sync"
	"time"

	"github.com/alvidir/go-util"
)

const (
	DEFAULT_KEYWORD = "default"
	YAML_REGEX      = "^\\w*\\.(yaml|yml|YAML|YML)*$"
)

type CacheFile struct {
	Timeout string   `yaml:"timeout"`
	Methods []string `yaml:"methods"`
	Enabled bool     `yaml:"enabled"`
}

type MethodFile struct {
	Name    string            `yaml:"name"`
	Enabled *bool             `yaml:"enabled"`
	Cached  bool              `yaml:"cached"`
	Headers map[string]string `yaml:"headers"`
}

type RequestFile struct {
	Methods []MethodFile      `yaml:"methods"`
	Headers map[string]string `yaml:"headers"`
}

type RouterFile struct {
	Endpoints []string          `yaml:"endpoints"`
	Headers   map[string]string `yaml:"headers"`
	Methods   []MethodFile      `yaml:"methods"`
	Cached    bool              `yaml:"cached"`
}

// ConfigFile represents a configuration file for the webcache service
type ConfigFile struct {
	Cache   CacheFile    `yaml:"cache"`
	Request RequestFile  `yaml:"request"`
	Router  []RouterFile `yaml:"router"`
}

// Config represents a set of settings to apply over http requests and responses' cache
type Config struct {
	files sync.Map
}

type endpointConfig struct {
	router  RouterFile
	request RequestFile
	cache   CacheFile
}

func (browser *Config) getEndpointConfig(endpoint string) (config *endpointConfig, exists bool) {
	config = new(endpointConfig)

	browser.files.Range(func(key, value interface{}) bool {
		file, ok := value.(*ConfigFile)
		if !ok || file == nil {
			log.Printf("TYPE_ASSERT %s - want *File", endpoint)
			browser.files.Delete(key)
			return true
		}

		config.cache = file.Cache
		config.request = file.Request

		for _, route := range file.Router {
			for _, regex := range route.Endpoints {
				comp, err := regexp.Compile(regex)
				if err != nil {
					log.Printf("REGEX_COMP %s - %s", regex, err)
					continue
				}

				if exists = comp.MatchString(endpoint); exists {
					config.router = route
					return false
				}
			}
		}

		return true
	})

	return
}

// ApplyFile reads the file located at the provided path p and stores its content
func (browser *Config) ApplyFile(p string) (err error) {
	var f ConfigFile
	if err = util.YamlEncoder.Unmarshaler().Path(p, &f); err != nil {
		return
	}

	filename := path.Base(p)
	browser.files.Store(filename, &f)
	log.Printf("READ_FILE %s successfully applied", p)
	return
}

// UnapplyFile removes the configuration from the given file path p
func (browser *Config) UnapplyFile(p string) {
	filename := path.Base(p)
	browser.files.Delete(filename)
	log.Printf("REMOVE_FILE %s successfully applied", p)
}

// ReadPath reads all files located at the provided path p that matches the browser's regex
func (browser *Config) ReadPath(p string) error {
	files, err := ioutil.ReadDir(p)
	if err != nil {
		return browser.ApplyFile(p)
	}

	for _, f := range files {
		if browser.regex != nil && !browser.regex.MatchString(f.Name()) {
			continue
		}

		fullpath := path.Join(p, f.Name())
		if err := browser.ApplyFile(fullpath); err != nil {
			return err
		}
	}

	return nil
}

// IsEndpointAllowed returns true if, and only if, the given enpoint is allowed
func (browser *Config) IsEndpointAllowed(endpoint string) (ok bool) {
	_, ok = browser.getEndpointConfig(endpoint)
	return
}

// IsMethodAllowed returns true if, and only if, the given method is allowed for the given endpoint
func (browser *Config) IsMethodAllowed(endpoint string, method string) bool {
	config, ok := browser.getEndpointConfig(endpoint)
	if !ok {
		return false
	}

	if config == nil {
		return false
	}

	for _, rmethod := range config.router.Methods {
		if rmethod.Name == DEFAULT_KEYWORD || rmethod.Name == method {
			return rmethod.Enabled == nil || *rmethod.Enabled
		}
	}

	for _, rmethod := range config.request.Methods {
		if rmethod.Name == DEFAULT_KEYWORD || rmethod.Name == method {
			return rmethod.Enabled == nil || *rmethod.Enabled
		}
	}

	return false
}

// IsMethodCached returns true if, and only if, is allowed to cach responses for the given endpoint and method
func (browser *Config) IsMethodCached(endpoint string, method string) bool {
	config, ok := browser.getEndpointConfig(endpoint)
	if !ok {
		return false
	}

	return config.cache.Enabled && config.router.Cached
}

// ResponseLifetime returns the duration a response is considered valid
func (browser *Config) ResponseLifetime(endpoint string) time.Duration {
	config, ok := browser.getEndpointConfig(endpoint)
	if !ok {
		return 0
	}

	lifetime, err := time.ParseDuration(config.cache.Timeout)
	if err != nil {
		log.Printf("PARSE_TIME %s - %s", endpoint, err.Error())
		return 0
	}

	return lifetime
}

// Headers returns all those headers to add to any request with the given endpoint and method
func (browser *Config) Headers(endpoint string, method string) map[string]string {
	headers := make(map[string]string)
	config, ok := browser.getEndpointConfig(endpoint)
	if !ok {
		return headers
	}

	for key, value := range config.request.Headers {
		headers[key] = value
	}

	for _, rmethod := range config.request.Methods {
		if rmethod.Name == DEFAULT_KEYWORD || rmethod.Name == method {
			for key, value := range rmethod.Headers {
				headers[key] = value
			}
		}
	}

	for key, value := range config.router.Headers {
		headers[key] = value
	}

	for _, rmethod := range config.router.Methods {
		if rmethod.Name == DEFAULT_KEYWORD || rmethod.Name == method {
			for key, value := range rmethod.Headers {
				headers[key] = value
			}
		}
	}

	return headers
}
