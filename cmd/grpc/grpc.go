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

	"github.com/CSKU-Lab/config-server/domain/models/compare"
	"github.com/CSKU-Lab/config-server/domain/models/language"
	"github.com/CSKU-Lab/config-server/domain/services"
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

	compareRepo := mongodb.NewCompareRepo(db)
	compareService := services.NewCompareService(compareRepo)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		log.Fatalln("failed to listen: ", err)
	}

	s := grpc.NewServer()
	pb.RegisterConfigServiceServer(s, newServer(languageService, compareService))
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
	langService    services.LanguageService
	compareService services.CompareService
}

func newServer(langService services.LanguageService, compareService services.CompareService) *configServiceServer {
	return &configServiceServer{
		langService:    langService,
		compareService: compareService,
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

func (c *configServiceServer) AddCompare(ctx context.Context, req *pb.AddCompareRequest) (*pb.CompareResponse, error) {
	compare := compare.New(&compare.Option{
		Name:        req.GetName(),
		Script:      req.GetScript(),
		ScriptName:  req.GetScriptName(),
		BuildScript: req.GetBuildScript(),
		RunScript:   req.GetRunScript(),
		RunName:     req.GetRunName(),
		Description: req.GetDescription(),
	})

	err := c.compareService.Add(ctx, compare)
	if err != nil {
		return nil, err
	}

	return &pb.CompareResponse{
		Id:          compare.ID,
		Name:        compare.Name,
		Script:      compare.Script,
		ScriptName:  compare.ScriptName,
		BuildScript: compare.BuildScript,
		RunScript:   compare.RunScript,
		RunName:     compare.RunName,
		Description: compare.Description,
	}, nil
}

func (c *configServiceServer) GetCompare(ctx context.Context, req *pb.GetCompareRequest) (*pb.CompareResponse, error) {
	if req.GetId() == "" {
		return nil, fmt.Errorf("Id is required!")
	}

	compare, err := c.compareService.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.CompareResponse{
		Id:          compare.ID,
		Name:        compare.Name,
		Script:      compare.Script,
		ScriptName:  compare.ScriptName,
		BuildScript: compare.BuildScript,
		RunScript:   compare.RunScript,
		RunName:     compare.RunName,
		Description: compare.Description,
	}, nil
}

func (c *configServiceServer) GetCompares(ctx context.Context, req *emptypb.Empty) (*pb.GetComparesResponse, error) {
	compare, err := c.compareService.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := []*pb.CompareResponse{}
	for _, compare := range compare {
		responses = append(responses, &pb.CompareResponse{
			Id:          compare.ID,
			Name:        compare.Name,
			Script:      compare.Script,
			ScriptName:  compare.ScriptName,
			BuildScript: compare.BuildScript,
			RunScript:   compare.RunScript,
			RunName:     compare.RunName,
			Description: compare.Description,
		})
	}

	return &pb.GetComparesResponse{
		Compares: responses,
	}, nil
}

func (c *configServiceServer) UpdateCompare(ctx context.Context, req *pb.UpdateCompareRequest) (*pb.CompareResponse, error) {
	if req.GetId() == "" {
		return nil, fmt.Errorf("Id is required!")
	}

	updated, err := c.compareService.UpdateByID(ctx, req.GetId(), &compare.PartialOption{
		Name:        req.Name,
		Script:      req.Script,
		ScriptName:  req.ScriptName,
		BuildScript: req.BuildScript,
		RunScript:   req.RunScript,
		RunName:     req.RunName,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CompareResponse{
		Id:          *updated.ID,
		Name:        *updated.Name,
		Script:      *updated.Script,
		ScriptName:  *updated.ScriptName,
		BuildScript: *updated.BuildScript,
		RunScript:   *updated.RunScript,
		RunName:     *updated.RunName,
		Description: *updated.Description,
	}, nil
}

func (c *configServiceServer) DeleteCompare(ctx context.Context, req *pb.DeleteCompareRequest) (*emptypb.Empty, error) {
	if req.GetId() == "" {
		return nil, fmt.Errorf("Id is required!")
	}

	err := c.compareService.DeleteByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return nil, nil
}
