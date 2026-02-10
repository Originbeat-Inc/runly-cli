# 使用 Runly CLI 部署逻辑资产

<p align="center">
<img src="https://assets.runly.pro/runly-logo.png" alt="Runly Protocol" width="600">


</p>

<p align="center">
  <a href="https://github.com/originbeat-inc/runly-cli/releases">
    <img src="https://img.shields.io/github/v/release/originbeat-inc/runly-cli?style=flat-square" alt="Release">
  </a>
  <a href="https://github.com/originbeat-inc/runly-cli/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/originbeat-inc/runly-cli?style=flat-square" alt="License">
  </a>
  <img src="https://img.shields.io/badge/Language-8--Supported-blue?style=flat-square" alt="I18n">
</p>

---

## 🌟 特性 (Features)

* **🔒 专家确权**: 基于 Ed25519 的数字签名，确保逻辑资产不可篡改，身份全网公认。
* **🌐 8 语支持**: 内置中、繁、英、日、韩、西、法、德八国语言，全球化开箱即用。
* **☁️ 身份同步**: 独有的“云端优先”密钥同步机制，确保 CLI 与 Web Console 身份无缝对齐。
* **🛠️ 配置先行**: 强制环境准入逻辑，规范化管理 Hub (资产) 与 Me (身份) 双服务器。
* **🚀 跨平台**: 完美支持 Windows, macOS (Intel/M1), Linux (AMD64/ARM64)。

---

## 📥 安装 (Installation)

使用一键脚本安装，系统将自动识别您的架构并配置基础环境：

```bash
curl -fsSL https://get.runly.pro/install.sh | sh

```

---

## 🚀 快速上手 (Quick Start)

### 1. 初始化配置

首次安装后，必须配置您的连接凭证和服务器地址：

```bash
runly-cli config setup

```

### 2. 同步开发者身份

关联您的专家账户，拉取或生成您的加密密钥对：

```bash
runly-cli keys generate "YourName"

```

### 3. 创建与发布资产

```bash
# 初始化项目模板
runly-cli init my-agent

# 构建并签署数字签名
runly-cli build my-agent.runly

# 推送至 Runly Hub
runly-cli publish dist.runly

```

---

## 📋 常用命令 (Command Index)

| 命令 | 描述 |
| --- | --- |
| `config setup` | **[准入]** 交互式设置服务器与 AccessToken |
| `keys generate` | **[核心]** 同步/创建身份密钥并开启云端备份 |
| `profile [name]` | 切换环境配置 (例如从 cloud 切换到 local) |
| `init [name]` | 生成符合协议标准的 `.runly` 资产模版 |
| `build [file]` | 执行哈希计算与私钥签名，生成发布级资产 |
| `publish [file]` | 将签署过的资产推送至资产中心 (Runly Hub) |
| `run [file]` | 在本地仿真引擎中测试执行逻辑 |

---

## 🗺️ 国际化 (Internationalization)

Runly CLI 会自动检测您的系统语言。您也可以通过环境变量或参数强制指定：

```bash
# 使用日语界面运行
runly-cli config setup --lang ja

```

---

## 🛠️ 开发者指南 (For Developers)

如果您希望从源代码构建项目：

```bash
# 克隆仓库
git clone https://github.com/originbeat-inc/runly-cli.git
cd runly-cli

# 使用 Makefile 一键编译
make build-all

```

---

## 📄 开源协议 (License)

本项目采用 [Apache-2.0 License](https://www.google.com/search?q=LICENSE) 协议。

---

<p align="center">
Built with ❤️ by <b>OriginBeat Inc.</b>
</p>
