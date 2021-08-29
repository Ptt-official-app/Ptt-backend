#!/bin/bash

# C3-4-1
# 放入 Access Token 之後， GET /v1/classes/1 要返回分類主目錄

ACCESS_TOKEN=`./get_sysop_token.sh`
curl -s http://localhost:8081/v1/classes/2 -H "Authorization: bearer $ACCESS_TOKEN" 

