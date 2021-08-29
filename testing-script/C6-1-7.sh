#!/bin/bash

# C6-1-7
# 登入後進入 SYSOP看板 將文章轉錄至 test 看板 以 POST /v1/boards/SYSOP/articles/{{article_id}}  轉發，原文章下面要顯示被轉錄，test看板要看到新文章



ACCESS_TOKEN=`./get_sysop_token.sh`

curl -s http://localhost:8081/v1/boards/SYSOP/articles -H "Authorization: bearer $ACCESS_TOKEN" -d 'action=add_article' --data-urlencode 'title=[測試] test' --data-urlencode 'article=12345'
NEW_FILENAME=`curl  -s http://localhost:8081/v1/boards/SYSOP/articles -H "Authorization: bearer $ACCESS_TOKEN" | jq -r '[.data.items[] | select(.title == "[測試] test")][-1] | .filename'`
echo "add success: $NEW_FILENAME"



curl http://localhost:8081/v1/boards/test/articles/$NEW_FILENAME -d 'action=forward_article'  --data-urlencode 'board_id=test' 
curl  -s http://localhost:8081/v1/boards/test/articles -H "Authorization: bearer $ACCESS_TOKEN" | jq -r '[.data.items[] | select(.title == "[測試] test")][-1] | .filename'
