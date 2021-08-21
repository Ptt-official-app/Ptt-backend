#!/bin/bash

# C2-1-5
# 使用 SYSOP 的 Access Token 登入後 GET /v1/users/SYSOP/information 應該要能夠取得 SYSOP 的資料，其中登入時間應該要是今天


ACCESS_TOKEN=`./get_sysop_token.sh`
curl http://localhost:8081/v1/users/SYSOP/information -H "Authorization: bearer $ACCESS_TOKEN"

