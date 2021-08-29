#!/bin/bash

# C2-1-2
# 取得 Access Token 後 GET /v1/users/NOTEXIST/information 不應該取得任何人的資料，回應的 JSON 要有 error 欄位


ACCESS_TOKEN=`./get_sysop_token.sh`
curl http://localhost:8081/v1/users/NOTEXIST/information -H "Authorization: bearer $ACCESS_TOKEN"

