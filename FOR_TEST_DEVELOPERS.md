# 給測試開發者的資訊

## 分支規則

只允許推送 Pull Request 到 `development_enable_lint` 。
在推送 Pull Request 時，請依照以下命名規則為您的分支命名

| 變更種類 | 分支的命名規則 |
| ---- | ---- |
|測試功能|`test/#{ISSUE_ID}-#{branch_title_name}`|

### 基本規則

1. 不接受 go test 不通過的 pr
2. 不接受 route test 沒對 responseData assert 
3. 不接受跟 issue 無相關的測試，可以先開 issue 再來討論
4. 暫時接受 lint 不通過的測試，2月過後會改為需要通過 lint

### 開發流程

1. 在 slack 或者 github 上認領 issue(ex:https://github.com/Ptt-official-app/Ptt-backend/issues/73) 
2. fork 專案
3. 研究 test/#45-add_route_token_test 做法
4. 在 issue 底下訂下完成時間，完成時間訂在何時可以 slack 討論
5. 進行開發
6. 開發完後推送 pr
7. pr assign 給 y2468101216 跟你自己

### Q&A

1. Q:有測試相關問題想請教怎辦? A: slack 找 陳昱廷
2. 待補