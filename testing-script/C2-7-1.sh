#!/bin/bash

# C2-7-1
# 不放入 Access Token 的情況下 GET /v1/users/SYSOP/drafts/0 不應該看到任何資料，應該要看到 error

curl http://localhost:8081/v1/users/SYSOP/drafts/0 
