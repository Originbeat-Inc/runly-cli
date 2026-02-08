package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/olekukonko/tablewriter"
	"github.com/originbeat-inc/runly-cli/internal/config"
	"github.com/originbeat-inc/runly-cli/internal/i18n"
	"github.com/originbeat-inc/runly-cli/internal/ui"
	"github.com/originbeat-inc/runly-cli/pkg/compiler"
	"github.com/originbeat-inc/runly-cli/pkg/executor/adapter"
	"github.com/originbeat-inc/runly-cli/pkg/protocol"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var hubCmd = &cobra.Command{
	Use:   "hub",
	Short: "🌐 Explore and manage assets on Runly Hub",
}

// hubTemplatesCmd: 协议版本查询
var hubTemplatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "📋 List all available Runly Protocol versions from Hub",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintHeader("cmd.hub_header")

		cfg, _ := config.LoadConfig()
		profile := cfg.GetActive()

		client := adapter.NewClient()
		client.BaseURL = profile.HubServer

		resp, err := client.Post("/v1/hub/templates/list", nil)
		if err != nil {
			ui.PrintError("common.failure", err)
			return
		}

		versions, _ := resp["versions"].([]interface{})
		latest, _ := resp["latest"].(string)

		fmt.Printf("\n%s:\n", i18n.T("cmd.hub_templates_list_title"))
		table := tablewriter.NewWriter(os.Stdout)
		// 多语言表头
		table.SetHeader([]string{i18n.T("common.version"), i18n.T("common.status")})
		table.SetBorder(false)

		for _, v := range versions {
			vStr := v.(string)
			status := "-"
			if vStr == latest {
				status = i18n.T("common.status_latest")
			}
			table.Append([]string{vStr, status})
		}
		table.Render()
		fmt.Printf("\n💡 %s\n", i18n.T("cmd.hub_templates_tip"))
	},
}

// hubSearchCmd: 资产搜索
var hubSearchCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: "Search for SOP assets on the active Hub",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := ""
		if len(args) > 0 {
			keyword = args[0]
		}

		ui.PrintHeader("cmd.hub_header")

		cfg, _ := config.LoadConfig()
		profile := cfg.GetActive()
		ui.PrintStep("executor.skill_calling", profile.HubServer)

		client := adapter.NewClient()
		client.BaseURL = profile.HubServer

		payload := map[string]interface{}{"keyword": keyword}
		data, err := client.Post("/v1/hub/search", payload)
		if err != nil {
			ui.PrintError("common.failure", err)
			return
		}

		results, ok := data["results"].([]interface{})
		if !ok || len(results) == 0 {
			ui.PrintWarning("errors.no_assets_found")
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		// 多语言表头：URN, 标题, 版本, 作者, 价格
		table.SetHeader([]string{"URN", i18n.T("common.title"), i18n.T("common.version"), i18n.T("common.creator"), i18n.T("common.price")})
		table.SetAutoWrapText(false)
		table.SetBorder(false)
		table.SetTablePadding("\t")

		for _, item := range results {
			row := item.(map[string]interface{})
			table.Append([]string{
				fmt.Sprintf("%v", row["urn"]),
				fmt.Sprintf("%v", row["title"]),
				fmt.Sprintf("%v", row["version"]),
				fmt.Sprintf("%v", row["creator"]),
				fmt.Sprintf("%v %v", row["price"], row["currency"]),
			})
		}
		fmt.Println()
		table.Render()
	},
}

// hubPullCmd: 资产拉取与校验
var hubPullCmd = &cobra.Command{
	Use:   "pull [URN]",
	Short: "Download and verify an asset from Runly Hub",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		urn := args[0]
		ui.PrintHeader("cmd.hub_pull_header") // 使用 hub 专用的 Header

		cfg, _ := config.LoadConfig()
		profile := cfg.GetActive()
		ui.PrintStep("executor.skill_calling", profile.HubServer)

		client := adapter.NewClient()
		client.BaseURL = profile.HubServer
		payload := map[string]interface{}{"urn": urn}

		data, err := client.Post("/v1/hub/pull", payload)
		if err != nil {
			ui.PrintError("common.failure", err)
			return
		}

		assetRaw, _ := data["content"].(string)
		var proto protocol.RunlyProtocol
		if err := yaml.Unmarshal([]byte(assetRaw), &proto); err != nil {
			ui.PrintError("errors.yaml_unmarshal_fail", err)
			return
		}

		// 执行多语言指纹校验步骤
		ui.PrintStep("cmd.signing_step")
		isValid, err := compiler.VerifyIntegrity(&proto)
		if err != nil || !isValid {
			ui.PrintError("errors.sign_verify_fail")
			return
		}

		fileName := fmt.Sprintf("%s.runly", proto.Manifest.URN)
		savePath := filepath.Clean(fileName)

		err = os.WriteFile(savePath, []byte(assetRaw), 0644)
		if err != nil {
			ui.PrintError("common.failure", err)
			return
		}

		ui.PrintSuccess("cmd.pull_success")
		// 多语言保存成功反馈
		fmt.Printf("\n%s: %s\n", i18n.T("cmd.verified_saved"), savePath)
	},
}

func init() {
	hubCmd.AddCommand(hubTemplatesCmd)
	hubCmd.AddCommand(hubSearchCmd)
	hubCmd.AddCommand(hubPullCmd)
	rootCmd.AddCommand(hubCmd)
}
