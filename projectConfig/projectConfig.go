package projectConfig

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type ProjectConfigProject struct {
	Name    string `toml:"name"`
	Author  string `toml:"author"`
	Version string `toml:"version"`
	Entry   string `toml:"entry"`
}

type ProjectConfigDependency struct {
	Version string
}

type ProjectConfig struct {
	Project      ProjectConfigProject               `toml:"project"`
	Dependencies map[string]ProjectConfigDependency `toml:"dependencies"`
}

func (c *ProjectConfig) String() (string, error) {
	b, err := toml.Marshal(c)

	return string(b), err
}

func ParseFile(path string) (*ProjectConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config ProjectConfig
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	err = Validate(config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
