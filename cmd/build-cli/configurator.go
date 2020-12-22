package main

import "gopkg.in/yaml.v2"

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
