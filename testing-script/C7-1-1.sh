#!/bin/bash

# C7-1-1
# 在不放入Access Token 的情況下 GET /v1/popular-articles 應該出現熱門文章列表，上限為 100 個文章

ACCESS_TOKEN=`./get_sysop_token.sh`
curl -s http://localhost:8081/v1/popular-articles -H "Authorization: bearer $ACCESS_TOKEN" 
echo ""
echo "上限為 100 個看板"
curl -s http://localhost:8081/v1/popular-articles -H "Authorization: bearer $ACCESS_TOKEN" | jq '.data[] | length '

