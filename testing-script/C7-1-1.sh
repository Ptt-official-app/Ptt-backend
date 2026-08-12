#!/bin/bash

# C7-1-1
# 在不放入 Access Token 的情況下 GET /v1/popular-articles 應該出現熱門文章列表，上限為 100 個文章

curl -s http://localhost:8081/v1/popular-articles
echo ""
echo "熱門文章上限為 100 篇"
curl -s http://localhost:8081/v1/popular-articles | jq '.data.items | length'
