package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/CSKU-Lab/config-server/genproto/config/v1"
	languages "github.com/CSKU-Lab/config-server/models/language"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		fmt.Println("You forget to set the MONGO_URI environment variable!")
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx, cancel := context.WithCancel(context.Background())

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalln("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterConfigServiceServer(s, newServer(ctx, uri))
	reflection.Register(s)
	log.Println("gRPC ConfigService registered")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Printf("Receive %s signal from OS, going to shutdown...\n", sig)
		timer := time.AfterFunc(10*time.Second, func() {
			log.Println("Server couldn't stop grafully in time. Doing force stop.")
			s.Stop()
		})
		defer timer.Stop()
		cancel()
		s.GracefulStop()

		log.Println("Successfully gracefully shutdown the server :D")
	}()

	if err := s.Serve(lis); err != nil {
		log.Fatalln("Cannot start grpc server : %v", err)
	}
}

type configServiceServer struct {
	pb.UnimplementedConfigServiceServer

	db  *mongo.Client
	ctx context.Context
}

func newServer(ctx context.Context, mongoURI string) *configServiceServer {
	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal("Can't disconnect mongodb : ", err)
		}
	}()

	return &configServiceServer{
		ctx: ctx,
		db:  client,
	}
}

func strPtr(s string) *string {
	return &s
}

func (c *configServiceServer) GetLanguages(ctx context.Context, request *pb.GetLanguagesRequest) (*pb.GetLanguagesResponse, error) {
	return &pb.GetLanguagesResponse{
		Languages: []*pb.Language{
			{
				Id:          "python_3.11.2",
				Name:        strPtr("Python"),
				Version:     strPtr("3.11.2"),
				BuildScript: "no script",
				RunScript:   "no script",
			},
		},
	}, nil

}

// func getAll(ctx context.Context, uri string) {
//
// 	coll := client.Database("configs").Collection("languages")
//
// 	cursor, err := coll.Find(ctx, bson.D{})
// 	if err != nil {
// 		log.Fatalln("There is an error :", err)
// 	}
//
// 	var languages []languages.Language
// 	if err := cursor.All(ctx, &languages); err != nil {
// 		log.Fatalln("There is an error :", err)
// 	}
//
// 	log.Println(languages)
// }

func work(ctx context.Context, uri string) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal("Can't disconnect mongodb : ", err)
		}
	}()

	pythonLang := languages.New(languages.Options{
		Name:      "Python",
		Version:   "3.11.2",
		RunScript: "python3 main.py",
	})

	coll := client.Database("configs").Collection("languages")

	res, err := coll.InsertOne(ctx, &pythonLang)
	if err != nil {
		log.Fatalln("Cannot insert the document :", err)
	}

	log.Println("Language added : ", res)
}
