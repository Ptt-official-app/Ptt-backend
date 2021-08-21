#!/bin/bash

# C5-1-1
# 在不放入 Access Token 的情況下 GET /v1/boards/SYSOP/articles 不應該出現文章列表，應該出現 error

curl -s http://localhost:8081/v1/boards/SYSOP/articles 