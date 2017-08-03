package rconfig

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// VarDefaultsFromYAMLFile reads contents of a file and calls VarDefaultsFromYAML
func VarDefaultsFromYAMLFile(filename string) map[string]string {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return make(map[string]string)
	}

	return VarDefaultsFromYAML(data)
}

// VarDefaultsFromYAML creates a vardefaults map from YAML raw data
func VarDefaultsFromYAML(in []byte) map[string]string {
	out := make(map[string]string)
	err := yaml.Unmarshal(in, &out)
	if err != nil {
		return make(map[string]string)
	}
	return out
}
