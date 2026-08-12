#!/bin/bash

# C2-3-3
# 使用 action=add_favorite type=folder title=test 新增資料夾後，再次 GET favorites 時，
# title=test 的 folder 項目數應該增加 1。

if [ "$#" -lt 2 ]; then
	echo "usage: $0 [user_id] [access_token]"
	exit 1
fi

USER_ID="$1"
ACCESS_TOKEN="$2"
URL="http://localhost:8081/v1/users/$USER_ID/favorites"
RESPONSE_FILE="/tmp/ptt-favorite-folder-$$.json"
trap 'rm -f "$RESPONSE_FILE"' EXIT

BEFORE=`curl -s "$URL" -H "Authorization: bearer $ACCESS_TOKEN" | jq '[.data.items[] | select(.type == "folder" and .title == "test")] | length'`

HTTP_STATUS=`curl -s -o "$RESPONSE_FILE" -w '%{http_code}' "$URL" \
	-H "Authorization: bearer $ACCESS_TOKEN" \
	-d 'action=add_favorite' \
	-d 'type=folder' \
	-d 'title=test'`
if [ "$HTTP_STATUS" != "200" ]; then
	echo "add favorite folder failed: HTTP $HTTP_STATUS"
	cat "$RESPONSE_FILE"
	exit 1
fi

AFTER=`curl -s "$URL" -H "Authorization: bearer $ACCESS_TOKEN" | jq '[.data.items[] | select(.type == "folder" and .title == "test")] | length'`
if [ "$AFTER" -ne $((BEFORE + 1)) ]; then
	echo "favorite folder count did not increase: before=$BEFORE after=$AFTER"
	exit 1
fi

echo "favorite folder added: before=$BEFORE after=$AFTER"
