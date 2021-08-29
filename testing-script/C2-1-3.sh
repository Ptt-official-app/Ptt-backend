#!/bin/bash

# C2-1-3
# 不放入 Access Token 的時候 GET /v1/users/SYSOP/information 不應該取得任何人的資料，回應的 JSON 要有 error 欄位


curl http://localhost:8081/v1/users/SYSOP/information 

