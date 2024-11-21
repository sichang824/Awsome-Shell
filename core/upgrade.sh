#!/usr/bin/env bash

if [ -d "${AWESOME_SHELL_ROOT}" ]; then
    cd "${AWESOME_SHELL_ROOT}" || exit
    git pull
else
    echo "目录 ${AWESOME_SHELL_ROOT} 不存在。"
fi
