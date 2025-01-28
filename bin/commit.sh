#!/usr/bin/env bash

# 定义Ollama模型ID（请替换为您实际使用的模型）
MODEL_ID="deepseek-r1:14b"

# 获取Git变更记录（排除状态标识符，仅保留文件名）
changes=$(git status -s | cut -c3-)
contents=$(git diff --staged)

# 使用Ollama生成commit信息
commit_message=$(
    ollama run $MODEL_ID generate \
        "Changes:
    $(echo -e "$changes")

    Change Contents:
    $(echo -e "$contents")"
)

if [ $? -ne 0 ]; then
    echo "生成commit消息时出现错误。"
    exit 1
fi

# 显示生成的消息
echo -e "\n生成的commit消息:\n$commit_message"

# 提交代码（可选：直接执行git commit）
read -p "是否现在提交？(y/n) " confirm

if [ "$confirm" == "y" ]; then
    git add .
    git commit -m "$commit_message"
fi
