#!/bin/bash

# C1-1-2
# 輸入錯誤的帳號密碼 (SYSOP / 1231234) 登入系統不應該取得 access token 參數，同時要出現錯誤訊息

curl http://localhost:8081/v1/token -d 'grant_type=password&username=SYSOP&password=1231234'

