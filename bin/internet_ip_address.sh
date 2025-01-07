#!/usr/bin/env bash

# 获取 timeout 和 time 命令的完整路径
TIMEOUT_CMD=$(which timeout)
TIME_CMD=$(which time)

# Entry function to: test_curls
# Usage: ./network_ip.sh test_curls
# Description: Test all listed curl commands
entry_test_curls() {

    # Function to execute a command with timeout, clear environment variables, and print the command
    __execute_with_timeout() {
        local cmd="$1"
        echo "Executing: $cmd"
        # Clear environment variables that could interfere with the commands
        env -i "$TIME_CMD" "$TIMEOUT_CMD" 1s sh -c "$cmd"
        echo "--------------------------------"
        echo
    }

    # List of commands to execute
    local commands=(
        "curl ifconfig.me"
        "curl icanhazip.com"
        "curl myip.com"
        "curl ip.appspot.com"
        "curl ipinfo.io/ip"
        "curl ipecho.net/plain"
        "curl www.trackip.net/i"
        "curl ip.sb"
        "curl ip.6655.com/ip.aspx"
        "curl whatismyip.akamai.com"
        "wget -qO - ifconfig.co"
        "dig +short myip.opendns.com @resolver1.opendns.com"
        "curl ident.me"
        "curl v4.ident.me"
        "curl v6.ident.me"
        "curl inet-ip.info"
        "curl ip.6655.com/ip.aspx?area=1"
        "curl 1111.ip138.com/ic.asp"
        "curl ip.cn"
        "curl cip.cc"
    )
    for cmd in "${commands[@]}"; do
        __execute_with_timeout "$cmd"
    done
}

# Entry function to: ipv6
# Usage: ./internet_ip.sh ipv6 [url]
# Description: Get IPv6 address using ifconfig.me or provided URL
entry_ipv6() {
    local url="${1:-'ifconfig.me'}"
    env -i "$TIMEOUT_CMD" 1s sh -c "curl $url"
}

# Entry function to: ipv4
# Usage: ./internet_ip.sh ipv4 [url]
# Description: Get IPv4 address using ipinfo.io/ip or provided URL
entry_ipv4() {
    local url="${1:-'ipinfo.io/ip'}"
    env -i "$TIMEOUT_CMD" 1s sh -c "curl $url"
}

# Description: Get internet ip address
# Usage: Commands [ipv4|ipv6|test_curls]
main() {
    _usage "$@"
    # shellcheck disable=SC1091
    source "${AWESOME_SHELL_ROOT}/plugins/interactive.sh" && interactive_select "$0"
}

# 引入 usage.sh 并调用 usage 函数
# shellcheck disable=SC1091
source "${AWESOME_SHELL_ROOT}/core/usage.sh" && usage "$@" && exit 0
