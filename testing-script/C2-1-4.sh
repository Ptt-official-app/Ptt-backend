#!/bin/bash

# C2-1-4
# 取得 Access Token 後 GET /v1/users/{{自己的ID}}/information 應該取得自己的使用者資料，其中金錢數量應該能顯示正確數值

if [ "$#" -lt 2 ]; then
	echo "usage: $0 [user_id] [access_token]"
	exit -1
fi

echo "curl http://localhost:8081/v1/users/$1/information -H \"Authorization: bearer $2\""

curl http://localhost:8081/v1/users/$1/information -H "Authorization: bearer $2"

