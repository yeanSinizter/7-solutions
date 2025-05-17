.PHONY: grpc

grpc:
	protoc --go_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:. \
		proto/user.proto
