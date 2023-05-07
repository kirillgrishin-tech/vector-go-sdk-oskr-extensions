# vector-go-sdk-oskr-extensions
New gRPC server to run on Vector OSKR. It allows new commands to be implemented in the vector-go-sdk.

# PROTOBUF INTERFACE GENERATION
1. Download protobuf (https://github.com/protocolbuffers/protobuf/releases)
2. Generate go files from oskr.proto, from directory "proto":
   protoc --go_out=../pkg/oskrpb --go_opt=paths=source_relative --go-grpc_out=../pkg/oskrpb --go-grpc_opt=paths=source_relative oskr.proto
3. This will generate oskr.pb.go and oskr.grpc.go
   The first one is for the server, the 2nd one for the client. Copy them also in the client application (for example, 
   in vector-go-sdk)

# BUILDING AND RUNNING ON VECTOR
1. Run ./build.sh it will build the server executable in build/vic-oskr-server
2. On the OSKR Vector, do the following: remount the /data partition as executable, and add a firewall rule to allow
   connections on the port 50051 where the server listens to
   mount -o remount,exec /data
   iptables -A INPUT -p tcp --dport 50051 -m conntrack --ctstate NEW,ESTABLISHED -j ACCEPT
   iptables -A OUTPUT -p tcp --sport 50051 -m conntrack --ctstate ESTABLISHED -j ACCEPT
3. Upload the vic-oskr-server on Vector:
   scp -i <PATH_TO_id_RSA_key> build/vic-oskr-server root@192.168.43.216:/data
4. On the OSKR Vector, run the server
   cd /data
   ./vic-oskr-server

Now you can connect with a client (see https://github.com/kirillgrishin-tech/vector-go-sdk/blob/main/cmd/examples/navigator/main.go
for an example)
