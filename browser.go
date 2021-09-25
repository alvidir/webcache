package webcache

import (
	"log"
	fpath "path"
	"regexp"
	"sync"

	"github.com/alvidir/go-util"
	"github.com/fsnotify/fsnotify"
)

type CacheFile struct {
	Timeout int      `yaml:"timeout"`
	Methods []string `yaml:"methods"`
	Enabled bool     `yaml:"enabled"`
}

type MethodFile struct {
	Name    string `yaml:"name"`
	Enabled bool   `yaml:"enabled"`
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
	Cache   CacheFile   `yaml:"cache"`
	Request RequestFile `ỳaml:"request"`
	Router  RouterFile  `ỳaml:"router"`
}

// Browser represents a set of settings to apply over http requests and responses' cache
type Browser struct {
	regex   *regexp.Regexp
	decoder util.Unmarshaler
	files   sync.Map
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

// ReadFile reads the file located at the provided path and stores its content
func (browser *Browser) ReadFile(path string) (err error) {
	var f File
	if err = util.YamlEncoder.Unmarshaler().Path(path, &f); err != nil {
		return
	}

	base := fpath.Base(path)
	browser.files.Store(base, &f)
	return
}

func (browser *Browser) RemoveFile(path string) {
	base := fpath.Base(path)
	browser.files.Delete(base)
}

// WatchPath waits for any event on the provided path and stores any change on those files that matches the
// browser's regex
func (browser *Browser) WatchPath(path string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	defer watcher.Close()
	if err = watcher.Add(path); err != nil {
		return err
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			if browser.regex != nil && !browser.regex.MatchString(event.Name) {
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

		case err := <-watcher.Errors:
			return err
		}
	}
}
