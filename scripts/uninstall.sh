#!/usr/bin/env bash

source "$(dirname "$(dirname "$0")")/core/colors.sh"

get_message() {
    case "$1" in
    "UNINSTALLING") echo "${COLOR_YELLOW}正在卸载 Awesome-Shell...${COLOR_NONE}" ;;
    "UNINSTALL_SUCCESS") echo "${COLOR_GREEN}Awesome-Shell 卸载成功！${COLOR_NONE}" ;;
    esac
}

echo -e "$(get_message "UNINSTALLING")"

rm -rf ~/.Awesome-Shell

echo -e "$(get_message "UNINSTALL_SUCCESS")"
