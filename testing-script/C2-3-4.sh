#!/bin/bash

# C2-3-4
# 使用 action=add_favorite type=board board_id=SYSOP 新增看板後，再次 GET favorites 時，
# items 應該包含 type=board 且 board_id=SYSOP 的項目。

if [ "$#" -lt 2 ]; then
	echo "usage: $0 [user_id] [access_token]"
	exit 1
fi

USER_ID="$1"
ACCESS_TOKEN="$2"
URL="http://localhost:8081/v1/users/$USER_ID/favorites"
RESPONSE_FILE="/tmp/ptt-favorite-board-$$.json"
trap 'rm -f "$RESPONSE_FILE"' EXIT

HTTP_STATUS=`curl -s -o "$RESPONSE_FILE" -w '%{http_code}' "$URL" \
	-H "Authorization: bearer $ACCESS_TOKEN" \
	-d 'action=add_favorite' \
	-d 'type=board' \
	-d 'board_id=SYSOP'`
if [ "$HTTP_STATUS" != "200" ]; then
	echo "add favorite board failed: HTTP $HTTP_STATUS"
	cat "$RESPONSE_FILE"
	exit 1
fi

COUNT=`curl -s "$URL" -H "Authorization: bearer $ACCESS_TOKEN" | jq '[.data.items[] | select(.type == "board" and .board_id == "SYSOP")] | length'`
if [ "$COUNT" -lt 1 ]; then
	echo "favorite board SYSOP not found after addition"
	exit 1
fi

echo "favorite board SYSOP added"
