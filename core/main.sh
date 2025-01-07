#!/usr/bin/env bash

# 设置严格模式
set -euo pipefail

# 定义帮助信息函数
show_help() {
    cat <<EOF
用法: $(basename "$0") [选项] [命令]

选项:
    -h, --help     显示此帮助信息
    -l, --list     列出所有可用命令
    -v, --version  显示版本信息

如果没有提供参数，将显示可用命令列表。
EOF
}

# 定义版本信息
VERSION="1.0.0"

# 执行具体命令的函数
execute_command() {
    local command_name="$1"
    shift # 移除第一个参数

    local command_path="${AWESOME_SHELL_ROOT}/bin/${command_name}.sh"

    if [ ! -f "$command_path" ]; then
        echo "错误: 找不到命令 '$command_name'" >&2
        echo "使用 --list 查看所有可用命令" >&2
        exit 1
    fi

    if [ ! -x "$command_path" ]; then
        echo "错误: 命令文件没有执行权限 '$command_path'" >&2
        exit 1
    fi

    bash "$command_path" "$@"
}

# 检查必要的环境变量
if [ -z "${AWESOME_SHELL_ROOT:-}" ]; then
    echo "错误: AWESOME_SHELL_ROOT 环境变量未设置" >&2
    exit 1
fi

# Description: Welcome to use Awesome Shell
# Usage: Commands [-l]
main() {
    # 如果没有参数，显示命令列表
    if [ $# -eq 0 ]; then
        bash "${AWESOME_SHELL_ROOT}/plugins/list_executables.sh"
        exit 0
    fi

    # 解析命令行参数
    case "$1" in
    -h | --help)
        show_help
        exit 0
        ;;
    -v | --version)
        echo "Version: $VERSION"
        exit 0
        ;;
    -l | --list)
        bash "${AWESOME_SHELL_ROOT}/plugins/list_executables.sh" -l
        exit 0
        ;;
    -*)
        echo "错误: 未知选项 '$1'" >&2
        show_help
        exit 1
        ;;
    *)
        execute_command "$@"
        ;;
    esac
}

main "$@"
