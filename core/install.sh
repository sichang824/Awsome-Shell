#!/usr/bin/env bash

git clone https://e.coding.net/cloudbase-100009281119/Awesome-Shell/Awesome-Shell.git ~/.Awesome-Shell

echo echo '"export AWESOME_SHELL_ROOT=$HOME/.Awesome-Shell" >> ~/.config/fish/config.fish'
echo echo '"export AWESOME_SHELL_ROOT=$HOME/.Awesome-Shell" >> ~/.bashrc'
echo 'source ~/.bashrc'

