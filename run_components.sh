#!/usr/bin/env bash

cp oauth-demo.yml ~/.teamocil
sed -i '' s%CURRENT_DIR%$(pwd)%g oauth-demo.yml
itermocil oauth-demo