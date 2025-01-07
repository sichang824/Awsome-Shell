#!/usr/bin/env bash

# Description: 交互式列出并选择 bin 目录下的可执行文件
# Usage: ./list_executables.sh
entry_list_executables() {
    # 保存当前光标位置
    tput smcup

    # 检查 AWESOME_SHELL_ROOT 是否设置
    if [ -z "${AWESOME_SHELL_ROOT}" ]; then
        echo "错误: AWESOME_SHELL_ROOT 环境变量未设置"
        exit 1
    fi

    # 检查目录是否存在
    if [ ! -d "${AWESOME_SHELL_ROOT}/bin" ]; then
        echo "错误: ${AWESOME_SHELL_ROOT}/bin 目录不存在"
        exit 1
    fi

    # 获取所有 .sh 文件并确保它们可执行
    files=()
    while IFS= read -r file; do
        file_path="${AWESOME_SHELL_ROOT}/bin/$file"
        if [ ! -x "$file_path" ]; then
            echo "为 $file 添加执行权限..."
            chmod +x "$file_path"
        fi
        files+=("$file")
    done < <(cd "${AWESOME_SHELL_ROOT}/bin" && ls -1 *.sh 2>/dev/null)

    # 检查是否找到文件
    if [ ${#files[@]} -eq 0 ]; then
        echo "错误: 在 ${AWESOME_SHELL_ROOT}/bin 目录下没有找到 .sh 文件"
        exit 1
    fi
    
    # 初始化选择的索引
    current=0
    
    # 清屏函数
    clear_screen() {
        tput clear
    }
    
    # 显示文件列表
    show_files() {
        clear_screen
        echo "在 ${AWESOME_SHELL_ROOT}/bin 中找到的可执行文件："
        echo "使用方向键 ↑(k) ↓(j) 选择，回车确认，数字直接选择，q 退出"
        echo "----------------------------------------"
        for i in "${!files[@]}"; do
            if [ "$i" -eq "$current" ]; then
                echo -e "\033[32m$i > ${files[$i]}\033[0m"
            else
                echo "  $i   ${files[$i]}"
            fi
        done
    }
    
    # 主循环
    while true; do
        show_files
        
        # 读取按键
        read -rsn1 key
        
        # 检查是否是方向键的第一个字符
        if [[ $key == $'\x1b' ]]; then
            read -rsn2 key
            case $key in
                '[A') # 上箭头
                    ((current > 0)) && ((current--))
                    ;;
                '[B') # 下箭头
                    ((current < ${#files[@]}-1)) && ((current++))
                    ;;
            esac
        else
            case $key in
                [0-9]) # 数字键
                    num=$key
                    # 等待可能的第二个数字
                    read -t 0.5 -rsn1 second_key
                    if [[ $second_key =~ [0-9] ]]; then
                        num="${key}${second_key}"
                    fi
                    # 检查数字是否在有效范围内
                    if [ "$num" -lt "${#files[@]}" ]; then
                        current=$num
                        # 显示选择
                        show_files
                        # 短暂延迟后执行
                        sleep 0.2
                        selected_file="${files[$current]}"
                        selected_path="${AWESOME_SHELL_ROOT}/bin/$selected_file"
                        
                        # 检查文件是否存在且可执行
                        if [ ! -x "$selected_path" ]; then
                            tput rmcup
                            echo "错误: $selected_file 不可执行，正在添加执行权限..."
                            chmod +x "$selected_path"
                            if [ $? -ne 0 ]; then
                                echo "错误: 无法添加执行权限"
                                exit 1
                            fi
                            echo "已添加执行权限"
                        fi
                        
                        # 恢复光标位置
                        tput rmcup
                        echo "你选择了: $selected_file"
                        # 执行选中的脚本
                        "$selected_path"
                        exit 0
                    fi
                    ;;
                'k') # vim 上
                    ((current > 0)) && ((current--))
                    ;;
                'j') # vim 下
                    ((current < ${#files[@]}-1)) && ((current++))
                    ;;
                'q') # 退出
                    clear_screen
                    # 恢复光标位置
                    tput rmcup
                    exit 0
                    ;;
                '') # 回车
                    clear_screen
                    selected_file="${files[$current]}"
                    selected_path="${AWESOME_SHELL_ROOT}/bin/$selected_file"
                    
                    # 检查文件是否存在且可执行
                    if [ ! -x "$selected_path" ]; then
                        tput rmcup
                        echo "错误: $selected_file 不可执行，正在添加执行权限..."
                        chmod +x "$selected_path"
                        if [ $? -ne 0 ]; then
                            echo "错误: 无法添加执行权限"
                            exit 1
                        fi
                        echo "已添加执行权限"
                    fi
                    
                    # 恢复光标位置
                    tput rmcup
                    echo "你选择了: $selected_file"
                    # 执行选中的脚本
                    "$selected_path"
                    exit 0
                    ;;
            esac
        fi
    done
}

main() {
    entry_list_executables
}

# 引入 usage.sh 并调用 usage 函数
# shellcheck disable=SC1091
source "${AWESOME_SHELL_ROOT}/core/usage.sh" && usage "$@" 
