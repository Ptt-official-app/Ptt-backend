#!/bin/bash

# C2-7-2
# 放入 Access Token 之後，在沒有任何草稿的情況下 GET /v1/users/{{自己的ID}}/drafts/0 不應該看到任何資料

if [ "$#" -lt 2 ]; then
	echo "usage: $0 [user_id] [access_token]"
	exit -1
fi

curl http://localhost:8081/v1/users/$1/drafts/0 -H "Authorization: bearer $2" 
