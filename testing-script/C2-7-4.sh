#!/bin/bash

# C2-7-4
# 放入 Access Token 之後， POST /v1/users/{{自己的ID}}/drafts/0 參數 action=update_draft raw=MTIzC%3D%3D (123 的 base64編碼) 應該要成功

if [ "$#" -lt 1 ]; then
	echo "usage: $0 [user_id] [access_token]"
	exit -1
fi

curl http://localhost:8081/v1/users/$1/drafts/0 -v -H "Authorization: bearer $2" -d 'action=update_draft' --data-urlencode 'raw=MTIzC=='
