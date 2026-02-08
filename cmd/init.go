package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/originbeat-inc/runly-cli/internal/config"
	"github.com/originbeat-inc/runly-cli/internal/i18n"
	"github.com/originbeat-inc/runly-cli/internal/ui"
	"github.com/originbeat-inc/runly-cli/pkg/executor/adapter"
	"github.com/spf13/cobra"
)

var protoVersion string

var initCmd = &cobra.Command{
	Use:   "init [project_name]",
	Short: "🌟 Initialize a new Runly protocol project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		fileName := projectName + ".runly"

		ui.PrintHeader("cmd.init_header")

		// 1. 初始化 Client 并指向 HubServer (确保访问 https://api.runlyhub.com)
		client := adapter.NewClient().SetToHubServer()

		// 2. 远程获取模板逻辑优化
		// 构造请求载荷
		payload := map[string]interface{}{
			"version": protoVersion,
		}

		// 执行带进度的请求 (调用我们 http.go 里的 PostWithProgress)
		ui.PrintStep("cmd.init_pulling_remote")
		data, err := client.PostWithProgress("/v1/hub/templates/pull", payload, "common.downloading")
		if err != nil {
			ui.PrintError("errors.load_fail", err)
			os.Exit(1)
		}

		rawTemplate, ok := data["content"].(string)
		if !ok {
			ui.PrintError("errors.server_err")
			os.Exit(1)
		}

		// 3. 获取 Profile 进行确权注入
		cfg, _ := config.LoadConfig()
		profile := cfg.GetActive()

		// 4. 执行动态占位符替换
		replacer := strings.NewReplacer(
			"{{PROJECT_NAME}}", projectName,
			"{{ME_ID}}", fallback(profile.Name, "me_placeholder"), // 假设 Name 作为 ID 标识
			"{{CREATOR_NAME}}", fallback(profile.Name, "Solo Creator"),
			"{{PUB_KEY}}", fallback(profile.PublicKey, "ed25519:not_found"),
			"{{CREATED_AT}}", time.Now().Format(time.RFC3339),
			"{{PROTO_VERSION}}", protoVersion,
		)
		finalContent := replacer.Replace(rawTemplate)

		// 5. 写入本地文件
		if err := os.WriteFile(fileName, []byte(finalContent), 0644); err != nil {
			ui.PrintError("errors.load_fail", err)
			os.Exit(1)
		}

		ui.PrintSuccess("common.success")
		fmt.Printf("\n🚀 %s [%s]\n📄 %s: %s\n",
			i18n.T("cmd.init_success_msg"), projectName,
			i18n.T("cmd.init_file_path"), fileName)
	},
}

func fallback(val, def string) string {
	if val == "" {
		return def
	}
	return val
}

func init() {
	initCmd.Flags().StringVarP(&protoVersion, "proto", "p", "latest", "Force pull specific protocol version")
	rootCmd.AddCommand(initCmd)
}
