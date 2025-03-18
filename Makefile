# Global commands
run-all:
	@go run main.go run $(ARGS)
migrate-all:
	@go run main.go migrate $(ARGS)

# Orders commands
run-orders:
	@go run main.go orders run $(ARGS)
migrate-orders:
	@go run main.go orders migrate $(ARGS)

# Kitchen commands
run-kitchen:
	@go run main.go kitchen run $(ARGS)
migrate-kitchen:
	@go run main.go kitchen migrate $(ARGS)

# Commands for generating protobuf code from protobuf folder at different languages
gen-go:
	@protoc \
		--proto_path=./protobuf "orders.proto" \
		--go_out=services/common/genproto/orders --go_opt=paths=source_relative \
		--go-grpc_out=services/common/genproto/orders \
		--go-grpc_opt=paths=source_relative
gen-cpp:
	@protoc \
		--proto_path=./protobuf "orders.proto" \
		--cpp_out=services/common/genproto/orders