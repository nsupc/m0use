package config

import (
	"errors"
	"os"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
)

type Telegram struct {
	Template string `yaml:"template"`
	Id       int    `yaml:"id"`
	Key      string `yaml:"key"`
	Author   string `yaml:"author"`
}

func (t *Telegram) validate() error {
	if t.Template == "" && (t.Id == 0 || t.Key == "" || t.Author == "") {
		return errors.New("either telegram template or id, key, and author must be set")
	}

	return nil
}

type Eurocore struct {
	Url      string `yaml:"url"`
	Username string `json:"user"`
	Password string `json:"password"`
}

func (e *Eurocore) validate() error {
	if e.Url == "" {
		return errors.New("eurocore url not set")
	}

	if e.Username == "" {
		return errors.New("eurocore username not set")
	}

	if e.Password == "" {
		return errors.New("eurocore password not set")
	}

	return nil
}

type Cache struct {
	IsActive  bool   `yaml:"active"`
	CacheDir  string `yaml:"cache_dir"`
	CacheFile string `yaml:"cache_file"`
}

func (c *Cache) validate() error {
	if !c.IsActive {
		return nil
	}

	if c.CacheDir == "" {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		c.CacheDir = dir
	}

	if c.CacheFile == "" {
		c.CacheFile = "m0use_cache.txt"
	}

	return nil
}

type Log struct {
	Level    string `yaml:"level"`
	Token    string `yaml:"token"`
	Endpoint string `yaml:"endpoint"`
}

func (l *Log) validate() {
	l.Level = strings.ToLower(l.Level)

	if !slices.Contains([]string{"debug", "info", "warn", "error"}, l.Level) {
		l.Level = "info"
	}
}

type Config struct {
	User        string   `yaml:"user"`
	Region      string   `yaml:"region"`
	Telegram    Telegram `yaml:"telegram"`
	Eurocore    Eurocore `yaml:"eurocore"`
	RequestRate int      `yaml:"request_rate"`
	Cache       Cache    `yaml:"cache"`
	Log         Log      `yaml:"log"`
}

func ReadConfig(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}

	err = yaml.Unmarshal(bytes, config)
	if err != nil {
		return nil, err
	}

	err = config.validate()
	if err != nil {
		return config, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.User == "" {
		return errors.New("user not set")
	}

	if c.Region == "" {
		return errors.New("region not set")
	}

	err := c.Telegram.validate()
	if err != nil {
		return err
	}

	err = c.Eurocore.validate()
	if err != nil {
		return err
	}

	err = c.Cache.validate()
	if err != nil {
		return err
	}

	c.Log.validate()

	return nil
}
