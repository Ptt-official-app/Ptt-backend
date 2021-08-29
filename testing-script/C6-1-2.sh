#!/bin/bash

# C6-1-2
# 登入後進入 test 看板，發出一篇標題為 「[測試] test」內文為 12345 的文章，接下來在下面進行推文，推送 「推」「114422」 要能夠順利推文



ACCESS_TOKEN=`./get_sysop_token.sh`
# curl -v -s http://localhost:8081/v1/boards/test/articles -H "Authorization: bearer $ACCESS_TOKEN" -d 'action=add_article' --data-urlencode 'title=[測試] test' --data-urlencode 'article=12345'


NEW_FILENAME=`curl -v -s http://localhost:8081/v1/boards/test/articles -H "Authorization: bearer $ACCESS_TOKEN" | jq -r '[.data.items[] | select(.title == "[測試] test")][-1] | .filename'`


echo $NEW_FILENAME

curl -v http://localhost:8081/v1/boards/test/articles/$NEW_FILENAME -d 'action=append_comment'  --data-urlencode 'type＝推' -d 'text=test push'



