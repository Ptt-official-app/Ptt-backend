#!/bin/bash

# C2-1-1
# 取得 Access Token 後 GET /v1/users/SYSOP/information 應該要能夠取得 SYSOP 的資料，其中要有上次登入IP


ACCESS_TOKEN=`./get_sysop_token.sh`
curl http://localhost:8081/v1/users/SYSOP/information -H "Authorization: bearer $ACCESS_TOKEN"

