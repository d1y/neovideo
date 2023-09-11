package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var Instance *NeovideoConfig

func InitWithFile(path string) (*NeovideoConfig, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = Init(yamlFile)
	return Instance, err
}

func Init(bte []byte) error {
	return yaml.Unmarshal(bte, &Instance)
}

func Get() *NeovideoConfig {
	return Instance
}
