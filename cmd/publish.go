package cmd

import (
	"fmt"
	"os"

	"github.com/originbeat-inc/runly-cli/internal/config"
	"github.com/originbeat-inc/runly-cli/internal/i18n"
	"github.com/originbeat-inc/runly-cli/internal/ui"
	"github.com/originbeat-inc/runly-cli/pkg/compiler"
	"github.com/originbeat-inc/runly-cli/pkg/executor/adapter"
	"github.com/originbeat-inc/runly-cli/pkg/protocol"
	"github.com/spf13/cobra"
)

var publishCmd = &cobra.Command{
	Use:   "publish [file.runly]",
	Short: "📤 Publish asset to Runly Hub with progress tracking",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		// 1. 打印多语言 Header (📤 RUNLY HUB 资产发布)
		ui.PrintHeader("cmd.publish_header")

		// 2. 获取环境配置 (Profile)
		cfg, _ := config.LoadConfig()
		profile := cfg.GetActive()

		// 3. 加载协议资产
		proto, err := protocol.Load(file)
		if err != nil {
			ui.PrintError("errors.load_fail", err)
			os.Exit(1)
		}

		// 4. 安全指纹校验 (发布前最后一道关卡)
		ui.PrintStep("cmd.validate_step")

		// 检查是否有签名位
		if proto.Security.Signature == "" {
			ui.PrintError("errors.no_sig")
			os.Exit(1)
		}

		// 验证资产完整性与签名有效性
		isValid, err := compiler.VerifyIntegrity(proto)
		if err != nil || !isValid {
			ui.PrintError("errors.sign_verify_fail")
			os.Exit(1)
		}

		// 5. 执行多语言进度上传
		client := adapter.NewClient()
		client.BaseURL = profile.HubServer // 切换至当前 Profile 指定的 Hub 节点地址

		payload := map[string]interface{}{
			"asset": proto,
			"urn":   proto.Manifest.URN,
		}

		// 动态生成进度描述，例如："🚀 正在同步资产至 [Official] 节点..."
		progressMsg := fmt.Sprintf("%s [%s]...", i18n.T("cmd.publish_progress_msg"), profile.Name)

		// 调用带进度的 POST 请求，将动态生成的描述传入
		_, err = client.PostWithProgress("/v1/hub/publish", payload, progressMsg)
		if err != nil {
			ui.PrintError("common.failure", err)
			os.Exit(1)
		}

		// 6. 成功反馈与 URN 确认
		ui.PrintSuccess("common.success")

		// 生成展示用的规范 URN 路径
		finalURN := fmt.Sprintf("runly://hub/%s/%s", profile.MeID, proto.Manifest.URN)
		fmt.Printf("\n📦 %s: %s\n", i18n.T("cmd.publish_live_msg"), finalURN)
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
}
