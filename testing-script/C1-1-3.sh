#!/bin/bash

# C1-1-3
# 未輸入帳號或密碼登入系統，不應取得 access token 參數，同時要出現錯誤訊息

curl http://localhost:8081/v1/token -d 'grant_type=password&username=SYSOP&password='
curl http://localhost:8081/v1/token -d 'grant_type=password&username=&password=123123'

