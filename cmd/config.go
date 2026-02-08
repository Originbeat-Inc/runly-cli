package cmd

import (
	"github.com/originbeat-inc/runly-cli/internal/config"
	"github.com/originbeat-inc/runly-cli/internal/i18n"
	"github.com/originbeat-inc/runly-cli/internal/ui"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: i18n.T("cmd.config_short"),
}

// setCmd: 用于快捷设置 Token 或服务器地址
var setCmd = &cobra.Command{
	Use:     "set",
	Short:   i18n.T("cmd.config_set_short"),
	Example: "  runly-cli config set --token XXX\n  runly-cli config set --hub https://api.runlyhub.com",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		profile := cfg.Profiles[cfg.ActiveProfile]

		token, _ := cmd.Flags().GetString("token")
		hub, _ := cmd.Flags().GetString("hub")
		me, _ := cmd.Flags().GetString("me")

		modified := false

		if token != "" {
			profile.AccessToken = token
			modified = true
		}
		if hub != "" {
			profile.HubServer = hub
			modified = true
		}
		if me != "" {
			profile.MeServer = me
			modified = true
		}

		if modified {
			cfg.Profiles[cfg.ActiveProfile] = profile
			_ = cfg.SaveConfig()
			ui.PrintSuccess("common.success")
			// 修正点 1: 确保这里的提示文本已在语言包定义
			ui.PrintStep(i18n.T("cmd.config_updated_tip") + ": " + cfg.ActiveProfile)
		} else {
			// 修正点 2: 确保 common.warning 对应 UI 逻辑
			ui.PrintWarning("common.warning", i18n.T("errors.invalid_args"))
		}
	},
}

// setupCmd: 交互式引导配置
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: i18n.T("cmd.config_setup_short"),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		profile := cfg.Profiles[cfg.ActiveProfile]

		// 交互式输入：Hub Server 地址
		hub, _ := pterm.DefaultInteractiveTextInput.
			WithDefaultText(profile.HubServer).
			WithDefaultValue(profile.HubServer).
			Show(i18n.T("cmd.config_prompt_hub"))

		// 修正点 3: 补全缺失的 Me Server 交互配置 (因为 Profile 里有这个字段)
		me, _ := pterm.DefaultInteractiveTextInput.
			WithDefaultText(profile.MeServer).
			WithDefaultValue(profile.MeServer).
			Show(i18n.T("cmd.config_prompt_me"))

		// 交互式输入：Token
		token, _ := pterm.DefaultInteractiveTextInput.
			WithMask("*").
			Show(i18n.T("cmd.config_prompt_token"))

		profile.HubServer = hub
		profile.MeServer = me // 赋值
		if token != "" {
			profile.AccessToken = token
		}

		cfg.Profiles[cfg.ActiveProfile] = profile
		_ = cfg.SaveConfig()

		ui.PrintSuccess("common.success")
	},
}

func init() {
	setCmd.Flags().String("token", "", i18n.T("cmd.config_flag_token"))
	setCmd.Flags().String("hub", "", i18n.T("cmd.config_flag_hub"))
	setCmd.Flags().String("me", "", i18n.T("cmd.config_flag_me"))

	configCmd.AddCommand(setCmd)
	configCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(configCmd)
}
