#!/bin/bash

# C1-3-1
# 在沒有註冊過 user01 的狀況下輸入 GET /v1/register/precheck?type=username&value=user01 應該要回應 available

curl http://localhost:8081/v1/register/precheck?type=username&value=user01