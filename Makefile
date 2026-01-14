# ==============================
#  Go + Vue Petite Doubao Translator
# ==============================
# 变量定义
BINARY_NAME   := translator
GO            := go
GOFLAGS       := -ldflags="-s -w"
PORT          := 5000
GO_TEST_FLAGS := -v -race -cover
DIST_DIR      := dist
COVERAGE_FILE := coverage.out

# ---------- 颜色定义 ----------
GREEN  := \033[32m
YELLOW := \033[33m
BLUE   := \033[34m
RED    := \033[31m
RESET  := \033[0m

# ---------- 帮助 ----------
.PHONY: help
help:
	@echo -e "$(GREEN)=== Go + Vue Petite Doubao Translator ===$(RESET)"
	@echo ""
	@echo -e "$(BLUE)开发命令:$(RESET)"
	@echo "  make dev          - 启动开发服务器 (go run .)"
	@echo "  make setup        - 完整安装（依赖 + 前端库）"
	@echo "  make clean        - 清理构建文件"
	@echo ""
	@echo -e "$(BLUE)构建命令:$(RESET)"
	@echo "  make build        - 编译开发版本"
	@echo "  make build-prod   - 编译生产版本（优化）"
	@echo "  make build-all    - 编译多平台版本"
	@echo "  make serve        - 构建并运行生产版本"
	@echo ""
	@echo -e "$(BLUE)代码质量:$(RESET)"
	@echo "  make fmt          - 格式化代码 (go fmt)"
	@echo "  make lint         - 代码检查 (golangci-lint)"
	@echo "  make vet          - 静态分析 (go vet)"
	@echo "  make test         - 运行测试"
	@echo "  make test-cover   - 运行测试并生成覆盖率报告"
	@echo ""
	@echo -e "$(BLUE)系统服务 (Linux):$(RESET)"
	@echo "  make install-service - 安装为 systemd 服务"
	@echo "  make status          - 查看服务状态"
	@echo "  make logs           - 查看服务日志"
	@echo "  make uninstall      - 卸载服务"
	@echo ""
	@echo -e "$(BLUE)其他命令:$(RESET)"
	@echo "  make generate      - 运行 go generate"
	@echo "  make deps          - 更新依赖"
	@echo "  make version       - 显示版本信息"
	@echo ""

# ---------- 依赖管理 ----------
.PHONY: deps
deps:
	@echo -e "$(YELLOW)更新 Go 依赖...$(RESET)"
	$(GO) mod tidy
	$(GO) mod download

.PHONY: install
install: deps
	@echo -e "$(GREEN)依赖安装完成$(RESET)"

# ---------- 前端资源 ----------
.PHONY: libs
libs:
	@echo -e "$(YELLOW)下载前端依赖库...$(RESET)"
	@mkdir -p static/libs
	@if command -v curl >/dev/null; then \
		curl -L -s -o static/libs/petite-vue.js https://unpkg.com/petite-vue@0.4.1/dist/petite-vue.iife.js; \
		curl -L -s -o static/libs/marked.min.js https://cdn.jsdelivr.net/npm/marked/marked.min.js; \
		mkdir -p static/libs/mathjax; \
		curl -L -s -o static/libs/mathjax/tex-mml-chtml.js https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js; \
	elif command -v wget >/dev/null; then \
		wget -q -O static/libs/petite-vue.js https://unpkg.com/petite-vue@0.4.1/dist/petite-vue.iife.js; \
		wget -q -O static/libs/marked.min.js https://cdn.jsdelivr.net/npm/marked/marked.min.js; \
		mkdir -p static/libs/mathjax; \
		wget -q -O static/libs/mathjax/tex-mml-chtml.js https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js; \
	else \
		echo -e "$(RED)错误: 需要 curl 或 wget 命令$(RESET)"; \
		exit 1; \
	fi
	@echo -e "$(GREEN)前端依赖库下载完成$(RESET)"

# ---------- 开发 ----------
.PHONY: setup
setup: deps libs
	@echo -e "$(GREEN)项目设置完成，可以运行 make dev 启动开发服务器$(RESET)"

.PHONY: dev
dev: setup
	@echo -e "$(GREEN)启动开发服务器...$(RESET)"
	$(GO) run .

# ---------- 构建 ----------
.PHONY: build
build: deps
	@echo -e "$(GREEN)编译开发版本...$(RESET)"
	$(GO) build -o $(BINARY_NAME)

.PHONY: build-prod
build-prod: deps
	@echo -e "$(GREEN)编译生产版本...$(RESET)"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -o $(BINARY_NAME)

.PHONY: build-all
build-all: deps
	@echo -e "$(YELLOW)编译多平台版本...$(RESET)"
	@mkdir -p $(DIST_DIR)
	@for platform in "linux/amd64" "linux/arm64" "darwin/amd64" "darwin/arm64" "windows/amd64"; do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		output="$(DIST_DIR)/$(BINARY_NAME)-$$GOOS-$$GOARCH"; \
		if [ "$$GOOS" = "windows" ]; then output="$$output.exe"; fi; \
		echo "编译 $$GOOS/$$GOARCH -> $$output"; \
		CGO_ENABLED=0 GOOS=$$GOOS GOARCH=$$GOARCH $(GO) build $(GOFLAGS) -o $$output; \
	done
	@echo -e "$(GREEN)多平台编译完成，文件在 $(DIST_DIR)/ 目录$(RESET)"

# ---------- 运行 ----------
.PHONY: serve
serve: build-prod
	@echo -e "$(GREEN)运行生产版本...$(RESET)"
	./$(BINARY_NAME)

# ---------- 代码质量 ----------
.PHONY: fmt
fmt:
	@echo -e "$(YELLOW)格式化代码...$(RESET)"
	$(GO) fmt ./...

.PHONY: vet
vet:
	@echo -e "$(YELLOW)静态代码分析...$(RESET)"
	$(GO) vet ./...

.PHONY: lint
lint:
	@echo -e "$(YELLOW)运行代码检查...$(RESET)"
	@if command -v golangci-lint >/dev/null; then \
		golangci-lint run --timeout 5m; \
	else \
		echo -e "$(RED)golangci-lint 未安装，请先安装:$(RESET)"; \
		echo "  macOS: brew install golangci-lint"; \
		echo "  Linux: 使用包管理器安装或从 https://golangci-lint.run/usage/install/ 下载"; \
		echo "  Windows: scoop install golangci-lint 或从官网下载"; \
		exit 1; \
	fi

.PHONY: test
test: deps
	@echo -e "$(YELLOW)运行测试...$(RESET)"
	$(GO) test $(GO_TEST_FLAGS) ./...

.PHONY: test-cover
test-cover: deps
	@echo -e "$(YELLOW)运行测试并生成覆盖率报告...$(RESET)"
	$(GO) test $(GO_TEST_FLAGS) -coverprofile=$(COVERAGE_FILE) ./...
	$(GO) tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo -e "$(GREEN)覆盖率报告生成完成: coverage.html$(RESET)"

# ---------- 清理 ----------
.PHONY: clean
clean:
	@echo -e "$(YELLOW)清理构建文件...$(RESET)"
	rm -f $(BINARY_NAME) $(COVERAGE_FILE) coverage.html
	rm -rf $(DIST_DIR)

# ---------- 系统服务 ----------
.PHONY: install-service
install-service: build-prod
	@echo -e "$(YELLOW)安装 systemd 服务...$(RESET)"
	@if [ ! -f .env ]; then \
		echo -e "$(RED)错误: 请先创建 .env 文件并配置 API 密钥$(RESET)"; \
		exit 1; \
	fi
	@if [ ! -f $(BINARY_NAME) ]; then \
		echo -e "$(RED)错误: 请先运行 make build-prod 编译生产版本$(RESET)"; \
		exit 1; \
	fi
	@sudo mkdir -p /opt/translator
	@sudo cp $(BINARY_NAME) /opt/translator/
	@sudo cp .env /opt/translator/
	@sudo chmod 600 /opt/translator/.env
	@sudo chmod +x /opt/translator/$(BINARY_NAME)
	@sudo tee /etc/systemd/system/translator.service > /dev/null <<'EOF'
[Unit]
Description=Doubao Translator Service
Documentation=https://github.com/LouisLau-art/doubao-translator-go
After=network.target
Wants=network.target

[Service]
Type=simple
User=nobody
Group=nogroup
WorkingDirectory=/opt/translator
EnvironmentFile=/opt/translator/.env
ExecStart=/opt/translator/$(BINARY_NAME)
Restart=on-failure
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=translator
ProtectSystem=strict
ReadWritePaths=/opt/translator
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF
	@sudo systemctl daemon-reload
	@sudo systemctl enable translator.service
	@echo -e "$(GREEN)服务安装完成$(RESET)"
	@echo -e "$(YELLOW)启动服务:$(RESET) sudo systemctl start translator.service"
	@echo -e "$(YELLOW)查看状态:$(RESET) sudo systemctl status translator.service"
	@echo -e "$(YELLOW)查看日志:$(RESET) sudo journalctl -u translator.service -f"

.PHONY: status
status:
	@echo -e "$(YELLOW)服务状态:$(RESET)"
	@sudo systemctl status translator.service 2>/dev/null || echo "服务未安装"

.PHONY: logs
logs:
	@echo -e "$(YELLOW)服务日志 (最近 50 行):$(RESET)"
	@sudo journalctl -u translator.service -n 50 --no-pager 2>/dev/null || echo "服务未安装"

.PHONY: uninstall
uninstall:
	@echo -e "$(RED)停止并卸载服务...$(RESET)"
	@sudo systemctl stop translator.service 2>/dev/null || true
	@sudo systemctl disable translator.service 2>/dev/null || true
	@sudo rm -f /etc/systemd/system/translator.service
	@sudo rm -rf /opt/translator
	@sudo systemctl daemon-reload
	@echo -e "$(GREEN)服务卸载完成$(RESET)"

# ---------- 其他命令 ----------
.PHONY: generate
generate:
	@echo -e "$(YELLOW)运行代码生成...$(RESET)"
	$(GO) generate ./...

.PHONY: version
version:
	@echo -e "$(BLUE)项目信息:$(RESET)"
	@echo "名称: Doubao Translator - Go Edition"
	@echo "版本: 1.0.0"
	@echo "Go 版本: $$($(GO) version | cut -d' ' -f3)"
	@echo "架构: $$(uname -m)"
	@echo "操作系统: $$(uname -s)"

# ---------- 打包 ----------
.PHONY: prod
prod: build-prod
	@echo -e "$(GREEN)生产版本编译完成: ./$(BINARY_NAME)$(RESET)"
	@if command -v upx >/dev/null; then \
		echo -e "$(YELLOW)使用 UPX 压缩...$(RESET)"; \
		upx --best $(BINARY_NAME); \
		echo -e "$(GREEN)压缩完成，文件大小: $$(du -h $(BINARY_NAME) | cut -f1)$(RESET)"; \
	else \
		echo -e "$(YELLOW)UPX 未安装，跳过压缩$(RESET)"; \
		echo "安装 UPX: brew install upx (macOS) 或 sudo apt install upx (Ubuntu)"; \
	fi