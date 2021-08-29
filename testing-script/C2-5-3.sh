#!/bin/bash

# C2-5-2
# 放入 Access Token 之後， GET /v1/users/SYSOP/articles 應該看到 SYSOP 曾經發過的文章


ACCESS_TOKEN=`./get_sysop_token.sh`
curl http://localhost:8081/v1/users/SYSOP/articles -H "Authorization: bearer $ACCESS_TOKEN"

