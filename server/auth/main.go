package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	authpb "server/auth/api/gen/v1"
	"server/auth/appinfo"
	"server/auth/auth"
	"server/auth/dao"
	"server/auth/wechat"
)

//用户微服务，暴露8080接口给其他微服务调用
func main() {
	logger, err2 := newZapLogger()
	if err2 != nil {
		log.Fatalf("cannot create logger:%v", err2)
	}
	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		logger.Fatal("cannot listen ", zap.Error(err))
	}

	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://root:root@localhost:27017"))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}

	server := grpc.NewServer()
	//服务注册
	authpb.RegisterAuthServiceServer(server, &auth.Service{
		OpenIDResolver: &wechat.Service{
			AppID:     appinfo.AppID,
			AppSecret: appinfo.AppSecret,
		},
		Mongo:  dao.NewMongo(mongoClient.Database("test")),
		Logger: logger,
	})

	err2 = server.Serve(listen)
	if err2 != nil {
		logger.Fatal("cannot server", zap.Error(err2))
	}
}

func newZapLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.TimeKey = ""
	return cfg.Build()
}
