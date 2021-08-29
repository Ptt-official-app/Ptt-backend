# PTT APP 後端

## 這專案是做什麼的，跟什麼有關係？

Ptt-backend 是做 [BBS (Bulletin Board System)](https://zh.wikipedia.org/wiki/BBS) 的開源專案，
但 BBS 有很多種，[Ptt](https://www.ptt.cc/index.html) 只是 BBS 的其中一種。

[go-bbs](https://github.com/Ptt-official-app/go-bbs) 是與 Ptt-backend 有關的專案，解析 BBS 儲存在硬碟（或是共享記憶體）中的資料提供給 Ptt-backend 使用。

另外我們有個[測試站](https://pttapp.cc/)，可以透過這個測試站進行測試，你也可以在上面註冊新帳號。

請注意，**密碼以及資料請勿使用真實資料**，該測試站上的資料會被用作開發測試用途，因此任何人都可以拿到上面[雜湊 (Hash)](https://zh.wikipedia.org/wiki/%E6%95%A3%E5%88%97%E5%87%BD%E6%95%B8) 過後的密碼以及資料。

## 如何參與這個專案？

1. 要有個 [GitHub 帳號](https://github.com/join)
2. 參與 [issue 討論](https://github.com/Ptt-official-app/Ptt-backend/issues)

想要更加地投入可以參考下面的 [*貢獻者相關*](https://github.com/Ptt-official-app/Ptt-backend#%E8%B2%A2%E7%8D%BB%E8%80%85%E7%9B%B8%E9%97%9C) 跟 [*Q&A*](https://github.com/Ptt-official-app/Ptt-backend#qa)

## 如何編譯執行？

請參考[「給開發者」 編譯與執行章節](https://github.com/Ptt-official-app/Ptt-backend/blob/development/FOR_DEVELOPERS.md#%E7%B7%A8%E8%AD%AF%E8%88%87%E5%9F%B7%E8%A1%8C)

## 文件檔案部分？

請參閱 [RESTful API 文件](https://docs.google.com/document/d/18DsZOyrlr5BIl2kKxZH7P2QxFLG02xL2SO0PzVHVY3k/edit?usp=sharing)

另外我們也歡迎協助幫忙維護 Swagger 的志工。

## 軟體測試測試表：

測試表目前以 Google Sheet 寫成，這份測試表適合讓對於不熟悉程式語言的開發者了解目前各項功能的狀態。

測試表母表： [連結](https://docs.google.com/spreadsheets/d/1uo4AJuSi5xTXEht2o2EHogLivCJlJqlLaeqoj1RceDY/edit?usp=sharing)

歷次測試結果：

* [2021060401](https://docs.google.com/spreadsheets/d/1dyfmWZJaTiDrSMIFZ6ynmfeTWU3h45ScW1eyKFJU494/edit?usp=sharing)
* [2021062101](https://docs.google.com/spreadsheets/d/1RGIQPN6KfiCzWRQe-BeTXFUTvdwytOv0Yhd5a8ixCLk/edit?usp=sharing)
* [2021070501](https://docs.google.com/spreadsheets/d/1thxyY9jf2GkK3DMgGO1bAHWE3BmsFb4p5Ot5eAgbaII/edit?usp=sharing)
* [2021072001](https://docs.google.com/spreadsheets/d/1Dv0eZNTLU_NKiehR15qp8iPXAXZ85j2tWj_7h1rxUKg/edit?usp=sharing)
* [2021082101](https://docs.google.com/spreadsheets/d/1p-HTUe-x-6CMVUjmLXty8hPJuoNL_kn8YsCZlJrvRzM/edit?usp=sharing)

## 貢獻者相關：

如果想要參與開發與維護專案可參考: [給開發者](https://github.com/Ptt-official-app/Ptt-backend/blob/development/FOR_DEVELOPERS.md)

也可直接[發 issue](https://github.com/Ptt-official-app/Ptt-backend/issues/new/choose) 討論相關問題。

## Q&A：

* 如果參與這個專案要花多久時間

> 建議至少每周花1小時參與這個專案
> 如果想解決 issue 估計每個 issue 都落在 4\~8 小時左右，超過的話我們會拆分

* 協作方式是什麼

> 請參考[給開發者的開發流程](https://github.com/Ptt-official-app/Ptt-backend/blob/development/FOR_DEVELOPERS.md#%E9%96%8B%E7%99%BC%E6%B5%81%E7%A8%8B)

* 直接討論方式與軟體

> 根據專案之初的投票，目前 go-bbs 以及 Ptt-backend 比較即時的討論在 [g0v slack](https://join.g0v.tw) 的 [#bbs channel](https://g0v-tw.slack.com/archives/C01K6RAR17Y)

