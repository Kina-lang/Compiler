package projectConfig

import (
	"fmt"
)

func Validate(config ProjectConfig) error {
	err := validateProject(config.Project)
	if err != nil {
		return fmt.Errorf("invalid project config: %w", err)
	}

	return nil
}

func validateProject(config ProjectConfigProject) error {
	if config.Name == "" {
		return fmt.Errorf("project.name is required")
	}

	if config.Entry == "" {
		return fmt.Errorf("project.entry is required")
	}

	return nil
}
