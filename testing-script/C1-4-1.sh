#!/bin/bash

# C1-4-1
# 在剛註冊完成時 GET /v1/users/{{自己的ID}}/register-form 應該要是空的

curl http://localhost:8081/v1/register -d 'username=user01&password=pass01'
curl http://localhost:8081/v1/users/users01/register-form

