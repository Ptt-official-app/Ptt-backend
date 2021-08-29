#!/bin/bash

# C3-1-2
# 以 SYSOP 帳號登入之後， POST /v1/boards board_id=testboard01 title=TestBoard 後取得看板列表可以看到 testboard01 看板


ACCESS_TOKEN=`./get_sysop_token.sh`
result=`curl -s http://localhost:8081/v1/boards -H "Authorization: bearer $ACCESS_TOKEN"`
echo "SYSOP"
echo $result | jq '.data[] | select(.id=="SYSOP")' 
echo "ptt_app"
echo $result | jq '.data[] | select(.id=="ptt_app")' 

