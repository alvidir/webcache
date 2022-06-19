package webcache

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"sync"
	"time"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const (
	DEFAULT_KEYWORD = "default"
	YAML_REGEX      = "^\\w*\\.(yaml|yml|YAML|YML)*$"
)

var (
	DefaultTimeout = 10 * time.Minute
)

type Method struct {
	Name    string `yaml:"name"`
	Enabled *bool  `yaml:"enabled"`
	Cached  *bool  `yaml:"cached"`
	Timeout *string
	Headers map[string]string `yaml:"headers"`
}

type Router struct {
	Endpoints []string `yaml:"endpoints"`
	Methods   []Method `yaml:"methods"`
}

type ConfigFile struct {
	Methods []Method `yaml:"methods"`
	Router  []Router `yaml:"router"`
}

type options struct {
	enabled bool
	cached  bool
	timeout time.Duration
	headers map[string]string
}

func newOptions() *options {
	return &options{
		headers: make(map[string]string),
	}
}

func (ops *options) copy() *options {
	copy := &options{
		enabled: ops.enabled,
		cached:  ops.cached,
		timeout: ops.timeout,
		headers: make(map[string]string),
	}

	for key, value := range ops.headers {
		copy.headers[key] = value
	}

	return copy
}

type config struct {
	regex   []*regexp.Regexp
	options map[string]*options
}

// ConfigGroup represents a set of settings to apply over http requests and responses
type ConfigGroup struct {
	config []config
	logger *zap.Logger
	mu     sync.RWMutex
}

// NewConfigGroup reads the content of the provided path and returns the declared configuration
func NewConfigGroup(path string, logger *zap.Logger) (*ConfigGroup, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	group := &ConfigGroup{
		logger: logger,
	}

	if stat.IsDir() {
		err = group.ReadDir(path)
	} else {
		err = group.ReadFile(path)
	}

	return group, err
}

// ReadDir applies all configuration files inside the given directory into the current configuration
func (group *ConfigGroup) ReadDir(dir string) error {
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
		if err = group.ReadFile(filepath); err != nil {
			return err
		}
	}

	return nil
}

// ReadFile applies a configuration file into the current configuration
func (group *ConfigGroup) ReadFile(filepath string) error {
	file, err := group.read(filepath)
	if err != nil {
		return err
	}

	return group.AddConfig(file)
}

// AddConfig registers the given configuration into the config group
func (group *ConfigGroup) AddConfig(file *ConfigFile) error {
	globalOps, err := group.mapOptions(file.Methods, nil)
	if err != nil {
		return err
	}

	group.config = make([]config, len(file.Router))
	for ri, router := range file.Router {
		config := config{
			regex: make([]*regexp.Regexp, len(router.Endpoints)),
		}

		for ei, rx := range router.Endpoints {
			comp, err := regexp.Compile(rx)
			if err != nil {
				return err
			}

			config.regex[ei] = comp
		}

		config.options, err = group.mapOptions(router.Methods, globalOps)
		if err != nil {
			return err
		}

		group.config[ri] = config
	}

	return nil
}

func (group *ConfigGroup) read(filepath string) (*ConfigFile, error) {
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	var file ConfigFile
	if err = yaml.NewDecoder(f).Decode(&file); err != nil {
		return nil, err
	}

	return &file, nil
}

func (group *ConfigGroup) buildOptions(m *Method, base *options) (*options, error) {
	var ops *options
	if base == nil {
		ops = newOptions()
	} else {
		ops = base.copy()
	}

	if m == nil {
		return ops, nil
	}

	if m.Enabled != nil {
		ops.enabled = *m.Enabled
	}

	if m.Cached != nil {
		ops.cached = *m.Cached
	}

	if m.Timeout == nil {
		ops.timeout = DefaultTimeout
	} else {
		timeout, err := time.ParseDuration(*m.Timeout)
		if err != nil {
			return nil, err
		}

		ops.timeout = timeout
	}

	for key, value := range m.Headers {
		ops.headers[key] = value
	}

	return ops, nil
}

func (group *ConfigGroup) mapOptions(m []Method, base map[string]*options) (mops map[string]*options, err error) {
	mops = make(map[string]*options)

	if base == nil {
		mops[DEFAULT_KEYWORD] = newOptions()
	} else if globalDef, exists := base[DEFAULT_KEYWORD]; !exists {
		mops[DEFAULT_KEYWORD] = newOptions()
	} else {
		mops[DEFAULT_KEYWORD] = globalDef.copy()
	}

	var localDef *Method
	for _, method := range m {
		if method.Name != DEFAULT_KEYWORD {
			continue
		}

		localDef = &method
		ops, err := group.buildOptions(&method, mops[DEFAULT_KEYWORD])
		if err != nil {
			return nil, err
		}

		mops[DEFAULT_KEYWORD] = ops
		break
	}

	for _, method := range m {
		if method.Name == DEFAULT_KEYWORD {
			continue
		}

		var mbase *options
		if base == nil {
			mbase = mops[DEFAULT_KEYWORD]
		} else if globalDef, exists := base[method.Name]; !exists {
			mbase = mops[DEFAULT_KEYWORD]
		} else if mbase, err = group.buildOptions(localDef, globalDef); err != nil {
			return nil, err
		}

		if ops, err := group.buildOptions(&method, mbase); err != nil {
			return nil, err
		} else {
			mops[method.Name] = ops
		}
	}

	return mops, nil
}
