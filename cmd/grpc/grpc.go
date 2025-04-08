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

	"github.com/CSKU-Lab/config-server/domain/services"
	"github.com/CSKU-Lab/config-server/domain/models/language"
	pb "github.com/CSKU-Lab/config-server/genproto/config/v1"
	"github.com/CSKU-Lab/config-server/internal/adapters/mongodb"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
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

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln(err)
	}

	db := client.Database("configs")
	languageRepo := mongodb.NewLanguageRepo(db)
	languageService := services.NewLanguageService(languageRepo)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalln("failed to listen: ", err)
	}

	s := grpc.NewServer()
	pb.RegisterConfigServiceServer(s, newServer(languageService))
	reflection.Register(s)
	log.Println("gRPC ConfigService registered")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

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
	langService services.LanguageService
}

func newServer(langService services.LanguageService) *configServiceServer {
	return &configServiceServer{
		langService: langService,
	}
}

func (c *configServiceServer) GetLanguages(ctx context.Context, req *pb.GetLanguagesRequest) (*pb.GetLanguagesResponse, error) {
	responsesLangauges := []*pb.Language{}
	langauges, err := c.langService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

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

func (c *configServiceServer) GetLanguage(ctx context.Context, req *pb.GetLanguageRequest) (*pb.LanguageResponse, error) {
	if req.GetId() == "" {
		return nil, fmt.Errorf("Id is required!")
	}

	lang, err := c.langService.GetByID(ctx, req.GetId())
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

func (c *configServiceServer) AddLanguage(ctx context.Context, req *pb.AddLanguageRequest) (*pb.LanguageResponse, error) {
	lang := language.New(&language.Options{
		Name:        req.GetName(),
		Version:     req.GetVersion(),
		BuildScript: req.GetBuildScript(),
		RunScript:   req.GetRunScript(),
	})

	err := c.langService.Add(ctx, lang)
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

func (c *configServiceServer) UpdateLanguage(ctx context.Context, req *pb.UpdateLanguageRequest) (*pb.LanguageResponse, error) {
	if req.GetId() == "" {
		return nil, fmt.Errorf("Id is required!")
	}

	updated, err := c.langService.UpdateByID(ctx, req.GetId(), &language.PartialOptions{
		Name:        req.Name,
		Version:     req.Version,
		BuildScript: req.BuildScript,
		RunScript:   req.RunScript,
	})
	if err != nil {
		return nil, err
	}

	return &pb.LanguageResponse{
		Id:          updated.ID,
		Name:        updated.Name,
		Version:     updated.Version,
		BuildScript: &updated.BuildScript,
		RunScript:   updated.RunScript,
	}, nil
}

func (c *configServiceServer) DeleteLanguage(ctx context.Context, req *pb.DeleteLanguageRequest) (*emptypb.Empty, error) {
	if req.GetId() == "" {
		return nil, fmt.Errorf("Id is required!")
	}

	err := c.langService.DeleteByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
