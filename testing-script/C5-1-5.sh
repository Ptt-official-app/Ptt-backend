#!/bin/bash

# C5-1-4
# 在放入 Access Token 的情況下 POST /v1/boards/SYSOP/articles 後應該要在看板列表中看到這篇文章


if [ "$#" -lt 0 ]; then
	echo "usage: $0"
	exit -1
fi
ACCESS_TOKEN=`./get_sysop_token.sh`

curl -s http://localhost:8081/v1/boards/SYSOP/articles -H "Authorization: bearer $ACCESS_TOKEN" -d action=add_article --data-urlencode title=中文 --data-urlencode article=中文


curl -s http://localhost:8081/v1/boards/SYSOP/articles -H "Authorization: bearer $ACCESS_TOKEN" 