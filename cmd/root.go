package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/originbeat-inc/runly-cli/internal/config"
	"github.com/originbeat-inc/runly-cli/internal/i18n"
	"github.com/originbeat-inc/runly-cli/internal/ui"
	"github.com/spf13/cobra"
)

// 注入变量：首字母大写，确保 Makefile 注入成功
var (
	Version   = "1.0.1"
	GitCommit = "none"
	BuildTime = "unknown"
	userLang  string
	verbose   bool
)

// rootCmd 根命令定义
var rootCmd = &cobra.Command{
	Use: "runly-cli",
}

// Execute CLI 入口
func Execute() {
	// 1. 预探测语言 (必须在所有翻译调用前)
	preDetectLanguage()

	// 2. 初始化语言包
	i18n.Init(userLang)

	// 3. 拦截 -v 或 --version 标志并执行自定义打印
	// 这样做可以绕过 Cobra 默认的简单输出，实现你的 pterm 漂亮效果
	if isVersionRequest() {
		printPrettyVersion()
		return
	}

	// 4. 配置子命令和描述翻译
	rootCmd.Version = Version
	refreshI18nDescriptions()

	// 5. 初始化内置命令
	rootCmd.InitDefaultHelpCmd()
	rootCmd.InitDefaultCompletionCmd()

	// 6. 安全检查
	ensureConfigInitialized()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// isVersionRequest 检查是否请求了版本信息
func isVersionRequest() bool {
	for _, arg := range os.Args {
		if arg == "-v" || arg == "--version" {
			return true
		}
	}
	return false
}

// printPrettyVersion 原样搬迁并优化后的艺术化输出
func printPrettyVersion() {
	ui.PrintHeader(i18n.T("cmd.root_short"))
	fmt.Printf("   %-15s %s\n", "Version:", Version)
	fmt.Printf("   %-15s %s\n", "GitCommit:", GitCommit)
	fmt.Printf("   %-15s %s\n", "BuildTime:", BuildTime)
	fmt.Printf("   %-15s %s\n", "Architecture:", "1.0 (RSS-DSL Standard)")
	fmt.Printf("   %-15s %s\n", "Language:", i18n.GetLang())
	ui.PrintFooter("© " + time.Now().Format("2006") + " OriginBeat Inc. All Rights Reserved.")
}

// refreshI18nDescriptions 刷新翻译描述
func refreshI18nDescriptions() {
	rootCmd.Short = i18n.T("cmd.root_short")
	rootCmd.Long = i18n.T("cmd.root_long")

	if helpFlag := rootCmd.Flags().Lookup("help"); helpFlag != nil {
		helpFlag.Usage = i18n.T("common.help")
	}

	for _, c := range rootCmd.Commands() {
		key := "cmd." + c.Name() + "_short"
		translated := i18n.T(key)
		if translated != key {
			c.Short = translated
		}
	}
}

// preDetectLanguage 探测语言
func preDetectLanguage() {
	for i, arg := range os.Args {
		if arg == "-l" || arg == "--lang" {
			if i+1 < len(os.Args) {
				userLang = os.Args[i+1]
				break
			}
		} else if strings.HasPrefix(arg, "--lang=") {
			userLang = arg[7:]
			break
		}
	}
	if userLang == "" {
		userLang = os.Getenv("RUNLY_LANG")
	}
	if userLang == "" {
		userLang = detectSystemLanguage()
	}
}

// detectSystemLanguage 系统语言识别
func detectSystemLanguage() string {
	langEnv := os.Getenv("LANG")
	if langEnv == "" {
		langEnv = os.Getenv("LC_ALL")
	}
	if langEnv != "" {
		base := strings.ToLower(strings.Split(langEnv, ".")[0])
		if strings.HasPrefix(base, "zh_tw") || strings.HasPrefix(base, "zh_hk") {
			return "zh-TW"
		}
		prefixes := []string{"zh", "ja", "ko", "es", "fr", "de"}
		for _, p := range prefixes {
			if strings.HasPrefix(base, p) {
				return p
			}
		}
	}
	return "en"
}

// ensureConfigInitialized 配置初始化检查
func ensureConfigInitialized() {
	var subCmd string
	if len(os.Args) > 1 {
		subCmd = os.Args[1]
	}
	if subCmd == "config" || subCmd == "help" || subCmd == "completion" || subCmd == "-l" || subCmd == "-lang" ||
		subCmd == "-v" || subCmd == "--version" || subCmd == "-h" || subCmd == "--help" || subCmd == "" {
		return
	}
	if !config.Exists() {
		ui.PrintHeader("cmd.init_header")
		ui.PrintWarning("common.warning", "Missing configuration file (config.json)")
		ui.PrintStep("🚀 Step 1: Run 'runly-cli config setup' to initialize your environment.")
		ui.PrintStep("📂 Step 2: Then use 'runly-cli init <name>' to create your project.")
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&userLang, "lang", "l", "", "Force language")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose output")
}
