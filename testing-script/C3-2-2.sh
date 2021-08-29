#!/bin/bash

# C3-2-2
# 沒有放入Access Token 的時候，GET /v1/boards/SYSOP/information 不應該看到看板資訊，要回傳 error

curl -s http://localhost:8081/v1/boards/SYSOP/information 