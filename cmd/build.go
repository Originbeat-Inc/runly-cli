package cmd

import (
	"fmt"
	"os"

	"github.com/originbeat-inc/runly-cli/internal/config"
	"github.com/originbeat-inc/runly-cli/internal/i18n"
	"github.com/originbeat-inc/runly-cli/internal/ui"
	"github.com/originbeat-inc/runly-cli/pkg/compiler"
	"github.com/originbeat-inc/runly-cli/pkg/protocol"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var buildCmd = &cobra.Command{
	Use:   "build [file.runly]",
	Short: "🛠️  Compile, Sign and Solidify the SOP asset",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		// 1. 打印多语言 Header (🌟 BUILD 资产编译)
		ui.PrintHeader("cmd.build_header")

		// 2. 获取当前 Profile 身份用于签名
		cfg, _ := config.LoadConfig()
		profile := cfg.GetActive()
		if profile.SecretKey == "" {
			// 提示：未检测到有效密钥，请先运行 runly-cli keys generate
			ui.PrintError("errors.no_key")
			os.Exit(1)
		}

		// 3. 加载协议文件
		proto, err := protocol.Load(file)
		if err != nil {
			// 提示：加载协议文件失败
			ui.PrintError("errors.load_fail", err)
			os.Exit(1)
		}

		// 4. 静态语义与 7-Domain 校验
		ui.PrintStep("cmd.validate_step")
		if err := protocol.Validate(proto); err != nil {
			// 校验失败输出
			fmt.Printf("❌ [%s]: %v\n", i18n.T("common.failure"), err)
			os.Exit(1)
		}

		// 5. 强制确权：同步 Creator 信息为当前环境 MeID
		proto.Manifest.Creator.MeID = profile.MeID
		proto.Manifest.Creator.PubKey = profile.PublicKey

		// 6. 执行编译、哈希计算与数字签名
		ui.PrintStep("cmd.signing_step")
		hash, err := compiler.BuildArtifact(proto, profile.SecretKey)
		if err != nil {
			ui.PrintError("common.failure", err)
			os.Exit(1)
		}

		// 7. 导出固化资产 (dist.runly)
		finalData, _ := yaml.Marshal(proto)
		distFile := "dist.runly"
		if err := os.WriteFile(distFile, finalData, 0644); err != nil {
			ui.PrintError("errors.load_fail", err)
			os.Exit(1)
		}

		// 8. 成功反馈
		ui.PrintSuccess("cmd.export_success")
		fmt.Printf("📄 %s: %s\n", i18n.T("common.output"), distFile)
		fmt.Printf("🔐 %s:   %s\n", i18n.T("common.hash"), hash)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
