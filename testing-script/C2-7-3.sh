#!/bin/bash

# C2-7-3
# 放入非 SYSOP 的 Access Token 之後， POST /v1/users/SYSOP/drafts/0 不應該成功，應該要看到 error

if [ "$#" -lt 1 ]; then
	echo "usage: $0 [access_token]"
	exit -1
fi

curl http://localhost:8081/v1/users/SYSOP/drafts/0 -H "Authorization: bearer $1" 
