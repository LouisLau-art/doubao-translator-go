
# Doubao Translator - Go Edition

![Go](https://img.shields.io/badge/go-v1.25%2B-00ADD8.svg)
![Gin](https://img.shields.io/badge/gin-v1.11.0-009688.svg)
![Vue Petite](https://img.shields.io/badge/vue%20petite-v0.4.1-4FC08D.svg)
![License](https://img.shields.io/badge/license-MIT-blue.svg)

基于 **火山引擎豆包翻译 API** 的极简、高性能网页翻译器，使用 Go + Gin 后端和 Vue Petite 前端。

## ✨ 功能亮点

- 🌙 **深色主题**：默认暗色界面，长时间阅读不刺眼
- 🔄 **实时自动翻译**：输入停顿 0.5 秒自动触发，可开关控制
- 📝 **Markdown/LaTeX 支持**：内置 MathJax 3，支持 `$...$`、`$$...$$` 公式渲染
- 🌍 **多语言互译**：支持 14 种语言互译，源语言自动检测
- 📋 **增强复制体验**：复制后显示 ✔ 反馈，支持纯文本复制
- 🧹 **快速清空**：输入输出一键清空，方便重写
- 📜 **本地历史记录**：自动保存 50 条翻译记录到 localStorage
- 📐 **字体大小调节**：滑块调节输入/输出字体大小（12px–26px）
- 📄 **长文档翻译**：自动拆分超长文本，保留格式结构

## 🛠 技术栈

### 后端
- **语言**: Go 1.25+
- **框架**: Gin 1.11.0 (高性能 Web 框架)
- **中间件**: gin-contrib/cors (CORS 支持)
- **配置管理**: godotenv (环境变量加载)
- **速率限制**: golang.org/x/time/rate (令牌桶算法)

### 前端
- **框架**: Vue Petite 0.4.1 (轻量级 Vue 替代，仅 6KB)
- **Markdown 渲染**: Marked.js
- **数学公式渲染**: MathJax 3
- **样式**: 原生 CSS (响应式设计 + 暗色主题)

### 构建工具
- **自动化**: Makefile (开发、构建、部署)
- **压缩**: UPX (可选二进制压缩)

## 🚀 快速开始

### 环境要求
- Go 1.25 或更高版本
- 火山引擎 ARK API 密钥 ([申请地址](https://www.volcengine.com/))

### 安装步骤

1. **克隆项目**
   ```bash
   git clone https://github.com/LouisLau-art/doubao-translator-go.git
   cd doubao-translator-go
   ```

2. **完整安装**
   ```bash
   make setup
   ```
   此命令会：
   - 下载 Go 模块依赖
   - 下载前端库 (Vue Petite, Marked.js, MathJax)

3. **配置 API 密钥**
   ```bash
   cp .env.example .env
   # 编辑 .env 文件，填入真实的 API 密钥
   ```

4. **启动开发服务器**
   ```bash
   make dev
   ```

5. **访问应用**
   在浏览器中打开 `http://localhost:5000`

### 生产部署

```bash
# 编译生产版本
make build-prod

# 运行生产版本
make serve

# 或直接运行编译好的二进制
./translator
```

## ⚙️ 配置说明

### 环境变量 (.env)
```env
# 必需配置
ARK_API_KEY=your_ark_api_key_here        # 火山引擎 ARK API 密钥
ARK_API_URL=https://ark.cn-beijing.volces.com/api/v3/responses  # API 端点

# 可选配置
PORT=5000                                # 服务器端口 (默认: 5000)
GIN_MODE=release                         # Gin 运行模式: debug/release
CACHE_TTL=3600                          # 缓存有效期 (秒)
CACHE_MAX_SIZE=1000                     # 最大缓存条目数
MAX_TEXT_LENGTH=5000                    # 单次请求最大文本长度
RATE_LIMIT_RPM=30                       # API 速率限制 (每分钟请求数)
```

### API 端点
- `GET /` - 前端页面入口
- `GET /api/languages` - 获取支持的语言列表
- `POST /api/translate` - 翻译请求 (JSON)
- `GET /api/health` - 健康检查
- `GET /static/*` - 静态资源
- `GET /libs/*` - 前端库文件

## 📖 使用指南

1. **输入文本**: 在左侧输入框输入或粘贴要翻译的文本
2. **选择语言**:
   - 源语言: 可选 "自动检测" 或指定语言
   - 目标语言: 从下拉列表中选择
   - 点击 "⇄" 按钮可快速交换源语言和目标语言
3. **自动翻译**: 默认开启，输入停顿 0.5 秒后自动翻译
4. **字体调节**: 使用底部滑块调节输入/输出区域的字体大小
5. **复制结果**: 点击 "📋" 按钮复制翻译结果
6. **查看历史**: 展开历史记录面板查看之前的翻译

### 特殊功能
- **数学公式**: 支持 LaTeX 公式，使用 `$...$` (行内) 或 `$$...$$` (独立行)
- **Markdown 渲染**: 翻译结果会自动渲染 Markdown 格式
- **长文档**: 系统会自动拆分长文本，翻译后重新组合

## 🛠 开发命令

### 常用命令
```bash
make help              # 查看所有可用命令
make setup             # 完整安装 (依赖 + 前端库)
make dev               # 启动开发服务器
make build-prod        # 编译生产版本
make serve             # 构建并运行生产版本
make clean             # 清理构建文件
```

### 代码质量
```bash
make fmt               # 格式化代码 (go fmt)
make lint              # 代码检查 (golangci-lint)
make vet               # 静态分析 (go vet)
make test              # 运行测试
make test-cover        # 运行测试并生成覆盖率报告
```

### 系统服务 (Linux)
```bash
make install-service   # 安装为 systemd 服务
make status            # 查看服务状态
make logs              # 查看服务日志
make uninstall         # 卸载服务
```

## 📁 项目结构

```
go-translator/
├── main.go                      # Go 后端主入口
├── go.mod                       # Go 模块依赖
├── go.sum                       # 依赖校验和
├── Makefile                     # 构建自动化
├── .env.example                 # 环境变量示例
├── CLAUDE.md                    # Claude Code 开发指南
├── download-libs.sh             # 前端库下载脚本
├── static/                      # 前端资源
│   ├── index.html              # 主页面
│   ├── app.js                  # Vue Petite 应用逻辑
│   ├── style.css               # 完整样式 (暗色主题)
│   └── libs/                   # 前端依赖库
│       ├── petite-vue.js       # Vue Petite 框架
│       ├── marked.min.js       # Markdown 解析器
│       └── mathjax/            # MathJax 3
│           └── tex-mml-chtml.js # LaTeX 渲染
└── README.md                    # 项目文档
```

## 🔧 技术特性

### 后端特性
- **智能缓存**: MD5 哈希缓存键，可配置 TTL
- **速率限制**: 令牌桶算法，防止 API 滥用
- **文本分块**: 自动拆分超长文本 (800 字符/块)
- **错误处理**: 详细的错误响应和日志
- **健康检查**: 独立的健康检查端点

### 前端特性
- **响应式设计**: 适配桌面和移动设备
- **实时更新**: 防抖处理，避免频繁请求
- **本地存储**: 配置和历史记录持久化
- **公式渲染**: 动态重新渲染数学公式
- **无障碍**: 键盘导航和屏幕阅读器支持

## ❓ 常见问题

### API 相关问题
**Q: 提示 "ARK_API_KEY not set"**
A: 请确保已创建 `.env` 文件并正确配置 `ARK_API_KEY`

**Q: API 返回 401 错误**
A: API 密钥无效或过期，请检查火山引擎控制台

**Q: API 返回 429 错误**
A: 达到速率限制，请稍后重试或调整 `RATE_LIMIT_RPM` 配置

### 功能相关问题
**Q: 长文档翻译格式混乱**
A: 系统会按段落拆分，确保输入是标准 Markdown 格式

**Q: 数学公式不显示**
A: 检查公式语法是否正确，或刷新页面重新加载 MathJax

**Q: 历史记录丢失**
A: 浏览器清理缓存可能导致 localStorage 数据丢失

### 部署相关问题
**Q: 端口 5000 被占用**
A: 修改 `.env` 中的 `PORT` 配置为其他可用端口

**Q: 生产版本无法启动**
A: 确保已运行 `make build-prod` 编译生产版本

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'feat: add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

### 开发流程
```bash
# 1. 克隆并设置项目
git clone <your-fork-url>
cd doubao-translator-go
make setup

# 2. 创建并切换分支
git checkout -b feature/your-feature

# 3. 开发并测试
make dev
make test

# 4. 提交更改
git add .
git commit -m "feat: description of your feature"

# 5. 推送到远程
git push origin feature/your-feature
```

## 🔒 安全提示

- **切勿提交敏感信息**: 不要将 `.env` 文件提交到版本库
- **API 密钥保护**: 使用环境变量或密钥管理服务存储 API 密钥
- **访问控制**: 生产环境建议配置防火墙和访问控制
- **日志安全**: 避免在日志中记录敏感信息

## 📄 许可证

本项目基于 [MIT License](LICENSE) 开源。

## 🙏 致谢

- 火山引擎 ARK 团队提供的豆包翻译 API
- Gin 框架社区
- Vue Petite 项目
- MathJax 和 Marked.js 项目

## 👤 作者

- **Louis Lau** - [GitHub](https://github.com/LouisLau-art) - louis.shawn@qq.com

## ⭐ 支持

如果这个项目对你有帮助，请给个 Star ⭐！

---

**Made with ❤️ and ☕** 