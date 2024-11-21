#!/usr/bin/env sh

openssl rand -base64 32 | tr -d '+/=' | cut -c -32
