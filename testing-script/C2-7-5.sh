#!/bin/bash

# C2-7-5
# 放入 Access Token 之後，設定 draft 0 為 123 後， GET /v1/users/{{自己的ID}}/drafts/0 應該要看到 raw 有 MTIzC==  (123 的 base64 編碼)

if [ "$#" -lt 1 ]; then
	echo "usage: $0 [user_id] [access_token]"
	exit -1
fi

curl http://localhost:8081/v1/users/$1/drafts/0 -v -H "Authorization: bearer $2" -d 'action=update_draft' --data-urlencode 'raw=MTIzC=='


curl http://localhost:8081/v1/users/$1/drafts/0 -H "Authorization: bearer $2" 
