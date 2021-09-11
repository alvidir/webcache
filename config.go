package webcache

import (
	"io/fs"
	"log"
	"path"
	"regexp"
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

type Config struct {
	// cachesByEndpoint sync.Map
	// configByEndpoint sync.Map
	configByFilename sync.Map
}

var fregex, _ = regexp.Compile(`^.*\.(yaml|yml)`)

func (config *Config) applySettings(name string, file *ConfigFile) {
	log.Printf("%s: its being processed", name)

}

func (config *Config) removeSettings(name string, file *ConfigFile) {
	log.Printf("%s: its being removed", name)
}

// ApplyConfigFiles takes a set of files and applies these ones that matches with the configuration structure
func (config *Config) ApplyConfigFiles(files []fs.FileInfo, root string) error {
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !fregex.MatchString(file.Name()) {
			continue
		}

		fullpath := path.Join(root, file.Name())

		var configfile ConfigFile
		if err := util.YamlEncoder.Unmarshaler().Path(fullpath, &configfile); err != nil {
			log.Printf("%s: %s", fullpath, err)
			continue
		}

		config.applySettings(fullpath, &configfile)
	}

	return nil
}

// HandleConfigWatcher
func (config *Config) HandleConfigWatcher(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if !fregex.MatchString(event.Name) {
				continue
			}

			if event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Write == fsnotify.Write {

				var file ConfigFile
				if err := util.YamlEncoder.Unmarshaler().Path(event.Name, &file); err != nil {
					log.Printf("%s: %s", event.Name, err)
					continue
				}

				config.applySettings(event.Name, &file)

			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				v, ok := config.configByFilename.Load(event.Name)
				if !ok {
					log.Printf("%s: was not processed", event.Name)
					continue
				}

				if file, ok := v.(*ConfigFile); ok {
					config.removeSettings(event.Name, file)
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			log.Println(err)
		}
	}
}
