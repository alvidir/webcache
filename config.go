package webcache

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	DEFAULT_KEYWORD = "default"
	YAML_REGEX      = "^\\w*\\.(yaml|yml|YAML|YML)*$"
)

type CacheConfig struct {
	Timeout string   `yaml:"timeout"`
	Methods []string `yaml:"methods"`
	Enabled bool     `yaml:"enabled"`
}

type MethodConfig struct {
	Name    string            `yaml:"name"`
	Enabled *bool             `yaml:"enabled"`
	Cached  bool              `yaml:"cached"`
	Headers map[string]string `yaml:"headers"`
}

type RequestConfig struct {
	Methods []MethodConfig    `yaml:"methods"`
	Headers map[string]string `yaml:"headers"`
}

type RouterConfig struct {
	Endpoints []string          `yaml:"endpoints"`
	Headers   map[string]string `yaml:"headers"`
	Methods   []MethodConfig    `yaml:"methods"`
	Cached    bool              `yaml:"cached"`
}

// ConfigFile represents a configuration file for the webcache service
type ConfigFile struct {
	Cache   CacheConfig    `yaml:"cache"`
	Request RequestConfig  `yaml:"request"`
	Router  []RouterConfig `yaml:"router"`
}

// Config represents a set of settings to apply over http requests and responses' cache
type Config struct {
	files sync.Map
}

type endpointConfig struct {
	router  RouterConfig
	request RequestConfig
	cache   CacheConfig
}

func NewConfig(path string) (*Config, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if stat.IsDir() {
		err = config.ReadDir(path)
	} else {
		err = config.ReadFile(path)
	}

	return &config, err
}

// ReadDir applies all configuration files inside the given directory into the current configuration
func (config *Config) ReadDir(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	yamlRegex, err := regexp.Compile(YAML_REGEX)
	if err != nil {
		return err
	}

	for _, f := range files {
		if !yamlRegex.MatchString(f.Name()) {
			continue
		}

		filepath := path.Join(dir, f.Name())
		if err = config.ReadFile(filepath); err != nil {
			return err
		}
	}

	return nil
}

// ReadFile applies a configuration file into the current configuration
func (config *Config) ReadFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	var cfile ConfigFile
	if err = yaml.NewDecoder(file).Decode(&cfile); err != nil {
		return err
	}

	filename := path.Base(filepath)
	config.files.Store(filename, &file)
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
