#!/bin/bash

# C1-1-1
# 輸入正確的帳號密碼 (SYSOP / 123123) 以 /v1/token 登入系統應該要取得 access token 參數


if [ "$#" -lt 2 ]; then
	echo "usage: $0 [user_id] [password]"
	exit -1
fi

curl http://localhost:8081/v1/token  -q -d "grant_type=password&username=$1&password=$2" | jq -r '.access_token'

