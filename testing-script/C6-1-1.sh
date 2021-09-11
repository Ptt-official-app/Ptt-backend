#!/bin/bash

# C6-1-1
# 登入後進入 test 看板，發出一篇標題為「[測試] test」內文為 12345 的文章 ，要能夠順利發文



ACCESS_TOKEN=`./get_sysop_token.sh`
curl -v -s http://localhost:8081/v1/boards/test/articles -H "Authorization: bearer $ACCESS_TOKEN" -d 'action=add_article' --data-urlencode 'title=[測試] test0904' --data-urlencode 'article=1234中文5'


# curl -v -s http://localhost:8081/v1/boards/test/articles -H "Authorization: bearer $ACCESS_TOKEN" 
# curl -v -s http://localhost:8081/v1/boards/test/articles/M.1630766114.A.4BB -H "Authorization: bearer $ACCESS_TOKEN" 