package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Profile 环境配置文件结构
type Profile struct {
	Name        string `json:"name"`
	MeServer    string `json:"me_server"`
	HubServer   string `json:"hub_server"`
	AccessToken string `json:"access_token"`
	PublicKey   string `json:"public_key"`
	MeID        string `json:"me_id"`
	SecretKey   string `json:"secret_key"`
}

// CLIConfig 根配置结构
type CLIConfig struct {
	ActiveProfile string             `json:"active_profile"`
	Profiles      map[string]Profile `json:"profiles"`
}

// GetConfigPath 返回本地存储路径 (~/.runly/config.json)
func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".runly", "config.json")
	// 确保目录存在
	os.MkdirAll(filepath.Dir(path), 0755)
	return path
}

// Exists 检查配置文件是否存在
func Exists() bool {
	// 修正：统一使用 GetConfigPath 获取的路径，确保检查的是同一个 .json 文件
	path := GetConfigPath()
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// LoadConfig 加载本地配置，如果不存在则返回初始官方环境
func LoadConfig() (*CLIConfig, error) {
	path := GetConfigPath()
	data, err := os.ReadFile(path)

	// 如果文件不存在，返回内存中的默认结构，供 config setup 命令填充
	if err != nil {
		return &CLIConfig{
			ActiveProfile: "cloud",
			Profiles: map[string]Profile{
				"cloud": {
					Name:      "cloud",
					MeServer:  "https://api.runly.me",
					HubServer: "https://api.runlyhub.com",
				},
				"local": {
					Name:      "local",
					MeServer:  "http://localhost:8080",
					HubServer: "http://localhost:8081",
				},
			},
		}, nil
	}

	var cfg CLIConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveConfig 持久化配置到磁盘
func (c *CLIConfig) SaveConfig() error {
	path := GetConfigPath()
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// GetActive 获取当前活跃的环境配置
func (c *CLIConfig) GetActive() Profile {
	// 增加防御性编程：如果 Map 里没有找到 ActiveProfile，返回一个空的 Profile
	if p, ok := c.Profiles[c.ActiveProfile]; ok {
		return p
	}
	// 兜底返回 cloud
	return c.Profiles["cloud"]
}
