#!/bin/bash

# C3-2-1
# 放入 Access Token 之後， GET /v1/boards/SYSOP/information 可以看到 SYSOP 看板的看板資訊

ACCESS_TOKEN=`./get_sysop_token.sh`
curl -s http://localhost:8081/v1/boards/SYSOP/information -H "Authorization: bearer $ACCESS_TOKEN" 

