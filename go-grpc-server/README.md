1. Go side setup: https://grpc.io/docs/languages/go/quickstart/
2. The protoc may not work based on 2), instead, remove the installed protobuf compiler and reinstalled protobuf compiler with "Install pre-compiled binaries" option at https://grpc.io/docs/protoc-installation/
3. To run: go run main.go
   or go build main.go, then
     ./main