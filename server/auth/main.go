package main

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os"
	authpb "server/auth/api/gen/v1"
	"server/auth/appinfo"
	"server/auth/auth"
	"server/auth/dao"
	"server/auth/token"
	"server/auth/wechat"
	"server/shared/server"
	"time"
)

//用户微服务，暴露8080接口给其他微服务调用
func main() {
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf("cannot create logger:%v", err)
	}

	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://root:root@localhost:27017"))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}

	pkFile, err := os.Open("auth/private.key")
	if err != nil {
		logger.Fatal("cannot open private key", zap.Error(err))
	}

	pkBytes, err := ioutil.ReadAll(pkFile)
	if err != nil {
		logger.Fatal("cannot read private key", zap.Error(err))
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
	if err != nil {
		logger.Fatal("cannot parse private key", zap.Error(err))
	}

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:   "auth",
		Addr:   ":8081",
		Logger: logger,
		RegisterFunc: func(s *grpc.Server) {
			authpb.RegisterAuthServiceServer(s, &auth.Service{
				OpenIDResolver: &wechat.Service{
					AppID:     appinfo.AppID,
					AppSecret: appinfo.AppSecret,
				},
				Mongo:          dao.NewMongo(mongoClient.Database("test")),
				Logger:         logger,
				TokenExpire:    10 * time.Second,
				TokenGenerator: token.NewJWTTokenGen("server/auth", privKey),
			})
		},
	}))
}
