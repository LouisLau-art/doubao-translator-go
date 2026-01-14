# ==============================
#  Go Translator  Makefile
# ==============================
# 变量
BINARY_NAME   := translator
GO            := go
GOFLAGS       := -ldflags="-s -w"
PORT          := 5000
CARGO_TARGET  := linux/amd64   # 可根据需要改成你机器的架构

# ---------- 颜色 ----------
GREEN  := \033[32m
YELLOW := \033[33m
RED    := \033[31m
RESET  := \033[0m

# ---------- 帮助 ----------
.PHONY: help
help:
	@echo -e "$(GREEN)=== Go Translator Makefile ===$(RESET)"
	@echo "  make dev          - 启动开发服务器 (go run .)"
	@echo "  make build        - 编译普通二进制"
	@echo "  make build-prod   - 编译生产版（-ldflags -s -w）"
	@echo "  make install      - 安装依赖 (go mod tidy)"
	@echo "  make setup        - 安装依赖 + 下载前端 libs"
	@echo "  make clean        - 删除二进制"
	@echo "  make fmt          - go fmt ./..."
	@echo "  make lint         - 运行 golangci-lint (若未装需先 sudo pacman -S golangci-lint)"
	@echo "  make prod         - build-prod + 拷贝到 /opt (需要 sudo)"
	@echo "  make serve        - 启动已编译的二进制 (需要先 make prod)"
	@echo "  make status       - 查看 systemd 状态"
	@echo "  make logs         - 查看 journal 日志"
	@echo "  make uninstall    - 删除 systemd 服务 & /opt/translator"
	@echo ""

# ---------- 依赖 ----------
.PHONY: install
install:
	@echo "$(YELLOW)Downloading Go modules...$(RESET)"
	$(GO) mod tidy

# ---------- 前端资源 ----------
.PHONY: libs
libs:
	@echo "$(YELLOW)Downloading frontend libs...$(RESET)"
	@mkdir -p static/libs
	@wget -q -O static/libs/petite-vue.js https://unpkg.com/petite-vue@0.4.1/dist/petite-vue.iife.js
	@wget -q -O static/libs/marked.min.js https://cdn.jsdelivr.net/npm/marked/marked.min.js
	@mkdir -p static/libs/mathjax
	@wget -q -O static/libs/mathjax/tex-mml-chtml.js https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js
	@echo "$(GREEN)Frontend libs downloaded.$(RESET)"

# ---------- 开发 ----------
.PHONY: dev
dev: install
	@echo "$(GREEN)Starting development server...$(RESET)"
	$(GO) run .

# ---------- 构建 ----------
.PHONY: build
build: install
	@echo "$(GREEN)Building binary...$(RESET)"
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME)

.PHONY: build-prod
build-prod: install
	@echo "$(GREEN)Building production binary...$(RESET)"
	CGO_ENABLED=0 GOOS=linux $(GO) build $(GOFLAGS) -o $(BINARY_NAME)

# ---------- 运行 ----------
.PHONY: serve
serve: build-prod
	@echo "$(GREEN)Running translator...$(RESET)"
	./$(BINARY_NAME)

# ---------- 清理 ----------
.PHONY: clean
clean:
	@echo "$(YELLOW)Cleaning binary...$(RESET)"
	rm -f $(BINARY_NAME)

# ---------- 系统服务 ----------
.PHONY: install-service
install-service:
	@echo "$(YELLOW)Installing systemd service...$(RESET)"
	@sudo mkdir -p /opt/translator
	@sudo cp $(BINARY_NAME) /opt/translator/
	@sudo cp .env /opt/translator/
	@sudo cp $(BINARY_NAME) /opt/translator/
	@sudo chmod +x /opt/translator/$(BINARY_NAME)
	@sudo tee /etc/systemd/system/translator.service > /dev/null <<'EOF'
[Unit]
Description=Doubao Translator Service
After=network.target

[Service]
Type=simple
User=$(shell whoami)
WorkingDirectory=/opt/translator
EnvironmentFile=/opt/translator/.env
ExecStart=/opt/translator/$(BINARY_NAME)
Restart=on-failure
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF
	@sudo systemctl daemon-reload
	@sudo systemctl enable translator.service
	@echo "$(GREEN)Service installed. Enable with: sudo systemctl enable translator.service$(RESET)"
	@echo "Start it with: sudo systemctl start translator.service"

.PHONY: uninstall
uninstall:
	@echo "$(RED)Stopping and removing service...$(RESET)"
	@sudo systemctl stop translator.service || true
	@sudo systemctl disable translator.service || true
	@sudo rm -f /etc/systemd/system/translator.service
	@sudo rm -rf /opt/translator
	@sudo systemctl daemon-reload
	@echo "$(GREEN)Service removed.$(RESET)"

# ---------- 代码质量 ----------
.PHONY: fmt
fmt:
	@echo "$(YELLOW)Formatting code...$(RESET)"
	$(GO) fmt ./...

.PHONY: lint
lint:
	@if command -v golangci-lint >/dev/null; then \
		echo "$(YELLOW)Running golangci-lint...$(RESET)"; \
		golangci-lint run; \
	else \
		echo "$(RED)golangci-lint not installed. Install with: sudo pacman -S golangci-lint$(RESET)"; \
	fi

# ---------- 打包 ----------
.PHONY: prod
prod: build-prod
	@echo "$(GREEN)Production binary built: ./$(BINARY_NAME)$(RESET)"
	@if command -v upx >/dev/null; then \
		echo "$(YELLOW)Compressing with UPX...$(RESET)"; \
		upx --best $(BINARY_NAME); \
	else \
		echo "$(YELLOW)UPX not installed, skipping compression.$(RESET)"; \
	fi