#!/bin/bash

# C6-1-8
# 登入後進入 SYSOP 看板 將文章轉錄自自己的信箱 以POST /v1/boards/SYSOP/articles/{{article_id}} 轉發


if [ "$#" -lt 1 ]; then
	echo "usage: $0 [email]"
	exit -1
fi

ACCESS_TOKEN=`./get_sysop_token.sh`

curl -s http://localhost:8081/v1/boards/SYSOP/articles -H "Authorization: bearer $ACCESS_TOKEN" -d 'action=add_article' --data-urlencode 'title=[測試] test' --data-urlencode 'article=12345'
NEW_FILENAME=`curl  -s http://localhost:8081/v1/boards/SYSOP/articles -H "Authorization: bearer $ACCESS_TOKEN" | jq -r '[.data.items[] | select(.title == "[測試] test")][-1] | .filename'`
echo "add success: $NEW_FILENAME"



curl http://localhost:8081/v1/boards/test/articles/$NEW_FILENAME -H "Authorization: bearer $ACCESS_TOKEN" -d 'action=forward_article'  --data-urlencode "email=$1"
