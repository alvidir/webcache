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

type Options struct {
	enabled bool
	cached  bool
	timeout time.Duration
	headers map[string]string
}

func NewOptions() *Options {
	return &Options{
		timeout: DefaultTimeout,
		headers: make(map[string]string),
	}
}

func (ops *Options) copy() *Options {
	copy := &Options{
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
	options map[string]*Options
}

// ConfigGroup represents a set of settings to apply over http requests and responses
type ConfigGroup struct {
	configs []config
	base    map[string]*Options
	logger  *zap.Logger
	mu      sync.RWMutex
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

func overwrite(m *Method, base *Options) (*Options, error) {
	var ops *Options
	if base == nil {
		ops = NewOptions()
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

	if m.Timeout != nil {
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

func mapOptions(m []Method, base map[string]*Options) (mops map[string]*Options, err error) {
	mops = make(map[string]*Options)

	if base == nil {
		mops[DEFAULT_KEYWORD] = NewOptions()
	} else if globalDef, exists := base[DEFAULT_KEYWORD]; !exists {
		mops[DEFAULT_KEYWORD] = NewOptions()
	} else {
		mops[DEFAULT_KEYWORD] = globalDef.copy()
	}

	var localDef *Method
	for _, method := range m {
		if method.Name != DEFAULT_KEYWORD {
			continue
		}

		localDef = &method
		ops, err := overwrite(&method, mops[DEFAULT_KEYWORD])
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

		var mbase *Options
		if base == nil {
			mbase = mops[DEFAULT_KEYWORD]
		} else if globalDef, exists := base[method.Name]; !exists {
			mbase = mops[DEFAULT_KEYWORD]
		} else if mbase, err = overwrite(localDef, globalDef); err != nil {
			return nil, err
		}

		if ops, err := overwrite(&method, mbase); err != nil {
			return nil, err
		} else {
			mops[method.Name] = ops
		}
	}

	return mops, nil
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
func (group *ConfigGroup) AddConfig(file *ConfigFile) (err error) {
	group.base, err = mapOptions(file.Methods, nil)
	if err != nil {
		return err
	}

	group.configs = make([]config, len(file.Router))
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

		config.options, err = mapOptions(router.Methods, group.base)
		if err != nil {
			return err
		}

		group.configs[ri] = config
	}

	return nil
}

// RequestOptions returns the Options instance for the given endpoint and method
func (group *ConfigGroup) RequestOptions(endpoint string, method string) (*Options, bool) {
	for _, config := range group.configs {
		for _, regex := range config.regex {
			if !regex.MatchString(endpoint) {
				continue
			}

			if ops, exists := config.options[method]; exists {
				return ops, true
			}
		}
	}

	return group.base[method], false
}
