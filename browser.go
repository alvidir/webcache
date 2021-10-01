package webcache

import (
	"context"
	"io/ioutil"
	"log"
	"path"
	"regexp"
	"sync"
	"time"

	"github.com/alvidir/go-util"
	"github.com/fsnotify/fsnotify"
)

const (
	DEFAULT_KEYWORD = "default"
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

// File represents a configuration file for the webcache service
type File struct {
	Cache   CacheFile    `yaml:"cache"`
	Request RequestFile  `ỳaml:"request"`
	Router  []RouterFile `ỳaml:"router"`
}

// Browser represents a set of settings to apply over http requests and responses' cache
type Browser struct {
	regex   *regexp.Regexp
	decoder util.Unmarshaler
	files   sync.Map
}

type eConfig struct {
	router  RouterFile
	request RequestFile
	cache   CacheFile
}

func (browser *Browser) watch(ctx context.Context, w *fsnotify.Watcher) {
	defer w.Close()

	for ctx.Err() == nil {
		select {
		case event, ok := <-w.Events:
			if !ok {
				log.Printf("WATCHER event arrived as ko")
				return
			}

			filename := path.Base(event.Name)
			if browser.regex != nil && !browser.regex.MatchString(filename) {
				continue
			}

			if event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Write == fsnotify.Write {

				if err := browser.ReadFile(event.Name); err != nil {
					log.Printf("%s: %s", event.Name, err.Error())
				}

			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				browser.RemoveFile(event.Name)
			}

		case err := <-w.Errors:
			if err != nil {
				log.Printf("WATCHER %s", err.Error())
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

func (browser *Browser) findEndpointConfig(endpoint string) (config *eConfig, exists bool) {
	config = new(eConfig)

	browser.files.Range(func(key, value interface{}) bool {
		file, ok := value.(*File)
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

func NewBrowser(regex string, decoder util.Unmarshaler) (*Browser, error) {
	comp, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	browser := &Browser{
		regex:   comp,
		decoder: decoder,
	}

	return browser, err
}

// ReadFile reads the file located at the provided path p and stores its content
func (browser *Browser) ReadFile(p string) (err error) {
	var f File
	if err = util.YamlEncoder.Unmarshaler().Path(p, &f); err != nil {
		return
	}

	filename := path.Base(p)
	browser.files.Store(filename, &f)
	log.Printf("READ_FILE %s successfully applied", p)
	return
}

// RemoveFile removes the configuration from the given file path p
func (browser *Browser) RemoveFile(p string) {
	filename := path.Base(p)
	browser.files.Delete(filename)
	log.Printf("REMOVE_FILE %s successfully applied", p)
}

// ReadPath reads all files located at the provided path p that matches the browser's regex
func (browser *Browser) ReadPath(p string) error {
	files, err := ioutil.ReadDir(p)
	if err != nil {
		return browser.ReadFile(p)
	}

	for _, f := range files {
		if browser.regex != nil && !browser.regex.MatchString(f.Name()) {
			continue
		}

		fullpath := path.Join(p, f.Name())
		if err := browser.ReadFile(fullpath); err != nil {
			return err
		}
	}

	return nil
}

// WatchPath waits for any event on the provided path p and stores any change on those files that matches the
// browser's regex
func (browser *Browser) WatchPath(ctx context.Context, p string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go browser.watch(ctx, watcher)
	log.Printf("WATCHER %s ready to accept events", p)

	if err = watcher.Add(p); err != nil {
		return err
	}

	return nil
}

// IsEndpointAllowed returns true if, and only if, the given enpoint is allowed
func (browser *Browser) IsEndpointAllowed(endpoint string) (ok bool) {
	_, ok = browser.findEndpointConfig(endpoint)
	return
}

// IsMethodAllowed returns true if, and only if, the given method is allowed for the given endpoint
func (browser *Browser) IsMethodAllowed(endpoint string, method string) bool {
	config, ok := browser.findEndpointConfig(endpoint)
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
func (browser *Browser) IsMethodCached(endpoint string, method string) bool {
	config, ok := browser.findEndpointConfig(endpoint)
	if !ok {
		return false
	}

	return config.cache.Enabled && config.router.Cached
}

// ResponseLifetime returns the duration a response is considered valid
func (browser *Browser) ResponseLifetime(endpoint string) time.Duration {
	config, ok := browser.findEndpointConfig(endpoint)
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

// Headers returns all these headers to add to any request with the given endpoint and method
func (browser *Browser) Headers(endpoint string, method string) map[string]string {
	headers := make(map[string]string)
	config, ok := browser.findEndpointConfig(endpoint)
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
