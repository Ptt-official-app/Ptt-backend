# PTT APP 後端

## 這專案是做什麼的,跟什麼有關係

Ptt-backend 是做 BBS (Bulletin Board System) 的開源專案
但 BBS 有很多種，Ptt 只是 BBS 的其中一種
go-bbs 是與 Ptt-backend 有關的專案，解析儲存的檔案提供給 Ptt-backend 使用

## 如何參與這個專案

1. 要有個 GitHub 帳號
2. 參與 issue 討論

想要更加地投入可以參考下面的 *貢獻者相關* 跟 *Q&A*

## 文件檔案部分

* [RESTful API 文件](https://docs.google.com/document/d/18DsZOyrlr5BIl2kKxZH7P2QxFLG02xL2SO0PzVHVY3k/edit?usp=sharing)
* [go-bbs Package](https://github.com/Ptt-official-app/go-bbs)

[Openapi](http://localhost:8081/swagger) 使用 docker-compose 啟動，docker rootless 設定請參閱 [Manage Docker as a non-root user](https://docs.docker.com/engine/install/linux-postinstall/#manage-docker-as-a-non-root-user)

```sh
docker-compose up -d
```

## 軟體測試測試表

測試表目前以 Google Sheet 寫成，這份測試表適合讓對於不熟悉程式語言的開發者了解目前各項功能的狀態。

測試表： [連結](https://docs.google.com/spreadsheets/d/1uo4AJuSi5xTXEht2o2EHogLivCJlJqlLaeqoj1RceDY/edit?usp=sharing)

歷次測試結果：

* [2021060401](https://docs.google.com/spreadsheets/d/1dyfmWZJaTiDrSMIFZ6ynmfeTWU3h45ScW1eyKFJU494/edit?usp=sharing)
* [2021062101](https://docs.google.com/spreadsheets/d/1RGIQPN6KfiCzWRQe-BeTXFUTvdwytOv0Yhd5a8ixCLk/edit?usp=sharing)
* [2021070501](https://docs.google.com/spreadsheets/d/1thxyY9jf2GkK3DMgGO1bAHWE3BmsFb4p5Ot5eAgbaII/edit?usp=sharing)

## 貢獻者相關

如果想要參與開發與維護專案可參考: [FOR_DEVELOPERS.md](https://github.com/Ptt-official-app/Ptt-backend/blob/development/FOR_DEVELOPERS.md)

也可直接發 issue 討論相關問題

## Q&A

* 如果參與這個專案要花多久時間

建議至少每周花1小時參與這個專案

如果想解決 issue 估計每個 issue 都落在 4\~8 小時左右，超過的話我們會拆分

* 協作方式是什麼

可參考 [FOR_DEVELOPERS.md](https://github.com/Ptt-official-app/Ptt-backend/blob/development/FOR_DEVELOPERS.md) 的開發流程

* 直接討論方式與軟體

[slack](https://g0v-tw.slack.com/archives/C01K6RAR17Y)
