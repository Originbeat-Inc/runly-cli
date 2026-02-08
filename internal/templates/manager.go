package templates

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/originbeat-inc/runly-cli/internal/i18n"
	"github.com/originbeat-inc/runly-cli/pkg/executor/adapter"
)

// GetTemplate 强制从远程 Hub 获取协议模板
func GetTemplate(version string) (string, error) {
	if version == "" || version == "latest" {
		version = "latest"
	}

	home, _ := os.UserHomeDir()
	cachePath := filepath.Join(home, ".runly", "templates", version+".runly")

	// 1. 明确提示正在进行远程拉取
	fmt.Printf("📡 %s (v%s)...\n", i18n.T("cmd.init_pulling_remote"), version)

	client := adapter.NewClient()
	// 调用 Hub 接口
	resp, err := client.Post("/v1/hub/templates/pull", map[string]interface{}{
		"version": version,
	})

	if err != nil {
		// 远程失败时的兜底逻辑：尝试读取本地旧缓存
		if data, cacheErr := os.ReadFile(cachePath); cacheErr == nil {
			fmt.Printf("⚠️ %s\n", i18n.T("errors.remote_pull_failed_use_cache"))
			return string(data), nil
		}
		return "", fmt.Errorf("remote pull failed: %v", err)
	}

	content, ok := resp["content"].(string)
	if !ok {
		return "", fmt.Errorf("invalid template content received")
	}

	// 2. 更新本地缓存
	_ = os.MkdirAll(filepath.Dir(cachePath), 0755)
	_ = os.WriteFile(cachePath, []byte(content), 0644)

	return content, nil
}
