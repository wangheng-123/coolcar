package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	blobpb "server/blob/api/gen/v1"
	carpb "server/car/api/gen/v1"
	rentalpb "server/rental/api/gen/v1"
	"server/rental/profile"
	profiledao "server/rental/profile/dao"
	"server/rental/trip"
	"server/rental/trip/client/car"
	"server/rental/trip/client/poi"
	profClient "server/rental/trip/client/profile"
	tripdao "server/rental/trip/dao"
	"server/shared/server"
	"time"
)

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

	db := mongoClient.Database("test")

	blobConn, err := grpc.Dial("localhost:8083", grpc.WithInsecure())
	if err != nil {
		logger.Fatal("cannot connect blob service", zap.Error(err))
	}

	profService := &profile.Service{
		BlobClient:        blobpb.NewBlobServiceClient(blobConn),
		PhotoGetExpire:    5 * time.Second,
		PhotoUploadExpire: 10 * time.Second,
		Mongo:             profiledao.NewMongo(db),
		Logger:            logger,
	}

	carConn, err := grpc.Dial("localhost:8084", grpc.WithInsecure())
	if err != nil {
		logger.Fatal("cannot connect car service", zap.Error(err))
	}

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:              "rental",
		Addr:              ":8082",
		AuthPublicKeyFile: "shared/auth/public.key",
		Logger:            logger,
		RegisterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				CarManager: &car.Manager{
					CarService: carpb.NewCarServiceClient(carConn),
				},
				ProfileManager: &profClient.Manager{
					Fetcher: profService,
				},
				POIManager: &poi.Manager{},
				Mongo:      tripdao.NewMongo(db),
				Logger:     logger,
			})
			rentalpb.RegisterProfileServiceServer(s, profService)
		},
	}))
}
