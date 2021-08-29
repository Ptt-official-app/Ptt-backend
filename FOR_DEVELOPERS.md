# 給開發者的資訊

## 歷史緣由跟開發原因

這個專案主要的開發語言是 Golang 以及中文。

原有的 [PTT 程式碼](https://github.com/ptt/pttbbs) 是透過C語言進行開發的，
然而 C 語言開發的程式碼雖然效能高，但是可維護性稍低，以致於後續接手維護不易。

另外原有 C 程式碼編譯結果仰賴平台例如 LP32, LP32, ILP64, LLP64 等狀態，導致原有主機
故障後升級成 64 位元架構主機時可能會無法順利執行上線之問題。

再來目前 C 版本程式碼並沒有把顯示以及商業邏輯的程式碼分離，也就是不符合 MVC 架構，
使得需要修改前端顯示文字也需要一路找到後端才有辦法進行。

較長期的問題還有例如在 C 當中的 IPv6 支援需要另外處理，線程管理不易導致需要大量使用
行程間溝通 (IPC) 方式進行資料處理。這些部分需要整體架構翻新才有機會被改善。

另外是目前存放時間的方式大多是使用 32bit 的儲存單位為主，這會在 2038 年後產生問題。

選擇 Golang 來處理這些問題的理由在於 Golang 在一開始設計之初時就考慮到不同架構
編譯上的狀況，同時可以很簡單的進行跨平台編譯 (Cross compiling) ，因此未來系統要
在 Arm based 的平台上面或者是在 Windows 上面運行都是有機會的。

另一方面是 Golang 原生支援 UTF-8 編碼，因此可以透過此次大改版逐步將 Big5-UAO
編碼的資料進行轉換。

## 整體架構部分

目前整體而言會以 HTTP 的 [RESTful API](https://zh.wikipedia.org/wiki/%E8%A1%A8%E7%8E%B0%E5%B1%82%E7%8A%B6%E6%80%81%E8%BD%AC%E6%8D%A2) 作為對外設計的介面，未來也許會支援相容 VT100
（也就是傳統可以用 PCMAN 登入的 BBS 介面） 的SSH連線方式。

使用 [RESTful API](https://zh.wikipedia.org/wiki/%E8%A1%A8%E7%8E%B0%E5%B1%82%E7%8A%B6%E6%80%81%E8%BD%AC%E6%8D%A2) 
的主因在於便於其他開發者進行開發，HTTP 有許多現成的客戶端以及測試工具，在各種語言以及平台中皆有範例程式碼以及函式庫可以呼叫。
另外 HTTPS 也經過了數年來在資訊安全上的驗證，同時我們可以透過 HTTP 既有的快取機制設計來進一步節省維運伺服器所需要的流量。

這個專案採用 Golang 原生的 HTTP 解譯器不另外使用其他框架 (framework) 以降低未來框架修改時的維護成本。

在呼叫流程中如下列所示： （需要補圖）

```
==== 外面 （公有 IP 或是受防火牆保護的網段） =====

針對 IP 判斷流量以及流速決定是否發回 429 Too Many Request

檢查傳入的 Access Token 是否合法

解開傳入的 Access Token, 取出使用者必要資訊（例如 User ID, User Level）

透過使用者必要資訊以及 etag 和 if-not-modified 資訊來決定是否送回 304 Not Modified

檢查存取的項目是否為熱門項目（熱門文章或是看板），回傳資料

透過 Cache 模組 (Shared memory) 檢查資料是否存在 Shared memory 上，確認是否應該回傳 Shared memory 資料優先

透過 go-bbs 模組讀寫現有 BBS 檔案架構

將處理後的資料進行快取儲存後回傳資料
```


而較為即時性的部分會另外採用 [Event Stream](https://developer.mozilla.org/zh-TW/docs/Web/API/Server-sent_events/Using_server-sent_events) 的方式進行設計。使用 Event Stream 的主因為使用 WebSocket 的話開發者會需要另外維護一份 WebSocket 的函式庫。

以水球為例，在呼叫流程中如下列所示：（需要補圖）
```
====外面（公有IP或是受防火牆保護的網段）=====

針對 IP 判斷流量以及流速決定是否發回 429 Too Many Request

檢查傳入的 Access Token 是否合法

解開傳入的 Access Token, 取出使用者必要資訊（例如 User ID, User Level）

訂閱指定的 Go Channel

有人丟水球過來時回傳給使用者水球（不關閉連線）

有人丟第二顆水球時回傳第二個訊息給使用者

```

## 資料夾結構

主要的程式碼包含程式進入點會放在根目錄

其他功能以功能為分類分在不同資料夾

目前暫定的功能為下：
```
後端預定模組

流量控制
權限管控
	登入 CRL

後台版主部分

水球部分

寄信模組

認證模組
	TWID 等認證模組
		https://github.com/ptt/pttbbs/issues/65
	手機認證模組
		https://github.com/ptt/pttbbs/issues/63



視訊部分
推播模組
看板文章即時訂閱
看板文章 RSS 訂閱

系統狀態統計部分
	產生可讓 rrdtool 使用的資料或是直接寫入 rrd?

Open ID Provider


IPv6 Support
	https://github.com/ptt/pttbbs/issues/58

ASCII 圖片產圖引擎
```

模組的載入以及停用儘可能設計在一到三個註解就可以完成，此設計可以避免模組使用的套件未來發生問題時造成無法編譯的狀況。

請儘量避免使用 `unsafe` 以及 cgo 除非你確定你提供的足夠多的文件以及測試案例。

測試案例盡量避免只有一行 reflect.DeepEqual，可以使用 reflect.DeepEqual，
但是如果能夠多比對其他欄位輸出更多錯誤訊息會對未來的開發者更有幫助。

---

## 開發環境建置

安裝下列的應用程式來建構開發環境：

| 應用程式名稱                                               | 應用程式版本(有特定版本才填寫) | 安裝要求                                   |
| -------                                                    | -------                        | -------                                    |
| [Golang](https://golang.org/dl/)                           | 1.16 以上                      | 必要                                       |
| [GoLand](https://www.jetbrains.com/go/promo/)              |                                | 如果使用 GoLand * 推薦給新手               |
| [Sublime Text 3](https://classic.yarnpkg.com/zh-Hant/)     |                                | 如果使用 Sublime Text，記得安裝 Gofmt 套件 |
| [docker engine](https://docs.docker.com/engine/install/)   | 1.13.0+                        |                                            |
| [docker compose](https://docs.docker.com/compose/install/) | 1.10.0+                        | 使用 `docker compose` 執行本專案時         |

### Sublime Text 3 的套件

要使用 Sublime Text 3，請安裝以下套件。

| 套件                                                                                                                                                     | 安裝要求 |
| -------                                                                                                                                                  | -------  |
| [Gofmt](https://packagecontrol.io/packages/Gofmt)                                                                                                        | 必須     |
| [ConvertToUTF8](https://ephrain.net/sublime-text-%E8%AE%93-sublime-text-%E6%94%AF%E6%8F%B4-big5-%E7%B7%A8%E7%A2%BC%E7%9A%84%E6%96%87%E5%AD%97%E6%AA%94/) | 可選     |

---

## 安裝與設定

### 使用 docker-compose

* 請使用 1.10.0 以上的版本，可以使用 `docker-compose version` 來確認

```bash
$ docker-compose up
```

### 使用 docker

```bash
$ docker build -t ptt-backend
$ docker run -p 8081:8081 --name ptt-backend
```

### 手動安裝

#### Clone Ptt-backend 專案

```bash
$ git clone https://github.com/Ptt-official-app/Ptt-backend.git
```

#### 下載測試伺服器靜態資料與配置設定檔
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

- 設定檔預設讀取 `conf/config_default.toml`，如果希望改成自己的設定檔請直接改在 `config.toml` 中即可。

## 編譯與執行

請在工作目錄 `./Ptt-backend` 下執行此命令。

### 狀況 1: 直接使用 `go build`

```bash
$ ./make.bash build
$ ./Ptt-backend
```

如果是 windows 的話

```bat
$ .\make.bat build
$ .\Ptt-backend.exe
```

#### Troubleshoot

// TODO

### 狀況 2: 使用 Gitpod 開發

// TODO

---

## 生產環境/其他環境的判定

一但要上線到生產環境，請務必改掉 security 章節內的設定。

---

## Deploy 到 Staging 環境以及正式環境的方法

當下表左欄所列的分支更新後，分支和網站將會自動被更新。

| 分支 | 建置與更新的分支 | 對應站點 |
| ---- | ---- | ---- |
|`master` (尚未建置) |||
|`staging` (尚未建置)||
|`development`||https://pttapp.cc|

---

## 分支規則

只允許推送 Pull Request 到 `development` 。
在推送 Pull Request 時，請依照以下命名規則為您的分支命名

| 變更種類 | 分支的命名規則 |
| ---- | ---- |
|新增功能|`feature/#{ISSUE_ID}-#{branch_title_name}`|
|Hotfix commit|`hotfix/#{ISSUE_ID}-#{branch_title_name}`|

### 基本分支


| 目的 | 分支 | 預覽用 URL | 誰可以發 Pull Request | 備註 |
| ---- | ---- | ---- | ---- | ---- |
| 開發 | development | https://pttapp.cc | All developers | 基本上請推送 Pull Request 到這裡 |
| 正式版預覽 | staging |  | Only administrators | 對於正式版釋出前的最終確認，禁止管理員以外的人推送 Pull Request。 |
| 正式版 | master |  | Only administrators | 禁止管理員以外的人推送 Pull Request |


### 系統所使用的分支

| 目的 | 分支 | 預覽用 URL | 備註 |
| ---- | -------- | ---- | ---- |

### 基本規則

1. 不接受 go test  不通過的 pr
2. 不接受 route test 沒對 responseData assert 
3. 不接受跟 issue 無相關的 pr，可以先開 issue 再來討論
4. 不接受 golang-ci 不通過的 pr

### 開發流程

1. 在 slack 或者 github 上認領 issue
2. [fork](https://guides.github.com/activities/forking/) 專案
3. 參照他人開立 branch
4. 在 issue 底下訂下完成時間，完成時間訂在何時可以 slack 討論
5. 進行開發
6. 開發完後推送 [pull request](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/creating-a-pull-request)
7. pr assign 給你自己
8. 增加 reviewer


### Test Sheet 腳本

在 testing-script 裡面有可以簡單測試 Google Sheet 版本 Test Sheet 的腳本，不過目前只支援 Linux 或是 macOS 的系統 (是 bash 腳本)，歡迎多加利用。

### Q&A

* 有測試相關問題想請教怎辦? 

slack 找 陳昱廷


* 第一次上手搞不懂要怎麼寫

Ptt-backend 的部分通常是實作 delivery/usecase/repository/ 
- internal/repository/
- internal/usecase/
- internal/delivery/

以及這三個 layer 的測試。
如果 go-bbs 還沒實作可以先做 mock 假資料

* 我會寫程式，但是 query 找不到

        可以問一下，檔案儲存的方式有些不同，自己要找可能要找很久

* 開發到一半發現有地方還沒實作怎麼辦

        打 remark // TODO:123
        或是 // FIXME:123
