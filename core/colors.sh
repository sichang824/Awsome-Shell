#!/usr/bin/env bash

# 定义颜色
COLOR_RED='\033[0;31m'
COLOR_GREEN='\033[0;32m'
COLOR_YELLOW='\033[1;33m'
COLOR_BLUE='\033[0;34m'
COLOR_NONE='\033[0m'

# 定义符号
SYMBOL_SUCCESS="✓"
SYMBOL_ERROR="✗"
SYMBOL_WARNING="⚠"
SYMBOL_STAR="✨"

# Entry function to show all available colors and styles
# Usage: ./colors.sh
# Description: Displays all available ANSI colors and styles
main() {
    # 标准 8 色文本颜色
    printf "\033[30m黑色\033[0m "
    printf "\033[31m红色\033[0m "
    printf "\033[32m绿色\033[0m "
    printf "\033[33m黄色\033[0m "
    printf "\033[34m蓝色\033[0m "
    printf "\033[35m洋红\033[0m "
    printf "\033[36m青色\033[0m "
    printf "\033[37m白色\033[0m\n"

    # 高亮 8 色文本颜色
    printf "\033[1;30m亮黑色\033[0m "
    printf "\033[1;31m亮红色\033[0m "
    printf "\033[1;32m亮绿色\033[0m "
    printf "\033[1;33m亮黄色\033[0m "
    printf "\033[1;34m亮蓝色\033[0m "
    printf "\033[1;35m亮洋红\033[0m "
    printf "\033[1;36m亮青色\033[0m "
    printf "\033[1;37m亮白色\033[0m\n"

    # 其他 ANSI 转义码
    printf "\033[2mDIM\033[0m "
    printf "\033[3mItalic\033[0m "
    printf "\033[4mUnderline\033[0m "
    printf "\033[5mBlink\033[0m "
    printf "\033[7mReverse\033[0m "
    printf "\033[9mStrikethrough\033[0m\n"
}

# # shellcheck disable=SC1091
# source "${AWESOME_SHELL_ROOT}/core/usage.sh" && usage "${@}"
