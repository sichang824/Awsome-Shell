#!/bin/bash

filename=""

# Description: 创建 SSH 密钥对并添加到 SSH agent
# Usage: ./ssh_git_publib_key.sh create_ssh_key
entry_create_ssh_key() {
    echo "该脚本将创建一个新的 SSH 密钥对并添加到你的 SSH agent 中。"
    read -p "请输入你的邮箱: " email
    # 将生成的密钥文件路径赋值给全局变量 filename
    filename="$HOME/.ssh/id_rsa_$(date +%Y%m%d%H%M%S)"
    
    # 确保 ~/.ssh 目录存在
    mkdir -p ~/.ssh
    
    ssh-keygen -t rsa -b 4096 -C "$email" -f "$filename" -N ""
    eval "$(ssh-agent -s)"
    ssh-add "$filename"
}

# Description: 将公钥上传到 GitHub
# Usage: ./ssh_git_publib_key.sh upload_public_key
entry_upload_public_key() {
    # 直接使用全局变量 filename,不需要再传参数
    read -p "请输入 GitHub 域名: " github_domain
    cat "$filename" | ssh -T $github_domain
    echo "SSH 密钥对已创建并添加到你的 SSH agent 中。"
}

# Description: 配置 SSH config
# Usage: ./ssh_git_publib_key.sh configure_ssh_config
entry_configure_ssh_config() {
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
                        /Host '"$github_domain"'/,/IdentityFile/ {
                            if(/Host/) print "Host " github_domain; 
                            else if(/HostName/) print "    HostName " github_domain;
                            else if(/User/) print "    User " github_user;
                            else if(/Port/) print "    Port 22";
                            else if(/IdentityFile/) print "    IdentityFile " filename;
                            next;
                        }
                        1
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

# Description: 集成创建ssh key、上传公钥、配置SSH config
# Usage: ./ssh_git_publib_key.sh aio
entry_aio() {
    entry_create_ssh_key
    entry_upload_public_key
    entry_configure_ssh_config
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
