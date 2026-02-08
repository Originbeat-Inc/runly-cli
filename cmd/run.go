package cmd

import (
	"fmt"
	"os"

	"github.com/originbeat-inc/runly-cli/internal/i18n"
	"github.com/originbeat-inc/runly-cli/internal/ui"
	"github.com/originbeat-inc/runly-cli/pkg/executor"
	"github.com/originbeat-inc/runly-cli/pkg/protocol"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [file.runly]",
	Short: "🚀 Execute SOP in sandbox with full AI engine support",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]

		// 1. 打印多语言 Header (🚀 RUNLY 本地仿真运行)
		ui.PrintHeader("cmd.run_header")

		// 2. 加载协议资产 (自动处理环境变量注入)
		proto, err := protocol.Load(file)
		if err != nil {
			ui.PrintError("errors.load_fail", err)
			os.Exit(1)
		}

		// 3. 准备运行上下文：注入 Dictionary 定义的默认输入
		// 这确保了即使不传递外部参数，SOP 也能依靠默认配置运行
		inputs := make(map[string]interface{})
		for _, in := range proto.Dictionary.Inputs {
			if in.Default != nil {
				inputs[in.Name] = in.Default
			}
		}

		// 4. 初始化多语言执行引擎
		engine := executor.NewEngine(proto, inputs)

		// 提示：⚙️ RUNLY 执行引擎
		ui.PrintStep("executor.engine_header")

		if err := engine.Run(); err != nil {
			// 提示：❌ 失败
			ui.PrintError("common.failure", err)
			os.Exit(1)
		}

		// 5. 运行终点：输出生成的资产报告 (Artifacts)
		// 提示：🎁 生成资产报告 (ARTIFACTS)
		ui.PrintHeader("executor.artifact_header")

		if len(engine.Context.Artifacts) == 0 {
			// 提示：⚠️ 警告: 本次运行未产生任何交付资产
			ui.PrintWarning("common.warning", i18n.T("executor.no_artifacts"))
		} else {
			for id, data := range engine.Context.Artifacts {
				// 输出资产：✅ [asset_id]
				ui.PrintSuccess(fmt.Sprintf("[%s]", id))
				// 输出状态：📊 状态: {data}
				fmt.Printf("   %s: %v\n", i18n.T("common.status"), data)
			}
		}

		// 6. 成功结语：✨ SOP 执行链路已完整结束
		fmt.Printf("\n✨ %s\n", i18n.T("executor.execution_complete"))
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
