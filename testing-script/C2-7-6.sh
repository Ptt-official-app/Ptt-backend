#!/bin/bash

# C2-7-6
# 放入 Access Token 之後，設定 draft 0 為 「中文」(0xa4a4a4e5) 後， GET /v1/users/{{自己的ID}}/drafts/0 應該要看到 raw 有 「pKSk5Q==」

if [ "$#" -lt 1 ]; then
	echo "usage: $0 [user_id] [access_token]"
	exit -1
fi

curl http://localhost:8081/v1/users/$1/drafts/0 -v -H "Authorization: bearer $2" -d 'action=update_draft' --data-urlencode 'raw=MTIzC=='


curl http://localhost:8081/v1/users/$1/drafts/0 -H "Authorization: bearer $2" 
