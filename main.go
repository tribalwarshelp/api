package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/mode"

	httpdelivery "github.com/tribalwarshelp/api/graphql/delivery/http"
	"github.com/tribalwarshelp/api/graphql/resolvers"

	"github.com/tribalwarshelp/api/graphql/dataloaders"

	ennoblementrepo "github.com/tribalwarshelp/api/ennoblement/repository"
	ennoblementucase "github.com/tribalwarshelp/api/ennoblement/usecase"
	langversionrepo "github.com/tribalwarshelp/api/langversion/repository"
	langversionucase "github.com/tribalwarshelp/api/langversion/usecase"
	"github.com/tribalwarshelp/api/middleware"
	playerrepo "github.com/tribalwarshelp/api/player/repository"
	playerucase "github.com/tribalwarshelp/api/player/usecase"
	serverrepo "github.com/tribalwarshelp/api/server/repository"
	serverucase "github.com/tribalwarshelp/api/server/usecase"
	triberepo "github.com/tribalwarshelp/api/tribe/repository"
	tribeucase "github.com/tribalwarshelp/api/tribe/usecase"
	villagerepo "github.com/tribalwarshelp/api/village/repository"
	villageucase "github.com/tribalwarshelp/api/village/usecase"

	"github.com/go-pg/pg/v10"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func init() {
	os.Setenv("TZ", "UTC")

	if mode.Get() == mode.DevelopmentMode {
		godotenv.Load(".env.development")
	}
}

func main() {
	db := pg.Connect(&pg.Options{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		Addr:     os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
	})
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("Database disconnecting:", err)
		}
	}()
	//single server redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal(errors.Wrap(err, "cannot establish a connection with Redis"))
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	langversionRepo, err := langversionrepo.NewPGRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	serverRepo, err := serverrepo.NewPGRepository(db)
	if err != nil {
		log.Fatal(err)
	}
	tribeRepo := triberepo.NewPGRepository(db)
	playerRepo := playerrepo.NewPGRepository(db)
	villageRepo := villagerepo.NewPGRepository(db)
	ennoblementRepo := ennoblementrepo.NewPGRepository(db, redisClient)

	router := gin.Default()
	v1 := router.Group("")
	v1.Use(middleware.DataLoadersToContext(serverRepo, dataloaders.Config{
		PlayerRepo:  playerRepo,
		TribeRepo:   tribeRepo,
		VillageRepo: villageRepo,
	}))
	httpdelivery.Attach(httpdelivery.Config{
		RouterGroup: v1,
		Resolver: &resolvers.Resolver{
			LangVersionUcase: langversionucase.New(langversionRepo),
			ServerUcase:      serverucase.New(serverRepo),
			TribeUcase:       tribeucase.New(tribeRepo),
			PlayerUcase:      playerucase.New(playerRepo),
			VillageUcase:     villageucase.New(villageRepo),
			EnnoblementUcase: ennoblementucase.New(ennoblementRepo),
		},
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
