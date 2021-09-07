package webcache

import (
	"log"
	"sync"

	"github.com/alvidir/go-util"
	"github.com/fsnotify/fsnotify"
)

type methodConfigFile struct {
	Name    string            `yaml:"name"`
	Headers map[string]string `yaml:"headers"`
	Enabled bool              `yaml:"enabled"`
}

type cacheConfigFile struct {
	Timeout  int      `yaml:"timeout"`
	Capacity int      `yaml:"capacity"`
	Methods  []string `yaml:"methods"`
	Enabled  bool     `yaml:"enabled"`
}

// ConfigFile represents a configuration file for the webcache service
type ConfigFile struct {
	Endpoints []string           `yaml:"endpoints"`
	Headers   map[string]string  `yaml:"headers"`
	Methods   []methodConfigFile `yaml:"methods"`
	Cache     cacheConfigFile    `yaml:"cache"`
}

type config struct {
	CachesByEndpoint sync.Map
	ConfigByEndpoint sync.Map
}

var configuration config

func applySettings(file *ConfigFile, config *config) {
}

func removeSettings(file *ConfigFile, config *config) {
}

func HandleConfigWatcher(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			var file ConfigFile
			if err := util.YamlEncoder.Unmarshaler().Path(event.Name, &file); err != nil {
				log.Printf("%s: %s", event.Name, err)
			}

			if event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Write == fsnotify.Write {
				log.Printf("configuration updates: %s", event.Name)
				applySettings(&file, &configuration)

			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				log.Printf("config file removed: %s", event.Name)
				removeSettings(&file, &configuration)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			log.Println("error:", err)
		}
	}
}
