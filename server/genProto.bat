@REM protoc -I=D:\wx\coolcar\server\auth\api --go_out=paths=source_relative:gen/go trip.proto

protoc -I=D:\wx\coolcar\server\auth\api --go_out=plugins=grpc,paths=source_relative:gen/v1 auth.proto

protoc -I=D:\wx\coolcar\server\auth\api --grpc-gateway_out=paths=source_relative,grpc_api_configuration=auth.yaml:gen/v1 auth.proto
