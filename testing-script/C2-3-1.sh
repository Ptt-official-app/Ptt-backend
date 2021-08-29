#!/bin/bash

# C2-3-1
# 取得 Access Token 後在不新增任何最愛項目的情況 GET /v1/users/{{自己的ID}}/favorites 應該要看到 data 和 items 項目，其中items 應該要是空陣列

if [ "$#" -lt 2 ]; then
	echo "usage: $0 [user_id] [access_token]"
	exit -1
fi

echo "curl http://localhost:8081/v1/users/$1/favorites -H \"Authorization: bearer $2\""

curl http://localhost:8081/v1/users/$1/favorites -H "Authorization: bearer $2"

