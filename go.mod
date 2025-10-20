module github.com/Ptt-official-app/Ptt-backend

go 1.24.2

require (
	github.com/PichuChen/postgresql-gobbs v0.0.0-20231009120523-1f2a4b3c5d7e
	github.com/Ptt-official-app/go-bbs v0.12.0
	github.com/golang-jwt/jwt/v4 v4.5.2
	github.com/pelletier/go-toml v1.9.5
)

require (
	golang.org/x/net v0.37.0
	google.golang.org/grpc v1.72.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b // indirect
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/rvelhote/go-recaptcha v0.0.0-20170215232712-e143c6ea64e5 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/tools v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
)

replace github.com/Ptt-official-app/go-bbs => ../go-bbs

replace github.com/PichuChen/postgresql-gobbs => ../postgresql-gobbs

// replace github.com/ptt/pttweb => ../pttweb
