# 项目信息
BINARY_NAME=netflood
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# 源文件目录
CMD_DIR=./cmd
BUILD_DIR=./build

# Go 相关命令
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# 颜色输出
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
NC=\033[0m # No Color

.PHONY: all build clean test help linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64 cross-compile

# 默认目标
all: build

# 帮助信息
help:
	@echo "$(GREEN)NetFlood 编译工具$(NC)"
	@echo ""
	@echo "$(YELLOW)可用命令:$(NC)"
	@echo "  make build           - 编译当前平台版本"
	@echo "  make linux-amd64     - 编译 Linux AMD64 版本"
	@echo "  make linux-arm64     - 编译 Linux ARM64 版本"
	@echo "  make darwin-amd64    - 编译 macOS AMD64 版本"
	@echo "  make darwin-arm64    - 编译 macOS ARM64 (Apple Silicon) 版本"
	@echo "  make windows-amd64   - 编译 Windows AMD64 版本"
	@echo "  make cross-compile   - 编译所有平台版本"
	@echo "  make clean           - 清理编译产物"
	@echo "  make test            - 运行测试"
	@echo "  make tidy            - 整理依赖"
	@echo ""
	@echo "$(YELLOW)编译产物位置:$(NC) $(BUILD_DIR)/"

# 编译当前平台版本
build:
	@echo "$(GREEN)正在编译当前平台版本...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "$(GREEN)✓ 编译完成: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

# Linux AMD64
linux-amd64:
	@echo "$(GREEN)正在编译 Linux AMD64 版本...$(NC)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	@echo "$(GREEN)✓ 编译完成: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64$(NC)"

# Linux ARM64
linux-arm64:
	@echo "$(GREEN)正在编译 Linux ARM64 版本...$(NC)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)
	@echo "$(GREEN)✓ 编译完成: $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64$(NC)"

# macOS AMD64 (Intel)
darwin-amd64:
	@echo "$(GREEN)正在编译 macOS AMD64 (Intel) 版本...$(NC)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	@echo "$(GREEN)✓ 编译完成: $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64$(NC)"

# macOS ARM64 (Apple Silicon)
darwin-arm64:
	@echo "$(GREEN)正在编译 macOS ARM64 (Apple Silicon) 版本...$(NC)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)
	@echo "$(GREEN)✓ 编译完成: $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64$(NC)"

# Windows AMD64
windows-amd64:
	@echo "$(GREEN)正在编译 Windows AMD64 版本...$(NC)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)
	@echo "$(GREEN)✓ 编译完成: $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe$(NC)"

# 交叉编译所有平台
cross-compile: linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64
	@echo ""
	@echo "$(GREEN)========================================$(NC)"
	@echo "$(GREEN)所有平台编译完成！$(NC)"
	@echo "$(GREEN)========================================$(NC)"
	@ls -lh $(BUILD_DIR)

# 清理编译产物
clean:
	@echo "$(YELLOW)正在清理编译产物...$(NC)"
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME)
	rm -f speed
	@echo "$(GREEN)✓ 清理完成$(NC)"

# 运行测试
test:
	@echo "$(GREEN)正在运行测试...$(NC)"
	$(GOTEST) -v ./...

# 整理依赖
tidy:
	@echo "$(GREEN)正在整理依赖...$(NC)"
	$(GOMOD) tidy
	@echo "$(GREEN)✓ 依赖整理完成$(NC)"

# 安装到本地
install: build
	@echo "$(GREEN)正在安装到 GOPATH...$(NC)"
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/
	@echo "$(GREEN)✓ 安装完成$(NC)"

# 快速编译并运行（demo模式）
run-demo: build
	@echo "$(GREEN)运行 demo 模式...$(NC)"
	./$(BUILD_DIR)/$(BINARY_NAME) -demo

# 快速编译并运行（API模式）
run-api: build
	@echo "$(GREEN)运行 API 模式...$(NC)"
	./$(BUILD_DIR)/$(BINARY_NAME)

