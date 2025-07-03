package config

import (
	yaml "gopkg.in/yaml.v2"
	"os"
)

var Cfg *Config

// Config 结构体用于存储配置信息
type Config struct {
	Server    Server      `yaml:"server"`
	Huawei    Huawei      `yaml:"huawei"`
	Alerts    []Alerts    `yaml:"alerts"`
	Receivers []Receivers `yaml:"receivers"`
}

type Server struct {
	Port int `yaml:"port"`
}

type Huawei struct {
	AppKey         string `yaml:"app_key"`
	AppSecret      string `yaml:"app_secret"`
	ApiAddress     string `yaml:"api_address"`
	Sender         string `yaml:"sender"`
	TemplateId     string `yaml:"template_id"`
	Signature      string `yaml:"signature"`
	StatusCallBack string `yaml:"status_call_back"`
}

type Alerts struct {
	Type    string `yaml:"type"`
	Enabled bool   `yaml:"enabled"`
}

type Receivers struct {
	ContactNumbers []string `yaml:"contact_numbers"`
	Name           string   `yaml:"name"`
}

// LoadConfig 从指定的 YAML 文件加载配置
func LoadConfig(filePath string) (*Config, error) {
	config := &Config{}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func InitConfig(filePath string) error {
	cfg, err := LoadConfig(filePath)
	if err != nil {
		return err
	}
	Cfg = cfg
	return nil
}
