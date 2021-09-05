package webcache

type methodConfig struct {
	name    string            `yaml:"name"`
	headers map[string]string `yaml:"headers"`
	enabled bool              `yaml:"enabled"`
}

type cacheConfig struct {
	timeout  int      `yaml:"timeout"`
	capacity int      `yaml:"capacity"`
	methods  []string `yaml:"methods"`
	enabled  bool     `yaml:"enabled"`
}

// Config represents a configuration file for the webcache service
type Config struct {
	app       string            `yaml:"app"`
	endpoints []string          `yaml:"endpoints"`
	headers   map[string]string `yaml:"headers"`
	methods   []methodConfig    `yaml:"methods"`
	cache     cacheConfig       `yaml:"cache"`
}
