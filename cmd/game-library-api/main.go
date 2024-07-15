package main

import (
	"context"
	_ "expvar"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OutOfStack/game-library/internal/app/game-library-api/api"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/facade"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/model"
	"github.com/OutOfStack/game-library/internal/app/game-library-api/repo"
	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/OutOfStack/game-library/internal/auth"
	"github.com/OutOfStack/game-library/internal/client/igdb"
	"github.com/OutOfStack/game-library/internal/client/redis"
	"github.com/OutOfStack/game-library/internal/client/uploadcare"
	"github.com/OutOfStack/game-library/internal/pkg/cache"
	conf "github.com/OutOfStack/game-library/internal/pkg/config"
	"github.com/OutOfStack/game-library/internal/pkg/database"
	"github.com/OutOfStack/game-library/internal/taskprocessor"
	"github.com/go-co-op/gocron"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

// @title Game library API
// @version 0.2
// @description API for game library service
// @termsOfService http://swagger.io/terms/

// @host localhost:8000
// @BasePath /api
// @query.collection.format multi
// @schemes http
func main() {
	// load config
	var cfg appconf.Cfg
	if err := conf.Load(".", "app", "env", &cfg); err != nil {
		log.Fatalf("can't parse config: %v", err)
	}

	// init logger
	logger, err := initLogger(cfg)
	if err != nil {
		log.Fatalf("can't init logger: %v", err)
	}
	defer func(logger *zap.Logger) {
		if err = logger.Sync(); err != nil {
			log.Printf("can't sync logger: %v", err)
		}
	}(logger)

	// run
	if err = run(logger, cfg); err != nil {
		logger.Fatal("can't run app", zap.Error(err))
	}
}

func initLogger(cfg appconf.Cfg) (*zap.Logger, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	consoleWriter := zapcore.Lock(os.Stderr)
	cores := []zapcore.Core{
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), consoleWriter, zap.InfoLevel),
	}

	gelfWriter, err := gelf.NewTCPWriter(cfg.Graylog.Address)
	if err != nil {
		log.Printf("can't create gelf writer: %v", err)
	}
	if gelfWriter != nil {
		cores = append(cores,
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.AddSync(gelfWriter), zap.InfoLevel))
	}

	core := zapcore.NewTee(cores...)

	logger := zap.New(core, zap.WithCaller(false)).With(zap.String("service", appconf.ServiceName))

	return logger, nil
}

func run(logger *zap.Logger, cfg appconf.Cfg) error {
	// connect to database
	db, err := database.Open(database.Config{
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		RequireSSL: cfg.DB.RequireSSL,
	})
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer func(db *sqlx.DB) {
		if err = db.Close(); err != nil {
			logger.Error("call database close", zap.Error(err))
		}
	}(db)

	// create auth module
	authClient, err := auth.New(logger, cfg.Auth.SigningAlgorithm, cfg.Auth.VerifyTokenAPIURL)
	if err != nil {
		return fmt.Errorf("create Auth: %w", err)
	}

	// create IGDB client
	igdbClient, err := igdb.New(logger, cfg.IGDB)
	if err != nil {
		return fmt.Errorf("create IGDB client: %w", err)
	}

	// create uploadcare client
	uploadcareClient, err := uploadcare.New(logger, cfg.Uploadcare)
	if err != nil {
		return fmt.Errorf("create uploadcare client: %w", err)
	}

	// create redis client
	redisClient, err := redis.New(cfg.Redis)
	if err != nil {
		return fmt.Errorf("create redis client: %w", err)
	}

	// create redis cache service
	rCache := cache.NewRedisStore(redisClient, logger)

	// create storage
	storage := repo.New(db)

	// create game facade
	gameFacade := facade.NewProvider(logger, storage, rCache)

	// create api provider
	apiProvider := api.NewProvider(logger, rCache, gameFacade)

	// run background tasks
	taskProvider := taskprocessor.New(logger, storage, igdbClient, uploadcareClient)
	scheduler := gocron.NewScheduler(time.UTC)
	tasks := map[string]model.TaskInfo{
		taskprocessor.FetchIGDBGamesTaskName: {Schedule: cfg.Scheduler.FetchIGDBGames, Fn: taskProvider.StartFetchIGDBGames},
	}
	for name, task := range tasks {
		_, err = scheduler.Cron(task.Schedule).Name(name).Do(task.Fn)
		if err != nil {
			logger.Error("run task", zap.String("task", name), zap.Error(err))
			return fmt.Errorf("run task %s: %v", name, err)
		}
	}
	scheduler.StartAsync()

	// start debug service
	go func() {
		logger.Info("Debug service started", zap.String("address", cfg.Web.DebugAddress))
		debugService := http.Server{Addr: cfg.Web.DebugAddress, ReadTimeout: time.Second}
		err = debugService.ListenAndServe()
		if err != nil {
			logger.Error("Debug service stopped", zap.Error(err))
		}
	}()

	// start API service
	apiService, err := api.Service(logger, db, authClient, apiProvider, cfg)
	if err != nil {
		return fmt.Errorf("can't create service api: %w", err)
	}

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("API service started", zap.String("address", apiService.Addr))
		serverErrors <- apiService.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err = <-serverErrors:
		return fmt.Errorf("listening and serving: %w", err)
	case <-shutdown:
		logger.Info("Start shutdown")
		timeout := cfg.Web.ShutdownTimeout
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err = apiService.Shutdown(ctx); err != nil {
			logger.Error("Shutdown did not complete", zap.Duration("timeout", timeout), zap.Error(err))
			err = apiService.Close()
			if err != nil {
				return fmt.Errorf("shutdown: %w", err)
			}
		}
	}

	return nil
}
