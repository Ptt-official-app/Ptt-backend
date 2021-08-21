#!/bin/bash

# C3-1-1
# 放入 Access Token 之後， GET /v1/boards 可以看到 SYSOP 板


ACCESS_TOKEN=`./get_sysop_token.sh`
curl -s http://localhost:8081/v1/boards -H "Authorization: bearer $ACCESS_TOKEN" | jq '.data[] | select(.id=="SYSOP")' 

