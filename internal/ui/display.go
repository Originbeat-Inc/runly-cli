package ui

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/originbeat-inc/runly-cli/internal/i18n"
	"github.com/pterm/pterm"
)

// PrintHeader 打印带图标的多语言业务标题
func PrintHeader(key string) {
	title := i18n.T(key)

	pterm.DefaultHeader.
		WithFullWidth().
		WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)). // 使用蓝色作为命令操作的主色调
		WithTextStyle(pterm.NewStyle(pterm.FgLightWhite)).
		WithMargin(2).
		Println(title)
	pterm.Println() // 留白一行，更美观
}

func PrintFooter(key string) {
	footer := i18n.T(key)
	pterm.Println() // 留白一行，更美观
	pterm.DefaultHeader.
		WithFullWidth().
		WithBackgroundStyle(pterm.NewStyle(pterm.BgBlue)). // 使用蓝色作为命令操作的主色调
		WithTextStyle(pterm.NewStyle(pterm.FgLightWhite)).
		WithMargin(1).
		Println(footer)
}

// PrintStep 打印带图标的多语言执行步骤
func PrintStep(key string, args ...interface{}) {
	message := fmt.Sprintf(i18n.T(key), args...)
	pterm.Info.Println(message)
}

// PrintSuccess 打印多语言成功反馈
func PrintSuccess(key string) {
	pterm.Success.Println(i18n.T(key))
}

// PrintError 打印多语言错误反馈
func PrintError(key string, args ...interface{}) {
	message := fmt.Sprintf(i18n.T(key), args...)
	pterm.Error.Println(message)
}

// PrintWarning 打印多语言警告
func PrintWarning(key string, args ...interface{}) {
	message := fmt.Sprintf(i18n.T(key), args...)
	pterm.Warning.Println(message)
}

// StartLoading 启动一个多语言感知的加载动画
func StartLoading(key string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = " " + i18n.T(key)
	s.Start()
	return s
}

// PrintKV 打印多语言键值对 (修复点：确保此函数存在)
// 效果: 📊 状态: [Value]
func PrintKV(key string, value interface{}) {
	label := i18n.T(key)
	pterm.Printf("%s: %v\n", label, pterm.ThemeDefault.SecondaryStyle.Sprint(value))
}

// ShowProgress 展示多语言感知的进度条 (用于 Publish 或 Pull)
func ShowProgress(key string, total int) {
	// 获取多语言描述文本
	description := i18n.T(key)

	// 修复：使用 WithTitle 替换 WithDescription，并显式配置进度条
	p, err := pterm.DefaultProgressbar.
		WithTotal(total).
		WithTitle(description). // 最新版 pterm 使用 Title 替代 Description
		WithShowCount(true).
		WithShowPercentage(true).
		WithRemoveWhenDone(false).
		Start()

	if err != nil {
		// 容错处理：如果启动失败，退化为普通日志
		pterm.Info.Println(description)
		return
	}

	// 执行进度模拟
	for i := 0; i < total; i++ {
		p.Increment()
		time.Sleep(time.Millisecond * 30)
	}

	// 显式停止
	_, _ = p.Stop()

	// 打印完成后的换行
	pterm.Println()
}
