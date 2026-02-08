package main

import (
	"github.com/originbeat-inc/runly-cli/cmd"
)

func main() {
	// 所有的逻辑（包括 -v 的拦截和 i18n 初始化）都已经写在 cmd.Execute 内部了
	cmd.Execute()
}
