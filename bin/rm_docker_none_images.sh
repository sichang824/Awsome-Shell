#!/bin/bash

# Description: 删除 Docker 中所有 <none> 标签的镜像
# Usage: ./rm_docker_none_images.sh
entry_remove_docker_none_images() {
    # 获取所有 <none> 镜像的 ID
    none_images=$(docker images -f "dangling=true" -q)

    if [ -z "$none_images" ]; then
        echo "没有发现 <none> 标签的镜像"
        exit 0
    fi

    echo "发现以下 <none> 标签的镜像："
    docker images -f "dangling=true"

    echo -n "是否确认删除这些镜像？(y/n): "
    read confirm

    if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
        docker rmi $none_images
        echo "删除完成"
    else
        echo "操作已取消"
    fi
}

source "${AWESOME_SHELL_ROOT}/core/usage.sh" && usage "$@"
