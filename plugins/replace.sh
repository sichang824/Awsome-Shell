#!/usr/bin/env bash

# shellcheck disable=SC1091
source "$(dirname "$(dirname "$0")")/plugins/colors.sh"

# 替换或添加配置到文件
# $1: 目标文件路径
# $2: 匹配的字符串模式（用于grep）
# $3: 替换的内容
replace_line() {
    local config_file="$1"
    local match_pattern="$2"
    local replace_content="$3"

    # 如果文件不存在，创建它
    if [ ! -f "$config_file" ]; then
        touch "$config_file"
    fi

    # 读取文件内容到临时文件
    local temp_file
    temp_file=$(mktemp)

    # 处理配置
    if grep -q "${match_pattern}" "$config_file"; then
        # 替换已存在的配置
        sed "s|${match_pattern}.*|${replace_content}|" "$config_file" >"$temp_file"
    else
        # 添加新配置
        cat "$config_file" >"$temp_file"
        echo "$replace_content" >>"$temp_file"
    fi

    # 将临时文件移回原文件
    mv "$temp_file" "$config_file"

    # 确保文件权限正确
    chmod 644 "$config_file"
}
