package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/adrg/xdg"
)

type Config struct {
	Radios []Radio
}

type Radio struct {
	Name string
	Url  string
}

func NewConfig() *Config {
	return FixConfig(&Config{})
}

func FixConfig(conf *Config) *Config {
	return conf
}

func GetConfig() (*Config, error) {
	configFilePath, err := xdg.ConfigFile("qmradio/conf")
	if err != nil {
		return NewConfig(), err
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return NewConfig(), nil
		}
		return NewConfig(), err
	}

	var decoded map[string]any
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		return NewConfig(), err
	}

	conf := Config{Radios: make([]Radio, 0)}

	radios, ok := decoded["radios"]
	if ok {
		_, ok := radios.([]any)
		if ok {
			for _, radio := range radios.([]any) {
				_, ok := radio.(map[string]any)
				if ok {
					radio := radio.(map[string]any)
					name, nameok := radio["name"]
					url, urlok := radio["url"]
					_, nameisstring := name.(string)
					_, urlisstring := url.(string)
					if nameok && urlok && nameisstring && urlisstring {
						conf.Radios = append(conf.Radios, Radio{Name: name.(string), Url: url.(string)})
					}
				}
			}
		}
	}

	return &conf, nil
}
