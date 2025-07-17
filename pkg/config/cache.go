package config

import (
	"errors"
	"os"
	"path"
	"strings"
)

func (c *Config) ReadCache() ([]string, error) {
	if !c.Cache.IsActive {
		return nil, errors.New("caching is disabled")
	}

	cacheFilePath := path.Join(c.Cache.CacheDir, c.Cache.CacheFile)

	data, err := os.ReadFile(cacheFilePath)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(data), "\n"), nil
}

func (c *Config) WriteCache(input []string) error {
	cacheFilePath := path.Join(c.Cache.CacheDir, c.Cache.CacheFile)

	file, err := os.Create(cacheFilePath)
	if err != nil {
		return err
	}

	_, err = file.WriteString(strings.Join(input, "\n"))
	if err != nil {
		return err
	}

	return nil
}
