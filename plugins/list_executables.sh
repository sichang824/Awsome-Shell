#!/usr/bin/env bash

# 检查环境和目录
check_environment() {
    if [ -z "${AWESOME_SHELL_ROOT}" ]; then
        echo "错误: AWESOME_SHELL_ROOT 环境变量未设置"
        exit 1
    fi

    if [ ! -d "${AWESOME_SHELL_ROOT}/bin" ]; then
        echo "错误: ${AWESOME_SHELL_ROOT}/bin 目录不存在"
        exit 1
    fi
}

# 确保文件可执行
ensure_executable() {
    local file_path="$1"
    if [ ! -x "$file_path" ]; then
        echo "为 $(basename "$file_path") 添加执行权限..."
        chmod +x "$file_path" || {
            echo "错误: 无法添加执行权限"
            exit 1
        }
    fi
}

# 执行选中的文件
execute_selected_file() {
    local selected_file="$1"
    local selected_path="${AWESOME_SHELL_ROOT}/bin/$selected_file"

    ensure_executable "$selected_path"
    tput rmcup
    echo "你选择了: ${selected_file%.sh}"
    "$selected_path"
    exit 0
}

# 直接列出文件
list_files() {
    check_environment

    echo "在 ${AWESOME_SHELL_ROOT}/bin 中找到的可执行文件："
    echo "----------------------------------------"
    local index=0
    while IFS= read -r file; do
        echo "$index   ${file%.sh}"
        ((index++))
    done < <(cd "${AWESOME_SHELL_ROOT}/bin" && ls -1 *.sh 2>/dev/null)
}

# 显示交互式文件列表
show_files() {
    local current="$1"
    local -n files_ref="$2"
    local -n file_names_ref="$3"

    tput clear
    echo "在 ${AWESOME_SHELL_ROOT}/bin 中找到的可执行文件："
    echo "使用方向键 ↑(k) ↓(j) 选择，回车确认，数字直接选择，q 退出"
    echo "----------------------------------------"
    for i in "${!files_ref[@]}"; do
        if [ "$i" -eq "$current" ]; then
            echo -e "\033[32m$i > ${file_names_ref[$i]}\033[0m"
        else
            echo "  $i   ${file_names_ref[$i]}"
        fi
    done
}

# 处理数字键输入
handle_number_input() {
    local key="$1"
    local -n files_ref="$2"
    local current="$3"
    
    local num=$key
    read -t 0.5 -rsn1 second_key
    if [[ $second_key =~ [0-9] ]]; then
        num="${key}${second_key}"
    fi
    
    # 只返回有效的数字
    if [ "$num" -lt "${#files_ref[@]}" ]; then
        echo "$num"
    else
        echo "$current"
    fi
}

# Description: 交互式列出并选择 bin 目录下的可执行文件
# Usage: ./list_executables.sh [-l|--list]
entry_list_executables() {
    # 检查是否使用 list 参数
    if [ "$1" = "-l" ] || [ "$1" = "--list" ]; then
        list_files
        exit 0
    fi

    check_environment
    tput smcup

    # 获取文件列表
    local files=()
    local file_names=()
    while IFS= read -r file; do
        ensure_executable "${AWESOME_SHELL_ROOT}/bin/$file"
        files+=("$file")
        file_names+=("${file%.sh}")
    done < <(cd "${AWESOME_SHELL_ROOT}/bin" && ls -1 *.sh 2>/dev/null)

    [ ${#files[@]} -eq 0 ] && {
        echo "错误: 在 ${AWESOME_SHELL_ROOT}/bin 目录下没有找到 .sh 文件"
        exit 1
    }

    local current=0
    while true; do
        show_files "$current" files file_names
        read -rsn1 key

        if [[ $key == $'\x1b' ]]; then
            read -rsn2 key
            case $key in
            '[A' | 'k') ((current > 0)) && ((current--)) ;;
            '[B' | 'j') ((current < ${#files[@]} - 1)) && ((current++)) ;;
            esac
        else
            case $key in
            [0-9]) 
                new_current=$(handle_number_input "$key" files "$current")
                if [ "$new_current" -lt "${#files[@]}" ]; then
                    current=$new_current
                    show_files "$current" files file_names
                fi
                ;;
            'q')
                tput rmcup
                exit 0
                ;;
            '') execute_selected_file "${files[$current]}" ;;
            esac
        fi
    done
}

main() {
    entry_list_executables "$@"
}

# 引入 usage.sh 并调用 usage 函数
# shellcheck disable=SC1091
source "${AWESOME_SHELL_ROOT}/core/usage.sh" && usage "$@"
