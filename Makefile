generate-proto:
	protoc --go_out=grpc-go --go_opt=paths=source_relative    --go-grpc_out=grpc-go --go-grpc_opt=paths=source_relative    grpc-/proto/*.proto

