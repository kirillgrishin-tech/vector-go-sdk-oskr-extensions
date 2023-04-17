# vector-go-sdk-oskr-extensions
New gRPC server to run on Vector OSKR. It allows new commands to be implemented in the vector-go-sdk.
1. Download protobuf (https://github.com/protocolbuffers/protobuf/releases)
2. Generate go files from oskr.proto, from directory "proto":
   protoc --go_out=../pkg/oskrpb --go_opt=paths=source_relative --go-grpc_out=../pkg/oskrpb --go-grpc_opt=paths=source_relative oskr.proto

3. This will generate oskr.pb.go and oskr.grpc.go
