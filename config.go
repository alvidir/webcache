package webcache

import (
	"sync"

	"github.com/fsnotify/fsnotify"
)

type methodConfigFile struct {
	name    string            `yaml:"name"`
	headers map[string]string `yaml:"headers"`
	enabled bool              `yaml:"enabled"`
}

type cacheConfigFile struct {
	timeout  int      `yaml:"timeout"`
	capacity int      `yaml:"capacity"`
	methods  []string `yaml:"methods"`
	enabled  bool     `yaml:"enabled"`
}

// ConfigFile represents a configuration file for the webcache service
type ConfigFile struct {
	endpoints []string           `yaml:"endpoints"`
	headers   map[string]string  `yaml:"headers"`
	methods   []methodConfigFile `yaml:"methods"`
	cache     cacheConfigFile    `yaml:"cache"`
}

type config struct {
	CachesByEndpoint sync.Map
	ConfigByEndpoint sync.Map
}

func HandleConfigWatcher(watcher *fsnotify.Watcher) {
	log := Log.WithField("origin", "watcher")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			log.Println("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("modified file:", event.Name)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			log.Println("error:", err)
		}
	}
}
