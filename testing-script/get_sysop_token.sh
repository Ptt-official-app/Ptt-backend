#!/bin/bash

# C1-1-1
# 輸入正確的帳號密碼 (SYSOP / 123123) 以 /v1/token 登入系統應該要取得 access token 參數

curl http://localhost:8081/v1/token  -q -d 'grant_type=password&username=SYSOP&password=123123' | jq -r '.access_token'

