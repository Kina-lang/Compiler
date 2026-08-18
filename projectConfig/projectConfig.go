package projectConfig

import "github.com/pelletier/go-toml/v2"

type ProjectConfigProject struct {
	Name string `toml:"name"`
	Author string `toml:"author"`
	Version string `toml:"version"`
	Entry string `toml:"entry"`
}

type ProjectConfigDependency struct {
	Version string
}

type ProjectConfig struct {
	Project ProjectConfigProject `toml:"project"`
	Dependencies map[string]ProjectConfigDependency `toml:"dependencies"`
}

func (c *ProjectConfig) String() (string, error) {
	b, err := toml.Marshal(c)

	return string(b), err
}
