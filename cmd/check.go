package cmd

import (
	"fmt"
	"os"

	"github.com/originbeat-inc/runly-cli/internal/i18n"
	"github.com/originbeat-inc/runly-cli/internal/ui"
	"github.com/originbeat-inc/runly-cli/pkg/compiler"
	"github.com/originbeat-inc/runly-cli/pkg/protocol"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check [file.runly]",
	Short: "🔍 Validate protocol and signature",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		// 1. 打印多语言 Header (🔍 RUNLY 静态安全校验)
		ui.PrintHeader("cmd.check_header")

		// 2. 加载协议资产
		proto, err := protocol.Load(file)
		if err != nil {
			// 提示：读取协议文件失败
			ui.PrintError("errors.load_fail", err)
			os.Exit(1)
		}

		// 3. 7-Domain 逻辑校验
		ui.PrintStep("cmd.validate_step")
		if err := protocol.Validate(proto); err != nil {
			// 使用 i18n 翻译 "失败" 前缀并输出具体错误
			fmt.Printf("❌ [%s]: %v\n", i18n.T("common.failure"), err)
			os.Exit(1)
		}

		// 4. 数字签名与完整性验证
		ui.PrintStep("cmd.signing_step")
		isValid, err := compiler.VerifyIntegrity(proto)
		if err != nil || !isValid {
			// 提示：安全签名验证未通过
			ui.PrintError("errors.sign_verify_fail")
			os.Exit(1)
		}

		// 5. 最终反馈
		ui.PrintSuccess("common.success")

		// 额外打印资产基本信息，增加专业感
		fmt.Printf("\n%s: %s\n", i18n.T("manifest.title"), proto.Manifest.Title)
		fmt.Printf("%s: %s\n", i18n.T("manifest.urn"), proto.Manifest.URN)
		fmt.Printf("%s: %s\n", i18n.T("manifest.creator"), proto.Manifest.Creator.Name)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
