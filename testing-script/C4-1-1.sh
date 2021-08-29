#!/bin/bash

# C4-1-1
# 在不放入Access Token 的情況下 GET /v1/popular-boards 應該出現熱門看板列表，上限為 100 個看板

ACCESS_TOKEN=`./get_sysop_token.sh`
curl -s http://localhost:8081/v1/popular-boards -H "Authorization: bearer $ACCESS_TOKEN" 
echo ""
echo "上限為 100 個看板"
curl -s http://localhost:8081/v1/popular-boards -H "Authorization: bearer $ACCESS_TOKEN" | jq '.data[] | length '

