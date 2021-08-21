#!/bin/bash

# C5-1-2
# 在放入一般使用者的 Access Token 的情況下 GET /v1/boards/bm_only/articles (例如 SECURITY )不應該出現文章列表，因為這是版主以上才能看到的版，應該出現 error


if [ "$#" -lt 1 ]; then
	echo "usage: $0 [access_token]"
	exit -1
fi

curl -s http://localhost:8081/v1/boards/SECURITY/articles 