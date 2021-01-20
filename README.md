# PTT APP 後端

這個專案主要的開發語言是 Golang 以及中文。

原有的 [PTT 程式碼](https://github.com/ptt/pttbbs)是透過C語言進行開發的，
然而C語言開發的程式碼雖然效能高，但是可維護性稍低，以致於後續接手維護不易。

另外原有C程式碼編譯結果仰賴平台例如 LP32, LP32, ILP64, LLP64 等狀態，導致原有主機
故障後升級成 64 位元架構主機時可能會無法順利執行上線之問題。

再來目前C版本程式碼並沒有把顯示以及商業邏輯的程式碼分離，也就是不符合 MVC 架構，
使得需要修改前端顯示文字也需要一路找到後端才有辦法進行。

較長期的問題還有例如在C當中的IPv6支援需要另外處理，線程管理不易導致需要大量使用
行程間溝通 (IPC) 方式進行資料處理。這些部分需要整體架構翻新才有機會被改善。

另外是目前存放時間的方式大多是使用32bit的儲存單位為主，這會在2038年後產生問題。

選擇 Golang 來處理這些問題的理由在於 Golang 在一開始設計之初時就考慮到不同架構
編譯上的狀況，同時可以很簡單的進行跨平台編譯 (Cross compiling) ，因此未來系統要
在 Arm based 的平台上面或者是在 Windows 上面運行都是有機會的。

另一方面是 Golang 原生支援 UTF-8 編碼，因此可以透過此次大改版逐步將 Big5-UAO
編碼的資料進行轉換。

## 整體架構部分

目前整體而言會以 HTTP 的 [RESTful API](https://zh.wikipedia.org/wiki/%E8%A1%A8%E7%8E%B0%E5%B1%82%E7%8A%B6%E6%80%81%E8%BD%AC%E6%8D%A2) 作為對外設計的介面，未來也許會支援相容VT100
（也就是傳統可以用PCMAN登入的BBS介面）的SSH連線方式。

使用 [RESTful API](https://zh.wikipedia.org/wiki/%E8%A1%A8%E7%8E%B0%E5%B1%82%E7%8A%B6%E6%80%81%E8%BD%AC%E6%8D%A2) 
的主因在於便於其他開發者進行開發，HTTP有許多現成的客戶端以及測試工具，在各種語言以及平台中皆有範例程式碼以及函式庫可以呼叫。
另外HTTPS也經過了數年來在資訊安全上的驗證，同時我們可以透過HTTP既有的快取機制設計來進一步節省維運伺服器所需要的流量。

這個專案採用 Golang 原生的 HTTP 解譯器不另外使用其他框架(framework)以降低未來框架修改時的維護成本。

在呼叫流程中如下列所示： （需要補圖）

```
====外面（公有IP或是受防火牆保護的網段）=====

針對 IP 判斷流量以及流速決定是否發回 429 Too Many Request

檢查傳入的 Access Token 是否合法

解開傳入的 Access Token, 取出使用者必要資訊（例如User ID, User Level）

透過使用者必要資訊以及etag和if-not-modified資訊來決定是否送回304 Not Modified

檢查存取的項目是否為熱門項目（熱門文章或是看板），回傳資料

透過Cache模組(Shared memory)檢查資料是否存在Shared memory上，確認是否應該回傳Shared memory資料優先

透過 go-bbs 模組讀寫現有BBS檔案架構

將處理後的資料進行快取儲存後回傳資料
```


而較為即時性的部分會另外採用 [Event Stream](https://developer.mozilla.org/zh-TW/docs/Web/API/Server-sent_events/Using_server-sent_events) 的方式進行設計。使用 Event Stream 的主因為使用WebSocke的話開發者會需要另外維護一份WebSocket的函式庫。

以水球為例，在呼叫流程中如下列所示：（需要補圖）
```
====外面（公有IP或是受防火牆保護的網段）=====

針對 IP 判斷流量以及流速決定是否發回 429 Too Many Request

檢查傳入的 Access Token 是否合法

解開傳入的 Access Token, 取出使用者必要資訊（例如User ID, User Level）

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
	登入CRL

後台版主部分

水球部分

寄信模組

認證模組
	TWID等認證模組
		https://github.com/ptt/pttbbs/issues/65
	手機認證模組
		https://github.com/ptt/pttbbs/issues/63



視訊部分
推播模組
看板文章即時訂閱
看板文章RSS訂閱

系統狀態統計部分
	產生可讓rrdtool使用的資料或是直接寫入rrd?

Open ID Provider


IPv6 Support
	https://github.com/ptt/pttbbs/issues/58

ASCII 圖片產圖引擎
```

模組的載入以及停用儘可能設計在一到三個註解就可以完成，此設計可以避免模組使用的套件未來發生問題時造成無法編譯的狀況。

請儘量避免使用 `unsafe` 以及 cgo 除非你確定你提供的足夠多的文件以及測試案例。

測試案例盡量避免只有一行reflect.DeepEqual，可以使用 reflect.DeepEqual，
但是如果能夠多比對其他欄位輸出更多錯誤訊息會對未來的開發者更有幫助。


## 文件檔案部分

* [RESTful API 文件](https://docs.google.com/document/d/18DsZOyrlr5BIl2kKxZH7P2QxFLG02xL2SO0PzVHVY3k/edit?usp=sharing)
* [go-bbs Package](https://github.com/PichuChen/go-bbs)


