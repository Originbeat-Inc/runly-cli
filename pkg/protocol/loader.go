package protocol

import (
	"fmt"
	"os"
	"regexp"

	"github.com/originbeat-inc/runly-cli/internal/i18n"
	"gopkg.in/yaml.v3"
)

// envRegex 匹配占位符格式：{{env.VARIABLE_NAME}}
var envRegex = regexp.MustCompile(`\{\{\s*env\.([a-zA-Z0-9_]+)\s*\}\}`)

// Load 负责从磁盘加载协议并执行动态预处理（如环境变量注入）
func Load(path string) (*RunlyProtocol, error) {
	// 1. 读取文件原始字节流
	data, err := os.ReadFile(path)
	if err != nil {
		// 📂 读取协议文件失败: %v
		return nil, fmt.Errorf(i18n.T("errors.load_fail"), err)
	}

	// 2. 环境变量热注入 (Secret Injection)
	// 在解析结构化对象前，先替换掉内存中的敏感信息占位符，保护密钥安全
	processedData := injectSecrets(data)

	// 3. 执行 YAML 反序列化
	var proto RunlyProtocol
	if err := yaml.Unmarshal(processedData, &proto); err != nil {
		// 🧩 协议语法解析异常，请检查 YAML 格式: %v
		return nil, fmt.Errorf(i18n.T("errors.yaml_unmarshal_fail"), err)
	}

	return &proto, nil
}

// injectSecrets 查找并替换所有的环境变量占位符
func injectSecrets(input []byte) []byte {
	return envRegex.ReplaceAllFunc(input, func(match []byte) []byte {
		submatch := envRegex.FindSubmatch(match)
		if len(submatch) < 2 {
			return match
		}

		envKey := string(submatch[1])
		// 从系统环境变量中查询真实值
		realValue, exists := os.LookupEnv(envKey)

		if !exists {
			// 如果缺失环境变量，保留占位符，由后续的 Validator 进行语义层面的拦截提示
			return match
		}

		return []byte(realValue)
	})
}
