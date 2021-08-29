#!/bin/bash

# C3-3-1
# 放入 Access Token 之後， GET /v1/boards/SYSOP/settings 可以看到 SYSOP 看板的看板設定

ACCESS_TOKEN=`./get_sysop_token.sh`
curl -s http://localhost:8081/v1/boards/SYSOP/settings -H "Authorization: bearer $ACCESS_TOKEN" 

