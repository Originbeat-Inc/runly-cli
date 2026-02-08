# 项目基本信息
BINARY_NAME=runly-cli
VERSION=1.0.1
BUILD_DIR=build
DIST_DIR=dist
MAIN_FILE=main.go
MODULE_NAME=github.com/originbeat-inc/runly-cli

# 注入信息获取
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_TIME=$(shell date "+%Y-%m-%dT%H:%M:%S")

# 注入路径：模块名/包名
INJECT_PATH=$(MODULE_NAME)/cmd

# 编译参数：不换行，确保所有参数作为一个字符串传递给 ldflags
LDFLAGS=-ldflags "-s -w -X '$(INJECT_PATH).Version=$(VERSION)' -X '$(INJECT_PATH).GitCommit=$(GIT_COMMIT)' -X '$(INJECT_PATH).BuildTime=$(BUILD_TIME)'"

.PHONY: all clean build-local build-all build-linux build-darwin build-windows

# 默认编译
all: clean build-local

clean:
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	@echo "🧹 Cleaned old builds."

# 编译本地版本
build-local:
	@mkdir -p $(BUILD_DIR)
	@echo "🚀 Injecting into: $(INJECT_PATH)"
	@echo "   Commit: $(GIT_COMMIT)"
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "✅ Built binary: $(BUILD_DIR)/$(BINARY_NAME)"

# 跨平台一键编译打包
build-linux:
	@echo "🐧 Building Linux..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@mkdir -p $(DIST_DIR)
	tar -czf $(DIST_DIR)/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)

build-darwin:
	@echo "🍎 Building macOS..."
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@mkdir -p $(DIST_DIR)
	tar -czf $(DIST_DIR)/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)

build-windows:
	@echo "🪟 Building Windows..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME).exe $(MAIN_FILE)
	@mkdir -p $(DIST_DIR)
	zip -q -j $(DIST_DIR)/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BUILD_DIR)/$(BINARY_NAME).exe

build-all: clean build-linux build-darwin build-windows