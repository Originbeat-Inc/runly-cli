package executor

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/originbeat-inc/runly-cli/internal/i18n"
)

// RenderTemplate 执行变量插值，例如将 {{inputs.topic}} 替换为实际值
func RenderTemplate(tpl string, ctx *Context) string {
	// 匹配 {{...}} 格式的变量
	re := regexp.MustCompile(`\{\{\s*([\w\.]+)\s*\}\}`)

	return re.ReplaceAllStringFunc(tpl, func(match string) string {
		// 提取内部路径，如 "inputs.topic"
		path := re.FindStringSubmatch(match)[1]
		parts := strings.Split(path, ".")

		if len(parts) < 2 {
			// 提示：变量引用格式错误
			return fmt.Sprintf("<! %s: %s !>", i18n.T("errors.var_format_err"), path)
		}

		domain := parts[0]
		key := parts[1]

		// 从 Context.Vars 中检索数据
		if data, ok := ctx.Vars[domain]; ok {
			if domainMap, ok := data.(map[string]interface{}); ok {
				if val, exists := domainMap[key]; exists {
					return fmt.Sprintf("%v", val)
				}
			}
		}

		// 提示：引用的变量不存在
		return fmt.Sprintf("<! %s: %s !>", i18n.T("errors.input_ref_missing"), path)
	})
}
