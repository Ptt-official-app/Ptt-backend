#!/bin/bash

# C6-1-3
# 以 SYSOP 登入後於 test 看板發出一篇文章，接下來透過同一使用者刪除該篇文章。
# SAFE_ARTICLE_DELETE 會保留 .DIR tombstone，文章列表中應出現「(本文已被刪除)」。

ACCESS_TOKEN=`./get_sysop_token.sh`
curl -v -s http://localhost:8081/v1/boards/test/articles -H "Authorization: bearer $ACCESS_TOKEN" -d 'action=add_article' --data-urlencode 'title=[測試] test' --data-urlencode 'article=12345'

NEW_FILENAME=`curl -s http://localhost:8081/v1/boards/test/articles -H "Authorization: bearer $ACCESS_TOKEN" | jq -r '[.data.items[] | select(.title == "[測試] test")][-1] | .filename'`
echo "add success: $NEW_FILENAME"

HTTP_STATUS=`curl -s -o /tmp/ptt-delete-article-response.json -w '%{http_code}' http://localhost:8081/v1/boards/test/articles/$NEW_FILENAME -H "Authorization: bearer $ACCESS_TOKEN" -d 'action=delete'`
if [ "$HTTP_STATUS" != "200" ]; then
  echo "delete failed: HTTP $HTTP_STATUS"
  cat /tmp/ptt-delete-article-response.json
  exit 1
fi

DELETED_TITLE=`curl -s http://localhost:8081/v1/boards/test/articles -H "Authorization: bearer $ACCESS_TOKEN" | jq -r '[.data.items[] | select(.title | startswith("(本文已被刪除)"))][-1] | .title'`
if [[ "$DELETED_TITLE" != "(本文已被刪除)"* ]]; then
  echo "delete tombstone not found"
  exit 1
fi

echo "delete success: $DELETED_TITLE"
