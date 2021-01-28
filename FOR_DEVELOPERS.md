# 給開發者的資訊

## 1. 開發環境建置

安裝下列的應用程式來建構開發環境：

| 應用程式名稱 | 應用程式版本(有特定版本才填寫) | 安裝要求 |
| ------- | ------- | ------- |
|[Golang](https://golang.org/dl/)|1.14 以上|必要|
|[GoLand](https://www.jetbrains.com/go/promo/)| |如果使用 GoLand * 推薦給新手|
|[Sublime Text 3](https://classic.yarnpkg.com/zh-Hant/)| |如果使用 Sublime Text，記得安裝 Gofmt 套件|
|[docker compose](https://docs.docker.com/compose/install/)| |使用 `docker compose` 直襲本專案時|

### 1-1. Sublime Text 3 的套件

要使用 Sublime Text 3，請安裝以下套件。

| 套件 | 安裝要求 |
| ------- | ------- |
|[Gofmt](https://packagecontrol.io/packages/Gofmt)|可選|
|[ConvertToUTF8](https://ephrain.net/sublime-text-%E8%AE%93-sublime-text-%E6%94%AF%E6%8F%B4-big5-%E7%B7%A8%E7%A2%BC%E7%9A%84%E6%96%87%E5%AD%97%E6%AA%94/)|可選|

---

## 2. 下載測試伺服器靜態資料

```bash
# 下載 BBS Home
$ wget http://pttapp.cc/data-archives/bbs_backup_lastest.tar.xz
# 用 xz 進行解壓縮
$ tar -Jxvf bbs_backup_lastest.tar.xz

# 下載 SHM 測試資料
$ wget http://pttapp.cc/data-archives/dump.shm.lastest.tar.bz2
# 用 bzip2 進行解壓縮
$ tar -jxvf dump.shm.lastest.tar.bz2
```


## 3. 執行此專案

請在工作目錄(./Ptt-backend)下執行此命令。

### 3-1. 直接使用 `go build` 的狀況

#### 3-1-1. 安裝以及編譯

```bash
# 編譯
$ go build
```

#### 3-1-2. 執行

```bash
# 編譯
$ ./Ptt-backend
```

#### 3-1-3. 提示

- 設定檔預設讀取 config_default.toml，如果希望改成自己的設定檔請將他複製成 config.toml 即可。

#### 3-1-4. Troubleshoot

// TODO

### 3-2. 使用 Gitpod 開發的狀況

// TODO

---

## 4. 生產環境/其他環境的判定

一但要上線到生產環境，請務必改掉 security 章節內的設定。

---

## 5. Deploy 到 Staging 環境以及正式環境的方法

當下表左欄所列的分支更新後，分支和網站將會自動被更新。

| 分支 | 建置與更新的分支 | 對應站點 |
| ---- | ---- | ---- |
|`master` (尚未建置) |||
|`staging` (尚未建置)||
|`development`||https://pttapp.cc|

---

## 6. 分支規則

只允許推送 Pull Request 到 `development` 。
在推送 Pull Request 時，請依照以下命名規則為您的分支命名

| 變更種類 | 分支的命名規則 |
| ---- | ---- |
|新增功能|`feature/#{ISSUE_ID}-#{branch_title_name}`|
|Hotfix commit|`hotfix/#{ISSUE_ID}-#{branch_title_name}`|

### 6-1. 基本分支


| 目的 | 分支 | 預覽用 URL | 誰可以發 Pull Request | 備註 |
| ---- | ---- | ---- | ---- | ---- |
| 開發 | development | https://pttapp.cc | All developers | 基本上請推送 Pull Request 到這裡 |
| 正式版預覽 | staging |  | Only administrators | 對於正式版釋出前的最終確認，禁止管理員以外的人推送 Pull Request。 |
| 正式版 | master |  | Only administrators | 禁止管理員以外的人推送 Pull Request |


### 7-2. 系統所使用的分支

| 目的 | 分支 | 預覽用 URL | 備註 |
| ---- | -------- | ---- | ---- |
