#!/bin/bash

# C3-1-2
# 以 SYSOP 帳號登入之後，POST /v1/boards 新增看板，並確認不需重啟即可在列表中看到。

ACCESS_TOKEN=$(./get_sysop_token.sh)
BOARD_ID="testboard01"
TITLE="TestBoard"

curl -s http://localhost:8081/v1/boards \
	-H "Authorization: bearer $ACCESS_TOKEN" \
	-H "Content-Type: application/x-www-form-urlencoded" \
	--data-urlencode "board_id=${BOARD_ID}" \
	--data-urlencode "title=${TITLE}"
echo ""

curl -s http://localhost:8081/v1/boards \
	-H "Authorization: bearer $ACCESS_TOKEN" |
	jq --arg board_id "$BOARD_ID" '.data[] | select(.id == $board_id)'
