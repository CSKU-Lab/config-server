package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	pb "github.com/CSKU-Lab/config-server/genproto/config/v1"
	languages "github.com/CSKU-Lab/config-server/models/language"
	"go.mongodb.org/mongo-driver/v2/bson"
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

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalln("failed to listen: ", err)
	}

	s := grpc.NewServer()
	pb.RegisterConfigServiceServer(s, newServer(ctx, client))
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

		if err := client.Disconnect(ctx); err != nil {
			log.Fatalln("Can't disconnect mongodb : ", err)
		}

		log.Println("Successfully gracefully shutdown the server :D")
	}()

	if err := s.Serve(lis); err != nil {
		log.Fatalln("Cannot start grpc server :", err)
	}
}

type configServiceServer struct {
	pb.UnimplementedConfigServiceServer

	col *mongo.Collection
	ctx context.Context
}

func newServer(ctx context.Context, client *mongo.Client) *configServiceServer {
	return &configServiceServer{
		ctx: ctx,
		col: client.Database("configs").Collection("languages"),
	}
}

func (c *configServiceServer) GetLanguages(ctx context.Context, req *pb.GetLanguagesRequest) (*pb.GetLanguagesResponse, error) {

	cursor, err := c.col.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("Cannot get languages : %v", err)
	}

	var langauges []languages.Language
	err = cursor.All(ctx, &langauges)
	if err != nil {
		return nil, err
	}

	responsesLangauges := []*pb.Language{}

	for _, language := range langauges {
		var name *string = nil
		if req.IncludeName {
			name = &language.Name
		}
		var version *string = nil
		if req.IncludeVersion {
			version = &language.Version
		}
		responsesLangauges = append(responsesLangauges, &pb.Language{
			Id:          language.ID,
			Name:        name,
			Version:     version,
			BuildScript: language.BuildScript,
			RunScript:   language.RunScript,
		})
	}

	return &pb.GetLanguagesResponse{
		Languages: responsesLangauges,
	}, nil

}

func (c *configServiceServer) AddLanguage(ctx context.Context, req *pb.AddLanguageRequest) (*pb.LanguageResponse, error) {
	lang := languages.New(&languages.Options{
		Name:        req.GetName(),
		Version:     req.GetVersion(),
		BuildScript: req.GetBuildScript(),
		RunScript:   req.GetRunScript(),
	})
	_, err := c.col.InsertOne(ctx, lang)
	if err != nil {
		return nil, err
	}

	return &pb.LanguageResponse{
		Id:          lang.ID,
		Name:        lang.Name,
		Version:     lang.Version,
		BuildScript: &lang.BuildScript,
		RunScript:   lang.RunScript,
	}, nil
}

func getUpdatedFields(i any) bson.D {
	v := reflect.ValueOf(i).Elem()
	t := reflect.TypeOf(i).Elem()

	fields := bson.D{}
	for i := range v.NumField() {
		fieldVal := v.Field(i)
		fieldTyp := t.Field(i)
		bsonTag := fieldTyp.Tag.Get("bson")

		if fieldVal.IsNil() {
			continue
		}

		fields = append(fields, bson.E{Key: bsonTag, Value: fieldVal.Elem().Interface()})
	}

	return fields
}

func (c *configServiceServer) UpdateLanguage(ctx context.Context, req *pb.UpdateLanguageRequest) (*pb.LanguageResponse, error) {
	if req.GetId() == "" {
		return nil, fmt.Errorf("Id is required!")
	}
	modLang := languages.NewUpdate(&languages.PartialOptions{
		Name:        req.Name,
		Version:     req.Version,
		BuildScript: req.BuildScript,
		RunScript:   req.RunScript,
	})

	updatedFields := getUpdatedFields(modLang)
	_, err := c.col.UpdateOne(ctx, bson.M{"id": req.GetId()}, bson.D{{"$set", updatedFields}})
	if err != nil {
		return nil, err
	}
	return &pb.LanguageResponse{
		Id:          *modLang.ID,
		Name:        *modLang.Name,
		Version:     *modLang.Version,
		BuildScript: modLang.BuildScript,
		RunScript:   *modLang.RunScript,
	}, nil
}
