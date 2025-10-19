.PHONY: all

all:
	protoc --go_out=. --go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--go_opt=Mman.proto=github.com/Ptt-official-app/Ptt-backend/internal/proto/man man.proto \
		--go-grpc_opt=Mman.proto=github.com/Ptt-official-app/Ptt-backend/internal/proto/man man.proto
