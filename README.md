
# Doubao Translator - Go Edition

![Go](https://img.shields.io/badge/go-v1.21%2B-00ADD8.svg)
![Gin](https://img.shields.io/badge/gin-v1.9-009688.svg)
![Vue Petite](https://img.shields.io/badge/vue%20petite-v0.4.1-4FC08D.svg)
![License](https://img.shields.io/badge/license-MIT-blue.svg)

基于 **火山引擎豆包翻译 API** 的极简、高性能网页翻译器。

## 功能亮点 ✨

- 🌙 **深色主题**：自带暗色界面，长时间阅读不刺眼。
- 🔄 **实时自动翻译**：输入停顿 0.5 秒自动触发，可随时切换开关。
- 📝 **Markdown/LaTeX 支持**：内置 MathJax，对 `$...$`、`$$...$$` 公式友好。
- 🌍 **多语言互译**：内置 28 种语种选项，并支持源语言自动检测。
- 📋 **增强复制体验**：复制后按钮会短暂显示 ✔，反馈更及时。
- 🧹 **快速清空**：输入输出一键清空，重写更高效。
- 📜 **本地历史记录**：自动保存所有翻译记录，可随时展开重用。
- 📐 **字体滑杆**：通过拖动滑杆精细调节输入/输出面板字体大小（12px–26px）。
- 📄 **长文档翻译**：支持任意长度的 Markdown 文档翻译，自动拆分合并，保留格式。

## 技术栈 🛠

- **后端**：Go + Gin
- **前端**：Vue Petite + 原生 HTML / CSS / JavaScript（拆分静态资源）
- **API**：ARK Doubao Translation
- **渲染**：Python‑Markdown + MathJax（原项目已替换为前端渲染）

## 快速开始 🚀

### 环境要求

- Go 1.21+
- 火山引擎 ARK API 密钥（[申请入口](https://www.volcengine.com/)）

### 安装步骤

1. **克隆项目**
   ```bash
   git clone https://github.com/LouisLau-art/doubao-translator-go.git
   cd doubao-translator-go
   ```

2. **安装依赖**
   ```bash
   make setup
   ```

3. **配置 API 密钥**
   ```bash
   cp .env.example .env
   # 编辑 .env，填入真实的 API Key
   ```
   ```env
   ARK_API_KEY=your_actual_api_key
   ARK_API_URL=https://ark.cn-beijing.volces.com/api/v3/responses
   PORT=5000
   ```

4. **启动服务**
   ```bash
   make dev
   ```
   或生产模式：
   ```bash
   make build-prod
   ./translator
   ```

5. **访问页面**
   在浏览器打开 `http://127.0.0.1:5000`（或 `http://localhost:5000`）。

## 配置说明 ⚙️

- 应用启动时会自动加载同目录下的 `translator.env`。
- `python-dotenv` 负责读取环境变量，避免将密钥硬编码进代码。
- 如需部署到其他路径，请确保 `translator.env` 与 `main.go`（或 `web_translator.py`）在同一目录。

## 使用指南 📖

1. 在左侧输入框键入或粘贴文本。
2. 依据需求选择源语言与目标语言，或使用 “⇄” 按钮互换。
3. 保持 “自动翻译” 开启即可实时出结果，也可手动控制。
4. 拖动状态栏的字体滑杆，调整输入与输出区域的字号。
5. 点击 “📋” 按钮复制翻译结果的纯文本。
6. 若文本包含公式，MathJax 会在翻译完成后自动重新渲染。
7. **长文档翻译**：直接粘贴任意长度的 Markdown 文档，系统会自动拆分、翻译并合并结果，保留原格式。

## API 说明 📚

- 当前实现内置 28 种语言（详见 `app.py` 中的 `LANGUAGE_MAP`）。
- 源语言支持自动检测，目标语言需从列表中选择。
- 翻译请求默认超时 30 秒，超时会提示 “网络错误”。
- 对常见状态码（401/429 等）提供了更明确的错误提示。
- **长文档处理**：自动将超过 1000 个 token 的文档拆分为多个块，翻译后合并，保留 Markdown 格式。

## 常见问题 ❓

- **提示 “请创建 translator.env 文件并设置 ARK_API_KEY”**：检查文件是否存在且密钥填写正确。
- **返回错误信息**：留意接口响应，确认 Key 权限与余额。
- **网络错误**：检查本地网络环境，终端日志可帮助定位问题。
- **长文档翻译问题**：确保输入的是有效的 Markdown 格式，代码块会自动完整保留。

## 目录结构 🗂

```
doubao-translator-go/
├── app.py                   # Flask 应用主入口（原项目）
├── templates/
│   └── index.html           # 页面模板
├── static/
│   ├── style.css            # 全部样式
│   └── script.js            # 前端交互脚本
├── translator.env           # 私有环境变量（需手动创建）
├── translator.example.env   # 环境变量示例
├── requirements.txt         # 项目依赖
└── README.md
```

## 贡献指南 🤝

欢迎提交 Issue / PR：

1. Fork 本仓库。
2. 新建分支 `git checkout -b feature/your-feature`。
3. 提交修改 `git commit -m "feat: add your feature"`。
4. 推送分支 `git push origin feature/your-feature`。
5. 发起 Pull Request 并描述变更内容。

## 安全提示 🔒

- 切勿将真实 API Key 提交到版本库。
- 保持 `.gitignore` 中对 `translator.env` 等敏感文件的忽略规则。
- 部署到服务器时，推荐使用系统级环境变量或密钥管理服务。

## 许可证 📄

本项目遵循 [MIT License](LICENSE)。

## 致谢 🙏

- 火山引擎 ARK 团队提供的豆包翻译 API。
- Flask 社区提供的优秀 Web 框架。
- Python‑Markdown 与 MathJax 项目的支持。

## 作者 👤

- GitHub：[LouisLau-art](https://github.com/LouisLau-art)
- Email：louis.shawn@qq.com

## 支持 ⭐

如果这个项目对你有帮助，请点亮 Star ⭐，或分享给更多开发者。

Made with ❤️ and ☕

--- 