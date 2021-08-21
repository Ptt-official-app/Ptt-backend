#!/bin/bash

# C3-2-3
# 成功取回 /v1/boards/ptt_app/information 後，應該要在moderators 看到 SYSOP 為版主

ACCESS_TOKEN=`./get_sysop_token.sh`
curl -s http://localhost:8081/v1/boards/ptt_app/information -H "Authorization: bearer $ACCESS_TOKEN" 

