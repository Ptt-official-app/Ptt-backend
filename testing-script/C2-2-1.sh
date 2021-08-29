#!/bin/bash

# C2-2-1
# 使用非 SYSOP 帳號取得 Access Token 後 GET /v1/users/SYSOP/preferences 不應該取得任何資料，回應的 JSON 要有 error 欄位


if [ "$#" -lt 1 ]; then
	echo "usage: $0 [access_token]"
	exit -1
fi
ACCESS_TOKEN=$1

curl http://localhost:8081/v1/users/SYSOP/preferences -H "Authorization: bearer $ACCESS_TOKEN"

