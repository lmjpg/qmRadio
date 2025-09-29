package main

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/adrg/xdg"
)

type Config struct {
	Radios []*Radio
}

type Radio struct {
	Name string
	Url  string
}

func ConfigLocation() (string, error) {
	return xdg.ConfigFile("qmradio/conf")
}

func NewConfig() *Config {
	return FixConfig(&Config{})
}

func FixConfig(conf *Config) *Config {
	if conf.Radios == nil {
		conf.Radios = make([]*Radio, 0)
	}
	return conf
}

func GetConfig() (*Config, error) {
	configFilePath, err := ConfigLocation()
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

	conf := Config{Radios: make([]*Radio, 0)}

	radios, ok := decoded["Radios"]
	if ok {
		_, ok := radios.([]any)
		if ok {
			for _, radio := range radios.([]any) {
				_, ok := radio.(map[string]any)
				if ok {
					radio := radio.(map[string]any)
					name, nameok := radio["Name"]
					url, urlok := radio["Url"]
					_, nameisstring := name.(string)
					_, urlisstring := url.(string)
					if nameok && urlok && nameisstring && urlisstring {
						conf.Radios = append(conf.Radios, &Radio{Name: name.(string), Url: url.(string)})
					}
				}
			}
		}
	}

	return FixConfig(&conf), nil
}

func SaveConfig(conf *Config) error {
	configFilePath, err := ConfigLocation()
	if err != nil {
		return err
	}

	data, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFilePath, data, 0666)
	if err != nil {
		return err
	}

	return nil
}

func AddRadio(conf *Config, name string, url string) error {
	conf.Radios = append(conf.Radios, &Radio{Name: name, Url: url})
	return SaveConfig(conf)
}
