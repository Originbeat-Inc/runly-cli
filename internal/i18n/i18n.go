package i18n

import (
	"embed"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed locales/*.yaml
var localeFS embed.FS

var bundle map[string]interface{}
var currentLang = "en"

// Init 加载语言包
func Init(lang string) {
	if lang == "" {
		lang = "en"
	}

	// 处理特殊语言映射
	if strings.HasPrefix(lang, "zh-") {
		lang = "zh-TW"
	} else if strings.HasPrefix(lang, "zh") {
		lang = "zh"
	}

	data, err := localeFS.ReadFile(fmt.Sprintf("locales/%s.yaml", lang))
	if err != nil {
		// 回退到英文
		data, _ = localeFS.ReadFile("locales/en.yaml")
		lang = "en"
	}

	currentLang = lang
	yaml.Unmarshal(data, &bundle)
}

// T 翻译函数，支持 a.b.c 格式的 Key
func T(path string) string {
	keys := strings.Split(path, ".")
	var val interface{} = bundle

	for _, key := range keys {
		if m, ok := val.(map[string]interface{}); ok {
			val = m[key]
		} else {
			return path // 找不到则返回原始 Key
		}
	}

	if s, ok := val.(string); ok {
		return s
	}
	return path
}

// GetLang 获取当前活跃语言
func GetLang() string {
	return currentLang
}
