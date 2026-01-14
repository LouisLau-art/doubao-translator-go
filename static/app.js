// Vue Petite 应用
function TranslatorApp() {
    return {
        // 状态
        inputText: '',
        outputText: '',
        sourceLang: '',
        targetLang: 'zh',
        languages: {},
        autoTranslate: true,
        fontSize: 16,
        loading: false,
        error: '',
        cached: false,
        copyBtnText: '📋 复制',
        history: [],
        debounceTimer: null,

        // 初始化
        async init() {
            // 加载语言列表
            await this.loadLanguages();
            
            // 加载本地存储
            this.loadFromStorage();
            
            // 设置 MathJax 配置
            window.MathJax = {
                tex: {
                    inlineMath: [['$', '$']],
                    displayMath: [['$$', '$$']],
                },
                startup: {
                    pageReady: () => {
                        return MathJax.startup.defaultPageReady();
                    }
                }
            };
        },

        // 加载语言列表
        async loadLanguages() {
            try {
                const response = await fetch('/api/languages');
                const data = await response.json();
                if (data.success) {
                    this.languages = data.languages;
                }
            } catch (error) {
                console.error('Failed to load languages:', error);
                // 使用默认语言列表
                this.languages = {
                    'zh': '中文（简体）',
                    'en': '英语',
                    'ja': '日语',
                    'ko': '韩语',
                };
            }
        },

        // 处理输入
        handleInput() {
            this.error = '';
            
            if (this.autoTranslate) {
                // 防抖处理
                clearTimeout(this.debounceTimer);
                this.debounceTimer = setTimeout(() => {
                    if (this.inputText.trim()) {
                        this.translateNow();
                    }
                }, 500);
            }
        },

        // 立即翻译
        async translateNow() {
            if (!this.inputText.trim() || this.loading) return;
            
            this.loading = true;
            this.error = '';
            this.cached = false;

            try {
                const response = await fetch('/api/translate', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        text: this.inputText,
                        source: this.sourceLang,
                        target: this.targetLang,
                    }),
                });

                const data = await response.json();
                
                if (data.success) {
                    this.outputText = data.text;
                    this.cached = data.cached || false;
                    
                    // 保存到历史
                    this.saveToHistory();
                    
                    // 触发 MathJax 重新渲染
                    this.$nextTick(() => {
                        if (window.MathJax?.typesetPromise) {
                            window.MathJax.typesetPromise();
                        }
                    });
                } else {
                    this.error = data.error || '翻译失败';
                }
            } catch (error) {
                this.error = '网络错误，请检查连接';
                console.error('Translation error:', error);
            } finally {
                this.loading = false;
            }
        },

        // 交换语言
        swapLanguages() {
            if (this.sourceLang && this.targetLang) {
                [this.sourceLang, this.targetLang] = [this.targetLang, this.sourceLang];
                if (this.inputText && this.outputText) {
                    [this.inputText, this.outputText] = [this.outputText, this.inputText];
                }
            }
        },

        // 复制结果
        async copyResult() {
            if (!this.outputText) return;
            
            try {
                await navigator.clipboard.writeText(this.outputText);
                this.copyBtnText = '✔ 已复制';
                setTimeout(() => {
                    this.copyBtnText = '📋 复制';
                }, 2000);
            } catch (error) {
                this.error = '复制失败';
            }
        },

        // 粘贴文本
        async pasteText() {
            try {
                const text = await navigator.clipboard.readText();
                this.inputText = text;
                this.handleInput();
            } catch (error) {
                this.error = '粘贴失败，请检查权限';
            }
        },

        // 清空输入
        clearInput() {
            this.inputText = '';
            this.outputText = '';
            this.error = '';
        },

        // 清空输出
        clearOutput() {
            this.outputText = '';
        },

        // 保存到历史
        saveToHistory() {
            const item = {
                timestamp: Date.now(),
                source: this.sourceLang,
                target: this.targetLang,
                inputText: this.inputText,
                outputText: this.outputText,
                preview: this.inputText.substring(0, 100) + (this.inputText.length > 100 ? '...' : ''),
            };
            
            // 添加到开头，限制数量
            this.history.unshift(item);
            if (this.history.length > 50) {
                this.history = this.history.slice(0, 50);
            }
            
            // 保存到本地存储
            this.saveToStorage();
        },

        // 从历史加载
        loadFromHistory(item) {
            this.inputText = item.inputText;
            this.outputText = item.outputText;
            this.sourceLang = item.source;
            this.targetLang = item.target;
            
            // 触发渲染
            this.$nextTick(() => {
                if (window.MathJax?.typesetPromise) {
                    window.MathJax.typesetPromise();
                }
            });
        },

        // 清空历史
        clearHistory() {
            if (confirm('确定要清空所有历史记录吗？')) {
                this.history = [];
                this.saveToStorage();
            }
        },

        // 格式化时间
        formatTime(timestamp) {
            const date = new Date(timestamp);
            const now = new Date();
            const diff = now - date;
            
            if (diff < 60000) {
                return '刚刚';
            } else if (diff < 3600000) {
                return Math.floor(diff / 60000) + '分钟前';
            } else if (diff < 86400000) {
                return Math.floor(diff / 3600000) + '小时前';
            } else {
                return date.toLocaleDateString() + ' ' + date.toLocaleTimeString().slice(0, 5);
            }
        },

        // 保存到本地存储
        saveToStorage() {
            const data = {
                history: this.history,
                fontSize: this.fontSize,
                autoTranslate: this.autoTranslate,
                sourceLang: this.sourceLang,
                targetLang: this.targetLang,
            };
            localStorage.setItem('translator_data', JSON.stringify(data));
        },

        // 从本地存储加载
        loadFromStorage() {
            const stored = localStorage.getItem('translator_data');
            if (stored) {
                try {
                    const data = JSON.parse(stored);
                    this.history = data.history || [];
                    this.fontSize = data.fontSize || 16;
                    this.autoTranslate = data.autoTranslate !== false;
                    this.sourceLang = data.sourceLang || '';
                    this.targetLang = data.targetLang || 'zh';
                } catch (error) {
                    console.error('Failed to load storage:', error);
                }
            }
        },

        // 计算属性：渲染后的输出
        get renderedOutput() {
            if (!this.outputText) {
                return '<div class="placeholder">翻译结果将显示在这里</div>';
            }
            
            // 使用 marked 渲染 Markdown
            if (window.marked) {
                return marked.parse(this.outputText);
            }
            
            // 备用：简单文本显示
            return this.outputText.replace(/\n/g, '<br>');
        }
    };
}

// 启动应用
PetiteVue.createApp(TranslatorApp).mount('#app');