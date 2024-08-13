package main

import (
	"context"
	"errors"
	"fmt"
	"github/diegoHDCz/gopet/internal/api"
	"github/diegoHDCz/gopet/internal/api/spec"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGKILL)

	defer cancel()

	if err := run(ctx); err != nil {
		log.Println(os.Stderr, err.Error())
		os.Exit(1)
	}
	fmt.Println("GoodBye :)")
}

func run(ctx context.Context) error {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	var uri = os.Getenv("MONGO_URI")
	var db = os.Getenv("DB_NAME")

	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	logger = logger.Named("petstore_app")

	defer func() { _ = logger.Sync() }()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	var result bson.M
	if err := client.Database(db).RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		panic(err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	dbClient := client.Database(db)
	si := api.NewPetAPI(dbClient, logger)

	r := gin.Default()
	r.Use(requestid.New(), gin.Recovery(), ginLogger(logger))

	spec.RegisterHandlers(r, &si)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	defer func() {
		const timeout = 30 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("failed to shutdown server", zap.Error(err))
		}
	}()

	errChan := make(chan error, 1)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return nil
}

func ginLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("Incoming request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("request_id", requestid.Get(c)),
		)
		c.Next()

	}

}
