package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config 配置结构体
type Config struct {
	// API 接口地址
	API string `yaml:"api"`
	// 同时下载的协程数量
	Goroutines int `yaml:"goroutines"`
}

// Load 从指定路径加载配置文件
func Load(path string) (*Config, error) {
	// 读取配置文件
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// 解析配置文件
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
