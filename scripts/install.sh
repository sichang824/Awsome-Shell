#!/usr/bin/env bash

# shellcheck disable=SC1091
source "$(dirname "$(dirname "$0")")/plugins/colors.sh"
source "$(dirname "$(dirname "$0")")/plugins/replace.sh"

SHELL_ROOT_CONFIG='export AWESOME_SHELL_ROOT=$HOME/.Awesome-Shell'
SHELL_ALIAS_CONFIG='alias as=$HOME/.Awesome-Shell/bin/hello.sh'

configure_shell() {
    local config_file="$1"
    replace_line "$config_file" "^export AWESOME_SHELL_ROOT=" "$SHELL_ROOT_CONFIG"
    replace_line "$config_file" "^alias as=" "$SHELL_ALIAS_CONFIG"
}

get_message() {
    case "$1" in
    "START") echo "开始安装 Awesome-Shell..." ;;
    "MODE") printf "安装模式: %s" "$2" ;;
    "DIR_EXISTS") echo "检测到 ~/.Awesome-Shell 目录已存在" ;;
    "CONFIRM_DELETE") echo "是否删除已存在的目录？(y/N) " ;;
    "DELETING") echo "正在删除已存在的目录..." ;;
    "DELETE_FAILED") echo "${COLOR_RED}${SYMBOL_ERROR} 删除目录失败${COLOR_NONE}" ;;
    "INSTALL_CANCELLED") echo "${COLOR_YELLOW}安装已取消${COLOR_NONE}" ;;
    "COPYING") echo "正在复制本地文件..." ;;
    "COPY_SUCCESS") echo "${COLOR_GREEN}${SYMBOL_SUCCESS} 本地文件复制完成${COLOR_NONE}" ;;
    "COPY_FAILED") echo "${COLOR_RED}${SYMBOL_ERROR} 本地文件复制失败${COLOR_NONE}" ;;
    "CLONING") echo "正在从远程仓库克隆..." ;;
    "CLONE_SUCCESS") echo "${COLOR_GREEN}${SYMBOL_SUCCESS} 远程仓库克隆完成${COLOR_NONE}" ;;
    "CLONE_FAILED") echo "${COLOR_RED}${SYMBOL_ERROR} 远程仓库克隆失败${COLOR_NONE}" ;;
    "SHELL_DETECTED") printf "检测到当前 Shell: %s" "$2" ;;
    "CONFIG_START") echo "正在配置环境变量..." ;;
    "FISH_SUCCESS") echo "${COLOR_GREEN}${SYMBOL_SUCCESS} Fish shell 配置完成${COLOR_NONE}" ;;
    "BASH_SUCCESS") echo "${COLOR_GREEN}${SYMBOL_SUCCESS} Bash shell 配置完成${COLOR_NONE}" ;;
    "ZSH_SUCCESS") echo "${COLOR_GREEN}${SYMBOL_SUCCESS} Zsh shell 配置完成${COLOR_NONE}" ;;
    "UNSUPPORTED_SHELL") printf "${COLOR_YELLOW}${SYMBOL_WARNING} 警告：当前 shell (%s) 不支持自动配置，请手动添加以下配置：${COLOR_NONE}" "$2" ;;
    "INSTALL_COMPLETE") echo "${COLOR_GREEN}${SYMBOL_STAR} Awesome-Shell 安装完成！${COLOR_NONE}" ;;
    "USAGE_HINT") echo "${COLOR_BLUE}使用 'as' 命令开始使用${COLOR_NONE}" ;;
    esac
}

ENV=${1:-"remote"}

echo -e "$(get_message "START")"
echo -e "$(get_message "MODE" "$ENV")"

if [ "$ENV" = "local" ]; then
    if [ -d ~/.Awesome-Shell ]; then
        echo -e "$(get_message "DIR_EXISTS")"
        read -p "$(get_message "CONFIRM_DELETE")" confirm
        if [[ $confirm =~ ^[Yy]$ ]]; then
            echo -e "$(get_message "DELETING")"
            rm -rf ~/.Awesome-Shell
            if [ $? -ne 0 ]; then
                echo -e "$(get_message "DELETE_FAILED")"
                exit 1
            fi
        else
            echo -e "$(get_message "INSTALL_CANCELLED")"
            exit 0
        fi
    fi

    echo -e "$(get_message "COPYING")"
    cp -r "$(dirname "$(dirname "$0")")" ~/.Awesome-Shell
    if [ $? -eq 0 ]; then
        echo -e "$(get_message "COPY_SUCCESS")"
    else
        echo -e "$(get_message "COPY_FAILED")"
        exit 1
    fi
else
    echo -e "$(get_message "CLONING")"
    git clone https://e.coding.net/cloudbase-100009281119/Awesome-Shell/Awesome-Shell.git ~/.Awesome-Shell
    if [ $? -eq 0 ]; then
        echo -e "$(get_message "CLONE_SUCCESS")"
    else
        echo -e "$(get_message "CLONE_FAILED")"
        exit 1
    fi
fi

current_shell=$(basename "$SHELL")
echo -e "$(get_message "SHELL_DETECTED" "$current_shell")"

echo -e "$(get_message "CONFIG_START")"
case "$current_shell" in
fish)
    configure_shell ~/.config/fish/config.fish
    echo -e "$(get_message "FISH_SUCCESS")"
    ;;
bash)
    configure_shell ~/.bashrc
    source ~/.bashrc
    echo -e "$(get_message "BASH_SUCCESS")"
    ;;
zsh)
    configure_shell ~/.zshrc
    source ~/.zshrc
    echo -e "$(get_message "ZSH_SUCCESS")"
    ;;
*)
    echo -e "$(get_message "UNSUPPORTED_SHELL" "$current_shell")"
    echo "$SHELL_ROOT_CONFIG"
    echo "$SHELL_ALIAS_CONFIG"
    ;;
esac

echo -e "$(get_message "INSTALL_COMPLETE")"
echo -e "$(get_message "USAGE_HINT")"
