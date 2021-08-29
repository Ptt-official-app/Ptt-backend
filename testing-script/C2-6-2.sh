#!/bin/bash

# C2-6-2
# 放入 Access Token 之後， GET /v1/users/SYSOP/comments 應該看到 SYSOP 曾經留言過的文章


ACCESS_TOKEN=`./get_sysop_token.sh`
curl http://localhost:8081/v1/users/SYSOP/comments -H "Authorization: bearer $ACCESS_TOKEN"

