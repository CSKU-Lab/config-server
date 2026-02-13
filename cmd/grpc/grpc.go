package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/CSKU-Lab/cache"
	"github.com/CSKU-Lab/config-server/configs"
	"github.com/CSKU-Lab/config-server/domain/models"
	"github.com/CSKU-Lab/config-server/domain/requests"
	"github.com/CSKU-Lab/config-server/domain/services"
	pb "github.com/CSKU-Lab/config-server/genproto/config/v1"
	graderPB "github.com/CSKU-Lab/config-server/genproto/grader/v1"
	taskPB "github.com/CSKU-Lab/config-server/genproto/task/v1"
	"github.com/CSKU-Lab/config-server/internal/adapters/mongodb"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	env := configs.NewEnv()

	client, err := mongo.Connect(options.Client().ApplyURI(env.Get("MONGO_URI")))
	if err != nil {
		log.Fatalln(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalln("Failed to ping MongoDB: ", err)
	}

	db := client.Database(env.Get("DATABASE_NAME"))
	runnerRepo := mongodb.NewRunnerRepo(db)
	runnerService := services.NewRunnerService(runnerRepo)

	compareRepo := mongodb.NewCompareRepo(db)
	compareService := services.NewCompareService(compareRepo)

	taskGrpcClient, closeConn, err := initTaskGRPCClient(env.Get("TASK_SERVER_URL"))
	if err != nil {
		log.Fatal("Failed to connect to task gRPC server: ", err)
	}
	defer closeConn()

	graderGRPCClient, closeConn, err := initGraderGRPCClient(env.Get("GO_GRADER_SERVER_URL"))
	if err != nil {
		log.Fatal("Failed to connect to grader gRPC server: ", err)
	}
	defer closeConn()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", env.Get("PORT")))
	if err != nil {
		log.Fatalln("failed to listen: ", err)
	}

	redis, err := cache.NewRedis(&cache.RedisOptions{
		Addr:     env.Get("REDIS_SERVER_URL"),
		Password: env.Get("REDIS_PASSWORD"),
	})
	if err != nil {
		log.Fatalln("Cannot initialize cache repository", "error", err)
	}
	defer redis.Close()

	s := grpc.NewServer()
	pb.RegisterConfigServiceServer(s, newServer(runnerService, compareService, taskGrpcClient, graderGRPCClient, redis))
	reflection.Register(s)
	log.Println("gRPC ConfigService registered")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Go(func() {
		sig := <-sigs
		log.Printf("Receive %s signal from OS, going to shutdown...\n", sig)
		timer := time.AfterFunc(10*time.Second, func() {
			log.Println("Server couldn't stop grafully in time. Doing force stop.")
			s.Stop()
		})
		defer timer.Stop()
		s.GracefulStop()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := client.Disconnect(ctx); err != nil {
			log.Println("Can't disconnect mongodb : ", err)
		}
		log.Println("Successfully gracefully shutdown the server :D")
	})

	if err := s.Serve(lis); err != nil {
		log.Fatalln("Cannot start grpc server :", err)
	}

	wg.Wait()
}

type configServiceServer struct {
	pb.UnimplementedConfigServiceServer
	runnerService  services.RunnerService
	compareService services.CompareService
	taskClient     taskPB.TaskServiceClient
	graderClient   graderPB.GraderServiceClient
	runnerCache    cache.CacheBuild
	compareCache   cache.CacheBuild
	cacheApp       cache.CacheApp
}

func newServer(runnerService services.RunnerService, compareService services.CompareService, taskClient taskPB.TaskServiceClient, graderClient graderPB.GraderServiceClient, cacheApp cache.CacheApp) *configServiceServer {
	runnerCache := cacheApp.Build("runnerCache")
	compareCache := cacheApp.Build("compareCache")

	return &configServiceServer{
		runnerService:  runnerService,
		compareService: compareService,
		taskClient:     taskClient,
		graderClient:   graderClient,
		runnerCache:    runnerCache,
		compareCache:   compareCache,
		cacheApp:       cacheApp,
	}
}

func (c *configServiceServer) GetRunners(ctx context.Context, req *pb.GetRunnersRequest) (*pb.GetRunnersResponse, error) {
	cacheObj := c.runnerCache.All(time.Hour * 4)
	cacheInstance := cache.NewCacheInstance[*pb.GetRunnersResponse](cacheObj)

	runnerRes, err := cacheInstance.LazyCaching(ctx, func() (*pb.GetRunnersResponse, error) {
		responsesRunners := []*pb.Runner{}
		runners, err := c.runnerService.GetAll(ctx)
		if err != nil {
			return nil, err
		}

		for _, runner := range runners {
			var name *string = nil
			if req.IncludeName {
				name = &runner.Name
			}
			responsesRunners = append(responsesRunners, &pb.Runner{
				Id:          runner.ID,
				Name:        name,
				BuildScript: runner.BuildScript,
				RunScript:   runner.RunScript,
			})
		}

		return &pb.GetRunnersResponse{
			Runners: responsesRunners,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return runnerRes, nil
}

func (c *configServiceServer) GetRunnersPagination(ctx context.Context, req *pb.GetRunnersPaginationRequest) (*pb.GetRunnersPaginationResponse, error) {
	cacheObj := c.runnerCache.One(time.Hour*4, req.Pagination.String())
	cacheInstance := cache.NewCacheInstance[*pb.GetRunnersPaginationResponse](cacheObj)

	paginationRes, err := cacheInstance.LazyCaching(ctx, func() (*pb.GetRunnersPaginationResponse, error) {
		runners, total, err := c.runnerService.GetPagination(
			ctx,
			&requests.GetPagination{
				Page:      int(req.Pagination.GetPage()),
				PageSize:  int(req.Pagination.GetPageSize()),
				SortOrder: req.Pagination.GetSortOrder(),
				Search:    req.Pagination.GetSearch(),
			})
		if err != nil {
			return nil, err
		}

		responseRunners := []*pb.Runner{}
		for _, runner := range runners {
			var name *string = nil
			if req.IncludeName {
				name = &runner.Name
			}
			responseRunners = append(responseRunners, &pb.Runner{
				Id:          runner.ID,
				Name:        name,
				BuildScript: runner.BuildScript,
				RunScript:   runner.RunScript,
			})
		}

		return &pb.GetRunnersPaginationResponse{
			Runners: responseRunners,
			Count:   int32(total),
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return paginationRes, nil
}

func (c *configServiceServer) GetRunner(ctx context.Context, req *pb.GetRunnerRequest) (*pb.RunnerResponse, error) {
	if req.GetId() == "" {
		return nil, fmt.Errorf("Id is required!")
	}

	cacheObj := c.runnerCache.One(time.Hour*4, req.GetId())
	cacheInstance := cache.NewCacheInstance[*pb.RunnerResponse](cacheObj)

	runnerRes, err := cacheInstance.LazyCaching(ctx, func() (*pb.RunnerResponse, error) {
		runner, err := c.runnerService.GetByID(ctx, req.GetId())
		if err != nil {
			return nil, err
		}

		return &pb.RunnerResponse{
			Id:          runner.ID,
			Name:        runner.Name,
			BuildScript: &runner.BuildScript,
			RunScript:   runner.RunScript,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return runnerRes, nil
}

func (c *configServiceServer) CreateRunner(ctx context.Context, req *pb.CreateRunnerRequest) (*pb.CreateRunnerResponse, error) {
	runner := &requests.CreateRunner{
		Name:        req.GetName(),
		BuildScript: req.GetBuildScript(),
		RunScript:   req.GetRunScript(),
	}

	runnerID, err := c.runnerService.Create(ctx, runner)
	if err != nil {
		return nil, err
	}

	err = c.runnerCache.InvalidateAll(ctx)
	if err != nil {
		return nil, err
	}

	broadcastReq := &graderPB.BroadcastRequest{
		Action: graderPB.BroadcastAction_REFETCH_CONFIG,
	}
	_, err = c.graderClient.Broadcast(ctx, broadcastReq)
	if err != nil {
		return nil, err
	}

	return &pb.CreateRunnerResponse{
		Id: runnerID,
	}, nil
}

func (c *configServiceServer) UpdateRunner(ctx context.Context, req *pb.UpdateRunnerRequest) (*emptypb.Empty, error) {
	if req.GetId() == "" {
		return nil, fmt.Errorf("Id is required!")
	}

	err := c.runnerService.UpdateByID(ctx, req.GetId(), &requests.UpdateRunner{
		Name:        req.Name,
		BuildScript: req.BuildScript,
		RunScript:   req.RunScript,
	})
	if err != nil {
		return nil, err
	}

	err = c.runnerCache.InvalidateAll(ctx)
	if err != nil {
		return nil, err
	}

	broadcastReq := &graderPB.BroadcastRequest{
		Action: graderPB.BroadcastAction_REFETCH_CONFIG,
	}
	_, err = c.graderClient.Broadcast(ctx, broadcastReq)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *configServiceServer) DeleteRunner(ctx context.Context, req *pb.DeleteRunnerRequest) (*emptypb.Empty, error) {
	if req.GetId() == "" {
		return nil, fmt.Errorf("Id is required!")
	}

	err := c.runnerService.DeleteByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	err = c.runnerCache.InvalidateAll(ctx)
	if err != nil {
		return nil, err
	}

	taskReq := &taskPB.RemoveRunnerOnCascadeRequest{
		RunnerId: req.GetId(),
	}
	_, err = c.taskClient.RemoveRunnerOnCascade(ctx, taskReq)
	if err != nil {
		return nil, err
	}

	broadcastReq := &graderPB.BroadcastRequest{
		Action: graderPB.BroadcastAction_REFETCH_CONFIG,
	}
	_, err = c.graderClient.Broadcast(ctx, broadcastReq)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *configServiceServer) CreateCompare(ctx context.Context, req *pb.CreateCompareRequest) (*pb.CreateCompareResponse, error) {
	compareID, err := c.compareService.Create(ctx, &requests.CreateCompare{
		Name:        req.GetName(),
		Files:       models.PBFileToFile(req.GetFiles()),
		BuildScript: req.GetBuildScript(),
		RunScript:   req.GetRunScript(),
		RunName:     req.GetRunName(),
		Description: req.GetDescription(),
	})
	if err != nil {
		return nil, err
	}

	err = c.compareCache.InvalidateAll(ctx)
	if err != nil {
		return nil, err
	}

	broadcastReq := &graderPB.BroadcastRequest{
		Action: graderPB.BroadcastAction_REFETCH_CONFIG,
	}
	_, err = c.graderClient.Broadcast(ctx, broadcastReq)
	if err != nil {
		return nil, err
	}

	return &pb.CreateCompareResponse{
		Id: compareID,
	}, nil
}

func (c *configServiceServer) GetCompare(ctx context.Context, req *pb.GetCompareRequest) (*pb.CompareResponse, error) {
	if req.GetId() == "" {
		return nil, fmt.Errorf("Id is required!")
	}

	cacheObj := c.compareCache.One(time.Hour*4, req.GetId())
	cacheInstance := cache.NewCacheInstance[*pb.CompareResponse](cacheObj)

	compareRes, err := cacheInstance.LazyCaching(ctx, func() (*pb.CompareResponse, error) {
		compare, err := c.compareService.GetByID(ctx, req.GetId())
		if err != nil {
			return nil, err
		}

		return &pb.CompareResponse{
			Id:          compare.ID,
			Name:        compare.Name,
			Files:       models.FileToPBFile(compare.Files),
			BuildScript: compare.BuildScript,
			RunScript:   compare.RunScript,
			RunName:     compare.RunName,
			Description: compare.Description,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return compareRes, nil
}

func (c *configServiceServer) GetCompares(ctx context.Context, req *emptypb.Empty) (*pb.GetComparesResponse, error) {
	cacheObj := c.compareCache.All(time.Hour * 4)
	cacheInstance := cache.NewCacheInstance[*pb.GetComparesResponse](cacheObj)

	compareRes, err := cacheInstance.LazyCaching(ctx, func() (*pb.GetComparesResponse, error) {
		compare, err := c.compareService.GetAll(ctx)
		if err != nil {
			return nil, err
		}

		responses := []*pb.CompareResponse{}
		for _, compare := range compare {
			responses = append(responses, &pb.CompareResponse{
				Id:          compare.ID,
				Name:        compare.Name,
				Files:       models.FileToPBFile(compare.Files),
				BuildScript: compare.BuildScript,
				RunScript:   compare.RunScript,
				RunName:     compare.RunName,
				Description: compare.Description,
			})
		}

		return &pb.GetComparesResponse{
			Compares: responses,
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return compareRes, nil
}

func (c *configServiceServer) GetComparesPagination(ctx context.Context, req *pb.GetComparesPaginationRequest) (*pb.GetComparesPaginationResponse, error) {
	cacheObj := c.compareCache.One(time.Hour*4, req.Pagination.String())
	cacheInstance := cache.NewCacheInstance[*pb.GetComparesPaginationResponse](cacheObj)

	paginationRes, err := cacheInstance.LazyCaching(ctx, func() (*pb.GetComparesPaginationResponse, error) {
		compares, total, err := c.compareService.GetPagination(
			ctx,
			&requests.GetPagination{
				Page:      int(req.Pagination.GetPage()),
				PageSize:  int(req.Pagination.GetPageSize()),
				SortOrder: req.Pagination.GetSortOrder(),
				Search:    req.Pagination.GetSearch(),
			})
		if err != nil {
			return nil, err
		}

		responses := []*pb.CompareResponse{}
		for _, compare := range compares {
			responses = append(responses, &pb.CompareResponse{
				Id:          compare.ID,
				Name:        compare.Name,
				Files:       models.FileToPBFile(compare.Files),
				BuildScript: compare.BuildScript,
				RunScript:   compare.RunScript,
				RunName:     compare.RunName,
				Description: compare.Description,
			})
		}

		return &pb.GetComparesPaginationResponse{
			Compares: responses,
			Count:    int32(total),
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return paginationRes, nil
}

func (c *configServiceServer) UpdateCompare(ctx context.Context, req *pb.UpdateCompareRequest) (*emptypb.Empty, error) {
	if req.GetId() == "" {
		return nil, fmt.Errorf("Id is required!")
	}

	err := c.compareService.UpdateByID(ctx, req.GetId(), &requests.UpdateCompare{
		Name:        req.Name,
		Files:       models.PBFileToFile(req.GetFiles()),
		BuildScript: req.BuildScript,
		RunScript:   req.RunScript,
		RunName:     req.RunName,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	err = c.compareCache.InvalidateAll(ctx)
	if err != nil {
		return nil, err
	}

	broadcastReq := &graderPB.BroadcastRequest{
		Action: graderPB.BroadcastAction_REFETCH_CONFIG,
	}
	_, err = c.graderClient.Broadcast(ctx, broadcastReq)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *configServiceServer) DeleteCompare(ctx context.Context, req *pb.DeleteCompareRequest) (*emptypb.Empty, error) {
	if req.GetId() == "" {
		return nil, fmt.Errorf("Id is required!")
	}

	err := c.compareService.DeleteByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	err = c.compareCache.InvalidateAll(ctx)
	if err != nil {
		return nil, err
	}

	taskReq := &taskPB.RemoveCompareScriptOnCascadeRequest{
		CompareScriptId: req.GetId(),
	}
	_, err = c.taskClient.RemoveCompareScriptOnCascade(ctx, taskReq)
	if err != nil {
		return nil, err
	}

	broadcastReq := &graderPB.BroadcastRequest{
		Action: graderPB.BroadcastAction_REFETCH_CONFIG,
	}
	_, err = c.graderClient.Broadcast(ctx, broadcastReq)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func initTaskGRPCClient(clientAddr string) (taskPB.TaskServiceClient, func(), error) {
	conn, err := grpc.NewClient(clientAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	client := taskPB.NewTaskServiceClient(conn)

	return client, func() {
		conn.Close()
	}, nil
}

func initGraderGRPCClient(clientAddr string) (graderPB.GraderServiceClient, func(), error) {
	conn, err := grpc.NewClient(clientAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}

	client := graderPB.NewGraderServiceClient(conn)

	return client, func() {
		conn.Close()
	}, nil
}
