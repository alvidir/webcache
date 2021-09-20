package webcache

import (
	"log"
	"regexp"

	"github.com/alvidir/go-util"
	"github.com/fsnotify/fsnotify"
)

type cacheFile struct {
	Timeout int      `yaml:"timeout"`
	Methods []string `yaml:"methods"`
	Enabled bool     `yaml:"enabled"`
}

type methodFile struct {
	Name    string `yaml:"name"`
	Enabled bool   `yaml:"enabled"`
}

type requestFile struct {
	Timeout int          `yaml:"timeout"`
	Methods []methodFile `yaml:"methods"`
}

type routerFile struct {
	Endpoints []string          `yaml:"endpoints"`
	Headers   map[string]string `yaml:"headers"`
	Methods   []methodFile      `yaml:"methods"`
	Cached    bool              `yaml:"cached"`
}

// file represents a configuration file for the webcache service
type file struct {
	Cache   cacheFile   `yaml:"cache"`
	Request requestFile `ỳaml:"request"`
	Router  routerFile  `ỳaml:"router"`
}

// Browser represents a set of settings to apply over http requests and responses' cache
type Browser struct {
	regex   *regexp.Regexp
	decoder util.Unmarshaler
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

// AttachWatcher takes a fs watchers and waits for any update, create or delete event from it
func (browser *Browser) AttachPath(path string) error {
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

				var f file
				if err := util.YamlEncoder.Unmarshaler().Path(event.Name, &f); err != nil {
					log.Printf("%s: %s", event.Name, err)
					continue
				}

				//browser.applySettings(event.Name, &f)

			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				//config.removeSettings(event.Name, file)
			}

		case err := <-watcher.Errors:
			return err
		}
	}
}
