#!/usr/bin/env bash

# Description: 交互式命令选择模块

# 交互式命令选择
interactive_select() {
    local script_path="$1"
    local history_file="/tmp/.cmd_history_$$"
    
    # 确保退出时删除临时历史文件
    trap "rm -f $history_file" EXIT
    
    # 启用 readline 功能
    if [ -z "$HISTFILE" ]; then
        export HISTFILE="$history_file"
    fi
    
    # 设置历史记录大小
    HISTSIZE=1000
    HISTFILESIZE=2000
    
    # 启用历史记录
    set -o history
    
    # 只在交互式 shell 中设置按键绑定
    if [[ $- == *i* ]]; then
        bind '"\e[A": history-search-backward'
        bind '"\e[B": history-search-forward'
    fi
    
    while true; do
        # 显示提示信息
        echo -e "\n请输入命令参数（直接回车退出）："
        
        # 读取用户输入的参数（启用行编辑功能）
        read -e -p "> " params
        
        # 如果用户直接按回车，退出
        if [ -z "$params" ]; then
            return 0
        fi
        
        # 将有效命令添加到历史记录
        history -s "$params"
        
        # 执行命令并传入用户输入的参数
        "$script_path" $params
        echo  # 添加换行，让输出更清晰
    done
} 