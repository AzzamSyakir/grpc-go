generate-proto:
	protoc --proto_path=protos --go_out=paths=source_relative:./pb --go-grpc_out=paths=source_relative:./pb/pb-grpc protos/*.proto

