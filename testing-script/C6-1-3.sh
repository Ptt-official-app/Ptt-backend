#!/bin/bash

# C6-1-1
# 以user01登入後於test看板發出一篇文章，接下來透過 user01 可以刪除該篇文章



ACCESS_TOKEN=`./get_sysop_token.sh`
curl -v -s http://localhost:8081/v1/boards/test/articles -H "Authorization: bearer $ACCESS_TOKEN" -d 'action=add_article' --data-urlencode 'title=[測試] test' --data-urlencode 'article=12345'


NEW_FILENAME=`curl  -s http://localhost:8081/v1/boards/test/articles -H "Authorization: bearer $ACCESS_TOKEN" | jq -r '[.data.items[] | select(.title == "[測試] test")][-1] | .filename'`
echo "add success: $NEW_FILENAME"
curl  http://localhost:8081/v1/boards/test/articles/$NEW_FILENAME -d 'action=delete'
NEW_FILENAME2=`curl -s http://localhost:8081/v1/boards/test/articles -H "Authorization: bearer $ACCESS_TOKEN" | jq -r '[.data.items[] | select(.title == "[測試] test")][-1] | .filename'`

echo $NEW_FILENAME2




