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

type CacheFile struct {
	Timeout int      `yaml:"timeout"`
	Methods []string `yaml:"methods"`
	Enabled bool     `yaml:"enabled"`
}

type MethodFile struct {
	Name    string            `yaml:"name"`
	Enabled bool              `yaml:"enabled"`
	Cached  bool              `yaml:"cached"`
	Headers map[string]string `yaml:"headers"`
}

type RequestFile struct {
	Timeout int          `yaml:"timeout"`
	Methods []MethodFile `yaml:"methods"`
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

func (browser *Browser) watch(ctx context.Context, w *fsnotify.Watcher) {
	defer w.Close()

	for ctx.Err() == nil {
		select {
		case event, ok := <-w.Events:
			if !ok {
				log.Printf("watcher's event arrived as: ko")
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
				log.Printf("watcher has got errors: %s", err.Error())
				return
			}

		case <-ctx.Done():
			return
		}
	}
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
	log.Printf("%s: successfully applied", p)
	return
}

// RemoveFile removes the configuration from the given file path p
func (browser *Browser) RemoveFile(p string) {
	filename := path.Base(p)
	browser.files.Delete(filename)
	log.Printf("%s: successfully removed", p)
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
	log.Printf("watcher is ready to accept events from %s", p)

	if err = watcher.Add(p); err != nil {
		return err
	}

	return nil
}

// IsEndpointAllowed returns true if, and only if, the given enpoint is allowed
func (browser *Browser) IsEndpointAllowed(string) bool {
	return false
}

// IsMethodAllowed returns true if, and only if, the given method is allowed for the given endpoint
func (browser *Browser) IsMethodAllowed(string, string) bool {
	return false
}

// IsMethodCached returns true if, and only if, is allowed to cach responses for the given endpoint and method
func (browser *Browser) IsMethodCached(string, string) bool {
	return false
}

// ResponseTimeout returns the time a response is considered valid
func (browser *Browser) ResponseTimeout(string) time.Duration {
	return 0
}

// Headers returns all these headers to add to any request with the given endpoint and method
func (browser *Browser) Headers(string, string) map[string]string {
	return nil
}
