#!/bin/bash

config_content="""#	$OpenBSD: ssh_config,v 1.35 2020/07/17 03:43:42 dtucker Exp $

# This is the ssh client system-wide configuration file.  See
# ssh_config(5) for more information.  This file provides defaults for
# users, and the values can be changed in per-user configuration files
# or on the command line.

# Configuration data is parsed as follows:
#  1. command line options
#  2. user-specific file
#  3. system-wide file
# Any configuration value is only changed the first time it is set.
# Thus, host-specific definitions should be at the beginning of the
# configuration file, and defaults at the end.

# This Include directive is not part of the default ssh_config shipped with
# OpenSSH. Options set in the included configuration files generally override
# those that follow.  The defaults only apply to options that have not been
# explicitly set.  Options that appear multiple times keep the first value set,
# unless they are a multivalue option such as IdentityFile.
Include config.d/*

# Site-wide defaults for some commonly used options.  For a comprehensive
# list of available options, their meanings and defaults, please see the
# ssh_config(5) man page.

# Host *
#   ForwardAgent no
#   ForwardX11 no
#   PasswordAuthentication yes
#   HostbasedAuthentication no
#   GSSAPIAuthentication no
#   GSSAPIDelegateCredentials no
#   BatchMode no
#   CheckHostIP yes
#   AddressFamily any
#   ConnectTimeout 0
#   StrictHostKeyChecking ask
#   IdentityFile ~/.ssh/id_rsa
#   IdentityFile ~/.ssh/id_dsa
#   IdentityFile ~/.ssh/id_ecdsa
#   IdentityFile ~/.ssh/id_ed25519
#   Port 22
#   Ciphers aes128-ctr,aes192-ctr,aes256-ctr,aes128-cbc,3des-cbc
#   MACs hmac-md5,hmac-sha1,umac-64@openssh.com
#   EscapeChar ~
#   Tunnel no
#   TunnelDevice any:any
#   PermitLocalCommand no
#   VisualHostKey no
#   ProxyCommand ssh -q -W %h:%p gateway.example.com
#   RekeyLimit 1G 1h
#   UserKnownHostsFile ~/.ssh/known_hosts.d/%k
Host *
    SendEnv LANG LC_*
    ServerAliveInterval 10
    ServerAliveCountMax 99999
    TCPKeepAlive yes
    HashKnownHosts yes
    GSSAPIAuthentication yes"""



# Description: 创建 SSH 密钥对并添加到 SSH agent
# Usage: ./ssh_git_publib_key.sh create_ssh_key
entry_create_ssh_key() {
    echo "该脚本将创建一个新的 SSH 密钥对并添加到你的 SSH agent 中。"
    read -p "请输入你的邮箱: " email
    # 将生成的密钥文件路径赋值给全局变量 filename
    filename="$HOME/.ssh/id_rsa_$(date +%Y%m%d%H%M%S)"

    echo "密钥文件路径: $filename"
    
    # 确保 ~/.ssh 目录存在
    mkdir -p ~/.ssh
    
    ssh-keygen -t rsa -b 4096 -C "$email" -f "$filename" -N ""
    eval "$(ssh-agent -s)"
    ssh-add "$filename"
}

# Description: 将公钥上传到 GitHub
# Usage: ./ssh_git_publib_key.sh check_public_key
entry_check_public_key() {
    # 直接使用全局变量 filename,不需要再传参数
    read -p "请输入 GitHub 域名: " github_domain
    ssh -T $github_domain
}

# Description: 配置 SSH config
# Usage: ./ssh_git_publib_key.sh configure_ssh_config
entry_configure_ssh_config() {
    read -p "请输入密钥文件路径: " filename
    read -p "请输入 GitHub 用户名: " github_user
    read -p "请输入 GitHub 域名: " github_domain
    
    mkdir -p ~/.ssh/config.d
    ssh_config_file=~/.ssh/config.d/git
    
    if grep -q "Host $github_domain" "$ssh_config_file" 2>/dev/null && 
       grep -q "User $github_user" "$ssh_config_file" 2>/dev/null; then
        echo "检测到 $ssh_config_file 中已存在相同的 Host 和 User 配置:"
        awk "/Host $github_domain/,/IdentityFile/" "$ssh_config_file"
        
        while true; do
            read -p "是否覆盖该配置? [y/n]: " yn
            case $yn in
                [Yy]* ) 
                    awk -v github_domain="$github_domain" -v github_user="$github_user" -v filename="$filename" '
                        /Host '"$github_domain"'/{
                            found=1
                            print "Host " github_domain
                            print "    HostName " github_domain
                            print "    User " github_user
                            print "    Port 22" 
                            print "    IdentityFile " filename
                            next
                        }
                        found==0 {print}
                    ' "$ssh_config_file" > tmp_config && mv tmp_config "$ssh_config_file"
                    echo "Host $github_domain User $github_user 的配置已被覆盖"
                    break
                    ;;
                [Nn]* )
                    echo "Host $github_domain User $github_user 的配置未被修改"
                    break
                    ;;
                * ) 
                    echo "请输入 y 或者 n"
                    ;;
            esac
        done
    else
        cat <<EOF >> "$ssh_config_file"
Host $github_domain
    HostName $github_domain
    User $github_user
    Port 22
    IdentityFile $filename
EOF
        echo "已添加 Host $github_domain User $github_user 到 $ssh_config_file"
    fi
}

# Description: 写入 SSH 配置文件 
# Usage: ./ssh_git_publib_key.sh write_ssh_config
entry_write_ssh_config() {
    ssh_config_file=~/.ssh/config
    
    echo "$config_content" > "$ssh_config_file"
    echo "已写入配置到 $ssh_config_file"
    
    echo "当前 $ssh_config_file 的内容为:"
    cat "$ssh_config_file"
}

# Description: 集成创建ssh key、上传公钥、配置SSH config
# Usage: ./ssh_git_publib_key.sh aio
entry_aio() {
    entry_create_ssh_key
    entry_configure_ssh_config
    entry_check_public_key
}

# Description: 查看ssh-add的秘钥列表
# Usage: ./ssh_git_publib_key.sh list_ssh_keys
entry_list_ssh_keys() {
    echo "当前 ssh-add 中的秘钥列表如下:"
    ssh-add -l
}

# Description: 集成创建ssh key、上传公钥、配置SSH config
main() {
    _usage
}

# 引入 usage.sh 并调用 usage 函数
# shellcheck disable=SC1091
source "${AWESOME_SHELL_ROOT}/core/usage.sh" && usage "$@"
