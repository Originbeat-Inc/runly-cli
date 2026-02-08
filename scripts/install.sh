#!/bin/bash

set -e

# --- 配置区 ---
VERSION="1.0.1"
BASE_URL="https://get.runly.pro/dist"
BINARY_NAME="runly-cli"
CONF_DIR="$HOME/.runly"
CONF_FILE="$CONF_DIR/config.json"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Starting Runly CLI $VERSION Installation...${NC}"

# 1. 检测系统与架构
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo -e "${RED}❌ Unsupported architecture: $ARCH${NC}"; exit 1 ;;
esac

case "$OS" in
    linux) PLATFORM="linux-$ARCH" ;;
    darwin) PLATFORM="darwin-$ARCH" ;;
    msys*|mingw*|cygwin*) PLATFORM="windows-amd64"; BINARY_NAME="runly-cli.exe" ;;
    *) echo -e "${RED}❌ Unsupported OS: $OS${NC}"; exit 1 ;;
esac

# 2. 下载二进制文件
DOWNLOAD_URL="${BASE_URL}/${BINARY_NAME}-${VERSION}-${PLATFORM}.tar.gz"
if [[ "$PLATFORM" == *"windows"* ]]; then
    DOWNLOAD_URL="${BASE_URL}/${BINARY_NAME}-${VERSION}-${PLATFORM}.zip"
fi

echo -e "${BLUE}📥 Downloading from: $DOWNLOAD_URL${NC}"
TMP_DIR=$(mktemp -d)
curl -L "$DOWNLOAD_URL" -o "$TMP_DIR/runly_package"

# 3. 解压并安装
cd "$TMP_DIR"
if [[ "$DOWNLOAD_URL" == *.zip ]]; then
    unzip -q runly_package
else
    tar -xzf runly_package
fi

echo -e "${BLUE}🔧 Installing to /usr/local/bin...${NC}"
if [[ "$OS" == "linux" || "$OS" == "darwin" ]]; then
    chmod +x "$BINARY_NAME"
    # 使用 sudo 移动到系统目录
    if [ -w "/usr/local/bin" ]; then
        mv "$BINARY_NAME" /usr/local/bin/runly-cli
    else
        sudo mv "$BINARY_NAME" /usr/local/bin/runly-cli
    fi
else
    # Windows 模式下尝试放在环境变量路径，这里假设用户有 ~/bin
    mkdir -p "$HOME/bin"
    mv "$BINARY_NAME" "$HOME/bin/runly-cli.exe"
    echo -e "${RED}⚠️ Please ensure $HOME/bin is in your PATH environment variable.${NC}"
fi

# 4. 默认配置初始化 (Critical Step)
echo -e "${BLUE}⚙️ Initializing default config.json...${NC}"
mkdir -p "$CONF_DIR"

if [ ! -f "$CONF_FILE" ]; then
cat <<EOF > "$CONF_FILE"
{
  "active_profile": "cloud",
  "profiles": {
    "cloud": {
      "name": "cloud",
      "me_server": "https://api.runly.me",
      "hub_server": "https://api.runlyhub.com",
      "access_token": "",
      "public_key": "",
      "me_id": "",
      "secret_key": ""
    },
    "local": {
      "name": "local",
      "me_server": "http://localhost:8080",
      "hub_server": "http://localhost:8081",
      "access_token": "",
      "public_key": "",
      "me_id": "",
      "secret_key": ""
    }
  }
}
EOF
    echo -e "${GREEN}✅ Created default configuration at $CONF_FILE${NC}"
else
    echo -e "${BLUE}ℹ️ Configuration already exists, skipping initialization.${NC}"
fi

# 5. 完成提示
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✨ Runly CLI $VERSION installed successfully!${NC}"
echo -e "${BLUE}👉 Next step: Run 'runly-cli config setup' to set your token.${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

# 清理
rm -rf "$TMP_DIR"