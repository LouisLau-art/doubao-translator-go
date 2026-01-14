#!/bin/bash
# download-libs.sh

mkdir -p static/libs

# 下载 Vue Petite
curl -L https://unpkg.com/petite-vue@0.4.1/dist/petite-vue.iife.js -o static/libs/petite-vue.js

# 下载 Marked.js
curl -L https://cdn.jsdelivr.net/npm/marked/marked.min.js -o static/libs/marked.min.js

# 下载 MathJax
mkdir -p static/libs/mathjax
curl -L https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-mml-chtml.js -o static/libs/mathjax/tex-mml-chtml.js

echo "所有库文件已下载到 static/libs/"
