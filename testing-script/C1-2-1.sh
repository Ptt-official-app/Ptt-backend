#!/bin/bash

# C1-2-1
# 在沒有 user01 的狀況下輸入 POST /v1/register 參數 username=user01 password=pass01 要能註冊成功並且取得 Access Token

curl http://localhost:8081/v1/register -d 'username=user01&password=pass01'

