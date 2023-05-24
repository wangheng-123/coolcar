package server

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"server/shared/auth"
)

type GRPCConfig struct {
	Name              string
	Addr              string
	AuthPublicKeyFile string
	Logger            *zap.Logger
	RegisterFunc      func(*grpc.Server)
}

// RunGRPCServer runs a grpc server
func RunGRPCServer(c *GRPCConfig) error {
	nameField := zap.String("name", c.Name)
	listen, err := net.Listen("tcp", c.Addr)
	if err != nil {
		c.Logger.Fatal("cannot listen ", nameField, zap.Error(err))
	}

	var opts []grpc.ServerOption

	if c.AuthPublicKeyFile != "" {
		in, err := auth.Interceptor(c.AuthPublicKeyFile)
		if err != nil {
			c.Logger.Fatal("cannot create auth interceptor", nameField, zap.Error(err))
		}
		opts = append(opts, grpc.UnaryInterceptor(in))
	}

	server := grpc.NewServer(opts...)
	//服务注册
	c.RegisterFunc(server)

	c.Logger.Sugar().Info("server started", nameField, zap.String("addr", c.Addr))
	return server.Serve(listen)
}
