package cmd

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/originbeat-inc/runly-cli/internal/config"
	"github.com/originbeat-inc/runly-cli/internal/i18n"
	"github.com/originbeat-inc/runly-cli/internal/ui"
	"github.com/originbeat-inc/runly-cli/pkg/executor/adapter"
	"github.com/spf13/cobra"
)

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "🔑 Identity & Runly Me Management",
}

// generateCmd：云端优先逻辑
var generateCmd = &cobra.Command{
	Use:   "generate [username]",
	Short: "Sync or create identity keys with Runly Cloud",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintHeader("cmd.keys_header")

		// 1. 加载配置
		cfg, err := config.LoadConfig()
		if err != nil {
			ui.PrintError("common.failure", err)
			os.Exit(1)
		}
		profile := cfg.GetActive()

		// 2. 初始化 Client 并强制切换到 MeServer
		client := adapter.NewClient().SetToMeServer()
		ui.PrintStep("executor.syncing")

		// 3. 尝试从云端拉取
		cloudData, err := client.Post("/v1/me/keys/pull", nil)

		if err == nil {
			// --- 情况 A: 云端已有密钥 ---
			cloudPub, _ := cloudData["public_key"].(string)
			cloudPriv, _ := cloudData["secret_key"].(string)
			cloudMeID, _ := cloudData["me_id"].(string)

			// 一致性校验：如果本地已有密钥且与云端不符，报错拦截
			if profile.PublicKey != "" && profile.PublicKey != cloudPub {
				ui.PrintWarning("common.warning", "Identity Mismatch!")
				fmt.Printf("   Local Public Key: %s\n", profile.PublicKey)
				fmt.Printf("   Cloud Public Key: %s\n", cloudPub)
				ui.PrintError("errors.auth_failed", "Keys mismatch between local and cloud.")
				os.Exit(1)
			}

			// 保存拉取的密钥到本地
			saveKeys(cfg, cloudMeID, cloudPub, cloudPriv)
			ui.PrintSuccess("common.success")
			ui.PrintStep("Identity synced from Cloud Console.")

		} else {
			// --- 情况 B: 云端无密钥，执行本地生成并同步 ---
			ui.PrintStep("No identity found on cloud. Generating new keypair...")

			pub, priv, _ := ed25519.GenerateKey(rand.Reader)
			pubHex := hex.EncodeToString(pub)
			privSeedHex := hex.EncodeToString(priv.Seed())

			username := "anonymous"
			if len(args) > 0 {
				username = args[0]
			}

			// 构造同步载荷 (包含私钥用于 Web 可视化编辑签名)
			payload := map[string]interface{}{
				"username":   username,
				"public_key": pubHex,
				"secret_key": privSeedHex,
			}

			resp, err := client.Post("/v1/me/keys/sync", payload)
			if err != nil {
				ui.PrintError("common.failure", err)
				os.Exit(1)
			}

			meID, _ := resp["me_id"].(string)
			if meID == "" {
				meID = "me_0x" + pubHex[:12]
			}

			saveKeys(cfg, meID, pubHex, privSeedHex)
			ui.PrintSuccess("cmd.key_gen_success")
		}

		// 4. 打印最终身份状态
		printCurrentIdentity(cfg)
	},
}

// showCmd：显示本地身份 (100% 还原你的 Emoji 格式)
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current identity information",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, _ := config.LoadConfig()
		profile := cfg.GetActive()

		if profile.SecretKey == "" {
			ui.PrintError("errors.no_key")
			return
		}

		ui.PrintHeader("cmd.keys_header")

		fmt.Printf("📊 %s: %s\n", i18n.T("common.status"), cfg.ActiveProfile)
		fmt.Printf("👤 MeID:    %s\n", profile.MeID)
		fmt.Printf("🔑 %s:  %s\n", i18n.T("manifest.pub_key_label"), profile.PublicKey)
		fmt.Printf("🌐 %s:  %s\n", i18n.T("manifest.server_label"), profile.MeServer)
	},
}

// --- 辅助函数：修复了 cfg 的引用问题 ---

func saveKeys(cfg *config.CLIConfig, meID, pub, priv string) {
	p := cfg.Profiles[cfg.ActiveProfile]
	p.MeID = meID
	p.PublicKey = pub
	p.SecretKey = priv
	cfg.Profiles[cfg.ActiveProfile] = p
	_ = cfg.SaveConfig()
}

func printCurrentIdentity(cfg *config.CLIConfig) {
	profile := cfg.GetActive()
	fmt.Printf("\n🆔 MeID:   %s\n", profile.MeID)
	fmt.Printf("🔑 %s: %s\n", i18n.T("manifest.pub_key_label"), profile.PublicKey)
	fmt.Printf("🌐 %s: %s (%s)\n", i18n.T("common.status"), cfg.ActiveProfile, profile.MeServer)
	fmt.Printf("\n💡 %s\n", i18n.T("cmd.keys_info_tip"))
}

func init() {
	keysCmd.AddCommand(generateCmd)
	keysCmd.AddCommand(showCmd)
	rootCmd.AddCommand(keysCmd)
}
