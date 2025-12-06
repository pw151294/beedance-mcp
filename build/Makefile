# 项目根目录 (相对于 Makefile 的位置)
ROOT_DIR := ../

# Go 主程序入口
MAIN_GO := $(ROOT_DIR)/main.go

# 构建产物
# 二进制文件将输出到 build 目录下
BIN_NAME := beedance_mcp
BIN_PATH := ./$(BIN_NAME)

# 打包产物
# 压缩包将输出到 build 目录下
TAR_NAME := beedance_mcp.tar.gz
TAR_PATH := ./$(TAR_NAME)

# 需要打包的配置文件
ENV_FILE := $(ROOT_DIR)/.env
CONFIG_FILE := $(ROOT_DIR)/configs/config.toml

# 默认目标
.PHONY: all
all: package

# 1. 构建 Go 程序为 linux/amd64 二进制文件
.PHONY: build
build:
	@echo "==> Building binary for linux/amd64..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BIN_PATH) $(MAIN_GO)
	@echo "==> Build complete: $(BIN_PATH)"

# 2. 打包二进制文件和配置文件
.PHONY: package
package: build
	@echo "==> Packaging to $(TAR_PATH)..."
	@# 先切换到项目根目录，然后执行打包，所有路径都相对于根目录
	cd $(ROOT_DIR) && tar -czf build/$(TAR_NAME) \
		.env \
		configs/config.toml \
		build/$(BIN_NAME)
	@echo "==> Package complete: $(TAR_PATH)"

# 清理构建产物
.PHONY: clean
clean:
	@echo "==> Cleaning build artifacts..."
	@rm -f $(BIN_PATH) $(TAR_PATH)
	@echo "==> Clean complete."
