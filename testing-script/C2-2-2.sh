#!/bin/bash

# C2-2-2
# 取得 Access Token 後 GET /v1/users/{{自己的ID}}/preferences 應該要取得自己的資料

if [ "$#" -lt 2 ]; then
	echo "usage: $0 [user_id] [access_token]"
	exit -1
fi

echo "curl http://localhost:8081/v1/users/$1/preferences -H \"Authorization: bearer $2\""

curl http://localhost:8081/v1/users/$1/preferences -H "Authorization: bearer $2"

