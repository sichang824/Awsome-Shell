#!/usr/bin/env sh

# 默认密码长度为32
length=${1:-32}

# 生成随机密码并按指定长度截取
openssl rand -base64 48 | tr -d '+/=' | cut -c -${length}
