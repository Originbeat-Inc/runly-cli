# 🤝 Contributing to Runly CLI

首先，感谢你考虑为 Runly CLI 做出贡献！正是有了社区的支持，我们才能构建出更好的 AI Agent 部署工具。

在提交代码之前，请花几分钟时间阅读以下指南。

## 🛠️ 开发环境准备

在开始贡献之前，请确保你的本地环境已安装：
* **Go**: 1.21 或更高版本
* **Git**: 用于版本控制
* **Make**: 用于执行 Makefile 指令

### 克隆项目
```bash
git clone [https://github.com/originbeat-inc/runly-cli.git](https://github.com/originbeat-inc/runly-cli.git)
cd runly-cli
go mod download
```

### 🌐 国际化 (I18n) 指南

由于 Runly CLI 支持 8 国语言，如果你添加了新的 UI 提示或命令描述，请务必更新：

- internal/i18n/locales/ 下的各个语言文件。

- 即使你不会其他语言，也请至少更新 en.yaml 和 zh.yaml。

### 📜 代码规范

我们遵循标准的 Go 代码风格：

- 使用 gofmt 或 goimports 格式化代码。

- 错误处理：请使用 internal/ui 包中的 PrintError 或 PrintWarning 来向用户展示错误，确保支持国际化。

- 配置访问：始终通过 internal/config 加载配置，不要直接读写磁盘文件。

- 身份安全：严禁在代码中硬编码任何私钥或敏感信息。

### 🚀 提交流程 (Workflow)

1、**Fork** 本仓库并创建一个新的分支 (Branch)。
```bash
git checkout -b feat/your-feature-name
```

2、**编写代码** 并确保通过本地编译。
```bash
make build
```

3、**测试：** 确保你的修改没有破坏 `config setup` 或 `keys generate` 等核心逻辑。

4、**提交代码：**
```bash
git commit -m "feat: add support for new server protocols"
```

5、**发起 Pull Request (PR)：**
- 请在 PR 描述中详细说明你的修改动机和测试结果。
- 勾选对应的 `Issue` 编号（如果有）。

### 📝 提交信息规范 (Commit Messages)

我们推荐使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

- `feat:` 新功能

- `fix:` 修复 Bug

- `docs:` 文档修改

- `style:` 代码格式（不影响功能的修改）

- `refactor:` 代码重构

- `i18n:` 国际化词条更新

### ⚖️ 许可协议
通过向本项目贡献代码，你同意你的贡献将基于 [Apache-2.0 License](https://www.google.com/search?q=LICENSE) 进行授权。

---
再次感谢你对 **Runly CLI** 的支持！

<p align="center">
Built with ❤️ by <b>OriginBeat Inc.</b>
</p>