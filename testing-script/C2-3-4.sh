#!/bin/bash

# C2-3-4
# 使用參數 action=add_favorite type=board board_id=SYSOP 新增一個看板後，在只有一個看板時，取得 Access Token 後 GET /v1/users/{{自己的ID}}/favorites 應該要看到 items 裡面有一個 type 為 board 的元素

if [ "$#" -lt 2 ]; then
	echo "usage: $0 [user_id] [access_token]"
	exit -1
fi

curl http://localhost:8081/v1/users/$1/favorites -H "Authorization: bearer $2" -d 'action=add_favorite' -d 'type=board' -d 'board_id=SYSOP'

curl http://localhost:8081/v1/users/$1/favorites -H "Authorization: bearer $2"
