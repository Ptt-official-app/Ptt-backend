#!/bin/bash

# C2-7-7
# 放入 Access Token 之後，設定 draft 0 為 123 後 存取 draft 1 不應該看到任何資料

if [ "$#" -lt 1 ]; then
	echo "usage: $0 [user_id] [access_token]"
	exit -1
fi

curl http://localhost:8081/v1/users/$1/drafts/0 -v -H "Authorization: bearer $2" -d 'action=update_draft' --data-urlencode 'raw=MTIzC=='

curl http://localhost:8081/v1/users/$1/drafts/1 -H "Authorization: bearer $2" 
