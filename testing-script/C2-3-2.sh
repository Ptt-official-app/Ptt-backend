#!/bin/bash

# C2-3-2
# 使用 action=add_favorite type=line 新增分隔線後，再次 GET favorites 時，
# type=line 的項目數應該增加 1。

if [ "$#" -lt 2 ]; then
	echo "usage: $0 [user_id] [access_token]"
	exit 1
fi

USER_ID="$1"
ACCESS_TOKEN="$2"
URL="http://localhost:8081/v1/users/$USER_ID/favorites"
RESPONSE_FILE="/tmp/ptt-favorite-line-$$.json"
trap 'rm -f "$RESPONSE_FILE"' EXIT

BEFORE=`curl -s "$URL" -H "Authorization: bearer $ACCESS_TOKEN" | jq '[.data.items[] | select(.type == "line")] | length'`

HTTP_STATUS=`curl -s -o "$RESPONSE_FILE" -w '%{http_code}' "$URL" \
	-H "Authorization: bearer $ACCESS_TOKEN" \
	-d 'action=add_favorite' \
	-d 'type=line'`
if [ "$HTTP_STATUS" != "200" ]; then
	echo "add favorite line failed: HTTP $HTTP_STATUS"
	cat "$RESPONSE_FILE"
	exit 1
fi

AFTER=`curl -s "$URL" -H "Authorization: bearer $ACCESS_TOKEN" | jq '[.data.items[] | select(.type == "line")] | length'`
if [ "$AFTER" -ne $((BEFORE + 1)) ]; then
	echo "favorite line count did not increase: before=$BEFORE after=$AFTER"
	exit 1
fi

echo "favorite line added: before=$BEFORE after=$AFTER"
