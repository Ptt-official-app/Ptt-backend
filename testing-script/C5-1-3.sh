#!/bin/bash

# C5-1-3
# 在放入 Access Token 的情況下在 GET /v1/boards/test/articles 搜尋推文 recommend_count_gt 10 應該要出現「來測試看看轉錄看板文章」這篇文章


if [ "$#" -lt 0 ]; then
	echo "usage: $0"
	exit -1
fi
ACCESS_TOKEN=`./get_sysop_token.sh`

# curl -s http://localhost:8081/v1/boards/SYSOP/articles -H "Authorization: bearer $ACCESS_TOKEN" -d action=add_article --data-urlencode title=Test --data-urlencode article=Test
curl -s http://localhost:8081/v1/boards/test/articles?recommend_count_gt=10 -H "Authorization: bearer $ACCESS_TOKEN"
