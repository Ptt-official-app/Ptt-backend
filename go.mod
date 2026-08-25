module github.com/Ptt-official-app/Ptt-backend

go 1.25.0

require (
	github.com/golang-jwt/jwt/v4 v4.5.2
	github.com/pelletier/go-toml v1.9.5
)

require (
	github.com/PichuChen/postgresql-gobbs v0.0.0-20251019143819-aba6e40a12ae
	github.com/Ptt-official-app/go-bbs v0.12.1-0.20251019065628-76580d57f8e5
	github.com/gorilla/websocket v1.5.3
	github.com/lib/pq v1.10.9
	github.com/modelcontextprotocol/go-sdk v1.7.0
	golang.org/x/net v0.50.0
	golang.org/x/text v0.34.0
	google.golang.org/grpc v1.72.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/google/jsonschema-go v0.4.3 // indirect
	github.com/segmentio/asm v1.1.3 // indirect
	github.com/segmentio/encoding v0.5.4 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	golang.org/x/oauth2 v0.35.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/time v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
)

// replace github.com/Ptt-official-app/go-bbs => ../go-bbs

// replace github.com/PichuChen/postgresql-gobbs => ../postgresql-gobbs

// replace github.com/ptt/pttweb => ../pttweb
