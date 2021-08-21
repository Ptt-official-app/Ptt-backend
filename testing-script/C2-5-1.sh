#!/bin/bash

# C2-5-1
# 不放入 Access Token 的情況下 GET /v1/users/SYSOP/articles 不應該看到任何文章，應該要看到 error 元素

curl http://localhost:8081/v1/users/SYSOP/articles 
