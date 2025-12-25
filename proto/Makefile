.PHONY: generate
generate:
	@echo "Generating Protocol Buffer code..."
	mkdir -p ../apps/auth-service/pkg/proto
	mkdir -p ../apps/financial-service/pkg/proto
	protoc --proto_path=. \
	       --go_out=../apps/auth-service/pkg/proto --go_opt=paths=source_relative \
	       --go-grpc_out=../apps/auth-service/pkg/proto --go-grpc_opt=paths=source_relative \
	       auth/v1/auth.proto
	
	protoc --proto_path=. \
	       --go_out=../apps/financial-service/pkg/proto --go_opt=paths=source_relative \
	       --go-grpc_out=../apps/financial-service/pkg/proto --go-grpc_opt=paths=source_relative \
	       auth/v1/auth.proto
	@echo "Done!"

.PHONY: clean
clean:
	@echo "Cleaning generated files..."
	find ../apps -name "*.pb.go" -type f -delete
	@echo "Cleaned!"

.PHONY: install-tools
install-tools:
	@echo "Installing protoc plugins..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Done!"
