#!/bin/bash

# C3-4-2
# 放入 Access Token 之後， GET /v1/classes/-1 不應返回東西，應該出現 error

ACCESS_TOKEN=`./get_sysop_token.sh`
curl -s http://localhost:8081/v1/classes/-1 -H "Authorization: bearer $ACCESS_TOKEN" 

