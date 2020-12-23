package main

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// CompositionConfiguration yaml based configuration for windtunnel
type CompositionConfiguration struct {
	Schema struct {
		BaseDir    string `yaml:"baseDir"`
		Stage      string `yaml:"stage"`
		Definition string `yaml:"definition"`
	}
	Stages []struct {
		Name   string `yaml:"name"`
		RunCmd string `yaml:"runCmd"`
		Code   struct {
			BaseDir   string `yaml:"baseDir"`
			PreStage  string `yaml:"preStage"`
			Stage     string `yaml:"stage"`
			PostStage string `yaml:"postStage"`
		}
	}
}

func (config *CompositionConfiguration) marshal(data []byte) error {
	err := yaml.Unmarshal([]byte(data), config)

	return err
}

func (config *CompositionConfiguration) validate() error {
	stageSchemaPath := fmt.Sprintf("%s/%s", config.Schema.BaseDir, config.Schema.Stage)
	definitionSchemaPath := fmt.Sprintf("%s/%s", config.Schema.BaseDir, config.Schema.Definition)

	if _, err := os.Stat(stageSchemaPath); os.IsNotExist(err) {
		return errors.New("cannot locate stage schema file")
	}

	if _, err := os.Stat(definitionSchemaPath); os.IsNotExist(err) {
		return errors.New("cannot locate definition schema file")
	}

	return nil
}
