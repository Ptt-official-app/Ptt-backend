#!/bin/bash

# C6-1-5
# 在放入 user01 Access Token 的情況下，發出一篇文章，user02 可以用 「↑」進行推文，評價數要加一，user02第二次用 ↑ 推文時評價數不變 接著user02 用 「↓」推文時原先的上箭頭推文被刪去，評價數減一（收回），第二次用 ↓ 推文時評價數減一，第三次用↓推文時評價數不變



ACCESS_TOKEN=`./get_sysop_token.sh`

PICHU_2_TOKEN=`./get_user_token.sh pichu2 123123`
PICHU_3_TOKEN=`./get_user_token.sh pichu3 123123`

curl -s http://localhost:8081/v1/boards/test/articles -H "Authorization: bearer $PICHU_2_TOKEN" -d 'action=add_article' --data-urlencode 'title=[測試] test' --data-urlencode 'article=12345'
NEW_FILENAME=`curl  -s http://localhost:8081/v1/boards/test/articles -H "Authorization: bearer $ACCESS_TOKEN" | jq -r '[.data.items[] | select(.title == "[測試] test")][-1] | .filename'`
echo "add success: $NEW_FILENAME"
curl http://localhost:8081/v1/boards/test/articles/$NEW_FILENAME -d 'action=append_comment'  --data-urlencode 'type＝↑' 
curl -v -s http://localhost:8081/v1/boards/test/articles -H "Authorization: bearer $ACCESS_TOKEN" 