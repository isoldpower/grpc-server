# Commands for running the services
run-orders:
	@go run services/orders/*.go
run-kitchen:
	@go run services/kitchen/*.go

# Commands for generating protobuf code from protobuf folder at different languages
gen-go:
	@protoc \
		--proto_path=protobuf "protobuf/orders.proto" \
		--go_out=services/common/genproto/orders --go_opt=paths=source_relative \
		--go-grpc_out=services/common/genproto/orders \
		--go-grpc_opt=paths=source_relative
gen-cpp:
	@protoc \
		--proto_path=protobuf "protobuf/orders.proto" \
		--cpp_out=services/common/genproto/orders